package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" bson:"_id,omitempty"`
	Name         string    `json:"name" bson:"name"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"password" bson:"password"`
	RegisteredAt time.Time `json:"registered_at" bson:"registered_at"`
	LastVisitAt  time.Time `json:"last_visit_at" bson:"last_visit_at"`
}
