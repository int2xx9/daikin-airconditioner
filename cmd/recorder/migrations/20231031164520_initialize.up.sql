create table devices(
  -- 16bytes(128bits)
  id char(32) primary key not null,
  name varchar(128)
);

create type operation_mode_type as enum('auto', 'cooling', 'heating', 'dehumidification', 'ventilating', 'other');

create table records(
  device_id char(32) not null,
  time timestamp not null,

  operation_status boolean,
  instantaneous_power_consumption int,
  cumulative_power_consumption int,
  fault_status boolean,
  operation_mode operation_mode_type,
  airflowrate_auto boolean,
  airflowrate_setting int,
  temperature_setting int,
  humidity_setting int,
  room_humidity int,
  room_temperature int,
  outdoor_temperature int,

  primary key (device_id, time)
);

select create_hypertable('records', 'time');
