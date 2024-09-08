package usecase

import "fmt"

type GetIpDTO struct {
	ClientIP string
	ApiKey   string
}

type GetIpUseCase struct {
}

func NewGetIpUseCase() *GetIpUseCase {
	return &GetIpUseCase{}
}

func (c *GetIpUseCase) Execute(i *GetIpDTO) {
	fmt.Printf("GET /ip called from IP: %s with API_KEY: %s\n", i.ClientIP, i.ApiKey)
}
