package util

// IsExportToAllNamespaces returns true if export to applies to all namespaces
// and false if it is set to namespace local.
func IsExportToAllNamespaces(exportTos []string) bool {
	exportedToAll := true
	for _, e := range exportTos {
		if e == ExportToNamespaceLocal {
			exportedToAll = false
		} else if e == ExportToAllNamespaces {
			// give preference to "*" and stop iterating
			exportedToAll = true
			break
		}
	}
	return exportedToAll
}
