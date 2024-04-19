package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/gocql/gocql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *bun.DB
}

type CassandraDB struct {
	Session *gocql.Session
}

type MongoDB struct {
	DB     *mongo.Database
	Cancel context.CancelFunc
}

func Init(dsn string) *Database {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return &Database{
		Client: db,
	}
}

func InitCassandra(h string, ks string) (*CassandraDB, error) {
	cluster := gocql.NewCluster(h)
	cluster.Keyspace = ks
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &CassandraDB{
		Session: session,
	}, nil
}

func InitMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoDB{
		DB:     db,
		Cancel: cancel,
	}, nil
}
