package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/yomorun/yomo/serverless"
)

var ImgDir = "./"

// Implement DataTags() to observe data with the given tags
func DataTags() []uint32 {
	return []uint32{0x11}
}

// Implement Init() for state initialization, such as loading LLM Model to GPU memory.
func Init() error {
	if v, ok := os.LookupEnv("IMG_DIR"); ok {
		ImgDir = v
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
	return "A function that gets an image from the /capture endpoint."
}

// Implement InputSchema() to define the input schema of the function
func InputSchema() any {
	return &Parameter{}
}

// Implement Handler() to handle the function call
func Handler(ctx serverless.Context) {
	result := "please see the picture image.png"
	defer ctx.WriteLLMResult(result)

	img, err := get("http://localhost:30080/deviceshifu-camera/capture")
	if err != nil {
		fmt.Println(err)
		result = err.Error()
		return
	}

	err = os.WriteFile(path.Join(ImgDir, "image.png"), img, 0644)
	if err != nil {
		fmt.Println(err)
		result = err.Error()
		return
	}
}

func get(url string) ([]byte, error) {
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
