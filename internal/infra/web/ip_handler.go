package web

import (
	"net/http"

	"github.com/carlosmeds/rate-limiter/internal/usecase"
)

type WebIpHandler struct {
}

func NewWebIpHandler() *WebIpHandler {
	return &WebIpHandler{}
}

func (h *WebIpHandler) Get(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	apiKey := r.Header.Get("API_KEY")

	uc := usecase.NewGetIpUseCase()
	uc.Execute(&usecase.GetIpDTO{
		ClientIP: clientIP,
		ApiKey:   apiKey,
	})
}
