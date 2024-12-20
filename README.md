# shelly-exporter

## tl;dr

This project is a prometheus exporter targeting shelly devices to collect metrics from them and providing them to prometheus

While this project started as a mainly focused on Shelly Plug S devices, the codebase turned out generic enough to be easily extended.

Currently it also contains support for the [Homewizard P1](https://www.homewizard.com/p1-meter/)

## Configuration

See [config.yaml](config.yaml) for some config examples.

The configuration file supports environment variables for `username`, `password` and `frequency` similar to bash variables: `${env:VARIABLE_NAME:-default_value}` but only as the only value of a field. Mixed usage of variables and strings are currently not supported.

### Supported devices

* SHPLG-S - Shelly Plug S (Tested only with the old dual color non bluetooth variant)
* HWE-P1 - [Homewizard P1](https://www.homewizard.com/p1-meter/)

### Planned suppport

* [Shelly Pro 3EM](https://www.shelly.com/products/shelly-pro-3em-x1)
* [Shelly Plus H&T](https://www.shelly.com/products/shelly-plus-h-t)
* [Shelly H&T](https://www.shelly.com/products/shelly-h-t-white)
* [Shelly FLood](https://www.shelly.com/products/shelly-flood)

## How to run

``` bash
docker run -v "$(pwd)/config.yaml:/config.yaml" ghcr.io/gentoomaniac/shelly-exporter:latest --config-file /config.yaml -vv
```

## Planned features

I'm currently working on a webhook that allows Shelly sensors to send their current measurements to the exporter.

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