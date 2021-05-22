package repo

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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
		if mongodb.IsDuplicate(err) {
			return [12]byte{}, ErrUserAlreadyExists
		}

		return [12]byte{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User

	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) UpdateLastLogin(ctx context.Context, id primitive.ObjectID, lastLogin time.Time) error {
	if _, err := r.db.UpdateByID(ctx, id, bson.M{"$set": bson.M{"lastLogin": lastLogin}}); err != nil {
		return err
	}

	return nil
}
