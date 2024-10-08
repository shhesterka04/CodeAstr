package llm

import (
	"github.com/pkg/errors"

	"glaphyra/internal/llm/handlers"
)

// Orchestrator управляет всеми ручками для работы с разными апи
type Orchestrator struct {
	handlers map[string]handlers.Handler
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		handlers: make(map[string]handlers.Handler),
	}
}

// AddHandler добавляет нового клиента
func (o *Orchestrator) AddHandler(clientName string, handler handlers.Handler) {
	o.handlers[clientName] = handler
}

func (o *Orchestrator) CallAPI(clientName string, request interface{}) (interface{}, error) {
	h, exists := o.handlers[clientName]
	if !exists {
		return nil, errors.Wrap(errors.New("handler not found"), clientName)
	}

	return h.CallAPI(request)
}
