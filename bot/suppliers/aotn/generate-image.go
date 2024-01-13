package aotn

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func generateImage(answer, authToken string) (string, error) {
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
