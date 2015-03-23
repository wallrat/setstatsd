# setstatsd
A small daemon to collect unique members of sets and send their cardinality as metrics to [InfluxDb](http://influxdb.com/).

### Usage
```bash
$ ./setstatsd -h
Usage of ./setstatsd:
  -p="9010": port to listen to
  -host="192.168.10.10": InfluxDB Host
  -port="8086": InfluxDB Port
  -db="metrics": InfluxDB Database
  -password="metrics": InfluxDB Password
  -user="metrics": InfluxDB User
```

### Sending data to setstatsd
The daemon listens for HTTP POSTs to /[set-name] and expects values separated by \n.

```
curl -i --data $'Value1\nValue2\nValue3' localhost:9010/my-set-of-uniques
```

To peek at collected metrics since last sent report:
```
$ curl localhost:9010/dump
Sets (and their size) seen since last report to InfluxDB

my-set-of-uniques: 3
```

### InfluxDB
Currently supports InfluxDb 0.8, will be updated to the new protocol (with tagsets) as soon as 0.9 is released.

## License
The MIT License (MIT), Copyright (c) 2015 Andreas Bielk

See the file [LICENSE](LICENSE) for full text and copyright notice.
