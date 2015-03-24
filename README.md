# setstatsd
A small daemon to count things unique. Collects values in sets and periodically send their cardinality to [InfluxDb](http://influxdb.com/).

Typical use cases are aggregating unique IPs or user Ids seen across a cluster of servers. For more general aggregation of counters etc,
see [statsd](https://github.com/etsy/statsd/).

### Usage
```bash
$ ./setstatsd -h
Usage of ./setstatsd:
  -p="9010": port to listen to
  -host="localhost": InfluxDB Host
  -port="8086": InfluxDB Port
  -db="metrics": InfluxDB Database
  -password="metrics": InfluxDB Password
  -user="metrics": InfluxDB User
  -interval=10s: Interval between reports to InfluxDB
```

### Sending data to setstatsd
The daemon listens for HTTP POSTs to /[set-name] and expects values separated by '\n'.

```
curl --data $'Value1\nValue2\nValue3' localhost:9010/my-set-of-uniques
```

To peek at metrics collected since last snapshot sent report:
```
$ curl localhost:9010/dump
Sets (and their size) seen since last report to InfluxDB

my-set-of-uniques: 3
```

### InfluxDB
Currently supports InfluxDb 0.8, will be updated to the new protocol (with tagsets) as soon as 0.9 is released.

## License
The MIT License (MIT)

Copyright (c) 2015 Andreas Bielk

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
