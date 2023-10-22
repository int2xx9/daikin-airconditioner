package main

import (
	"fmt"
	"net"

	"github.com/int2xx9/daikin-airconditioner/daikin"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/constraints"
)

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

func updateNumberMetrics[T constraints.Signed | constraints.Unsigned](addr net.UDPAddr, id string, getter func() (value T, err error), gaugeVec *prometheus.GaugeVec) error {
	value, err := getter()
	if err != nil {
		return err
	}
	gaugeVec.WithLabelValues(addr.String(), id).Set(float64(value))
	return nil
}

func updateNumberWithAutoMetrics[T constraints.Signed | constraints.Unsigned](addr net.UDPAddr, id string, getter func() (value T, auto bool, err error), gaugeVec *prometheus.GaugeVec, autoGaugeVec *prometheus.GaugeVec) error {
	value, auto, err := getter()
	if err != nil {
		return err
	}
	if auto {
		autoGaugeVec.WithLabelValues(addr.String(), id).Set(1)
	} else {
		autoGaugeVec.WithLabelValues(addr.String(), id).Set(0)
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
