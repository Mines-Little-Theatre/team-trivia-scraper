package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type Config struct {
	WebhookURL          string
	FreeAnswerMessage   string
	NoFreeAnswerMessage string
}

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "config", "team-trivia-scraper-config.json", "Location of configuration file")
	flag.Parse()

	var config Config
	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatalln("failed to open config file:", err)
	}
	configData, err := io.ReadAll(configFile)
	configFile.Close()
	if err != nil {
		log.Fatalln("failed to read config file:", err)
	}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalln("failed to unmarshal config file:", err)
	}

	body, bodyWriter := io.Pipe()
	mpartWriter := multipart.NewWriter(bodyWriter)
	mpartBoundary := mpartWriter.Boundary()

	go writeMessage(&config, bodyWriter, mpartWriter)
	executeWebhook(&config, body, mpartBoundary)
}
