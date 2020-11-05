package cmd

import (
	"istio.io/istio/tools/istio-iptables/pkg/constants"
	dep "istio.io/istio/tools/istio-iptables/pkg/dependencies"
)

func flushAndDeleteChains(ext dep.Dependencies, cmd string, table string, chains []string) {
	for _, chain := range chains {
		ext.RunQuietlyAndIgnore(cmd, "-t", table, "-F", chain)
		ext.RunQuietlyAndIgnore(cmd, "-t", table, "-X", chain)
	}
}

func removeOldChains(ext dep.Dependencies, cmd string) {
	for _, table := range []string{constants.NAT, constants.MANGLE} {
		// Remove the old chains
		ext.RunQuietlyAndIgnore(cmd, "-t", table, "-D", constants.PREROUTING, "-p", constants.TCP, "-j", constants.ISTIOINBOUND)
	}
	ext.RunQuietlyAndIgnore(cmd, "-t", constants.NAT, "-D", constants.OUTPUT, "-p", constants.TCP, "-j", constants.ISTIOOUTPUT)

	// Flush and delete the istio chains from NAT table.
	chains := []string{constants.ISTIOOUTPUT, constants.ISTIOINBOUND}
	flushAndDeleteChains(ext, cmd, constants.NAT, chains)
	// Flush and delete the istio chains from MANGLE table.
	chains = []string{constants.ISTIOINBOUND, constants.ISTIODIVERT, constants.ISTIOTPROXY}
	flushAndDeleteChains(ext, cmd, constants.MANGLE, chains)

	// Must be last, the others refer to it
	chains = []string{constants.ISTIOREDIRECT, constants.ISTIOINREDIRECT}
	flushAndDeleteChains(ext, cmd, constants.NAT, chains)
}

func cleanup(dryRun bool) {
	var ext dep.Dependencies
	if dryRun {
		ext = &dep.StdoutStubDependencies{}
	} else {
		ext = &dep.RealDependencies{}
	}

	defer func() {
		for _, cmd := range []string{constants.IPTABLESSAVE, constants.IP6TABLESSAVE} {
			// iptables-save is best efforts
			_ = ext.Run(cmd)
		}
	}()

	for _, cmd := range []string{constants.IPTABLES, constants.IP6TABLES} {
		removeOldChains(ext, cmd)
	}
}
