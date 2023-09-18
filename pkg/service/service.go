package service

import "fetch-rewards.com/mx/receipt-processor/pkg/handlers/rest"

type Service struct {
	receipts []rest.ProcessReceiptRequest
}

func New() *Service {
	return &Service{
		receipts: []rest.ProcessReceiptRequest{},
	}
}
