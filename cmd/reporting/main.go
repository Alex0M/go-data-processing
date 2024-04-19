package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"out/internal/controllers/auth"
	"out/internal/controllers/reporting"
	"out/internal/db"
	"out/internal/middleware"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	mongoHost       = ""
	mongoDatabase   = ""
	mongoCollection = ""
	mongoUser       = ""
	mongoPassword   = ""
)

func init() {
	flag.StringVar(&mongoHost, "mongohost", os.Getenv("MONGODB_HOST"), "MongoDB hosts to connect to")
	flag.StringVar(&mongoDatabase, "mongodatabase", os.Getenv("MONGODB_DB"), "MongoDB database")
	flag.StringVar(&mongoCollection, "mongocollection", os.Getenv("MONGODB_COLLECTION"), "MongoDB collection")
	flag.StringVar(&mongoUser, "mongouser", os.Getenv("MONGODB_USER"), "MongoDB username")
	flag.StringVar(&mongoPassword, "mongopassword", os.Getenv("MONGODB_PASSWORD"), "MongoDB password")

	flag.Parse()
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error creating logger: %s", err)
	}

	db, err := db.InitMongoDB(fmt.Sprintf("mongodb://%s:%s@%s/admin", mongoUser, mongoPassword, mongoHost), mongoDatabase)
	if err != nil {
		panic(err)
	}
	defer db.Cancel()

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(middleware.ErrorHandler(logger))

	reporting.RegisterRoutes(r, db, mongoCollection)
	auth.RegisterRoutes(r)

	log.Fatal(r.Run(fmt.Sprintf("0.0.0.0:%d", 9098)))
}
