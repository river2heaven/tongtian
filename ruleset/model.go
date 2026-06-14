// Package ruleset 是 通天河 的统一规则模型 —— 编译器把各上游解析成这个中立模型，
// 再扇出成各客户端的 rule-set 语法。
//
// ⚠️ 注意本模型**只承载「匹配条件」（MatchType + Value），不含「动作 / 策略」（Action）**：
// 编译产物是 *headless* rule-set（无 outbound / policy），策略由订阅模板在引用
// `RULE-SET` 时施加（`RULE-SET,netflix,🎬 Streaming`）。这正是它与「内联用户规则」
// （那种带 PROXY/DIRECT/REJECT 的）的本质区别 —— 后者不在本仓库。
package ruleset

import "unicode/utf8"

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

// SafeValue 校验一条规则的 Value 是否能安全写入文本类 rule-set 产物。
//
// 各文本方言（clash classical `.list` / surge RULE-SET / Anywhere `.arrs`）都用
// 「换行分隔规则、逗号分隔字段」的语法。若上游某次坏 commit 让 Value 里混入换行 / 回车
// / 逗号 / 其它控制字符，渲染时就会被注入成额外的伪造规则行（甚至带 policy 的行）。
// 因此扇出前必须在统一模型出口把这些字符**拒绝**（而非转义——本编译器消费的全是
// 结构化域名 / 关键字 / CIDR，合法值里永不含这些字符，出现即视为被污染，直接丢弃）。
//
// 规则：
//   - 拒绝非法 UTF-8（防畸形字节序列）。
//   - 拒绝任何 ASCII 控制字符（含 \n \r \t）与 DEL。
//   - 拒绝空格 / 逗号（裸 value 字段 + classical/arrs 字段分隔符）。
//   - 空值拒绝。
//
// DOMAIN-REGEX 允许正则元字符（`. * \ [ ] ( ) | ^ $ + ?` 等），但同样禁换行 / 逗号 /
// 控制字符——参数 m 暂未按类型放宽，保留以备后续细化。
func SafeValue(m MatchType, v string) bool {
	_ = m
	if v == "" || !utf8.ValidString(v) {
		return false
	}
	for _, r := range v {
		if r < 0x20 || r == 0x7f { // ASCII 控制字符（含 \n \r \t）+ DEL
			return false
		}
		if r == ' ' || r == ',' { // 字段分隔符：空格 / 逗号
			return false
		}
	}
	return true
}
