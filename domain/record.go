package domain

import (
	"context"
	"time"
)

type Record struct {
	Key   string
	Value string
	Ttl   time.Duration
}

type RecordService interface {
	Set(ctx context.Context, record *Record) error
	Get(ctx context.Context, key string) (*Record, error)
	GetAll(ctx context.Context) []*Record
	SetTtl(ctx context.Context, req *Record) (*Record, error)
}

type RecordRepository interface {
	Set(ctx context.Context, record *Record) error
	Get(ctx context.Context, key string) (*Record, error)
	GetAll(ctx context.Context) []*Record
	Delete(ctx context.Context, key string)
}

func (r *Record) IsExpired() bool {
	return r.Ttl < 0
}
