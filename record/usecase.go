package record

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache/v3"
	"log"
	"storage/domain"
	"time"
)

type service struct {
	repo  domain.RecordRepository
	cache *bigcache.BigCache
}

func NewRecordService(repo domain.RecordRepository) domain.RecordService {
	cache, _ := bigcache.New(context.Background(), getCacheConfig())

	s := service{
		repo:  repo,
		cache: cache,
	}

	go printCacheStats(cache)
	go s.removeExpiredRecordJob(10 * time.Minute)

	return &s
}

func getCacheConfig() bigcache.Config {
	config := bigcache.DefaultConfig(10 * time.Minute)
	config.MaxEntriesInWindow = 10000
	config.HardMaxCacheSize = 32 // MB
	return config
}

func (s *service) Set(ctx context.Context, record *domain.Record) error {
	return s.repo.Set(ctx, record)
}

func (s *service) Get(ctx context.Context, key string) (*domain.Record, error) {
	if value := s.cacheGet(key); value != nil {
		return value, nil
	}

	record, err := s.repo.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if record.IsExpired() {
		go s.repo.Delete(context.Background(), key)
		return nil, errors.New("record expired")
	}

	s.cacheSet(key, record)

	return record, nil
}

func (s *service) GetAll(ctx context.Context) []*domain.Record {
	records := s.repo.GetAll(ctx)
	var notExpiredRecords []*domain.Record
	var expiredKeys []string
	for _, r := range records {
		if r.IsExpired() {
			expiredKeys = append(expiredKeys, r.Key)
		} else {
			notExpiredRecords = append(notExpiredRecords, r)
		}
	}

	if len(expiredKeys) > 0 {
		go s.repo.Delete(context.Background(), expiredKeys...)
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

func (s *service) removeExpiredRecordJob(per time.Duration) {
	for range time.Tick(per) {
		records := s.repo.GetAll(context.Background())
		var expiredKeys []string
		for _, record := range records {
			if record.IsExpired() {
				expiredKeys = append(expiredKeys, record.Key)
			}
		}

		s.repo.Delete(context.Background(), expiredKeys...)
	}
}

func (s *service) cacheGet(key string) *domain.Record {
	if value, err := s.cache.Get(key); err == nil {
		var record domain.Record
		json.Unmarshal(value, &record)

		if record.IsExpired() {
			s.cache.Delete(key)
			return nil
		}

		return &record
	}

	return nil
}

func (s *service) cacheSet(key string, value *domain.Record) {
	v, _ := json.Marshal(value)
	s.cache.Set(key, v)
}

func printCacheStats(cache *bigcache.BigCache) {
	for range time.Tick(time.Hour) {
		log.Printf("cache stats: %+v,	length: %d\n", cache.Stats(), cache.Len())
	}
}
