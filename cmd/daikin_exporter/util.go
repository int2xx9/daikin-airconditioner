package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/int2xx9/daikin-airconditioner/daikin"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/constraints"
)

type metricRange struct {
	name string
	min  float64
	max  float64
}

var (
	airflowRateSettingRange = metricRange{name: "AirflowRateSetting", min: 1, max: 8}
	temperatureSettingRange = metricRange{name: "TemperatureSetting", min: 0, max: 50}
	humiditySettingRange    = metricRange{name: "HumiditySetting", min: 0, max: 100}
	roomTemperatureRange    = metricRange{name: "RoomTemperature", min: -127, max: 125}
	roomHumidityRange       = metricRange{name: "RoomHumidity", min: 0, max: 100}
	outdoorTemperatureRange = metricRange{name: "OutdoorTemperature", min: -127, max: 125}
)

func inRangeOrWarn(value float64, addr net.UDPAddr, id string, validRange *metricRange) bool {
	if validRange == nil {
		return true
	}

	if value < validRange.min || value > validRange.max {
		slog.Warn("[updateMetrics] value out of range; skip metric", "address", addr.String(), "id", id, "property", validRange.name, "accepted_range", fmt.Sprintf("%.0f-%.0f", validRange.min, validRange.max), "actual_value", value)
		return false
	}

	return true
}

func bytesToString(data []byte) string {
	s := "0x"
	for _, b := range data {
		s += fmt.Sprintf("%02x", b)
	}
	return s
}

func updateBoolMetrics(addr net.UDPAddr, id string, getter func() (value bool, err error), gaugeVec *prometheus.GaugeVec) error {
	value, err := getter()
	if err != nil {
		return err
	}
	if value {
		gaugeVec.WithLabelValues(addr.String(), id).Set(1)
	} else {
		gaugeVec.WithLabelValues(addr.String(), id).Set(0)
	}
	return nil
}

func updateNumberMetrics[T constraints.Signed | constraints.Unsigned](addr net.UDPAddr, id string, getter func() (value T, err error), gaugeVec *prometheus.GaugeVec, validRange *metricRange) error {
	value, err := getter()
	if err != nil {
		return err
	}
	if !inRangeOrWarn(float64(value), addr, id, validRange) {
		return nil
	}
	gaugeVec.WithLabelValues(addr.String(), id).Set(float64(value))
	return nil
}

func updateNumberWithAutoMetrics[T constraints.Signed | constraints.Unsigned](addr net.UDPAddr, id string, getter func() (value T, auto bool, err error), gaugeVec *prometheus.GaugeVec, autoGaugeVec *prometheus.GaugeVec, validRange *metricRange) error {
	value, auto, err := getter()
	if err != nil {
		return err
	}
	if auto {
		autoGaugeVec.WithLabelValues(addr.String(), id).Set(1)
		return nil
	}
	autoGaugeVec.WithLabelValues(addr.String(), id).Set(0)
	if !inRangeOrWarn(float64(value), addr, id, validRange) {
		return nil
	}
	gaugeVec.WithLabelValues(addr.String(), id).Set(float64(value))
	return nil
}

func updateOperationMode(addr net.UDPAddr, id string, getter func() (value daikin.OperationMode, err error), gaugeVec *prometheus.GaugeVec) error {
	value, err := getter()
	if err != nil {
		return err
	}
	gaugeVec.WithLabelValues(addr.String(), id, "auto").Set(0)
	gaugeVec.WithLabelValues(addr.String(), id, "cooling").Set(0)
	gaugeVec.WithLabelValues(addr.String(), id, "heating").Set(0)
	gaugeVec.WithLabelValues(addr.String(), id, "dehumidification").Set(0)
	gaugeVec.WithLabelValues(addr.String(), id, "ventilation").Set(0)
	gaugeVec.WithLabelValues(addr.String(), id, "other").Set(0)
	switch value {
	case daikin.OperationModeAuto:
		gaugeVec.WithLabelValues(addr.String(), id, "auto").Set(1)
	case daikin.OperationModeCooling:
		gaugeVec.WithLabelValues(addr.String(), id, "cooling").Set(1)
	case daikin.OperationModeHeating:
		gaugeVec.WithLabelValues(addr.String(), id, "heating").Set(1)
	case daikin.OperationModeDehumidification:
		gaugeVec.WithLabelValues(addr.String(), id, "dehumidification").Set(1)
	case daikin.OperationModeVentilating:
		gaugeVec.WithLabelValues(addr.String(), id, "ventilation").Set(1)
	case daikin.OperationModeOther:
		gaugeVec.WithLabelValues(addr.String(), id, "other").Set(1)
	}
	return nil
}
