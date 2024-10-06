package yagpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"glaphyra/config"
	"glaphyra/internal/llm/yagpt/dto"
)

// пока хардкод мб в конфиг надо
const url = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"

type YaGPTClient struct {
	config *config.Config
}

func NewYaGPTClient(config *config.Config) *YaGPTClient {
	return &YaGPTClient{config: config}
}

func (c *YaGPTClient) CallAPI(reqDTO dto.RequestDTO) (dto.ResponseDTO, error) {
	log.Printf("Calling YaGPT API with request: %v", reqDTO)

	requestBodyBytes, err := c.buildRequestBody(reqDTO)
	if err != nil {
		log.Printf("build request body: %s", err)
		return dto.ResponseDTO{}, err
	}

	req, err := c.buildHTTPRequest(requestBodyBytes)
	if err != nil {
		log.Printf("create request: %s", err)
		return dto.ResponseDTO{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("send request: %s", err)
		return dto.ResponseDTO{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("response code: %v Response body: %s", resp.StatusCode, string(body))
		return dto.ResponseDTO{}, errors.New("received non-200 response from YaGPT API")
	}

	return c.processResponse(resp)
}

func (c *YaGPTClient) buildRequestBody(reqDTO dto.RequestDTO) ([]byte, error) {
	requestBody := map[string]interface{}{
		"modelUri": reqDTO.Model,
		"completionOptions": map[string]interface{}{
			"temperature": reqDTO.Temperature,
			"maxTokens":   reqDTO.MaxTokens,
		},
		"messages": []map[string]string{
			{"role": "system", "text": reqDTO.SystemMessage},
			{"role": "user", "text": reqDTO.UserMessage},
		},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	return requestBodyBytes, nil
}

func (c *YaGPTClient) buildHTTPRequest(body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.YagptKey))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *YaGPTClient) processResponse(resp *http.Response) (dto.ResponseDTO, error) {
	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("decode response: %s", err)
		return dto.ResponseDTO{}, err
	}

	result := apiResponse["result"].(map[string]interface{})
	usage := result["usage"].(map[string]interface{})
	alternatives := result["alternatives"].([]interface{})[0].(map[string]interface{})
	message := alternatives["message"].(map[string]interface{})

	responseDTO := dto.ResponseDTO{
		Result:           message["text"].(string),
		InputTextTokens:  usage["inputTextTokens"].(string),
		CompletionTokens: usage["completionTokens"].(string),
		TotalTokens:      usage["totalTokens"].(string),
		ModelVersion:     result["modelVersion"].(string),
	}

	log.Printf("Received response from YaGPT: %v", responseDTO)
	return responseDTO, nil
}
