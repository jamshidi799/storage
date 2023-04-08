package record

import (
	"context"
	"gorm.io/gorm"
	"log"
	"storage/domain"
	"time"
)

type record struct {
	Key      string `gorm:"primaryKey"`
	Value    string
	ExpireAt time.Time `gorm:"index"`
}

type postgresRepo struct {
	db *gorm.DB
}

func NewPostgresRecordRepository(db *gorm.DB) domain.RecordRepository {
	if err := db.AutoMigrate(record{}); err != nil {
		log.Println(err)
	}

	return &postgresRepo{db: db}
}

func (p *postgresRepo) Set(ctx context.Context, record *domain.Record) error {
	return p.db.WithContext(ctx).Save(convertToModel(record)).Error
}

func (p *postgresRepo) Get(ctx context.Context, key string) (*domain.Record, error) {
	var r record
	err := p.db.WithContext(ctx).Where("key = ?", key).First(&r).Error
	return r.toRecord(), err
}

func (p *postgresRepo) GetAll(ctx context.Context) []*domain.Record {
	var rows []record
	p.db.WithContext(ctx).Find(&rows)

	var records []*domain.Record
	for _, r := range rows {
		records = append(records, r.toRecord())
	}
	return records
}

func (p *postgresRepo) Delete(ctx context.Context, keys ...string) {
	p.db.WithContext(ctx).Delete(record{}, keys)
}

func convertToModel(r *domain.Record) *record {
	var expireAt time.Time
	if r.Ttl != 0 {
		expireAt = time.Now().Add(r.Ttl)
	}
	return &record{
		Key:      r.Key,
		Value:    r.Value,
		ExpireAt: expireAt,
	}
}

func (r *record) toRecord() *domain.Record {
	var ttl time.Duration
	if !r.ExpireAt.IsZero() {
		ttl = r.ExpireAt.Sub(time.Now())
	}
	return &domain.Record{
		Key:   r.Key,
		Value: r.Value,
		Ttl:   ttl,
	}
}
