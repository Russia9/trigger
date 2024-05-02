package domain

import (
	"context"
	"gopkg.in/telebot.v3"
)

type Trigger struct {
	ID      string `bson:"_id"`
	Trigger string
	Chat    int64

	Object   []byte
	Type     string
	Entities telebot.Entities
}

type TriggerRepository interface {
	Create(ctx context.Context, object *Trigger) error
	Get(ctx context.Context, trigger string, chat int64) ([]*Trigger, error)
	Delete(ctx context.Context, trigger string, chat int64) (int64, error)
	List(ctx context.Context, chat int64) ([]*Trigger, error)
}
