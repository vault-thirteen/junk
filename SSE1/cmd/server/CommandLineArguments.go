package main

// Settings of the Command Line Interface including:
//   - Flag Names,
//   - Default Flag Values,
//   - Flag brief Descriptions.
//
// To see all the available Flags, start the Server with '-h' Flag.
const (
	PathToConfigurationFileArgName      = "cfg"
	PathToConfigurationFileDefaultValue = `..\..\configs\default.xml`
	PathToConfigurationFileUsageHint    = "Path to a Configuration File"
)

// Arguments of the Command Line Interface (Console, Terminal, etc.).
type CommandLineArguments struct {
	// Configuration File's full Path.
	PathToConfigurationFile string
}
