package dalle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	apiResp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         answer,
		Model:          openai.CreateImageModelDallE2,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		Size:           openai.CreateImageSize256x256,
	})
	if err != nil {
		return nil, err
	}
	if len(apiResp.Data) < 1 {
		return nil, errors.New("got no images for some reason")
	}

	imageReq, err := http.NewRequestWithContext(ctx, "GET", apiResp.Data[0].URL, nil)
	if err != nil {
		return nil, err
	}
	imageResp, err := http.DefaultClient.Do(imageReq)
	if err != nil {
		return nil, err
	}
	defer imageResp.Body.Close()
	if imageResp.StatusCode != 200 {
		return nil, fmt.Errorf("got status %d %s while retrieving image", imageResp.StatusCode, imageResp.Status)
	}
	return io.ReadAll(imageResp.Body)
}
