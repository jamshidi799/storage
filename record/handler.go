package record

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"storage/domain"
	"time"
)

type handler struct {
	service domain.RecordService
}

func NewRecordController(rg *gin.RouterGroup, rs domain.RecordService) {
	h := &handler{service: rs}

	rg.POST("", h.set)
	rg.GET("", h.getAll)
	rg.GET(":key", h.get)
	rg.POST("ttl", h.setTtl)
}

func (h *handler) set(c *gin.Context) {
	var req setRecordRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Set(c.Request.Context(), req.toRecord()); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) getAll(c *gin.Context) {
	records := h.service.GetAll(c.Request.Context())

	var res []*response
	for _, record := range records {
		res = append(res, toResponse(record))
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) get(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, errors.New("key slug not found"))
		return
	}
	record, err := h.service.Get(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, toResponse(record))
}

func (h *handler) setTtl(c *gin.Context) {
	var req setRecordTtlRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	record, err := h.service.SetTtl(c.Request.Context(), req.toRecord())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, toResponse(record))
}

type setRecordRequest struct {
	Key   string        `json:"key" binding:"required"`
	Value string        `json:"value" binding:"required"`
	Ttl   time.Duration `json:"ttl"`
}

func (s *setRecordRequest) toRecord() *domain.Record {
	return &domain.Record{
		Key:   s.Key,
		Value: s.Value,
		Ttl:   s.Ttl,
	}
}

type response struct {
	Key   string        `json:"key"`
	Value string        `json:"value"`
	Ttl   time.Duration `json:"ttl,omitempty"`
}

func toResponse(r *domain.Record) *response {
	return &response{
		Key:   r.Key,
		Value: r.Value,
		Ttl:   r.Ttl,
	}
}

type setRecordTtlRequest struct {
	Key string        `json:"key" binding:"required"`
	Ttl time.Duration `json:"ttl" binding:"required"`
}

func (s *setRecordTtlRequest) toRecord() *domain.Record {
	return &domain.Record{
		Key: s.Key,
		Ttl: s.Ttl,
	}
}
