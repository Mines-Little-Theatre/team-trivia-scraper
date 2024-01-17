package dalle

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

var ErrNoToken error = errors.New("TRIVIA_BOT_OPENAI_TOKEN not set")

func GenerateImage(ctx context.Context, answer string) ([]byte, error) {
	authToken, ok := os.LookupEnv("TRIVIA_BOT_OPENAI_TOKEN")
	if !ok {
		return nil, ErrNoToken
	}

	client := openai.NewClient(authToken)
	resp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         answer,
		Model:          openai.CreateImageModelDallE2,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		Size:           openai.CreateImageSize256x256,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) < 1 {
		return nil, errors.New("got no images for some reason")
	}
	return base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
}
