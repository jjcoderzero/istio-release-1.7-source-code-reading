package cmd

import "strings"

func FormatIptablesCommands(commands [][]string) []string {
	output := make([]string, 0)
	for _, cmd := range commands {
		output = append(output, strings.Join(cmd, " "))
	}
	return output
}
