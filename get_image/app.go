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

var apiKey string

// Implement DataTags() to observe data with the given tags
func DataTags() []uint32 {
	return []uint32{0x11}
}

// Implement Init() for state initialization, such as loading LLM Model to GPU memory.
func Init() error {
	if v, ok := os.LookupEnv("VIVGRID_TOKEN"); ok {
		apiKey = v
	}
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
	fmt.Println("start running handler")
	ch := make(chan string)

	go func() {
		body, err := httpGet("http://localhost:30080/deviceshifu-camera/capture")
		if err != nil {
			ch <- "error: " + err.Error()
			return
		}

		image := base64.StdEncoding.EncodeToString(body)

		config := openai.DefaultConfig(apiKey)
		config.BaseURL = "https://openai.vivgrid.com/v1"
		client := openai.NewClientWithConfig(config)

		response, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Messages: []openai.ChatCompletionMessage{
					{
						Role: "user",
						MultiContent: []openai.ChatMessagePart{
							{
								Type: "text",
								Text: "Thanks! Can you tell me what is in the image? Specifically what is the number on the display and does the PLC has 4 output lights on? Please return in json format, like {\"display_number\":2929,\"plc_switch\":true}.",
							},
							{
								Type: "image_url",
								ImageURL: &openai.ChatMessageImageURL{
									URL: "data:image/jpeg;base64," + image,
								},
							},
						},
					},
				},
				// ResponseFormat: &openai.ChatCompletionResponseFormat{
				// 	Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				// 	JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				// 		Name:        "image_info",
				// 		Description: "the detail info from the captured image",
				// 		Schema: jsonschema.Definition{
				// 			Type: jsonschema.Object,
				// 			Properties: map[string]jsonschema.Definition{
				// 				"number": {
				// 					Type:        jsonschema.Integer,
				// 					Description: "the number on the display",
				// 				},
				// 				"switch": {
				// 					Type:        jsonschema.Boolean,
				// 					Description: "whether the PLC outputs lights are ON or OFF",
				// 				},
				// 			},
				// 			Required:             []string{"number", "switch"},
				// 			AdditionalProperties: false,
				// 		},
				// 		Strict: true,
				// 	},
				// },
			},
		)
		if err != nil {
			ch <- "error: " + err.Error()
			return
		}

		res := response.Choices[0].Message.Content

		ch <- "As a virtual camera assitant, I have taken an image for you. And the captured image shows the following information:\n" + res
	}()

	for res := range ch {
		fmt.Println("res:", res)
		ctx.WriteLLMResult(res)
	}
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status)
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
