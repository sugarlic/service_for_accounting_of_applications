package dto

import (
	"encoding/json"
	"time"
)

type errorResponse struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request body"`
}

type articleErrorEnvelope struct {
	Error errorResponse `json:"error"`
}

type articleBlockResponse struct {
	ID        string          `json:"id" example:"b7d0f4d5-8d5c-4d76-b12e-3f4fd16c2a01"`
	ArticleID string          `json:"article_id" example:"e39e3994-3ebf-4bad-b861-dd9e2787e1e6"`
	Position  int             `json:"position" example:"0"`
	Type      string          `json:"type" example:"title"`
	Payload   json.RawMessage `json:"payload" swaggertype:"object"`
	CreatedAt time.Time       `json:"created_at" example:"2026-04-11T15:59:56Z"`
	UpdatedAt time.Time       `json:"updated_at" example:"2026-04-11T16:10:00Z"`
}

type articleResponse struct {
	ID          string                 `json:"id" example:"e39e3994-3ebf-4bad-b861-dd9e2787e1e6"`
	AuthorID    string                 `json:"author_id" example:"af9e4c44-2d9c-4e08-9f30-985c646563b5"`
	Status      string                 `json:"status" example:"draft"`
	Title       string                 `json:"title" example:"Как инвестировать в крипту"`
	Description string                 `json:"description" example:"Краткое описание статьи"`
	CoverURL    string                 `json:"cover_url" example:"https://example.com/cover.jpg"`
	Keywords    []string               `json:"keywords" example:"крипта,инвестиции"`
	Tags        []string               `json:"tags" example:"крипта,инвестиции"`
	ViewsCount  int64                  `json:"views_count" example:"0"`
	LikesCount  int64                  `json:"likes_count" example:"0"`
	Blocks      []articleBlockResponse `json:"blocks"`
	CreatedAt   time.Time              `json:"created_at" example:"2026-04-11T15:59:56Z"`
	UpdatedAt   time.Time              `json:"updated_at" example:"2026-04-11T16:10:00Z"`
	PublishedAt *time.Time             `json:"published_at,omitempty" example:"2026-04-11T17:00:00Z"`
}

type validationErrorResponse struct {
	Field   string `json:"field" example:"title"`
	Code    string `json:"code" example:"TITLE_REQUIRED"`
	Message string `json:"message" example:"Title is required for publication"`
}

type createArticleSuccessResponse struct {
	Article articleResponse `json:"article"`
}

type listArticlesSuccessResponse struct {
	Items []articleResponse `json:"items"`
	Total int               `json:"total" example:"2"`
}

type getArticleSuccessResponse struct {
	Article articleResponse `json:"article"`
}

type updateArticleSuccessResponse struct {
	Article          articleResponse           `json:"article"`
	ValidationErrors []validationErrorResponse `json:"validation_errors"`
}

type publishArticleSuccessResponse struct {
	Article          articleResponse           `json:"article"`
	ValidationErrors []validationErrorResponse `json:"validation_errors"`
}

type unpublishArticleSuccessResponse struct {
	Article articleResponse `json:"article"`
}

type deleteArticleSuccessResponse struct {
	Deleted bool `json:"deleted" example:"true"`
}

// Эти структуры только для документации request body.

type createArticleRequestSwagger struct {
	Title       string   `json:"title" example:"Как инвестировать в крипту"`
	Description string   `json:"description" example:"Краткое описание статьи"`
	CoverURL    string   `json:"cover_url" example:"https://example.com/cover.jpg"`
	Keywords    []string `json:"keywords" example:"крипта,инвестиции"`
	Tags        []string `json:"tags" example:"крипта,инвестиции"`
}

type updateArticleBlockItemSwagger struct {
	Type    string          `json:"type" example:"paragraph"`
	Payload json.RawMessage `json:"payload" swaggertype:"object"`
}

type updateArticleRequestSwagger struct {
	Title       string                          `json:"title" example:"Как инвестировать в крипту"`
	Description string                          `json:"description" example:"Краткое описание статьи"`
	CoverURL    string                          `json:"cover_url" example:"https://example.com/cover.jpg"`
	Keywords    []string                        `json:"keywords" example:"крипта,инвестиции"`
	Tags        []string                        `json:"tags" example:"крипта,инвестиции"`
	Blocks      []updateArticleBlockItemSwagger `json:"blocks"`
}
