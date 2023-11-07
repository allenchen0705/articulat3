package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (h *PromptHandler) CreatePrompt(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ctx := r.Context()
	logger := log.Ctx(ctx)
	logger.Trace().Msg("create prompt request received")

	var req PromptRequest
	reqId := requestIdFromContext(ctx)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, r, http.StatusBadRequest,
			fmt.Errorf("error decoding the request: %v", err))
		return
	}

	prompt, _ := req.ToPrompt()
	newPrompt, _ := h.ctrl.PromptCreate(ctx, prompt)
	res := promptResponseFromPrompt(newPrompt, reqId)
	writeResponse(w, r, http.StatusCreated, res)
}
