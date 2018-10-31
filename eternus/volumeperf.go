package eternus

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

type VolumePerfStat struct {
	Name            string
	IopsRead        int64
	IopsWrite       int64
	ThroughputRead  int64
	ThroughputWrite int64
	ResptimeRead    int64
	ResptimeWrite   int64
	ProctimeRead    int64
	ProctimeWrite   int64
}

type VolumePerfCollector struct {
	stat []VolumePerfStat
}

func (c *VolumePerfCollector) GetMetricTypes(plugin.Config) []plugin.Metric {
	var mts []plugin.Metric
	for _, m := range []string{"iops_read", "iops_write", "throughput_read", "throughput_write", "resptime_read", "resptime_write", "proctime_read", "proctime_write"} {
		namespace := plugin.NewNamespace(PluginVedor, PluginName, "volume")
		namespace = namespace.AddDynamicElement("name", "component name")
		namespace = namespace.AddStaticElement(m)
		mts = append(mts, plugin.Metric{Namespace: namespace})
	}
	return mts
}

func (c *VolumePerfCollector) CollectMetrics(exec Executor, cliprefix string, metric plugin.Metric) ([]plugin.Metric, error) {
	var mts []plugin.Metric

	now := time.Now()

	if metric.Namespace[2].Value != "volume" {
		return nil, nil
	}

	out, err := exec.Execute("show performance -type host-io")
	if err != nil {
		return nil, err
	}
	if c.stat == nil {
		err = c.parseStat(out, cliprefix)
		if err != nil {
			return nil, err
		}
	}

	for _, stat := range c.stat {
		ns := plugin.NewNamespace()
		tags := make(map[string]string)
		ns = append(ns, metric.Namespace...)
		for k, v := range metric.Tags {
			tags[k] = v
		}
		m := plugin.Metric{Namespace: ns, Timestamp: now, Tags: tags}
		m.Namespace[3].Value = stat.Name
		switch m.Namespace[4].Value {
		case "iops_read":
			m.Data = stat.IopsRead
		case "iops_write":
			m.Data = stat.IopsWrite
		case "throughput_read":
			m.Data = stat.ThroughputRead
		case "throughput_write":
			m.Data = stat.ThroughputWrite
		case "resptime_read":
			m.Data = stat.ResptimeRead
		case "resptime_write":
			m.Data = stat.ResptimeWrite
		case "proctime_read":
			m.Data = stat.ProctimeRead
		case "proctime_write":
			m.Data = stat.ProctimeWrite
		}
		mts = append(mts, m)
	}
	return mts, nil
}

func (c *VolumePerfCollector) Reset() {
	c.stat = nil
}

func (c *VolumePerfCollector) parseStat(out, prefix string) error {
	c.stat = []VolumePerfStat{}
	re := regexp.MustCompile("-{2,}")
	lines := strings.Split(out, "\n")
	var header [][]int
	gotHeader := false
	for _, line := range lines {
		if !gotHeader {
			header = re.FindAllStringIndex(line, -1)
			if len(header) == 0 {
				continue
			}
			if len(header) == 13 {
				gotHeader = true
				continue
			} else {
				return fmt.Errorf("volume performance: looking for 13 fields, got %v", len(header))
			}
		}
		if strings.TrimSpace(line) == "" || strings.Contains(line, prefix) {
			continue
		}
		name := strings.TrimSpace(line[header[1][0]:header[1][1]])
		ioread, _ := strconv.ParseInt(strings.TrimSpace(line[header[2][0]:header[2][1]]), 10, 32)
		iowrite, _ := strconv.ParseInt(strings.TrimSpace(line[header[3][0]:header[3][1]]), 10, 32)
		thrread, _ := strconv.ParseInt(strings.TrimSpace(line[header[4][0]:header[4][1]]), 10, 32)
		thrwrite, _ := strconv.ParseInt(strings.TrimSpace(line[header[5][0]:header[5][1]]), 10, 32)
		resptimeread, _ := strconv.ParseInt(strings.TrimSpace(line[header[6][0]:header[6][1]]), 10, 32)
		resptimewrite, _ := strconv.ParseInt(strings.TrimSpace(line[header[7][0]:header[7][1]]), 10, 32)
		proctimeread, _ := strconv.ParseInt(strings.TrimSpace(line[header[8][0]:header[8][1]]), 10, 32)
		proctimewrite, _ := strconv.ParseInt(strings.TrimSpace(line[header[9][0]:header[9][1]]), 10, 32)
		c.stat = append(c.stat, VolumePerfStat{
			Name:            strings.TrimSpace(name),
			IopsRead:        ioread,
			IopsWrite:       iowrite,
			ThroughputRead:  thrread,
			ThroughputWrite: thrwrite,
			ResptimeRead:    resptimeread,
			ResptimeWrite:   resptimewrite,
			ProctimeRead:    proctimeread,
			ProctimeWrite:   proctimewrite,
		})
	}
	return nil
}
