package repo

import (
	"context"
	"github.com/mebr0/tiny-url/internal/domain"
	"github.com/mebr0/tiny-url/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLsRepo struct {
	db *mongo.Collection
}

func newURLsRepo(db *mongo.Database) *URLsRepo {
	return &URLsRepo{
		db: db.Collection(urlsCollection),
	}
}

func (r *URLsRepo) ListByOwner(ctx context.Context, userId primitive.ObjectID) ([]domain.URL, error) {
	urls := make([]domain.URL, 0)

	cur, err := r.db.Find(ctx, bson.M{"owner": userId})

	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &urls)

	return urls, err
}

func (r *URLsRepo) Create(ctx context.Context, url domain.URL) (string, error) {
	res, err := r.db.InsertOne(ctx, url)

	if err != nil {
		if mongodb.IsDuplicate(err) {
			return "", ErrURLAlreadyExists
		}

		return "", err
	}

	return res.InsertedID.(string), nil
}

func (r *URLsRepo) Get(ctx context.Context, alias string) (domain.URL, error) {
	var url domain.URL

	if err := r.db.FindOne(ctx, bson.M{"_id": alias}).Decode(&url); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.URL{}, ErrURLNotFound
		}

		return domain.URL{}, err
	}

	return url, nil
}

func (r *URLsRepo) GetByOriginalAndOwner(ctx context.Context, original string, owner primitive.ObjectID) (domain.URL, error) {
	var url domain.URL

	if err := r.db.FindOne(ctx, bson.M{"original": original, "owner": owner}).Decode(&url); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.URL{}, ErrURLNotFound
		}

		return domain.URL{}, err
	}

	return url, nil
}

func (r *URLsRepo) Delete(ctx context.Context, alias string) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": alias})

	return err
}
