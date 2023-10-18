package repository_user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	RegisteredAt time.Time          `bson:"registered_at"`
	LastVisitAt  time.Time          `bson:"last_visit_at"`
}
