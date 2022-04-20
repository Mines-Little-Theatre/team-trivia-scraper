package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func executeWebhook(config *Config, body *io.PipeReader, mpartBoundary string) {
	req, err := http.NewRequest("POST", config.WebhookURL, body)
	if err != nil {
		log.Fatalln("failed to create Discord request:", err)
	}

	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data;boundary=%s", mpartBoundary))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("failed to do request:", err)
	} else if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatalln("API call failed with code", resp.StatusCode, "and reading the response also failed:", err)
		} else {
			log.Fatalln("API call failed with code", resp.StatusCode, "and responded:", string(body))
		}
	}
	resp.Body.Close()
}
