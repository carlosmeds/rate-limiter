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
	fmt.Println("GET /ip called")
}
