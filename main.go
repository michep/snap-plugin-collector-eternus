package main

import (
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	"github.com/michep/snap-plugin-collector-eternus/eternus"
)

func main() {
	coll := eternus.NewCollector(
		&eternus.DiskHealthCollector{},
		&eternus.DiskBusyCollector{},
		&eternus.VolumePerfCollector{},
		&eternus.ControllerBusyCollector{},
		&eternus.PortPerfCollector{},
	)
	plugin.StartCollector(coll, eternus.PluginName, eternus.PluginVersion, plugin.RoutingStrategy(plugin.StickyRouter))
}
