package imagegen

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/genai"
)

const model = "gemini-3.1-flash-image-preview"

// GeminiClient implements Client using the Gemini image-generation API.
// Credentials are read from the GEMINI_API_KEY environment variable.
type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(ctx context.Context) (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("imagegen: GEMINI_API_KEY is not set")
	}
	c, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("imagegen: create Gemini client: %w", err)
	}
	return &GeminiClient{client: c}, nil
}

// Generate sends imageData to Gemini and returns the generated image bytes.
func (g *GeminiClient) Generate(ctx context.Context, imageData []byte, mimeType string) ([]byte, error) {

	userDescription := "I am a lizard making lizard things"

	contents := []*genai.Content{
		genai.NewContentFromParts(
			[]*genai.Part{
				genai.NewPartFromText("Generate a new similar version of this image. It is very important that the image stays the same. Add elements and modify the image to give room for this description: " + userDescription + ". The image is used for making two parallell images that you can toggle between by a slider. Therefore you cant change the picture that much. But try to follow the description of this user."),
				{InlineData: &genai.Blob{MIMEType: mimeType, Data: imageData}},
			},
			genai.RoleUser,
		),
	}

	result, err := g.client.Models.GenerateContent(ctx, model, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("imagegen: GenerateContent: %w", err)
	}

	if len(result.Candidates) == 0 || result.Candidates[0].Content == nil {
		return nil, errors.New("imagegen: no candidates returned")
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			return part.InlineData.Data, nil
		}
	}

	return nil, errors.New("imagegen: response contained no image data")
}
