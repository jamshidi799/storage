package record

import (
	"context"
	"gorm.io/gorm"
	"storage/domain"
	"time"
)

type recordModel struct {
	Key      string `gorm:"primaryKey"`
	Value    string
	ExpireAt time.Time `gorm:"index"`
}

type postgresRepo struct {
	db *gorm.DB
}

func NewPostgresRecordRepository(db *gorm.DB) domain.RecordRepository {
	return &postgresRepo{db: db}
}

func (p *postgresRepo) Set(ctx context.Context, record *domain.Record) error {
	return p.db.WithContext(ctx).Create(convertToModel(record)).Error
}

func (p *postgresRepo) Get(ctx context.Context, key string) (*domain.Record, error) {
	var r recordModel
	err := p.db.WithContext(ctx).Where("key = ?", key).First(&r).Error
	return r.toRecord(), err
}

func (p *postgresRepo) GetAll(ctx context.Context) []*domain.Record {
	var rows []recordModel
	p.db.WithContext(ctx).Find(&rows)

	var records []*domain.Record
	for _, r := range rows {
		records = append(records, r.toRecord())
	}
	return records
}

func (p *postgresRepo) Delete(ctx context.Context, key string) {
	p.db.WithContext(ctx).Delete(key)
}

func convertToModel(r *domain.Record) *recordModel {
	return &recordModel{
		Key:      r.Key,
		Value:    r.Value,
		ExpireAt: time.Now().Add(r.Ttl),
	}
}

func (r *recordModel) toRecord() *domain.Record {
	return &domain.Record{
		Key:   r.Key,
		Value: r.Value,
		Ttl:   r.ExpireAt.Sub(time.Now()),
	}
}
