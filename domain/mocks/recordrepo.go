package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"storage/domain"
)

type MockRecordRepository struct {
	mock.Mock
}

func (m *MockRecordRepository) Set(ctx context.Context, record *domain.Record) error {
	ret := m.Called(ctx, record)
	return ret.Error(0)
}

func (m *MockRecordRepository) Get(ctx context.Context, key string) (*domain.Record, error) {
	ret := m.Called(ctx, key)

	err := ret.Error(1)
	if r, ok := ret.Get(0).(*domain.Record); ok {
		return r, err
	} else {
		return nil, err
	}
}

func (m *MockRecordRepository) GetAll(ctx context.Context) []*domain.Record {
	ret := m.Called(ctx)
	if records, ok := ret.Get(0).([]*domain.Record); ok {
		return records
	} else {
		return nil
	}
}

func (m *MockRecordRepository) Delete(ctx context.Context, keys ...string) {
	_ = m.Called(ctx, keys)
}
