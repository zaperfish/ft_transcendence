package event

import (
	// Std
	"context"
	"fmt"
	"net/http"
	"time"

	// Extern
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type LabelHandler struct {
	db *gorm.DB
}

func RegisterLabelsApi(api huma.API, db *gorm.DB) {
	h := LabelHandler{db: db}

	db.AutoMigrate(&Label{})

	// Register POST /labels
	huma.Register(api, huma.Operation{
		OperationID:   "create-label",
		Method:        http.MethodPost,
		Path:          "/api/labels",
		Summary:       "Create a label",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusCreated,
	}, h.HandlePostLabel)

	// Register GET /labels/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "get-label",
		Method:        http.MethodGet,
		Path:          "/api/labels/{id}",
		Summary:       "Get a label by ID",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetLabel)

	// Register GET /labels
	huma.Register(api, huma.Operation{
		OperationID:   "get-labels",
		Method:        http.MethodGet,
		Path:          "/api/labels",
		Summary:       "Get a list of labels",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusOK,
	}, h.HandleGetLabels)

	// Register DELETE /labels/{id}
	huma.Register(api, huma.Operation{
		OperationID:   "delete-label",
		Method:        http.MethodDelete,
		Path:          "/api/labels/{id}",
		Summary:       "Delete a a label",
		Tags:          []string{"Labels"},
		DefaultStatus: http.StatusOK,
	}, h.HandleDeleteLabel)
}

type Label struct {
	gorm.Model

	// Core
	Name string `gorm:"not null;uniqueIndex:idx_labels_name,WHERE:deleted_at IS NULL;check:length(name) >= 1"`
}

type LabelDTO struct {
	ID        uint      `json:"id" doc:"ID of the event"`
	CreatedAt time.Time `json:"created_at" doc:"Time the event got created"`
	UpdatedAt time.Time `json:"updated_at" doc:"Time the event got updated"`
	Name      string    `json:"name" doc:"Name of the label"`
}

type LabelListDTO struct {
	Data     []LabelDTO `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
}

type CreateLabelDTO struct {
	Name string `json:"name" minLenght:"1" maxLength:"15" example:"Go" doc:"Name of the label"`
}

type CreateLabelInput struct {
	Body CreateLabelDTO
}

type GetLabelsInput struct {
	Page     int `query:"page" minimum:"1" default:"1" doc:"Filter by page"`
	PageSize int `query:"page_size" minimum:"1" default:"10" doc:"Page size"`
}

type LabelOutput struct {
	Body LabelDTO
}

type LabelsOutput struct {
	Body LabelListDTO
}

func (l *Label) toDTO() LabelDTO {
	labelDTO := LabelDTO{
		Name:      l.Name,
		ID:        l.ID,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}

	return labelDTO
}

func (h *LabelHandler) HandlePostLabel(ctx context.Context, input *CreateLabelInput) (*LabelOutput, error) {
	label := Label{
		Name: input.Body.Name,
	}

	err := gorm.G[Label](h.db.Debug()).Create(ctx, &label)
	if err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	return &LabelOutput{Body: label.toDTO()}, nil
}

func (h *LabelHandler) HandleGetLabel(ctx context.Context, input *GetEventInput) (*LabelOutput, error) {
	label, err := gorm.G[Label](h.db.Debug()).Where("id = ?", input.ID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	return &LabelOutput{Body: label.toDTO()}, nil
}

func (h *LabelHandler) HandleGetLabels(ctx context.Context, input *GetLabelsInput) (*LabelsOutput, error) {
	base := gorm.G[Label](h.db.Debug())
	q := base.Limit(input.PageSize)
	offset := (input.Page - 1) * input.PageSize
	q = q.Offset(offset)

	labels, err := q.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	total := len(labels)
	labelsDTO := make([]LabelDTO, total)
	for i, label := range labels {
		labelsDTO[i] = label.toDTO()
	}

	labelsOutput := &LabelsOutput{
		Body: LabelListDTO{
			Data:     labelsDTO,
			Page:     input.Page,
			PageSize: input.PageSize,
			Total:    total,
		},
	}

	return labelsOutput, nil
}

func (h *LabelHandler) HandleDeleteLabel(ctx context.Context, input *GetEventInput) (*struct{}, error) {
	rows, err := gorm.G[Label](h.db.Debug()).Where("id = ?", input.ID).Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete label: %w", err)
	}

	if rows == 0 {
		return nil, fmt.Errorf("failed to delete label: record not found")
	}

	return nil, nil
}
