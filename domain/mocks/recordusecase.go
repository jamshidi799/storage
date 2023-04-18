package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"storage/domain"
)

type MockRecordService struct {
	mock.Mock
}

func (m *MockRecordService) Set(ctx context.Context, record *domain.Record) error {
	ret := m.Called(ctx, record)
	return ret.Error(0)
}

func (m *MockRecordService) Get(ctx context.Context, key string) (*domain.Record, error) {
	ret := m.Called(ctx, key)

	err := ret.Error(1)
	if r, ok := ret.Get(0).(*domain.Record); ok {
		return r, err
	}
	return nil, err
}

func (m *MockRecordService) GetAll(ctx context.Context) []*domain.Record {
	ret := m.Called(ctx)
	if r, ok := ret.Get(0).([]*domain.Record); ok {
		return r
	}
	return nil
}

func (m *MockRecordService) SetTtl(ctx context.Context, req *domain.Record) (*domain.Record, error) {
	ret := m.Called(ctx, req)

	err := ret.Error(1)
	if r, ok := ret.Get(0).(*domain.Record); ok {
		return r, err
	}
	return nil, err
}
