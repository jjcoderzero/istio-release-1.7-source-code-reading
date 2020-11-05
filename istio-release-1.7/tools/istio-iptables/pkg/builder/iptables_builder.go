package builder

// IptablesProducer is an interface for adding iptables rules
type IptablesProducer interface {
	// AppendRuleV4 appends an IPv4 rule into the given iptables chain
	AppendRuleV4(chain string, table string, params ...string) IptablesProducer
	// AppendRuleV6 appends an IPv6 rule into the given iptables chain
	AppendRuleV6(chain string, table string, params ...string) IptablesProducer
	// InsertRuleV4 inserts IPv4 rule at a particular position in the chain
	InsertRuleV4(chain string, table string, position int, params ...string) IptablesProducer
	// InsertRuleV6 inserts IPv6 rule at a particular position in the chain
	InsertRuleV6(chain string, table string, position int, params ...string) IptablesProducer
}
