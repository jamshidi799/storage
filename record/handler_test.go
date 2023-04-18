package record

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"net/url"
	"storage/domain"
	"storage/domain/mocks"
	"storage/util"
	"testing"
	"time"
)

func Test_handler_set(t *testing.T) {
	mockRecord := &domain.Record{
		Key:   "key",
		Value: "value",
		Ttl:   time.Hour,
	}

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockRecordService)
		mockService.
			On("Set", mock.Anything, mockRecord).
			Return(nil).Once()

		setReq := setRecordRequest{
			Key:   mockRecord.Key,
			Value: mockRecord.Value,
			Ttl:   mockRecord.Ttl,
		}

		w := httptest.NewRecorder()
		ctx := util.GetTestGinContext(w)
		util.MockJsonPost(ctx, setReq)
		h := handler{service: mockService}
		h.set(ctx)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		mockService := new(mocks.MockRecordService)
		setReq := setRecordRequest{
			Value: mockRecord.Value,
			Ttl:   mockRecord.Ttl,
		}

		w := httptest.NewRecorder()
		ctx := util.GetTestGinContext(w)
		util.MockJsonPost(ctx, setReq)

		h := handler{service: mockService}
		h.set(ctx)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "validation for 'Key' failed")
	})
}

func Test_handler_getAll(t *testing.T) {
	records := []*domain.Record{
		{
			Key:   "key1",
			Value: "val1",
			Ttl:   0,
		},
		{
			Key:   "key2",
			Value: "val2",
			Ttl:   time.Hour,
		},
	}

	mockService := new(mocks.MockRecordService)
	mockService.On("GetAll", mock.Anything).
		Return(records).Once()

	w := httptest.NewRecorder()
	ctx := util.GetTestGinContext(w)
	util.MockJsonGet(ctx, []gin.Param{}, url.Values{})

	h := handler{service: mockService}
	h.getAll(ctx)

	var res []*response
	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
	assert.Equal(t, toResponse(records[0]), res[0])
	assert.Equal(t, toResponse(records[1]), res[1])
}

func Test_handler_get(t *testing.T) {
	mockRecord := &domain.Record{
		Key:   "key",
		Value: "value",
		Ttl:   time.Hour,
	}

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockRecordService)
		mockService.On("Get", mock.Anything, mockRecord.Key).
			Return(mockRecord, nil).Once()

		w := httptest.NewRecorder()
		ctx := util.GetTestGinContext(w)
		util.MockJsonGet(ctx, []gin.Param{{Key: "key", Value: mockRecord.Key}}, url.Values{})

		h := handler{service: mockService}
		h.get(ctx)

		var res response
		err := json.Unmarshal(w.Body.Bytes(), &res)

		assert.Equal(t, 200, w.Code)
		assert.NoError(t, err)
		assert.Equal(t, toResponse(mockRecord), &res)
	})

	t.Run("key slug not found", func(t *testing.T) {
		mockService := new(mocks.MockRecordService)

		w := httptest.NewRecorder()
		ctx := util.GetTestGinContext(w)
		util.MockJsonGet(ctx, []gin.Param{}, url.Values{})

		h := handler{service: mockService}
		h.get(ctx)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("record not found", func(t *testing.T) {
		mockService := new(mocks.MockRecordService)
		mockService.On("Get", mock.Anything, mockRecord.Key).
			Return(nil, errors.New("")).Once()

		w := httptest.NewRecorder()
		ctx := util.GetTestGinContext(w)
		util.MockJsonGet(ctx, []gin.Param{{Key: "key", Value: mockRecord.Key}}, url.Values{})

		h := handler{service: mockService}
		h.get(ctx)

		assert.Equal(t, 400, w.Code)
	})
}

func Test_handler_setTtl(t *testing.T) {
	mockRecord := &domain.Record{
		Key: "key",
		Ttl: time.Hour,
	}

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockRecordService)
		mockService.
			On("SetTtl", mock.Anything, mockRecord).
			Return(mockRecord, nil).Once()

		setReq := setRecordTtlRequest{
			Key: mockRecord.Key,
			Ttl: mockRecord.Ttl,
		}

		w := httptest.NewRecorder()
		ctx := util.GetTestGinContext(w)
		util.MockJsonPost(ctx, setReq)
		h := handler{service: mockService}
		h.setTtl(ctx)

		var res response
		err := json.Unmarshal(w.Body.Bytes(), &res)

		assert.Equal(t, 200, w.Code)
		assert.NoError(t, err)
		assert.Equal(t, toResponse(mockRecord), &res)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockService := new(mocks.MockRecordService)

		setReq := setRecordTtlRequest{
			Ttl: mockRecord.Ttl,
		}

		w := httptest.NewRecorder()
		ctx := util.GetTestGinContext(w)
		util.MockJsonPost(ctx, setReq)
		h := handler{service: mockService}
		h.setTtl(ctx)

		assert.Equal(t, 400, w.Code)
	})
}
