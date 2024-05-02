package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"github.com/russia9/trigger/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
)

type triggerRepository struct {
	coll *mongo.Collection
}

func NewTriggerRepository(db *mongo.Database) domain.TriggerRepository {
	return &triggerRepository{
		coll: db.Collection("triggers"),
	}
}

func (t triggerRepository) Create(ctx context.Context, object *domain.Trigger) error {
	_, err := t.coll.InsertOne(ctx, object)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}

	return nil
}

func (t triggerRepository) Get(ctx context.Context, trigger string, chat int64) ([]*domain.Trigger, error) {
	// Do MongoDB query
	cursor, err := t.coll.Find(ctx, bson.M{
		"trigger": bson.M{
			"$regex": primitive.Regex{Pattern: "^" + regexp.QuoteMeta(trigger) + "$", Options: "i"},
		},
		"chat": chat,
	})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}

	// Decode result
	var result []*domain.Trigger
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}

	return result, nil
}

func (t triggerRepository) Delete(ctx context.Context, trigger string, chat int64) (int64, error) {
	// Do MongoDB query
	count, err := t.coll.DeleteOne(ctx, bson.M{
		"trigger": bson.M{
			"$regex": primitive.Regex{Pattern: "^" + regexp.QuoteMeta(trigger) + "$", Options: "i"},
		},
		"chat": chat,
	})
	if err != nil {
		return 0, errors.Wrap(err, "mongo")
	}

	return count.DeletedCount, nil
}

func (t triggerRepository) List(ctx context.Context, chat int64) ([]*domain.Trigger, error) {
	// Do MongoDB query
	cursor, err := t.coll.Find(ctx, bson.M{
		"chat": chat,
	})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}

	// Decode results
	var result []*domain.Trigger
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}

	return result, nil
}
