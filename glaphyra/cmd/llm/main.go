package main

import (
	"fmt"
	"log"
	"reflect"

	yaml_config "glaphyra/config"
	"glaphyra/internal/llm"
	"glaphyra/internal/llm/handlers"
	"glaphyra/internal/llm/yagpt/dto"
)

// пример использования, потом удалить
func main() {
	config, err := yaml_config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	orchestrator := llm.NewOrchestrator()

	// Добавляем хэндлеры
	orchestrator.AddHandler("YaGPT", handlers.NewYaGPTHandler(config))

	yaRequest := dto.RequestDTO{
		Model:         "gpt://b1g788lsg3ebu9vohhji/yandexgpt-lite", // хардкод, но вроде нормально
		Temperature:   0.6,
		MaxTokens:     2000,
		SystemMessage: "Ты умный ассистент",
		UserMessage:   "Какими науками занимался Альберт Эйнштейн?",
	}

	yaResponse, err := orchestrator.CallAPI("YaGPT", yaRequest)
	if err != nil {
		log.Printf("API: %v", err)
	}

	response, ok := yaResponse.(dto.ResponseDTO)
	if !ok {
		log.Printf("Invalid response type %v", reflect.TypeOf(yaResponse))
	}

	fmt.Printf("YaGPT Response: %s\n", response.Result)
}
