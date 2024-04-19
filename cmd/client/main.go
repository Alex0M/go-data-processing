package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"out/internal/client"
	"out/internal/models"
	"syscall"
)

var (
	clientNumber = 0
	randomMode   = true
	configFile   = ""
	frequency    = 1
	reportingAPI = ""
	token        = ""
)

func init() {
	flag.StringVar(&reportingAPI, "reportingAPI", os.Getenv("REPORTING_API"), "reporting API URL")
	flag.StringVar(&configFile, "configFile", os.Getenv("CONFIG_FILE"), "path to config file")
	flag.StringVar(&token, "token", os.Getenv("REPORTING_API_TOKEN"), "token to conect to reporting API")
	flag.IntVar(&clientNumber, "clientNumber", 0, "set number of clients")
	flag.IntVar(&frequency, "frequency", 1, "set frequency to send client report")
	flag.BoolVar(&randomMode, "randomMode", true, "mode to run")
	flag.Parse()

	if len(token) == 0 {
		panic("no tokec for reporting API URL defined, please set the -token flag")
	}

	if randomMode && len(configFile) == 0 {
		configFile = "./config/random_data.json"
	}

	if !randomMode && len(configFile) == 0 {
		configFile = "./config/client_data.json"
	}

	if len(reportingAPI) == 0 {
		panic("no reporting API URL defined, please set the -reportingAPI flag")
	}
}

func main() {
	configData := []*models.Stream{}
	c := client.NewClient(clientNumber, frequency, configData, randomMode, reportingAPI, token)

	if randomMode {
		err := client.ReadConfigRandomData(configFile)
		if err != nil {
			log.Fatalf("cannot read config file: %s", err)
		}
		c.Run()
	} else {
		d, err := client.ReadClientDataFromFile(configFile)
		if err != nil {
			log.Fatalf("cannot read config file: %s", err)
		}
		c.ClientNumber = len(d)
		c.ConfigData = d
		c.Run()
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	fmt.Println("done...")
}
