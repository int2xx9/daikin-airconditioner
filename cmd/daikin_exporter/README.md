daikin_exporter
==================================================

A prometheus exporter for daikin air conditioners.

## Tested devices

This exporter is tested on these devices.

- AN223ARS-W (FW: 2.11.0-g4)
- AN253ARS-W (FW: 2.11.0-g4)

## Usage

An exporter will automatically detect all air conditioners on the network. You just need to build and run the exporter.

```
git clone https://github.com/int2xx9/daikin-airconditioner
cd daikin-airconditioner/cmd/daikin_exporter
go build .
./daikin_exporter
```

### Options

- `--port` (default: 2112)
  - a port number to access to an exporter

## Metrics

### Labels

All metrics have these labels:

| label | summary |
|-|-|
| address | a device's IP address |
| id | a device's identification number obtained through the echonet lite property 0x83 |

### Metrics

All available metrics:

| metric | echonet lite epc | summary |
|-|-|-|
| operation_status                | 80 | operation status (1:on, 0:off) |
| instantaneous_power_consumption | 84 | instantaneous power consumption (unit:W) |
| cumulative_power_consumption    | 85 | cumulative power consumption (unit:Wh) |
| fault_status                    | 88 | fault status (1:on, 0:off) |
| airflow_rate_auto               | a0 | airflow rate (1:auto, 0:manual) |
| airflow_rate_setting            | a0 | airflow rate (1-8) |
| operation_mode_setting          | b0 | operation mode (1:on, 0:off) |
| temperature_setting             | b3 | temperature setting (0-50 degree(s) Celsius) |
| humidity_setting                | b4 | humidity setting (0-100%) |
| room_humidity                   | ba | room temperature (-127 to 125 degree(s) Celsius) |
| room_temperature                | bb | room humidity (0-100%) |
| outdoor_temperature             | be | outdoor temperature (-127 to 125 degree(s) Celsius) |
