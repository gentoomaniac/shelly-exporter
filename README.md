# shelly-exporter

## tl;dr

This project is a prometheus exporter targeting shelly devices to collect metrics from them and providing them to prometheus

While this project started as a mainly focused on Shelly Plug S devices, the codebase turned out generic enough to be easily extended.

Currently it also contains support for the [Homewizard P1](https://www.homewizard.com/p1-meter/)

### Supported Devices

Not all metrics are generated right yet but the base support for the listed devices is available

#### Pull based

* SHPLG-S - Shelly Plug S (Tested only with the old dual color non bluetooth variant)
* HWE-P1 - [Homewizard P1](https://www.homewizard.com/p1-meter/)
* [Shelly Pro 3EM](https://www.shelly.com/products/shelly-pro-3em-x1)

#### Push based

The below devices are mostly in a sleep state and because of that can't be querried reliably.

The exporter offers webhooks for these that can be configured in the devices to retrieve the sensor data once the device wakes up.

* [Shelly Plus H&T](https://www.shelly.com/products/shelly-plus-h-t)
* WIP: [Shelly H&T](https://www.shelly.com/products/shelly-h-t-white)

### Planned suppport

* [Shelly FLood](https://www.shelly.com/products/shelly-flood)
* ShellyBLU devices

## How to run

### Configuration

See [config.yaml](config.yaml) for some config examples.

The configuration file supports environment variables for `username`, `password` and `frequency` similar to bash variables: `${env:VARIABLE_NAME:-default_value}` but only as the only value of a field. Mixed usage of variables and strings are currently not supported.

#### Devices with API

TODO: config example

#### Webhook

You can send arbitrary data to the exporter to allow for sleep state devices to send their data.

The URL is constructed like this, some [mandatory tags](https://github.com/gentoomaniac/shelly-exporter/blob/dbdcdcf266652e45f9bd85b1009ebbb22e45102d/pkg/exporter/webhook.go#L14) have to be specified for the exporter to function properly

```
https://exporter/webhook?tag=value&tag2=value2&metric=<metric_name>&value=<value>
```

example:

```
https://<exporter>:<port>/webhook?building=main&room=bedroom&type=PLUSHT&name=PlusHT%20GH%20Downstairs&deviceid=08B61FCEA4BC&namespace=shelly&metric=ambient_temperature_celsius&value=${ev.tC}
```

#### Legacy Webhook

This webhook is for old shelly devices that have a fixed list of parameters they send.

Below is an example of the path you cna configure to pass on arbitrary labels.

``` bash
http://127.0.0.1:8080/legacywebhook/location=test/label=fizz/label2=buzz/
# will cause a call like
http://127.0.0.1:8080/legacywebhook/location=test/label=fizz/label2/buzz/?hum=40.0&temp=22.1&id=afsefa
```

the resulting metrics will look like this

```
shelly_humidity{deviceId="afsefa",ip="::1",lable="fizz",lable2="buzz",type="SHHT-1",userAgent="Shelly/20230913-112531/v1.14.0-gcb84623 (SHHT-1)"} 40
shelly_temperature{deviceId="afsefa",ip="::1",lable="fizz",lable2="buzz",type="SHHT-1",userAgent="Shelly/20230913-112531/v1.14.0-gcb84623 (SHHT-1)"} 22.1
```

### Run the exporter

``` bash
docker run -v "$(pwd)/config.yaml:/config.yaml" ghcr.io/gentoomaniac/shelly-exporter:latest --config-file /config.yaml -vv
```

## How to extend the exporter

* add a new package to the codebase implementing the `Device` interface:

```go
type Device interface {
    Collectors() ([]prometheus.Collector, error)
    Name() string
    Refresh() error
    RefreshDeviceinfo() error
}
```

* add the device type to the [`config package`](https://github.com/gentoomaniac/shelly-exporter/blob/69c63f8b3b413b9e60ab968e17a687fdbafcc849/pkg/config/config.go#L15-L31)

* add the instantiation code to the exporter [setup function](https://github.com/gentoomaniac/shelly-exporter/blob/69c63f8b3b413b9e60ab968e17a687fdbafcc849/pkg/exporter/exporter.go#L115-L126)
