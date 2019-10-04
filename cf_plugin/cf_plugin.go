package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
	token_store "github.com/uaa_token_store/token_store"
)

type UaaTsPlugin struct {
}

func (c *UaaTsPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	fmt.Println("Running the UAA Token Store")
	token_store.EntryPoint()
}

func (c *UaaTsPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "UAA Token Store",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "uaats",
				HelpText: "UAA Token Store plugin records and refreshes the UAA Token",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "uaats\n   cf uaats --add \n cf uaats --add-from-file \n cf uaats",
				},
			},
		},
	}
}
func main() {
	// Any initialization for your plugin can be handled here
	//
	// Note: to run the plugin.Start method, we pass in a pointer to the struct
	// implementing the interface defined at "code.cloudfoundry.org/cli/plugin/plugin.go"
	//
	// Note: The plugin's main() method is invoked at install time to collect
	// metadata. The plugin will exit 0 and the Run([]string) method will not be
	// invoked.
	plugin.Start(new(UaaTsPlugin))
	// Plugin code should be written in the Run([]string) method,
	// ensuring the plugin environment is bootstrapped.
}
