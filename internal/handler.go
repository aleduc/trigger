package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"test_trigger/internal/call"
	"test_trigger/internal/logger"
)

//go:generate go run github.com/golang/mock/mockgen --source=handler.go --destination=handler_mock.go --package=internal

// CallSaver is responsible for saving calls for later processing.
type CallSaver interface {
	AddToQueueBack(_ context.Context, meta call.Meta) error
}

// TriggerResponse response struct for /trigger request.
type TriggerResponse struct {
	CallID string `json:"call_id"`
}

// Server is responsible for handling requests.
type Server struct {
	callSaver CallSaver
	getUUID   func() string // decided to save time there.
	logger    logger.Logger
}

func NewServer(callSaver CallSaver, getUUID func() string, logger logger.Logger) *Server {
	return &Server{callSaver: callSaver, getUUID: getUUID, logger: logger}
}

// Trigger processes http request, save correct body to storage for later processing.
func (s *Server) Trigger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	callBody := &call.Body{}
	err := json.NewDecoder(r.Body).Decode(callBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if callBody.VirtualAgentID == "" || callBody.PhoneNumber == "" {
		http.Error(w, "virtual_agent_id or phone_number can't be empty", http.StatusBadRequest)
		return
	}
	callID := s.getUUID()
	err = s.callSaver.AddToQueueBack(r.Context(), call.Meta{
		PhoneNumber:    callBody.PhoneNumber,
		VirtualAgentID: callBody.VirtualAgentID,
		ID:             call.ID(callID),
	})
	if err != nil {
		s.logger.Error(fmt.Errorf("trigger: AddToQueueBack: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := TriggerResponse{CallID: callID}
	respBody, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error(fmt.Errorf("trigger: marshall: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		s.logger.Error(fmt.Errorf("trigger: write bytes: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
