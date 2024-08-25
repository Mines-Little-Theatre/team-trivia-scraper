package imagegen

import (
	"context"
	"errors"
	"log"
)

type GeneratedImage struct {
	PNGData   []byte
	ModelName string
}

type imageGenerator func(ctx context.Context, prompt string) (*GeneratedImage, error)

var imageGenerators = [...]imageGenerator{
	generateImageWithOpenAI,
	generateImageWithCloudflare,
}

var errNotConfigured error = errors.New("not configured")

func GenerateImage(ctx context.Context, prompt string) (*GeneratedImage, error) {
	var errs []error
	for _, gen := range imageGenerators {
		result, err := gen(ctx, prompt)
		if err == nil {
			return result, nil
		} else if !errors.Is(err, errNotConfigured) {
			log.Println(err)
			errs = append(errs, err)
		}
	}
	return nil, errors.Join(errs...)
}
