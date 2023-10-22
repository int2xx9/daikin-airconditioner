package main

import (
	"net/http"
	"os"

	"github.com/int2xx9/daikin-airconditioner/daikin"
	"github.com/int2xx9/daikin-airconditioner/echonetlite"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	handler := newDaikinPrometheusHandler()
	handler.controller.Logger = logger
	handler.controller.Start()

	http.Handle("/metrics", handler)
	http.ListenAndServe(":2112", nil)
}

type daikinPrometheusHandler struct {
	registry          *prometheus.Registry
	prometheusHandler http.Handler
	metrics           daikinMetrics
	controller        echonetlite.Controller
	daikin            daikin.Daikin
}

func newDaikinPrometheusHandler() *daikinPrometheusHandler {
	handler := &daikinPrometheusHandler{
		registry:   prometheus.NewRegistry(),
		controller: echonetlite.NewController(),
	}
	handler.metrics = newDaikinMetrics(handler.registry)
	handler.prometheusHandler = promhttp.HandlerFor(handler.registry, promhttp.HandlerOpts{Registry: handler.registry})
	handler.daikin = daikin.NewDaikin(&handler.controller)
	return handler
}

func (handler *daikinPrometheusHandler) Close() error {
	return handler.controller.Close()
}

func (handler *daikinPrometheusHandler) updateMetrics() error {
	resps, err := handler.daikin.Request().
		IdentificationNumber().
		OperationStatus().
		InstantaneousPowerConsumption().
		CumulativePowerConsumption().
		FaultStatus().
		AirflowRate().
		OperationMode().
		TemperatureSetting().
		HumiditySetting().
		RoomTemperature().
		RoomHumidity().
		OutdoorTemperature().
		DehumidifyingSetting().
		Query()
	if err != nil {
		return err
	}

	handler.metrics.operationStatus.Reset()
	handler.metrics.instantaneousPowerConsumption.Reset()
	handler.metrics.cumulativePowerConsumption.Reset()
	handler.metrics.faultStatus.Reset()
	handler.metrics.airflowRateAuto.Reset()
	handler.metrics.airflowRateSetting.Reset()
	handler.metrics.operationModeSetting.Reset()
	handler.metrics.temperatureSetting.Reset()
	handler.metrics.humiditySetting.Reset()
	handler.metrics.dehumidifyingSetting.Reset()
	handler.metrics.roomHumidity.Reset()
	handler.metrics.roomTemperature.Reset()
	handler.metrics.outdoorTemperature.Reset()

	slog.Debug("[updateMetrics] responses retrieved", "device_count", len(resps))
	for _, resp := range resps {
		id, err := resp.IdentificationNumber()
		if err != nil {
			return err
		}
		idstr := bytesToString(id)

		if err := updateBoolMetrics(resp.Address, idstr, resp.OperationStatus, handler.metrics.operationStatus); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "OperationStatus", "error", err)
		}

		if err := updateNumberMetrics(resp.Address, idstr, resp.InstantaneousPowerConsumption, handler.metrics.instantaneousPowerConsumption); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "InstantaneousPowerConsumption", "error", err)
		}

		if err := updateNumberMetrics(resp.Address, idstr, resp.CumulativePowerConsumption, handler.metrics.cumulativePowerConsumption); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "CumulativePowerConsumption", "error", err)
		}

		if err := updateBoolMetrics(resp.Address, idstr, resp.FaultStatus, handler.metrics.faultStatus); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "FaultStatus", "error", err)
		}

		if err := updateNumberWithAutoMetrics(resp.Address, idstr, resp.AirflowRate, handler.metrics.airflowRateAuto, handler.metrics.airflowRateSetting); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "AirflowRate", "error", err)
		}

		if err := updateOperationMode(resp.Address, idstr, resp.OperationMode, handler.metrics.operationModeSetting); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "OperationMode", "error", err)
		}

		if err := updateNumberMetrics(resp.Address, idstr, resp.TemperatureSetting, handler.metrics.temperatureSetting); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "TemperatureSetting", "error", err)
		}

		if err := updateNumberMetrics(resp.Address, idstr, resp.HumiditySetting, handler.metrics.humiditySetting); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "HumiditySetting", "error", err)
		}

		if err := updateNumberMetrics(resp.Address, idstr, resp.RoomTemperature, handler.metrics.roomTemperature); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "RoomTemperature", "error", err)
		}

		if err := updateNumberMetrics(resp.Address, idstr, resp.RoomHumidity, handler.metrics.roomHumidity); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "RoomHumidity", "error", err)
		}

		if err := updateNumberMetrics(resp.Address, idstr, resp.OutdoorTemperature, handler.metrics.outdoorTemperature); err != nil {
			slog.Info("[updateMetrics] update failed", "id", idstr, "property", "OutdoorTemperature", "error", err)
		}
	}
	return nil
}

type daikinMetrics struct {
	operationStatus               *prometheus.GaugeVec
	instantaneousPowerConsumption *prometheus.GaugeVec
	cumulativePowerConsumption    *prometheus.GaugeVec
	faultStatus                   *prometheus.GaugeVec
	airflowRateAuto               *prometheus.GaugeVec
	airflowRateSetting            *prometheus.GaugeVec
	operationModeSetting          *prometheus.GaugeVec
	temperatureSetting            *prometheus.GaugeVec
	humiditySetting               *prometheus.GaugeVec
	dehumidifyingSetting          *prometheus.GaugeVec
	roomTemperature               *prometheus.GaugeVec
	roomHumidity                  *prometheus.GaugeVec
	outdoorTemperature            *prometheus.GaugeVec
}

func newDaikinMetrics(reg prometheus.Registerer) daikinMetrics {
	commonLabels := []string{"address", "id"}

	metrics := daikinMetrics{
		operationStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "operation_status", Help: "operation status (1:on, 0:off)"},
			commonLabels,
		),
		instantaneousPowerConsumption: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "instantaneous_power_consumption", Help: "instantaneous power consumption (unit:W)"},
			commonLabels,
		),
		cumulativePowerConsumption: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "cumulative_power_consumption", Help: "cumulative power consumption (unit:Wh)"},
			commonLabels,
		),
		faultStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "fault_status", Help: "fault status (1:on, 0:off)"},
			commonLabels,
		),
		airflowRateAuto: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "airflow_rate_auto", Help: "airflow rate (1:auto, 0:manual)"},
			commonLabels,
		),
		airflowRateSetting: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "airflow_rate_setting", Help: "airflow rate (1-8)"},
			commonLabels,
		),
		operationModeSetting: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "operation_mode_setting", Help: "operation mode (1:on, 0:off)"},
			append(commonLabels, "mode"),
		),
		temperatureSetting: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "temperature_setting", Help: "temperature setting (0-50 degree(s) Celsius)"},
			commonLabels,
		),
		humiditySetting: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "humidity_setting", Help: "humidity setting (0-100%)"},
			commonLabels,
		),
		dehumidifyingSetting: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "dehumidifying_setting", Help: "dehumidifying setting (0-100%)"},
			commonLabels,
		),
		roomTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "room_temperature", Help: "room temperature (-127 to 125 degree(s) Celsius)"},
			commonLabels,
		),
		roomHumidity: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "room_humidity", Help: "room humidity (0-100%)"},
			commonLabels,
		),
		outdoorTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "outdoor_temperature", Help: "outdoor temperature (-127 to 125 degree(s) Celsius)"},
			commonLabels,
		),
	}

	reg.MustRegister(metrics.operationStatus)
	reg.MustRegister(metrics.instantaneousPowerConsumption)
	reg.MustRegister(metrics.cumulativePowerConsumption)
	reg.MustRegister(metrics.faultStatus)
	reg.MustRegister(metrics.airflowRateAuto)
	reg.MustRegister(metrics.airflowRateSetting)
	reg.MustRegister(metrics.operationModeSetting)
	reg.MustRegister(metrics.temperatureSetting)
	reg.MustRegister(metrics.humiditySetting)
	reg.MustRegister(metrics.dehumidifyingSetting)
	reg.MustRegister(metrics.roomHumidity)
	reg.MustRegister(metrics.roomTemperature)
	reg.MustRegister(metrics.outdoorTemperature)

	return metrics
}

func (handler *daikinPrometheusHandler) ServeHTTP(response http.ResponseWriter, req *http.Request) {
	handler.updateMetrics()
	handler.prometheusHandler.ServeHTTP(response, req)
}
