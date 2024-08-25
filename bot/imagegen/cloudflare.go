package imagegen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cloudflareAPIRequest struct {
	Prompt string `json:"prompt"`
}

func generateImageWithCloudflare(ctx context.Context, prompt string) (*GeneratedImage, error) {
	accountId, ok := os.LookupEnv("TRIVIA_BOT_CLOUDFLARE_ACCOUNT_ID")
	if !ok {
		return nil, errNotConfigured
	}

	apiToken, ok := os.LookupEnv("TRIVIA_BOT_CLOUDFLARE_API_TOKEN")
	if !ok {
		return nil, errNotConfigured
	}

	requestBody, err := json.Marshal(cloudflareAPIRequest{
		Prompt: prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("cloudflare: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/@cf/stabilityai/stable-diffusion-xl-base-1.0", accountId), bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("cloudflare: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloudflare: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cloudflare: got status %s and couldn't read body: %w", resp.Status, err)
		}
		return nil, fmt.Errorf("cloudflare: got status %s: %s", resp.Status, content)
	}

	pngData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cloudflare: %w", err)
	}

	return &GeneratedImage{PNGData: pngData, ModelName: "stable-diffusion-xl-base-1.0"}, nil
}
