package repository

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type RecordRepository struct {
	db *sql.DB
}

func NewRecordRepository(db *sql.DB) *RecordRepository {
	return &RecordRepository{
		db: db,
	}
}

func (r *RecordRepository) Add(id string, timestamp time.Time, data map[string]any) (int64, error) {
	columns := []string{"device_id", "time"}
	values := []any{id, timestamp.UTC()}
	placeholders := []string{"$1", "$2"}
	i := 3
	for key, value := range data {
		columns = append(columns, key)
		values = append(values, value)
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
		i++
	}

	sql := "insert into records (" + strings.Join(columns, ", ") + ") values (" + strings.Join(placeholders, ", ") + ")"
	result, err := r.db.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (r *RecordRepository) DeleteByDeviceId(deviceId string) (int64, error) {
	result, err := r.db.Exec("delete from records where device_id=$1", deviceId)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}
