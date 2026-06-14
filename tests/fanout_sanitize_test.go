// 通天河 编译器：扇出注入字符过滤 + 方言可移植性测试。
//
// 背景：上游某次坏 commit 可能让规则 value 里混入换行 / 回车 / 逗号 / 控制字符。
// 文本类 rule-set（clash classical `.list` / surge RULE-SET / Anywhere `.arrs`）用
// 「换行分隔规则、逗号分隔字段」语法，被污染的 value 会被注入成额外的伪造规则行
// （甚至带 policy）。扇出前必须在统一模型出口拒绝这些 value，并把不可移植的
// DOMAIN-REGEX 从共享 `.list` 里剔除（surge 不支持）。
package compiler_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SafeValue 拒换行 / 回车 / 逗号 / 制表 / 控制字符 / 空格 / 空值 / 非法 UTF-8；
// 放行正常域名 / 关键字 / CIDR / 正则元字符。
func TestSafeValue(t *testing.T) {
	t.Parallel()

	// 合法：放行
	for _, v := range []string{
		"netflix.com",
		"www.example.com",
		"netflix",            // keyword
		"10.0.0.0/8",         // cidr
		"2001:db8::/32",      // cidr6
		`.*\.nflxext\.com`,   // regex 元字符
		"xn--fiqs8s.example", // punycode
	} {
		assert.True(t, ruleset.SafeValue(ruleset.MatchDomainSuffix, v), "应放行: %q", v)
	}

	// 非法：拒绝（注入向量）
	for name, v := range map[string]string{
		"换行":   "evil.com\nDOMAIN,bank.com",
		"回车":   "evil.com\rDOMAIN,bank.com",
		"逗号":   "evil.com,REJECT",
		"制表":   "evil\t.com",
		"空格":   "evil .com",
		"NUL":  "evil\x00.com",
		"控制字符": "evil\x1b.com",
		"DEL":  "evil\x7f.com",
		"空值":   "",
	} {
		assert.False(t, ruleset.SafeValue(ruleset.MatchDomainSuffix, v), "应拒绝 (%s): %q", name, v)
	}
}

// ClashClassicalList 丢弃含注入字符的 value，并计入 droppedUnsafe；合法行不受影响。
func TestClashClassicalList_DropsInjectedValue(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "good.com"},
		{Match: ruleset.MatchDomainSuffix, Value: "evil.com\nDOMAIN,bank.com"}, // 换行注入
		{Match: ruleset.MatchDomain, Value: "x.com,REJECT"},                    // 逗号注入
		{Match: ruleset.MatchDomainKeyword, Value: "ok"},
	}
	got, droppedRegex, droppedUnsafe := compiler.ClashClassicalList(rules)

	assert.Equal(t, "DOMAIN-SUFFIX,good.com\nDOMAIN-KEYWORD,ok\n", got)
	assert.Equal(t, 0, droppedRegex)
	assert.Equal(t, 2, droppedUnsafe)
	// 注入的伪造规则行（含 policy）绝不出现在产物里
	assert.NotContains(t, got, "bank.com")
	assert.NotContains(t, got, "REJECT")
	// 每行严格 `<matcher>,<value>`：行数 == 合法规则数（无注入额外行）
	assert.Equal(t, 2, strings.Count(strings.TrimSpace(got), "\n")+1)
}

// ClashClassicalList 把不可移植的 DOMAIN-REGEX 从 .list 剔除（surge 不支持），并计入 droppedRegex。
func TestClashClassicalList_DropsRegexForPortability(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "netflix.com"},
		{Match: ruleset.MatchDomainRegex, Value: `.*\.nflxext\.com`},
		{Match: ruleset.MatchIPCIDR, Value: "1.0.1.0/24"},
	}
	got, droppedRegex, droppedUnsafe := compiler.ClashClassicalList(rules)

	assert.Equal(t, 1, droppedRegex)
	assert.Equal(t, 0, droppedUnsafe)
	assert.NotContains(t, got, "DOMAIN-REGEX", ".list 不含 surge 不支持的 regex 行")
	assert.NotContains(t, got, "nflxext")
	assert.Contains(t, got, "DOMAIN-SUFFIX,netflix.com")
	assert.Contains(t, got, "IP-CIDR,1.0.1.0/24")
}

// ArrsList 丢弃含注入字符的 value，计入 dropped；合法行照常。
func TestArrsList_DropsInjectedValue(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "good.com"},
		{Match: ruleset.MatchDomainSuffix, Value: "evil.com\n2, bank.com"}, // 换行注入伪造 .arrs 行
	}
	content, kept, dropped := compiler.ArrsList("x", rules)

	assert.Equal(t, "name = x\n2, good.com\n", content)
	assert.Equal(t, 1, kept)
	assert.Equal(t, 1, dropped)
	assert.NotContains(t, content, "bank.com")
}

// SingboxSource 丢弃含注入字符的 value，但保留 DOMAIN-REGEX（sing-box 原生支持 domain_regex）。
func TestSingboxSource_DropsInjectedKeepsRegex(t *testing.T) {
	t.Parallel()
	src, err := compiler.SingboxSource([]ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "good.com"},
		{Match: ruleset.MatchDomainSuffix, Value: "evil.com\nbad"}, // 注入
		{Match: ruleset.MatchDomainRegex, Value: `.*\.nflxext\.com`},
	})
	require.NoError(t, err)

	assert.Contains(t, src, "good.com")
	assert.Contains(t, src, "domain_regex", "sing-box 原生支持 regex，保留")
	assert.Contains(t, src, "nflxext")
	// 注入的换行 value 不进产物（被污染整条丢弃）
	assert.NotContains(t, src, "evil.com")
}

// WriteCategory 端到端：含注入 + regex 的混合规则集，.list 既不含注入也不含 regex，
// .singbox.json 保留 regex，所有产物语法合法。
func TestWriteCategory_SanitizesAndFiltersDialect(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "good.com"},
		{Match: ruleset.MatchDomainSuffix, Value: "evil.com\nDOMAIN,bank.com,REJECT"},
		{Match: ruleset.MatchDomainRegex, Value: `.*\.nflxext\.com`},
	}
	require.NoError(t, compiler.WriteCategory(dir, "cat", rules, compiler.Tools{}))

	list, err := os.ReadFile(filepath.Join(dir, "cat.list"))
	require.NoError(t, err)
	assert.Equal(t, "DOMAIN-SUFFIX,good.com\n", string(list), ".list 仅留可移植 + 干净行")

	sb, err := os.ReadFile(filepath.Join(dir, "cat.singbox.json"))
	require.NoError(t, err)
	assert.Contains(t, string(sb), "nflxext", "sing-box 源保留 regex")
	assert.NotContains(t, string(sb), "bank.com", "注入未进 sing-box 源")
}
