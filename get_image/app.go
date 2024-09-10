package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/yomorun/yomo/serverless"
)

// Implement DataTags() to observe data with the given tags
func DataTags() []uint32 {
	return []uint32{0x11}
}

// Implement Init() for state initialization, such as loading LLM Model to GPU memory.
func Init() error {
	return nil
}

// Parameters needed for OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
type Parameter struct {
}

// Implement Description() to define the description of OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
func Description() string {
	return "A function that gets an image from the capture camera."
}

// Implement InputSchema() to define the input schema of the function
func InputSchema() any {
	return &Parameter{}
}

// Implement Handler() to handle the function call
func Handler(ctx serverless.Context) {
	var result string
	defer ctx.WriteLLMResult(result)

	body, err := httpGet("http://localhost:30080/deviceshifu-camera/capture")
	if err != nil {
		result = "an error occurred: " + err.Error()
		return
	}

	image := base64.StdEncoding.EncodeToString(body)

	c := openai.DefaultConfig(os.Getenv("VIVGRID_TOKEN"))
	c.BaseURL = "https://openai.vivgrid.com/v1"
	client := openai.NewClientWithConfig(c)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			// Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "user",
					MultiContent: []openai.ChatMessagePart{
						{
							Type: "text",
							Text: "please describe the given image",
						},
						{
							Type: "img_url",
							ImageURL: &openai.ChatMessageImageURL{
								URL: image,
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		result = "an error occurred: " + err.Error()
		return
	}

	result = response.Choices[0].Message.Content
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("http status:", resp.Status)
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
