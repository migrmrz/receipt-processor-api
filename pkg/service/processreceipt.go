package service

import "fetch-rewards.com/mx/receipt-processor/pkg/handlers/rest"

// ProcessReceipts saves receipt to a receipts array and returns its generated UUID
func (s *Service) ProcessReceipts(receipt rest.ProcessReceiptRequest) string {
	s.receipts = append(s.receipts, receipt)

	return receipt.ID
}
