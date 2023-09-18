package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	gokithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type getPointsRequest struct {
	receiptID string
}

type getPointsResponse struct {
	statusCode int
	Error      string      `json:"error,omitempty"`
	Points     interface{} `json:"points"`
}

// StatusCode func to implement StatusCoder interface
func (gr getPointsResponse) StatusCode() int {
	return gr.statusCode
}

func makeGetPointsEndpoint(s ReceiptProcessorHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getPointsRequest)

		points, _ := s.GetPoints(req.receiptID)

		return points, nil
	}
}

// decode request function
func decodeGetPointsRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	receiptID, err := uuid.Parse(chi.URLParam(req, "id"))
	if err != nil {
		return nil, ErrBadRequest{message: "invalid receipt id"}
	}

	getPointsRequest := getPointsRequest{
		receiptID: receiptID.String(),
	}

	return getPointsRequest, nil
}

// encode response and error functions
func encodeProcessReceiptResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := processReceiptResponse{
		statusCode: http.StatusOK,
		ID:         response,
	}

	return gokithttp.EncodeJSONResponse(ctx, w, resp)
}

func encodeGetPointsError(ctx context.Context, err error, w http.ResponseWriter) {
	resp := getPointsResponse{
		Points:     nil,
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
