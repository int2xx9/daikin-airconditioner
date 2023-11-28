package service

import (
	"database/sql"
)

type DeviceQueryService struct {
	db *sql.DB
}

type Device struct {
	ID   string
	Name string
}

func NewDeviceQueryService(db *sql.DB) *DeviceQueryService {
	return &DeviceQueryService{
		db: db,
	}
}

func (s DeviceQueryService) GetDevices() ([]Device, error) {
	rows, err := s.db.Query("select id, name from devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := []Device{}
	for rows.Next() {
		var device Device
		rows.Scan(&device.ID, &device.Name)
		devices = append(devices, device)
	}

	return devices, nil
}
