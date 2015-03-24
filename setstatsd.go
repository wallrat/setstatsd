package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

import _ "expvar"

//TODO: expose and report some metrics about the service itself

var port = flag.String("p", "9010", "port to listen to")

var influxdbUrl string
var influxdbHost = flag.String("host", "localhost", "InfluxDB Host")
var influxdbPort = flag.String("port", "8086", "InfluxDB Port")
var influxdbDb = flag.String("db", "metrics", "InfluxDB Database")
var influxdbUser = flag.String("user", "metrics", "InfluxDB User")
var influxdbPassword = flag.String("password", "metrics", "InfluxDB Password")

var reportPeriod = flag.Duration("interval", 10*time.Second, "Interval between reports to InfluxDB")

var metrics map[string]map[string]bool
var mutex = &sync.Mutex{}

var OK = []byte("OK\n")

func init() {
	metrics = make(map[string]map[string]bool)
	flag.Parse()
	influxdbUrl = "http://" + *influxdbHost + ":" + *influxdbPort + "/db/" + *influxdbDb + "/series?u=" + *influxdbUser + "&p=" + *influxdbPassword
}

func main() {
	fmt.Println("Starting set stats daemon listening for HTTP on port " + *port)
	fmt.Printf("Posting metrics to %s each %v\n", influxdbUrl, *reportPeriod)

	// reporter
	ticker := time.NewTicker(*reportPeriod)
	go func() {
		for _ = range ticker.C {
			// snapshot metrics
			mutex.Lock()
			snapshot := metrics
			metrics = make(map[string]map[string]bool)
			mutex.Unlock()

			// send snapshot
			go storeMetrics(snapshot)
		}
	}()

	// HTTP stuff
	// expvar metrics is available at /debug/vars
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
	http.HandleFunc("/dump", dumpMetrics)
	http.HandleFunc("/", metricPostHandler)
	http.ListenAndServe(":"+*port, nil)
}

func storeMetrics(snapshot map[string]map[string]bool) {
	if len(snapshot) > 0 {
		// create some JSON
		buf := "["
		i := 0
		for k, v := range snapshot {
			if i > 0 {
				buf = buf + ","
			}
			m := fmt.Sprintf("{\"name\":\"%s\",\"columns\":[\"value\"],\"points\":[[%d]]}", k, len(v))
			buf = buf + m
			i++
		}
		buf = buf + "]"

		// POST to influx
		resp, err := http.Post(influxdbUrl, "application/json", strings.NewReader(buf))
		if err != nil {
			fmt.Printf("Error sending report to influx db error='%v'\n", err)
			return
		}
		if resp.StatusCode != 200 {
			fmt.Printf("Error sending report to influx db status='%s'\n", resp.Status)
			return
		}
	}
}

func dumpMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Sets (and their size) seen since last report to InfluxDB\n\n")
	for k, v := range metrics {
		fmt.Fprintf(w, "%s: %d\n", k, len(v))
	}
	fmt.Fprintf(w, "\n")
}

func metricPostHandler(w http.ResponseWriter, r *http.Request) {
	// make sure we got a POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request method (only POST allowed) "+r.Method)
		return
	}

	// read metric name
	metricName := r.URL.Path[1:]
	if len(metricName) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid metric name "+metricName)
		return
	}

	// read body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request body "+r.Method)
		return
	}

	// split body
	values := strings.Split(string(body), "\n")

	// update state
	mutex.Lock()
	for _, v := range values {
		set := metrics[metricName]
		if set == nil {
			set = make(map[string]bool)
			metrics[metricName] = set
		}
		set[v] = true
	}
	mutex.Unlock()

	// send a short response
	w.Write(OK)
}
