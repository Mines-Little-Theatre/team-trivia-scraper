package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func executeWebhook(config *Config, body *io.PipeReader, mpartBoundary string) {
	req, err := http.NewRequest("POST", config.WebhookURL, body)
	if err != nil {
		log.Fatalf("failed to create Discord request: %v", err)
	}

	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data;boundary=%s", mpartBoundary))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("failed to do request:", err)
	} else if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatalf("API call failed (%s) and reading the response also failed: %v", resp.Status, err)
		} else {
			log.Fatalf("API call failed (%s) and responded: %s", resp.Status, body)
		}
	}
	resp.Body.Close()
}
