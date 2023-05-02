# metrics-exporter

[![Build status](https://github.com/BlackLotus-SMP/metrics-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/BlackLotus-SMP/metrics-exporter/actions/workflows/build.yml)

Simple exporter for the minecraft metrics protocol mod: https://github.com/BlackLotus-SMP/metrics-protocol

## Docker deploy

**You can find the docs on how to deploy the service [here](https://github.com/BlackLotus-SMP/MetricsExporter/blob/master/docs/README.md).**

## Start

```bash
./metrics-exporter -p 8855 -interval 15 -mcAddress 127.0.0.1 -mcPort 25565
```
- `-p`
  - **port the service will be listening to**
  - **default** = 8462
- `-interval`
  - **in seconds, the interval of time the service will check for new minecraft metrics**
  - **default** = 30
- `-mcAddress`
  - **address (ip/dns) of the mc server**
  - **default** = 127.0.0.1
- `-mcPort`
    - **port of the mc server**
    - **default** = 25565

## Endpoints

### Metrics
- **endpoint**: `/metrics`
- **command**: `curl 127.0.0.1:{port}/metrics`
- **run**: returns the minecraft server metrics in a prometheus compatible format

### HealthCheck
- **endpoint**: `/healthcheck`
- **command**: `curl 127.0.0.1:{port}/healthcheck`
- **run**: Just returns 200, this is for docker/kubernetes integration