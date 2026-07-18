package health

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type HealthResponse struct {
	Body struct {
		Status string    `json:"status" example:"ok"`
		Time   time.Time `json:"time"`
	}
}

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "health-check",
		Method:        http.MethodGet,
		Path:          "/api/health",
		Summary:       "Health check",
		Tags:          []string{"Health"},
		DefaultStatus: http.StatusOK,
	}, healthCheck)
}

func healthCheck(ctx context.Context, input *struct{}) (*HealthResponse, error) {
	resp := &HealthResponse{}
	resp.Body.Status = "ok"
	resp.Body.Time = time.Now().UTC()
	return resp, nil
}
