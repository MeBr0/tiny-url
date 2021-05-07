package repo

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) (primitive.ObjectID, error)
}

type Repos struct {
	Users Users
}

func NewRepos(db *mongo.Database) *Repos {
	return &Repos{
		Users: newUsersRepo(db),
	}
}
