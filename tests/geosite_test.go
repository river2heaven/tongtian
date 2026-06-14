// 通天河 编译器：geosite DSL 解析 + chinalist + fanout 单元测试。
package compiler_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFixtureGeosite(t *testing.T) *compiler.Geosite {
	t.Helper()
	g, err := compiler.LoadGeosite(filepath.Join("fixtures", "geosite"))
	require.NoError(t, err)
	return g
}

func hasRule(rules []ruleset.Rule, m ruleset.MatchType, v string) bool {
	for _, r := range rules {
		if r.Match == m && r.Value == v {
			return true
		}
	}
	return false
}

func countRule(rules []ruleset.Rule, v string) int {
	n := 0
	for _, r := range rules {
		if r.Value == v {
			n++
		}
	}
	return n
}

// 各前缀翻译 + include 递归
func TestGeosite_Prefixes(t *testing.T) {
	t.Parallel()
	g := loadFixtureGeosite(t)
	rules, err := g.Resolve("netflix", nil)
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "netflix.com"), "裸域→DOMAIN-SUFFIX")
	assert.True(t, hasRule(rules, ruleset.MatchDomain, "www.netflix.com"), "full→DOMAIN")
	assert.True(t, hasRule(rules, ruleset.MatchDomainKeyword, "netflix"), "keyword→DOMAIN-KEYWORD")
	assert.True(t, hasRule(rules, ruleset.MatchDomainRegex, `.*\.nflxext\.com`), "regexp→DOMAIN-REGEX")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "fast.com"), "include 递归展开")
}

// @attr 过滤（manifest 级 exclude）
func TestGeosite_ExcludeAttrs(t *testing.T) {
	t.Parallel()
	g := loadFixtureGeosite(t)

	with, _ := g.Resolve("netflix", nil)
	assert.True(t, hasRule(with, ruleset.MatchDomainSuffix, "ads.netflix.com"), "无 exclude 时 @ads 条目在")

	without, _ := g.Resolve("netflix", []string{"ads"})
	assert.False(t, hasRule(without, ruleset.MatchDomainSuffix, "ads.netflix.com"), "exclude ads 后 @ads 条目被排除")
	assert.True(t, hasRule(without, ruleset.MatchDomainSuffix, "netflix.com"), "其余条目仍在")
}

// include 携带 @attr（穿透嵌套）：streaming = include:netflix @cdn → 仅 @cdn 条目
func TestGeosite_IncludeWithAttr(t *testing.T) {
	t.Parallel()
	g := loadFixtureGeosite(t)
	rules, err := g.Resolve("streaming", nil)
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "nflxvideo.net"), "仅 @cdn 的 nflxvideo.net")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "netflix.com"), "非 @cdn 的 netflix.com 不在")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "fast.com"), "嵌套 include 也受 @cdn 约束")
}

// 去重（netflix 与 netflix-extra 都有 netflix.com）
func TestGeosite_Dedup(t *testing.T) {
	t.Parallel()
	g := loadFixtureGeosite(t)
	rules, _ := g.Resolve("netflix", nil)
	assert.Equal(t, 1, countRule(rules, "netflix.com"), "重复域名应去重")
}

// 环保护：cycle-a ↔ cycle-b 互相 include，不死循环
func TestGeosite_CycleGuard(t *testing.T) {
	t.Parallel()
	g := loadFixtureGeosite(t)
	rules, err := g.Resolve("cycle-a", nil)
	require.NoError(t, err)
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "a.com"))
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "b.com"))
}

// dnsmasq-china-list 解析
func TestChinaList(t *testing.T) {
	t.Parallel()
	rules, err := compiler.ParseChinaList(filepath.Join("fixtures", "chinalist", "china.conf"))
	require.NoError(t, err)
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "baidu.com"))
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "taobao.com"))
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "qq.com"), "ipset= 行也解析")
	assert.Equal(t, 1, countRule(rules, "baidu.com"), "重复行去重")
}

// dnsmasq-china-list 收紧：只接受 server=/ipset=，其它指令（address=/local=/cache-size=）丢弃
func TestChinaList_OnlyServerAndIpset(t *testing.T) {
	t.Parallel()
	rules, err := compiler.ParseChinaList(filepath.Join("fixtures", "chinalist", "mixed.conf"))
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "good.com"), "server= 收")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "good2.com"), "ipset= 收")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "good3.com"), "server= 收")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "evil.com"), "address= 拒")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "evil2.com"), "local= 拒")
	assert.Len(t, rules, 3, "仅 3 条合法 server=/ipset= 域名")
}

// fanout .list 文本（golden）
func TestFanout_ClashList(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "netflix.com"},
		{Match: ruleset.MatchDomain, Value: "www.netflix.com"},
		{Match: ruleset.MatchDomainKeyword, Value: "netflix"},
	}
	got, droppedRegex, droppedUnsafe := compiler.ClashClassicalList(rules)
	want := "DOMAIN-SUFFIX,netflix.com\nDOMAIN,www.netflix.com\nDOMAIN-KEYWORD,netflix\n"
	assert.Equal(t, want, got)
	assert.Equal(t, 0, droppedRegex)
	assert.Equal(t, 0, droppedUnsafe)
}

// fanout sing-box 源 JSON（按 type 聚合，version 2）
func TestFanout_SingboxSource(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "a.com"},
		{Match: ruleset.MatchDomainSuffix, Value: "b.com"},
		{Match: ruleset.MatchDomain, Value: "c.com"},
	}
	src, err := compiler.SingboxSource(rules)
	require.NoError(t, err)
	assert.Contains(t, src, `"version": 2`)
	assert.Contains(t, src, `"domain_suffix"`)
	assert.Contains(t, src, `"a.com"`)
	assert.Contains(t, src, `"domain"`)
	// 合法 JSON
	assert.True(t, strings.HasPrefix(strings.TrimSpace(src), "{"))
}
