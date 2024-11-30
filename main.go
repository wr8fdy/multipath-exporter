package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type multipathCollector struct {
	mapInfo        *prometheus.Desc
	mapState       *prometheus.Desc
	mapFaults      *prometheus.Desc
	groupState     *prometheus.Desc
	pathInfo       *prometheus.Desc
	pathState      *prometheus.Desc
	pathDevState   *prometheus.Desc
	pathCheckState *prometheus.Desc
}

func newMultipathCollector() *multipathCollector {
	return &multipathCollector{
		mapInfo: prometheus.NewDesc("multipath_map_info",
			"Shows multipath map info",
			[]string{"name", "uuid", "vendor", "paths", "sysfs"}, nil,
		),
		mapState: prometheus.NewDesc("multipath_map_state",
			"Shows multipath map state",
			[]string{"uuid", "state"}, nil,
		),
		mapFaults: prometheus.NewDesc("multipath_map_faults",
			"Shows multipath map number of faults",
			[]string{"uuid"}, nil,
		),
		groupState: prometheus.NewDesc("multipath_group_state",
			"Shows multipath group info and state",
			[]string{"map_uuid", "group_id", "state"}, nil,
		),
		pathInfo: prometheus.NewDesc("multipath_path_info",
			"Shows multipath path info",
			[]string{"map_uuid", "group_id", "dev", "target_wwnn", "host_adapter"}, nil,
		),
		pathState: prometheus.NewDesc("multipath_path_state",
			"Shows multipath path state",
			[]string{"map_uuid", "group_id", "dev", "state"}, nil,
		),
		pathDevState: prometheus.NewDesc("multipath_path_device_state",
			"Shows multipath path device state",
			[]string{"map_uuid", "group_id", "dev", "state"}, nil,
		),
		pathCheckState: prometheus.NewDesc("multipath_path_check_state",
			"Shows multipath path check state",
			[]string{"map_uuid", "group_id", "dev", "state"}, nil,
		),
	}
}

func (collector *multipathCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.mapInfo
	ch <- collector.mapState
	ch <- collector.mapFaults
	ch <- collector.groupState
	ch <- collector.pathInfo
	ch <- collector.pathState
	ch <- collector.pathDevState
	ch <- collector.pathCheckState
}

func (collector *multipathCollector) Collect(ch chan<- prometheus.Metric) {
	out, err := exec.Command("multipathd", "show", "maps", "json").Output()
	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
		case *exec.ExitError:
			fmt.Println("command exit rc =", e.ExitCode())
		default:
			panic(err)
		}
	}

	var parsedOut ShowMapsOutput
	err = json.Unmarshal(out, &parsedOut)
	if err != nil {
		fmt.Println("failed parsing output:", err)
	}

	for _, device := range parsedOut.Maps {
		mapInfo := prometheus.MustNewConstMetric(collector.mapInfo, prometheus.GaugeValue, 1.0,
			device.Name, device.UUID, device.Vend, strconv.Itoa(device.Paths), device.Sysfs)
		ch <- mapInfo

		mapState := prometheus.MustNewConstMetric(collector.mapState, prometheus.GaugeValue, 1.0,
			device.Name, device.DmSt)
		ch <- mapState

		mapFaults := prometheus.MustNewConstMetric(collector.mapFaults, prometheus.GaugeValue,
			float64(device.PathFaults), device.UUID)
		ch <- mapFaults

		for _, group := range device.PathGroups {
			groupState := prometheus.MustNewConstMetric(collector.groupState, prometheus.GaugeValue, 1.0,
				device.UUID, strconv.Itoa(group.Group), group.DmSt)
			ch <- groupState

			for _, path := range group.Paths {
				pathInfo := prometheus.MustNewConstMetric(collector.pathInfo, prometheus.GaugeValue, 1.0,
					device.UUID, strconv.Itoa(group.Group), path.Dev, path.TargetWwnn, path.HostAdapter)
				ch <- pathInfo

				pathState := prometheus.MustNewConstMetric(collector.pathState, prometheus.GaugeValue, 1.0,
					device.UUID, strconv.Itoa(group.Group), path.Dev, path.DmSt)
				ch <- pathState

				pathDevState := prometheus.MustNewConstMetric(collector.pathDevState, prometheus.GaugeValue, 1.0,
					device.UUID, strconv.Itoa(group.Group), path.Dev, path.DevSt)
				ch <- pathDevState

				pathCheckState := prometheus.MustNewConstMetric(collector.pathCheckState, prometheus.GaugeValue, 1.0,
					device.UUID, strconv.Itoa(group.Group), path.Dev, path.ChkSt)
				ch <- pathCheckState
			}
		}
	}

}

var listenAddress = flag.String("web.listen-address", "9101", "Listen address")

func main() {
	flag.Parse()

	c := newMultipathCollector()
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *listenAddress), nil))
}
