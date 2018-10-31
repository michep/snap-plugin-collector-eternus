package eternus

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

type DiskBusyStat struct {
	Name        string
	BusyPercent int64
}

type DiskBusyCollector struct {
	stat []DiskBusyStat
}

func (c *DiskBusyCollector) GetMetricTypes(plugin.Config) []plugin.Metric {
	var mts []plugin.Metric
	for _, m := range []string{"busy_percent"} {
		namespace := plugin.NewNamespace(PluginVedor, PluginName, "disk")
		namespace = namespace.AddDynamicElement("name", "component name")
		namespace = namespace.AddStaticElement(m)
		mts = append(mts, plugin.Metric{Namespace: namespace})
	}
	return mts
}

func (c *DiskBusyCollector) CollectMetrics(exec Executor, cliprefix string, metric plugin.Metric) ([]plugin.Metric, error) {
	var mts []plugin.Metric

	now := time.Now()

	if !(metric.Namespace[2].Value == "disk" && metric.Namespace[4].Value == "busy_percent") {
		return nil, nil
	}

	out, err := exec.Execute("show performance -type disks")
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
		m.Data = stat.BusyPercent
		mts = append(mts, m)
	}
	return mts, nil
}

func (c *DiskBusyCollector) Reset() {
	c.stat = nil
}

func (c *DiskBusyCollector) parseStat(out, prefix string) error {
	c.stat = []DiskBusyStat{}
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
			if len(header) == 2 {
				gotHeader = true
				continue
			} else {
				return fmt.Errorf("disk busy: looking for 2 fields, got %v", len(header))
			}
		}
		if strings.TrimSpace(line) == "" || strings.Contains(line, prefix) {
			continue
		}
		name := strings.TrimSpace(line[header[0][0]:header[0][1]])
		percent, _ := strconv.ParseInt(strings.TrimSpace(line[header[1][0]:header[1][1]]), 10, 64)
		c.stat = append(c.stat, DiskBusyStat{Name: name, BusyPercent: percent})
	}
	return nil
}
