package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/yomorun/yomo/serverless"
)

// Implement DataTags() to observe data with the given tags
func DataTags() []uint32 {
	return []uint32{0x13}
}

// Implement Init() for state initialization, such as loading LLM Model to GPU memory.
func Init() error {
	return nil
}

// Parameters needed for OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
type Parameter struct {
	Number int `json:"num" jsonschema:"description=a number to display on the LED, between 0 and 9999"`
}

// Implement Description() to define the description of OpenAI Function Calling
// ref: https://platform.openai.com/docs/guides/function-calling
func Description() string {
	return "A function that sets the display number of the LED."
}

// Implement InputSchema() to define the input schema of the function
func InputSchema() any {
	return &Parameter{}
}

// Implement Handler() to handle the function call
func Handler(ctx serverless.Context) {
	var result string
	defer ctx.WriteLLMResult(result)

	var msg Parameter
	err := ctx.ReadLLMArguments(&msg)
	if err != nil {
		result = "an error occurred: " + err.Error()
		return
	}

	// buf, err := json.Marshal(struct)

	_, err = httpPost("http://localhost:30080/deviceshifu-led/number")
	if err != nil {
		result = "an error occurred: " + err.Error()
		return
	}

	result = ""
}

func httpPost(url string, body string) ([]byte, error) {
	resp, err := http.Post(url, "application/json", nil)
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
