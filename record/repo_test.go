package record

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"storage/domain"
	"time"

	"testing"
)

func initDB() (sqlmock.Sqlmock, error, *postgresRepo) {
	db, mock, err := sqlmock.New()
	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	repo := postgresRepo{db: gormDb}
	return mock, err, &repo
}

func TestPostgresRepo_Set(t *testing.T) {
	r := &domain.Record{
		Key:   "key",
		Value: "val",
		Ttl:   0,
	}
	model := convertToModel(r)

	mock, err, repo := initDB()
	assert.NoError(t, err)

	mock.ExpectBegin()
	query := `UPDATE "records" SET`
	mock.ExpectExec(query).
		WithArgs(model.Value, model.ExpireAt, model.Key).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Set(context.TODO(), r)
	assert.NoError(t, err)
}

func TestPostgresRepo_Get(t *testing.T) {
	r := &domain.Record{
		Key:   "key",
		Value: "val",
		Ttl:   0,
	}
	model := convertToModel(r)

	mock, err, repo := initDB()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"key", "value", "expire_at"}).
		AddRow(model.Key, model.Value, model.ExpireAt)

	query := `SELECT \* FROM "records"`
	mock.ExpectQuery(query).WithArgs(model.Key).WillReturnRows(rows)
	mock.ExpectCommit()

	actual, err := repo.Get(context.TODO(), model.Key)
	assert.NoError(t, err)
	assert.Equal(t, actual, r)
}

func TestPostgresRepo_GetAll(t *testing.T) {
	records := []*domain.Record{
		{
			Key:   "key",
			Value: "val",
			Ttl:   0,
		},
		{
			Key:   "key",
			Value: "val",
			Ttl:   time.Hour,
		},
	}

	firstRow := convertToModel(records[0])
	secondRow := convertToModel(records[1])
	rows := sqlmock.NewRows([]string{"key", "value", "expire_at"}).
		AddRow(firstRow.Key, firstRow.Value, firstRow.ExpireAt).
		AddRow(secondRow.Key, secondRow.Value, secondRow.ExpireAt)

	mock, err, repo := initDB()
	assert.NoError(t, err)

	query := `SELECT \* FROM "records"`
	mock.ExpectQuery(query).WithArgs().WillReturnRows(rows)

	result := repo.GetAll(context.TODO())
	assert.Equal(t, *records[0], *result[0])

	assert.Equal(t, records[1].Key, result[1].Key)
	assert.Equal(t, records[1].Value, result[1].Value)
	assert.Equal(t, records[1].Value, result[1].Value)
	assert.NotEqual(t, records[1].Ttl, result[1].Ttl)
}

func TestPostgresRepo_Delete(t *testing.T) {
	keys := []string{"key1", "key2"}
	mock, err, repo := initDB()
	assert.NoError(t, err)

	mock.ExpectBegin()
	query := `DELETE FROM "records"`
	mock.ExpectExec(query).WithArgs(keys[0], keys[1]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo.Delete(context.TODO(), keys...)
}
