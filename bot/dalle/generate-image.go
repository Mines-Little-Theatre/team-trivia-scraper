package dalle

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

var ErrNoToken error = errors.New("TRIVIA_BOT_OPENAI_TOKEN not set")

func GenerateImage(ctx context.Context, answer string) (string, error) {
	authToken, ok := os.LookupEnv("TRIVIA_BOT_OPENAI_TOKEN")
	if !ok {
		return "", ErrNoToken
	}

	client := openai.NewClient(authToken)
	resp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         answer,
		Model:          openai.CreateImageModelDallE2,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		Size:           openai.CreateImageSize256x256,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Data) != 1 {
		return "", fmt.Errorf("got %d images for some reason", len(resp.Data))
	}
	return resp.Data[0].URL, nil
}
