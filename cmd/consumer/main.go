package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"out/internal/db"
	"out/internal/kafka"
	"out/internal/models"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

var (
	brokers         = ""
	version         = ""
	group           = ""
	topics          = ""
	assignor        = ""
	oldest          = true
	verbose         = false
	mongoHost       = ""
	mongoDatabase   = ""
	mongoCollection = ""
	mongoUser       = ""
	mongoPassword   = ""
)

func init() {
	flag.StringVar(&brokers, "brokers", os.Getenv("KAFKA_BROKERS"), "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&group, "group", os.Getenv("KAFKA_CONSUMER_GROUP"), "Kafka consumer group definition")
	flag.StringVar(&version, "version", sarama.DefaultVersion.String(), "Kafka cluster version")
	flag.StringVar(&topics, "topics", os.Getenv("KAFKA_TOPIC"), "Kafka topics to be consumed, as a comma separated list")
	flag.StringVar(&assignor, "assignor", os.Getenv("KAFKA_CONSUMER_STRATEGY"), "Consumer group partition assignment strategy (range, roundrobin, sticky)")
	flag.BoolVar(&oldest, "oldest", true, "Kafka consumer consume initial offset from oldest")
	flag.BoolVar(&verbose, "verbose", false, "Sarama logging")
	flag.StringVar(&mongoHost, "mongohost", os.Getenv("MONGODB_HOST"), "MongoDB hosts to connect to")
	flag.StringVar(&mongoDatabase, "mongodatabase", os.Getenv("MONGODB_DB"), "MongoDB database")
	flag.StringVar(&mongoCollection, "mongocollection", os.Getenv("MONGODB_COLLECTION"), "MongoDB collection")
	flag.StringVar(&mongoUser, "mongouser", os.Getenv("MONGODB_USER"), "MongoDB username")
	flag.StringVar(&mongoPassword, "mongopassword", os.Getenv("MONGODB_PASSWORD"), "MongoDB password")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topics) == 0 {
		panic("no topics given to be consumed, please set the -topics flag")
	}

	if len(group) == 0 {
		panic("no Kafka consumer group defined, please set the -group flag")
	}

}

func main() {

	db, err := db.InitMongoDB(fmt.Sprintf("mongodb://%s:%s@%s/admin", mongoUser, mongoPassword, mongoHost), mongoDatabase)
	if err != nil {
		panic(err)
	}
	defer db.Cancel()

	keepRunning := true
	log.Println("Starting a new Sarama consumer")

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumer := kafka.Consumer{
		Ready: make(chan bool),
		Fn: func(msg []byte) error {
			var s *models.Stream
			err := json.Unmarshal(msg, &s)
			if err != nil {
				log.Println(err)
				return err
			}

			err = db.AddStream(context.TODO(), mongoCollection, s)
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}

			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
