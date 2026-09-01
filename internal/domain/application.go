package domain

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Comment string    `json:"comment"`
	Source  string    `json:"source"`

	Status ApplicationStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
