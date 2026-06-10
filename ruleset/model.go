// Package ruleset 是 通天河 的统一规则模型 —— 编译器把各上游解析成这个中立模型，
// 再扇出成各客户端的 rule-set 语法。
//
// ⚠️ 注意本模型**只承载「匹配条件」（MatchType + Value），不含「动作 / 策略」（Action）**：
// 编译产物是 *headless* rule-set（无 outbound / policy），策略由订阅模板在引用
// `RULE-SET` 时施加（`RULE-SET,netflix,🎬 Streaming`）。这正是它与「内联用户规则」
// （那种带 PROXY/DIRECT/REJECT 的）的本质区别 —— 后者不在本仓库。
package ruleset

// MatchType 是规则 matcher 类型。
type MatchType string

const (
	MatchDomain        MatchType = "DOMAIN"
	MatchDomainSuffix  MatchType = "DOMAIN-SUFFIX"
	MatchDomainKeyword MatchType = "DOMAIN-KEYWORD"
	MatchIPCIDR        MatchType = "IP-CIDR"
	MatchIPCIDR6       MatchType = "IP-CIDR6"
	MatchDomainRegex   MatchType = "DOMAIN-REGEX"
)

// Rule 是一条规则的匹配条件。**无 Policy/Action** —— 见包注释。
type Rule struct {
	Match MatchType
	Value string
}

// clashKeyword 把 MatchType 映射成 mihomo/clash classical rule-provider 关键字
// （clash `.list`（behavior: classical）与 surge `RULE-SET` text 共用此关键字）。
var clashKeyword = map[MatchType]string{
	MatchDomain:        "DOMAIN",
	MatchDomainSuffix:  "DOMAIN-SUFFIX",
	MatchDomainKeyword: "DOMAIN-KEYWORD",
	MatchIPCIDR:        "IP-CIDR",
	MatchIPCIDR6:       "IP-CIDR6",
	MatchDomainRegex:   "DOMAIN-REGEX",
}

// ClashKeyword 返回某 matcher 的 clash classical 关键字；不支持的 matcher 返回 false。
func ClashKeyword(m MatchType) (string, bool) {
	kw, ok := clashKeyword[m]
	return kw, ok
}
