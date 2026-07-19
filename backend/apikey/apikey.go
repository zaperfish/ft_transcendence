package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"ft_transcendence/backend/middleware"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type GormApiKeyModel struct {
	gorm.Model

	KeyHash string
	Revoked bool
}

func (GormApiKeyModel) TableName() string {
	return "api_keys"
}

func GenerateApiKey() (string, error) {
	prefix := "ft_transcendence"
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	secret := base64.RawURLEncoding.EncodeToString(b)

	return fmt.Sprintf("%s.%s", prefix, secret), nil
}

func HashApiKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func StoreApiKey(key string, ctx context.Context, db *gorm.DB) error {
	model := GormApiKeyModel{
		KeyHash: HashApiKey(key),
		Revoked: false,
	}

	err := gorm.G[GormApiKeyModel](db.Debug()).Create(ctx, &model)
	if err != nil {
		return huma.Error409Conflict("api key exists")
	}

	return nil
}

func ValidateApiKey(key string, ctx context.Context, db *gorm.DB) bool {
	hash := HashApiKey(key)

	var exists bool
	db.Model(&GormApiKeyModel{}).
		Select("count(*) > 0").
		Where("key_hash = ? AND revoked = false", hash).
		Find(&exists)

	return exists
}

type ApiKeyHandler struct {
	db *gorm.DB
}

type CreateApiKeyInput struct {
}

type CreateApiKeyOutput struct {
	Body CreateApiKeyOutputBody
}

type CreateApiKeyOutputBody struct {
	Key string `json:"key"`
}

func (h *ApiKeyHandler) CreateApiKeyHandler(ctx context.Context, input *CreateApiKeyInput) (*CreateApiKeyOutput, error) {
	key, err := GenerateApiKey()
	if err != nil {
		return nil, err
	}

	model := GormApiKeyModel{
		KeyHash: HashApiKey(key),
		Revoked: false,
	}

	err = gorm.G[GormApiKeyModel](h.db.Debug()).Create(ctx, &model)
	if err != nil {
		return nil, err
	}

	return &CreateApiKeyOutput{
		Body: CreateApiKeyOutputBody{
			Key: key,
		},
	}, nil
}

type DeleteApiKeyInput struct {
	ID uint `path:"id"`
}

type DeleteApiKeyOutput struct{}

func (h *ApiKeyHandler) DeleteApiKeyHandler(ctx context.Context, input *DeleteApiKeyInput) (*DeleteApiKeyOutput, error) {
	rows, err := gorm.G[GormApiKeyModel](h.db.Debug()).
		Where("id = ?", input.ID).
		Update(ctx, "revoked", true)

	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, huma.Error404NotFound("record not found")
	}

	return &DeleteApiKeyOutput{}, nil
}

type ListApiKeysInput struct{}

type ListApiKeysOutput struct {
	Body ListApiKeysOutputBody
}

type ListApiKeysOutputBody struct {
	Data []ApiKeyDTO `json:"data"`
}

type ApiKeyDTO struct {
	ID      uint `json:"id"`
	Revoked bool `json:"revoked"`
}

func (h *ApiKeyHandler) ListApiKeysHandler(ctx context.Context, input *ListApiKeysInput) (*ListApiKeysOutput, error) {
	keys, err := gorm.G[GormApiKeyModel](h.db.Debug()).Find(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ApiKeyDTO, len(keys))
	for i, k := range keys {
		out[i] = ApiKeyDTO{
			ID:      k.ID,
			Revoked: k.Revoked,
		}
	}

	return &ListApiKeysOutput{
		Body: ListApiKeysOutputBody{
			Data: out,
		},
	}, nil
}

func RegisterRoutes(api huma.API, db *gorm.DB) {
	db.AutoMigrate(&GormApiKeyModel{})

	h := ApiKeyHandler{
		db: db,
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is not set")
	}
	admin := middleware.PasswordVerifier(api, password)

	huma.Register(api, huma.Operation{
		OperationID:   "create-api-key",
		Method:        "POST",
		Path:          "/api/api-keys",
		Summary:       "Create API key",
		Tags:          []string{"ApiKey"},
		Middlewares:   huma.Middlewares{admin},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"AdminPassword": {}},
		},
	}, h.CreateApiKeyHandler)

	huma.Register(api, huma.Operation{
		OperationID:   "list-api-keys",
		Method:        "GET",
		Path:          "/api/api-keys",
		Summary:       "List API keys",
		Tags:          []string{"ApiKey"},
		Middlewares:   huma.Middlewares{admin},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"AdminPassword": {}},
		},
	}, h.ListApiKeysHandler)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-api-key",
		Method:        "DELETE",
		Path:          "/api/api-keys/{id}",
		Summary:       "Revoke API key",
		Tags:          []string{"ApiKey"},
		Middlewares:   huma.Middlewares{admin},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"AdminPassword": {}},
		},
	}, h.DeleteApiKeyHandler)
}

func ApiKeyVerifier(api huma.API, db *gorm.DB) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		auth := ctx.Header("Authorization")
		if auth == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid Authorization format")
			return
		}

		provided := strings.TrimPrefix(auth, "Bearer ")

		valid := ValidateApiKey(provided, ctx.Context(), db)
		if !valid {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid api key")
			return
		}

		next(ctx)
	}
}
