package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/VladislavZhr/highload-workflow/connector/internal/model"
	"github.com/VladislavZhr/highload-workflow/connector/internal/service"
)

type HTTPHandler struct {
	connectorService *service.ConnectorService
}

func NewHTTPHandler(connectorService *service.ConnectorService) *HTTPHandler {
	return &HTTPHandler{
		connectorService: connectorService,
	}
}

func (h *HTTPHandler) HandleProcessRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST method is allowed")
		return
	}

	defer r.Body.Close()

	var req model.Request

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		handleDecodeError(w, err)
		return
	}

	if decoder.More() {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must contain only one JSON object")
		return
	}

	_, err := h.connectorService.Process(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			writeError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		case errors.Is(err, service.ErrKafkaProduce):
			writeError(w, http.StatusInternalServerError, "kafka_produce_error", "failed to send message to kafka")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			return
		}
	}

	resp := model.SuccessResponse{
		Status:  "success",
		Message: "message was successfully sent to kafka",
	}

	writeJSON(w, http.StatusOK, resp)
}

func handleDecodeError(w http.ResponseWriter, err error) {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError

	switch {
	case errors.Is(err, io.EOF):
		writeError(w, http.StatusBadRequest, "empty_body", "request body is empty")
	case errors.As(err, &syntaxErr):
		writeError(w, http.StatusBadRequest, "invalid_json", "request body contains malformed JSON")
	case errors.As(err, &unmarshalTypeErr):
		writeError(w, http.StatusBadRequest, "invalid_json_type", "request body contains invalid field type")
	default:
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	resp := model.ErrorResponse{
		Code:    code,
		Message: message,
	}

	writeJSON(w, status, resp)
}
