package main

const (
	defInterceptRuleMgrType = "iptables"
)

// InterceptRuleMgr配置网络表(例如iptables或nftables)来将流量重定向到Istio代理。
type InterceptRuleMgr interface {
	Program(netns string, redirect *Redirect) error
}

type InterceptRuleMgrCtor func() InterceptRuleMgr

var (
	InterceptRuleMgrTypes = map[string]InterceptRuleMgrCtor{
		"iptables": IptablesInterceptRuleMgrCtor,
	}
)

// 已知类型的InterceptRuleMgr的构造函数工厂
func GetInterceptRuleMgrCtor(interceptType string) InterceptRuleMgrCtor {
	return InterceptRuleMgrTypes[interceptType]
}

// iptables的构造函数InterceptRuleMgr
func IptablesInterceptRuleMgrCtor() InterceptRuleMgr {
	return newIPTables()
}
