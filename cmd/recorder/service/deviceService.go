package service

import (
	"fmt"

	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/adapter"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/config"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/repository"
)

type DeviceService struct {
	config             config.Configuration
	daikinAdapter      *adapter.DaikinAdapter
	deviceQueryService *DeviceQueryService
	deviceRepository   *repository.DeviceRepository
	recordRepository   *repository.RecordRepository
}

func NewDeviceService(config config.Configuration, daikinAdapter *adapter.DaikinAdapter, deviceQueryService *DeviceQueryService, deviceRepository *repository.DeviceRepository, recordRepository *repository.RecordRepository) DeviceService {
	return DeviceService{
		config:             config,
		daikinAdapter:      daikinAdapter,
		deviceQueryService: deviceQueryService,
		deviceRepository:   deviceRepository,
		recordRepository:   recordRepository,
	}
}

func (s DeviceService) Add(id string, name string) (int64, error) {
	affected, err := s.deviceRepository.Add(id, name)
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (s DeviceService) Remove(id string) error {
	recordsAffected, err := s.recordRepository.DeleteByDeviceId(id)
	if err != nil {
		return err
	}

	devicesAffected, err := s.deviceRepository.Delete(id)
	if err != nil {
		return err
	}

	fmt.Printf("%d device(s), %d record(s) removed\n", devicesAffected, recordsAffected)
	return nil
}

func (s DeviceService) Rename(id string, name string) error {
	affected, err := s.deviceRepository.Rename(id, name)
	if err != nil {
		return err
	}

	fmt.Printf("%d device(s) renamed\n", affected)
	return nil
}

func (s DeviceService) List() error {
	devices, err := s.deviceQueryService.GetDevices()
	if err != nil {
		return err
	}

	fmt.Printf("ID\t\t\t\t\tName\n")
	for _, device := range devices {
		fmt.Printf("%s\t%s\n", device.ID, device.Name)
	}

	return nil
}

func (s DeviceService) Discover() error {
	devices, err := s.daikinAdapter.Discover()
	if err != nil {
		return err
	}

	fmt.Printf("%d device(s) are discovered.\n", len(devices))
	fmt.Printf("IPAddress\tID\n")

	for _, device := range devices {
		fmt.Printf("%s\t%s\n", device.Address.IP, device.ID)
	}

	return nil
}
