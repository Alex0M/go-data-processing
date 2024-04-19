package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"out/internal/controllers/auth"
	"out/internal/controllers/tracking"
	"out/internal/kafka"
	"out/internal/middleware"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	brokers = ""
	topic   = ""
)

func init() {
	flag.StringVar(&brokers, "brokers", os.Getenv("KAFKA_BROKERS"), "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&topic, "topics", os.Getenv("KAFKA_TOPIC"), "Kafka topic")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if topic == "" {
		panic("no topics given to be consumed, please set the -topics flag")
	}
}

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error creating logger: %s", err)
	}

	p, err := kafka.NewDataCollector(strings.Split(brokers, ","))
	if err != nil {
		log.Fatalf("failed to start Sarama producer: %s", err)
	}

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(middleware.ErrorHandler(logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	tracking.RegisterRoutes(r, p, topic)
	auth.RegisterRoutes(r)

	//log.Fatal(r.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))))
	log.Fatal(r.Run(fmt.Sprintf("0.0.0.0:%d", 9099)))
}
