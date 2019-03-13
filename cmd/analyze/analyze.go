package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
)

var (
	target  = flag.String("target", "", "Should be either debezium or streams")
	logfile = flag.String("logfile", "", "Path to the log file to analyze")
)

type Metric struct {
	AvgLatency                  float64
	NinetyPercentileLatency     float64
	NinetyNinePercentileLatency float64
	MaxLatency                  int64
	MinLatency                  int64
	sumLatency                  int64
	eventCount                  int64
}

func main() {
	flag.Parse()
	if *target == "" || *logfile == "" {
		flag.Usage()
		return
	}

	var analysisFunc func(string, *Metric) (int64, error)
	switch *target {
	case "debezium":
		analysisFunc = debezium
	case "streams":
		analysisFunc = streams
	default:
		panic("target should be one of debezium or streams")
	}

	file, err := os.Open(*logfile)
	if err != nil {
		panic(err)
	}
	defer func() {
		file.Close()
	}()

	reader := bufio.NewReader(file)
	var line string
	m := Metric{
		MaxLatency: math.MinInt64,
		MinLatency: math.MaxInt64,
	}
	var latencies []int64

	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		var latency int64
		latency, err = analysisFunc(line, &m)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		latencies = append(latencies, latency)
	}
	if err != io.EOF {
		panic(err)
	}

	// Calculate 90 percentile
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] > latencies[j]
	})

	m.AvgLatency = float64(m.sumLatency / m.eventCount)
	ninetyPercentileIndex := float64(m.eventCount) * float64(90.0/100.0)
	if ninetyPercentileIndex == float64(int64(ninetyPercentileIndex)) {
		i := int(ninetyPercentileIndex)
		m.NinetyPercentileLatency = float64(latencies[i-1])
	} else if ninetyPercentileIndex > 1 {
		i := int(ninetyPercentileIndex)
		m.NinetyPercentileLatency = float64((latencies[i-1] + latencies[i]) / 2.0)
	}

	ninetyNinePercentileIndex := float64(m.eventCount) * float64(99.0/100.0)
	if ninetyNinePercentileIndex == float64(int64(ninetyNinePercentileIndex)) {
		i := int(ninetyNinePercentileIndex)
		m.NinetyNinePercentileLatency = float64(latencies[i-1])
	} else if ninetyNinePercentileIndex > 1 {
		i := int(ninetyNinePercentileIndex)
		m.NinetyNinePercentileLatency = float64((latencies[i-1] + latencies[i]) / 2.0)
	}

	out, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}

func debezium(line string, m *Metric) (int64, error) {
	var jsonNode map[string]interface{}
	err := json.Unmarshal([]byte(line), &jsonNode)
	if err != nil {
		return 0, fmt.Errorf("%s %s", err, line)
	}
	debeziumTime, ok := jsonNode["ts_ms"].(float64)
	if !ok {
		// Skip
		return 0, fmt.Errorf("%s %#v", line, jsonNode)
	}
	debeziumTimeMs := int64(debeziumTime)
	innerNode, ok := jsonNode["source"].(map[string]interface{})
	if !ok {
		// Skip
		return 0, fmt.Errorf("%s %#v", line, innerNode)
	}
	transactionTimeMs := int64(innerNode["ts_usec"].(float64) / 1000000)

	lag := debeziumTimeMs - transactionTimeMs
	m.sumLatency += lag
	m.eventCount++
	if lag < m.MinLatency {
		m.MinLatency = lag
	}
	if lag > m.MaxLatency {
		m.MaxLatency = lag
	}
	return lag, nil
}

func streams(line string, m *Metric) (int64, error) {
	var jsonNode map[string]interface{}
	err := json.Unmarshal([]byte(line), &jsonNode)
	if err != nil {
		return 0, fmt.Errorf("%s %s", err, line)
	}

	payload, ok := jsonNode["payload"].(map[string]interface{})
	if !ok {
		// Skip
		return 0, fmt.Errorf("%s %#v", line, jsonNode)
	}

	debeziumTimeMs := int64(payload["debezium"].(float64))
	streamsTimeMs := int64(payload["streams"].(float64))

	lag := streamsTimeMs - debeziumTimeMs
	m.sumLatency += lag
	m.eventCount++
	if lag < m.MinLatency {
		m.MinLatency = lag
	}
	if lag > m.MaxLatency {
		m.MaxLatency = lag
	}
	return lag, nil
}
