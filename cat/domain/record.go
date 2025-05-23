package domain

import (
	"github.com/google/uuid"
)

type CatRecord struct {
	ID        uuid.UUID
	Timestamp string
	Cat       string
	Weight    float32
	Notes     *string
}
