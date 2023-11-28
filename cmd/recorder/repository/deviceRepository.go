package repository

import (
	"database/sql"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) GetDevices() ([]string, error) {
	rows, err := r.db.Query("select id from devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := []string{}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		devices = append(devices, id)
	}

	return devices, nil
}

func (r *DeviceRepository) Add(id string, name string) (int64, error) {
	result, err := r.db.Exec("insert into devices (id, name) values ($1, $2)", id, name)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, err
}

func (r *DeviceRepository) Delete(id string) (int64, error) {
	result, err := r.db.Exec("delete from devices where id=$1", id)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (r *DeviceRepository) Rename(id string, name string) (int64, error) {
	result, err := r.db.Exec("update devices set name=$1 where id=$2", name, id)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, err
}
