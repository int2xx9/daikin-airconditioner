recorder
==================================================

## Usage

### install

```shell
go install github.com/int2xx9/daikin-airconditioner/cmd/recorder@latest

# or

git clone https://github.com/int2xx9/daikin-airconditioner
cd cmd/recorder
go build -o recorder .
```

### setup

1. Prepare config.yaml (copy and edit config.example.yaml)
2. Run `recorder migrate`

```shell
recorder migrate
```

3. Add devices

```shell
$ recorder device discover
2 device(s) are discovered.
IPAddress       ID
192.168.0.3     XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
192.168.0.4     YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY
$ recorder device add XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX RoomA
$ recorder device add YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY RoomB
```

### Start recording

```shell
recorder record
```
