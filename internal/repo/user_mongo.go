package repo

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	db *mongo.Collection
}

func newUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) List(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	cur, err := r.db.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &users)

	return users, err
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, user)

	if err != nil {
		return [12]byte{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}
