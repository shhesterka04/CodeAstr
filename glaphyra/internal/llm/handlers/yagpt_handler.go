package handlers

import (
	"github.com/pkg/errors"

	"glaphyra/config"
	yagpt "glaphyra/internal/llm/yagpt/client"
	"glaphyra/internal/llm/yagpt/dto"
)

type YaGPTHandler struct {
	client *yagpt.YaGPTClient
}

func NewYaGPTHandler(config *config.Config) *YaGPTHandler {
	return &YaGPTHandler{
		client: yagpt.NewYaGPTClient(config),
	}
}

func (h *YaGPTHandler) CallAPI(request interface{}) (interface{}, error) {
	req, ok := request.(dto.RequestDTO)
	if !ok {
		return nil, errors.New("invalid request type for YaGPT")
	}

	return h.client.CallAPI(req)
}
