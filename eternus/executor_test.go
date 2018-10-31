package eternus

import (
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"testing"
)

func Test_Execute(t *testing.T) {
	exec, err := NewExecutor("172.16.18.67:22", "root", "NamEUd09af", "3des-cbc", 200)
	if err != nil {
		t.Fatal(err)
	}
	defer exec.Disconnect()
	out, err := exec.Execute("show disks")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
}

func Test_DiskHealthCollectMetrics(t *testing.T) {
	exec, err := NewExecutor("172.16.18.67:22", "root", "NamEUd09af", "3des-cbc", 200)
	if err != nil {
		t.Fatal(err)
	}
	defer exec.Disconnect()

	c := DiskHealthCollector{}
	namespace := plugin.NewNamespace(PluginVedor, PluginName, "disk")
	namespace = namespace.AddDynamicElement("name", "component name")
	namespace = namespace.AddStaticElement("health_percent")
	mts, err := c.CollectMetrics(exec, "CLI>", plugin.Metric{Namespace: namespace})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", mts)
}

func Test_DiskBusyCollectMetrics(t *testing.T) {
	exec, err := NewExecutor("172.16.18.67:22", "root", "NamEUd09af", "3des-cbc", 200)
	if err != nil {
		t.Fatal(err)
	}
	defer exec.Disconnect()

	c := DiskBusyCollector{}
	namespace := plugin.NewNamespace(PluginVedor, PluginName, "disk")
	namespace = namespace.AddDynamicElement("name", "component name")
	namespace = namespace.AddStaticElement("busy_percent")
	mts, err := c.CollectMetrics(exec, "CLI>", plugin.Metric{Namespace: namespace})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", mts)
}

func Test_VolumePerfCollectMetrics(t *testing.T) {
	exec, err := NewExecutor("172.16.18.67:22", "root", "NamEUd09af", "3des-cbc", 200)
	if err != nil {
		t.Fatal(err)
	}
	defer exec.Disconnect()

	c := VolumePerfCollector{}
	namespace := plugin.NewNamespace(PluginVedor, PluginName, "volume")
	namespace = namespace.AddDynamicElement("name", "component name")
	namespace = namespace.AddStaticElement("iops_read")
	mts, err := c.CollectMetrics(exec, "CLI>", plugin.Metric{Namespace: namespace})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", mts)
}

func Test_ControllerBusyCollectMetrics(t *testing.T) {
	exec, err := NewExecutor("172.16.18.67:22", "root", "NamEUd09af", "3des-cbc", 200)
	if err != nil {
		t.Fatal(err)
	}
	defer exec.Disconnect()

	c := ControllerBusyCollector{}
	namespace := plugin.NewNamespace(PluginVedor, PluginName, "controller")
	namespace = namespace.AddDynamicElement("name", "component name")
	namespace = namespace.AddStaticElement("busy_percent")
	mts, err := c.CollectMetrics(exec, "CLI>", plugin.Metric{Namespace: namespace})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", mts)
}

func Test_PortPerfCollectMetrics(t *testing.T) {
	exec, err := NewExecutor("172.16.18.67:22", "root", "NamEUd09af", "3des-cbc", 200)
	if err != nil {
		t.Fatal(err)
	}
	defer exec.Disconnect()

	c := PortPerfCollector{}
	namespace := plugin.NewNamespace(PluginVedor, PluginName, "port")
	namespace = namespace.AddDynamicElement("name", "component name")
	namespace = namespace.AddStaticElement("iops_read")
	mts, err := c.CollectMetrics(exec, "CLI>", plugin.Metric{Namespace: namespace})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", mts)
}
