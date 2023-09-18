package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fetch-rewards.com/mx/receipt-processor/pkg/handlers/rest"
)

func TestProcessReceipt(t *testing.T) {
	cases := []struct {
		name                string
		receipt             rest.ProcessReceiptRequest
		expectedResponse    string
		expectedReceiptsLen int
	}{
		{
			name: "process-target-ok",
			receipt: rest.ProcessReceiptRequest{
				ID:           "77b341b0-24db-4188-9954-ff353073f09a",
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Total:        "35.35",
				Items: []rest.Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            "6.49",
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            "12.25",
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            "1.26",
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            "3.35",
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            "12.00",
					},
				},
			},
			expectedResponse:    "77b341b0-24db-4188-9954-ff353073f09a",
			expectedReceiptsLen: 1,
		},
		{
			name: "process-m&m-corner-ok",
			receipt: rest.ProcessReceiptRequest{
				ID:           "6dd1a55a-48bc-460d-8643-4d616358acd1",
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Total:        "9.00",
				Items: []rest.Item{
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
				},
			},
			expectedResponse:    "6dd1a55a-48bc-460d-8643-4d616358acd1",
			expectedReceiptsLen: 2,
		},
	}

	// init service
	srv := New()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualResponse := srv.ProcessReceipts(tc.receipt)
			assert.Equal(t, tc.expectedResponse, actualResponse)
			assert.Equal(t, tc.expectedReceiptsLen, len(srv.receipts))
		})
	}
}
