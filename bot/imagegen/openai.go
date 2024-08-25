package imagegen

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
)

func generateImageWithOpenAI(ctx context.Context, answer string) (*GeneratedImage, error) {
	authToken, ok := os.LookupEnv("TRIVIA_BOT_OPENAI_TOKEN")
	if !ok {
		return nil, errNotConfigured
	}

	client := openai.NewClient(authToken)
	apiResp, err := client.CreateImage(ctx, openai.ImageRequest{
		Prompt:         answer,
		Model:          openai.CreateImageModelDallE3,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		Size:           openai.CreateImageSize1024x1024,
	})
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}
	if len(apiResp.Data) < 1 {
		return nil, errors.New("got no images for some reason")
	}

	log.Println("openai: image url is", apiResp.Data[0].URL)
	imageReq, err := http.NewRequestWithContext(ctx, "GET", apiResp.Data[0].URL, nil)
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}

	imageResp, err := http.DefaultClient.Do(imageReq)
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}
	defer imageResp.Body.Close()

	if imageResp.StatusCode != 200 {
		return nil, fmt.Errorf("openai: got status %s while retrieving image", imageResp.Status)
	}

	pngData, err := io.ReadAll(imageResp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}

	return &GeneratedImage{PNGData: pngData, ModelName: "DALL·E 3"}, nil
}
