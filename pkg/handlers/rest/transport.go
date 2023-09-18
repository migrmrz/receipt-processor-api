package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	gokithttp "github.com/go-kit/kit/transport/http"
)

// Defining endpoints for the service
type Endpoints struct {
	ProcessReceiptsEndpoint endpoint.Endpoint
	GetPointsEndpoint       endpoint.Endpoint
}

type ReceiptProcessorHandler interface {
	ProcessReceipts(ProcessReceiptRequest) string
	GetPoints(string) (int, []string)
}

func MakeHTTPHandlers(s ReceiptProcessorHandler) http.Handler {
	r := chi.NewRouter()

	version1RouterFunc := makeVersion1RouterFunc(s)
	r.Route("/receipt-processor/v1", version1RouterFunc)

	return r
}

func makeVersion1RouterFunc(s ReceiptProcessorHandler) func(chi.Router) {
	endpoints := makeServerEndpoints(s)

	receiptProcessorServer := gokithttp.NewServer(
		endpoints.ProcessReceiptsEndpoint,
		decodeProcessReceiptsRequest,
		encodeProcessReceiptResponse,
		gokithttp.ServerErrorEncoder(encodeProcessReceiptError),
	)

	pointsGetterServer := gokithttp.NewServer(
		endpoints.GetPointsEndpoint,
		decodeGetPointsRequest,
		encodeGetPointsResponse,
		gokithttp.ServerErrorEncoder(encodeGetPointsError),
	)

	return func(r chi.Router) {
		r.Get("/receipts/{id}/points", pointsGetterServer.ServeHTTP)
		r.Post("/receipts/process", receiptProcessorServer.ServeHTTP)
	}
}

func makeServerEndpoints(s ReceiptProcessorHandler) Endpoints {
	return Endpoints{
		ProcessReceiptsEndpoint: makeProcessReceiptsEndpoint(s),
		GetPointsEndpoint:       makeGetPointsEndpoint(s),
	}
}
