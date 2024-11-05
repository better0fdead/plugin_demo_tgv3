package context

import (
	"reflect"
	"strings"
)

type PluginCtx struct {
    // Plugin flags
	ServicePath        string

	// Flags for tg usage
    // DO NOT MODIFY
	Help        bool
	Doc         bool
	Source      bool
	Description bool
	Flags       bool
}

func SetPluginCtx(flags map[string]string) PluginCtx {
	// Plugin flags
	pluginCtx := PluginCtx{
		ServicePath:        flags["ServicePath"],
	}

	// Flags for tg usage
    // DO NOT MODIFY
	if _, exists := flags["help"]; exists {
		pluginCtx.Help = true
	}
	if _, exists := flags["h"]; exists {
		pluginCtx.Help = true
	}
	if _, exists := flags["doc"]; exists {
		pluginCtx.Doc = true
	}
	if _, exists := flags["source"]; exists {
		pluginCtx.Source = true
	}
	if _, exists := flags["desc"]; exists {
		pluginCtx.Description = true
	}
	if _, exists := flags["flags"]; exists {
		pluginCtx.Flags = true
	}
	return pluginCtx
}

// GetPluginFlags returns an array of field names of PluginCtx
func GetPluginFlags() []string {
	var fields []string

	ctx := PluginCtx{}

	val := reflect.TypeOf(ctx)
	serviceFlagsStr := "Help" + " Doc " + " Source " + " Description " + " Flags" + "Version"

	// Iterate over the struct fields and append their names to the fields slice
	for i := 0; i < val.NumField(); i++ {
		name := val.Field(i).Name
		if !strings.Contains(serviceFlagsStr, name) {
			fields = append(fields, name)
		}

	}

	return fields
}
