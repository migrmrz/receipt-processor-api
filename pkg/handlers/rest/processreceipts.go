package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	gokithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// processReceiptsRequest defines the request message to process receipts
type ProcessReceiptRequest struct {
	ID           string
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type processReceiptResponse struct {
	statusCode int
	Error      string      `json:"error,omitempty"`
	ID         interface{} `json:"id"`
}

// StatusCode func to implement StatusCoder interface
func (pr processReceiptResponse) StatusCode() int {
	return pr.statusCode
}

func makeProcessReceiptsEndpoint(s ReceiptProcessorHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ProcessReceiptRequest)

		receiptID := s.ProcessReceipts(req)

		return receiptID, nil
	}
}

// decode request function
func decodeProcessReceiptsRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, ErrBadRequest{"can't read body"}
	}

	var parsedReq ProcessReceiptRequest

	err = json.Unmarshal(reqBody, &parsedReq)
	if err != nil {
		return nil, ErrBadRequest{"can't parse body"}
	}

	// generating a UUID for the processed receipt
	parsedReq.ID = uuid.New().String()

	return parsedReq, nil
}

// encode response and error funcions
func encodeGetPointsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := getPointsResponse{
		statusCode: http.StatusOK,
		Points:     response,
	}

	return gokithttp.EncodeJSONResponse(ctx, w, resp)
}

func encodeProcessReceiptError(ctx context.Context, err error, w http.ResponseWriter) {
	resp := processReceiptResponse{
		ID:         nil,
		Error:      err.Error(),
		statusCode: http.StatusInternalServerError,
	}

	if statusCoder, ok := err.(gokithttp.StatusCoder); ok {
		resp.statusCode = statusCoder.StatusCode()
	}

	if encodeErr := gokithttp.EncodeJSONResponse(ctx, w, resp); encodeErr != nil {
		logrus.WithFields(logrus.Fields{
			"func":  "errorEncoder",
			"step":  "EncodeJSONRepsonse",
			"error": encodeErr,
		}).Error()

		gokithttp.DefaultErrorEncoder(ctx, encodeErr, w)
	}
}
