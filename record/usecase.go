package record

import (
	"context"
	"errors"
	"storage/domain"
)

type service struct {
	repo domain.RecordRepository
}

func NewRecordService(repo domain.RecordRepository) domain.RecordService {
	return &service{repo: repo}
}

func (s *service) Set(ctx context.Context, record *domain.Record) error {
	return s.repo.Set(ctx, record)
}

func (s *service) Get(ctx context.Context, key string) (*domain.Record, error) {
	record, err := s.repo.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if record.IsExpired() {
		go s.repo.Delete(context.Background(), key)
		return nil, errors.New("record expired")
	}

	return record, nil
}

func (s *service) GetAll(ctx context.Context) []*domain.Record {
	records := s.repo.GetAll(ctx)
	var notExpiredRecords []*domain.Record
	for _, record := range records {
		if record.IsExpired() {
			go s.repo.Delete(context.Background(), record.Key)
		} else {
			notExpiredRecords = append(notExpiredRecords, record)
		}
	}

	return notExpiredRecords
}

func (s *service) SetTtl(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	r, err := s.repo.Get(ctx, record.Key)
	if err != nil {
		return nil, err
	}

	r.Ttl = record.Ttl
	if err = s.repo.Set(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}
