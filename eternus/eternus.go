package eternus

import (
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"strings"
)

const (
	PluginVedor    = "mfms"
	PluginName     = "eternus"
	PluginVersion  = 1
	paramHost      = "host"
	paramUser      = "user"
	paramPassword  = "password"
	paramCipher    = "cipher"
	paramCLIPrefix = "cli_prefix"
	paramWaitTime  = "wait_time"
)

type EternusCollector interface {
	GetMetricTypes(plugin.Config) []plugin.Metric
	CollectMetrics(Executor, string, plugin.Metric) ([]plugin.Metric, error)
	Reset()
}

type Plugin struct {
	initialized bool
	host        string
	user        string
	password    string
	cipher      string
	cliprefix   string
	waittime    int64
	stats       []EternusCollector
	system      string
}

func NewCollector(stats ...EternusCollector) *Plugin {
	return &Plugin{stats: stats}
}

func (p *Plugin) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{PluginVedor, PluginName}, paramHost, true)
	policy.AddNewStringRule([]string{PluginVedor, PluginName}, paramUser, true)
	policy.AddNewStringRule([]string{PluginVedor, PluginName}, paramPassword, true)
	policy.AddNewStringRule([]string{PluginVedor, PluginName}, paramCipher, false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{PluginVedor, PluginName}, paramCLIPrefix, false, plugin.SetDefaultString("CLI>"))
	policy.AddNewIntRule([]string{PluginVedor, PluginName}, paramWaitTime, false, plugin.SetDefaultInt(200))
	return *policy, nil
}

func (p *Plugin) GetMetricTypes(config plugin.Config) ([]plugin.Metric, error) {
	var mts []plugin.Metric

	for _, stat := range p.stats {
		mts = append(mts, stat.GetMetricTypes(config)...)
	}

	return mts, nil
}

func (p *Plugin) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	var mts []plugin.Metric

	if !p.initialized {
		config := metrics[0].Config
		p.host, _ = config.GetString(paramHost)
		p.user, _ = config.GetString(paramUser)
		p.password, _ = config.GetString(paramPassword)
		p.cipher, _ = config.GetString(paramCipher)
		p.cliprefix, _ = config.GetString(paramCLIPrefix)
		p.waittime, _ = config.GetInt(paramWaitTime)
		p.system = strings.Split(p.host, ":")[0]
		p.initialized = true
	}

	exec, err := NewExecutor(p.host, p.user, p.password, p.cipher, p.waittime)
	if err != nil {
		return nil, err
	}
	defer exec.Disconnect()

	for _, stat := range p.stats {
		stat.Reset()
		for _, metric := range metrics {
			metric.Tags = map[string]string{"system": p.system}
			mt, err := stat.CollectMetrics(exec, p.cliprefix, metric)
			if err != nil {
				return nil, err
			}
			mts = append(mts, mt...)
		}
	}

	return mts, nil
}
