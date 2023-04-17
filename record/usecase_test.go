package record

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"storage/domain"
	"storage/domain/mocks"
	"testing"
	"time"
)

func Test_service_Set(t *testing.T) {
	repo := new(mocks.MockRecordRepository)
	mockRecord := domain.Record{
		Key:   "key",
		Value: "val",
		Ttl:   0,
	}

	t.Run("success", func(t *testing.T) {
		repo.
			On("Set", mock.Anything, &mockRecord).
			Return(nil).Once()

		u := NewRecordService(repo)
		err := u.Set(context.TODO(), &mockRecord)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})
}

func Test_service_Get(t *testing.T) {
	repo := new(mocks.MockRecordRepository)
	mockRecord := domain.Record{
		Key:   "key",
		Value: "val",
		Ttl:   0,
	}

	t.Run("success", func(t *testing.T) {
		repo.
			On("Get", mock.Anything, mockRecord.Key).
			Return(&mockRecord, nil).Once()

		u := NewRecordService(repo)
		r, err := u.Get(context.TODO(), mockRecord.Key)
		assert.Equal(t, mockRecord, *r)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("record not exist", func(t *testing.T) {
		expectedErr := errors.New("key not found")
		repo.
			On("Get", mock.Anything, mockRecord.Key).
			Return(nil, expectedErr).Once()

		u := NewRecordService(repo)
		r, err := u.Get(context.TODO(), mockRecord.Key)
		assert.Empty(t, r)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})

	t.Run("record was expired", func(t *testing.T) {
		mockRecord := mockRecord
		mockRecord.Ttl = time.Millisecond * -5

		repo.
			On("Get", mock.Anything, mockRecord.Key).Return(&mockRecord, nil).Once().
			On("Delete", mock.Anything, []string{mockRecord.Key}).Return().Once()

		s := NewRecordService(repo)
		r, err := s.Get(context.TODO(), mockRecord.Key)
		assert.Empty(t, r)
		if assert.Error(t, err) {
			assert.Equal(t, errors.New("record expired"), err)
		}

		time.Sleep(time.Millisecond)
		repo.AssertExpectations(t)
	})
}

func Test_service_GetAll(t *testing.T) {
	repo := new(mocks.MockRecordRepository)
	mockRecords := []*domain.Record{
		{
			Key:   "key1",
			Value: "val",
			Ttl:   0,
		},
		{
			Key:   "key2",
			Value: "val",
			Ttl:   -1,
		},
	}

	t.Run("get all record", func(t *testing.T) {
		repo.
			On("GetAll", mock.Anything).Return(mockRecords).Once().
			On("Delete", mock.Anything, []string{mockRecords[1].Key}).Return().Once()

		s := NewRecordService(repo)
		records := s.GetAll(context.TODO())
		expected := []*domain.Record{mockRecords[0]}
		assert.Equal(t, expected, records)

		time.Sleep(time.Millisecond * 1)
		repo.AssertExpectations(t)
	})
}

func Test_service_SetTtl(t *testing.T) {
	repo := new(mocks.MockRecordRepository)
	mockRecord := domain.Record{
		Key:   "key",
		Value: "val",
		Ttl:   0,
	}
	newTtl := 10
	mockRecordWithNewTtl := mockRecord
	mockRecordWithNewTtl.Ttl = time.Duration(newTtl)

	t.Run("success", func(t *testing.T) {
		repo.On("Get", mock.Anything, mockRecord.Key).Return(&mockRecord, nil).Once()

		repo.
			On("Set", mock.Anything, &mockRecordWithNewTtl).Return(nil)

		s := NewRecordService(repo)
		r, err := s.SetTtl(context.TODO(), &mockRecordWithNewTtl)
		assert.NoError(t, err)
		assert.Equal(t, &mockRecordWithNewTtl, r)

		repo.AssertExpectations(t)
	})

	t.Run("record not exist", func(t *testing.T) {
		repo.On("Get", mock.Anything, mockRecord.Key).
			Return(nil, errors.New("")).Once()

		s := NewRecordService(repo)
		r, err := s.SetTtl(context.TODO(), &mockRecordWithNewTtl)
		assert.Empty(t, r)
		assert.Error(t, err)

		repo.AssertExpectations(t)
	})
}

func Test_service_removeExpiredRecordJob(t *testing.T) {
	repo := new(mocks.MockRecordRepository)
	mockRecords := []*domain.Record{
		{
			Key:   "key1",
			Value: "val",
			Ttl:   0,
		},
		{
			Key:   "key2",
			Value: "val",
			Ttl:   -1,
		},
	}

	t.Run("get all record", func(t *testing.T) {
		repo.
			On("GetAll", mock.Anything).Return(mockRecords).
			On("Delete", mock.Anything, []string{mockRecords[1].Key}).Return()

		s := service{repo: repo}
		go s.removeExpiredRecordJob(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		repo.AssertExpectations(t)
	})
}
