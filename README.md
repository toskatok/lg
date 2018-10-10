# MQTT Load Generator :fire:

[![Travis branch](https://img.shields.io/travis/com/I1820/lg/master.svg?style=flat-square)](https://travis-ci.com/I1820/lg)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b6019269d2cf4b11897acd4340e78bdc)](https://www.codacy.com/app/i1820/lg?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=I1820/lg&amp;utm_campaign=Badge_Grade)

**Sometimes we tell lies, sometimes we prove we don't lie. Let's prove ourselves.**

## Introduction
This application gives a way for creating a load with MQTT, HTTP ...
MQTT is a messaging protocol and various platforms work with it so
creating a load with MQTT creates a way for testing platforms.

## Transport
LG have an awesome way for customizing portocol named `Transport`.
Transports give you a way for sending your generated data over your protocol with ease.

## Generator
LG is fully customizable so you can generate data with
your own structure and publish it on your own topic.
For this LG has the `Generator` interface is defined in `generator/`.

## Running
You can run LG with following syntax:

```sh
lg --destination http://127.0.0.1:1883 --rate 1ms
```

In `config.yml` you can specifies generator configuration. If your destination has scheme `http://` your transport is HTTP
and if it has `mqtt://` your transport is MQTT.

## Set I1820 on fire
### TTN (over HTTP)
With this load generator, we try to put a load on I1820 TTN Integration module in link component
and we get the following results.

```yml
generator:
  name: ttn
  info:
    applicationName: fan
    applicationID: 5ba3f19c87a142b0a840fae0
    devEUI: 000AE31955C049FC
    deviceName: agrinode
token: ttnIStheBEST
messages:
  - count: "{{.Count}}"
```

These results show generated parsed information ratio (number of parsed data / number of received data)
with send interval of data:

| Interval | Accept Ratio  |
|:---------|:-------------:|
| 1s       | 100%          |
| 100ms    | 100%          |
| 10ms     | 100%          |
| 1ms      | 100%          |
| 100us    | 100%          |
| 10us     | 100%          |

Please consider that HTTP requests cannot get their response in `100us` or lower interval so in these intervals there is no
difference with `1ms` interval.

### Fanco (over MQTT)
With this load generator, we try to put a load on I1820 MQTT service in link component
and we get the following results.

```yml
generator:
  name: fanco
  info:
    thingID: 5bbd104cefe940cb57dfeb76
token: 1BLzO2YYB1jH91pRB0cpeIdPMsM
messages:
  - count: "{{.Count}}"
```

These results show generated parsed information ratio (number of parsed data / number of received data)
with send interval of data:

| Interval | Accept Ratio  |
|:---------|:-------------:|
| 1s       | 100%          |
| 100ms    | 100%          |
| 10ms     | 100%          |
| 1ms      | 100%          |
| 100us    | 52.3809524%   |
| 10us     | 47.6190476%   |

### Test System

- 4 Core Intel(R) Core(TM) i7-5930K CPU @ 3.50GHz
- 8 Gb of RAM
- SSD
