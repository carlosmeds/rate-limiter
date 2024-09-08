package web

import (
	"fmt"
	"net/http"
)

type WebIpHandler struct {
}

func NewWebIpHandler() *WebIpHandler {
	return &WebIpHandler{}
}

func (h *WebIpHandler) Get(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	apiKey := r.Header.Get("API_KEY")

	fmt.Printf("GET /ip called from IP: %s with API_KEY: %s\n", clientIP, apiKey)
}
