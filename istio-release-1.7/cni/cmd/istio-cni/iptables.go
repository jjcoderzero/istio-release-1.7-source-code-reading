// 这是一个示例链式插件，支持多个CNI版本。它根据cniVersion解析prevResult
package main

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"

	"istio.io/pkg/log"
)

var (
	nsSetupProg = "istio-iptables"
)

type iptables struct {
}

func newIPTables() InterceptRuleMgr {
	return &iptables{}
}

// Program 定义基于重定向中提供的参数对iptables编程的方法.
func (ipt *iptables) Program(netns string, rdrct *Redirect) error {
	netnsArg := fmt.Sprintf("--net=%s", netns)
	nsSetupExecutable := fmt.Sprintf("%s/%s", nsSetupBinDir, nsSetupProg)
	nsenterArgs := []string{
		netnsArg,
		nsSetupExecutable,
		"-p", rdrct.targetPort,
		"-u", rdrct.noRedirectUID,
		"-m", rdrct.redirectMode,
		"-i", rdrct.includeIPCidrs,
		"-b", rdrct.includePorts,
		"-d", rdrct.excludeInboundPorts,
		"-o", rdrct.excludeOutboundPorts,
		"-x", rdrct.excludeIPCidrs,
		"-k", rdrct.kubevirtInterfaces,
	}
	log.Info("nsenter args",
		zap.Reflect("nsenterArgs", nsenterArgs))
	out, err := exec.Command("nsenter", nsenterArgs...).CombinedOutput()
	if err != nil {
		log.Error("nsenter failed",
			zap.String("out", string(out)),
			zap.Error(err))
		log.Infof("nsenter out: %s", out)
	} else {
		log.Infof("nsenter done: %s", out)
	}
	return err
}
