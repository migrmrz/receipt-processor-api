package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHTTPHandlers(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		url            string
		data           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "process-receipt-target-ok",
			method: http.MethodPost,
			url:    "/receipt-processor/v1/receipts/process",
			data: `{"retailer":"Target","purchaseDate":"2022-01-01","purchaseTime":"13:01","items":[{"` +
				`shortDescription":"Mountain Dew 12PK","price": "6.49"},{"shortDescription":"Emils Cheese Pizza",` +
				`"price":"12.25"},{"shortDescription":"Knorr Creamy Chicken","price":"1.26"},{"shortDescription":` +
				`"Doritos Nacho Cheese","price":"3.35"},{"shortDescription":"Klarbrunn 12-PK 12 FL OZ","price":` +
				`"12.00"}],"total": "35.35"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"ddbfe2da-b86b-49db-9aba-168e3397286a"}`,
		},
		{
			name:           "get-receipt-target-ok",
			method:         http.MethodGet,
			url:            "/receipt-processor/v1/receipts/ddbfe2da-b86b-49db-9aba-168e3397286a/points",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"points":28}`,
		},
	}

	mockHandler := mockHandler{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := MakeHTTPHandlers(mockHandler)
			assert.NotNil(t, handler)

			// create request
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.data))
			assert.NoError(t, err)

			reqRec := httptest.NewRecorder()

			handler.ServeHTTP(reqRec, req)
			assert.Equal(t, tc.expectedStatus, reqRec.Code)

			actualBody := strings.TrimRight(reqRec.Body.String(), "\n")
			assert.Equal(t, tc.expectedBody, actualBody)
		})
	}
}

type mockHandler struct{}

func (m mockHandler) ProcessReceipts(receipt ProcessReceiptRequest) string {
	switch receipt.Retailer {
	case "Target":
		return "ddbfe2da-b86b-49db-9aba-168e3397286a"
	case "M&M Corner Market":
		return "2a4ac2bc-ce84-4a62-bc47-1c3a4984cb02"
	default:
		return "31311cf8-f18a-4111-a6d6-356469af5605"
	}
}

func (m mockHandler) GetPoints(receiptID string) (int, []string) {
	switch receiptID {
	case "ddbfe2da-b86b-49db-9aba-168e3397286a":
		return 28, nil
	case "2a4ac2bc-ce84-4a62-bc47-1c3a4984cb02":
		return 109, nil
	case "31311cf8-f18a-4111-a6d6-356469af5605":
		return 1, nil
	default:
		return 0, nil
	}
}
