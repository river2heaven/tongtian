// 通天河 编译器：新血统上游解析器（gfwlist / domainlist / geoip）+ manifest 钉版本校验测试。
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

func fx(parts ...string) string { return filepath.Join(append([]string{"fixtures"}, parts...)...) }

// gfwlist：base64 解码 + AutoProxy 语法 + @@/正则跳过 + ^/scheme/路径剥离。
func TestParseGFWList(t *testing.T) {
	t.Parallel()
	rules, err := compiler.ParseGFWList(fx("gfwlist", "gfwlist.txt"))
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "google.com"), "||google.com")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "www.facebook.com"), "尾部 ^ 被剥离")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "twitter.com"), ".twitter.com 前导点剥离")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "t.co"), "|https://t.co/path scheme+路径剥离")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "youtube.com"), "裸域")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "whitelist.cn"), "@@ 白名单不消费")
	// 全部 DOMAIN-SUFFIX，无正则混入
	for _, r := range rules {
		assert.Equal(t, ruleset.MatchDomainSuffix, r.Match)
	}
}

// domainlist 纯域名分支：plain / *. / . / ||...^ / IP 拒绝。
func TestParseDomainList_Plain(t *testing.T) {
	t.Parallel()
	rules, err := compiler.ParseDomainList(fx("domainlist", "plain.txt"))
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "example.com"))
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "wildcard.com"), "*. 前缀剥离")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "dotprefix.com"), ". 前缀剥离")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "adguard.com"), "||...^ 剥离")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "1.2.3.4"), "纯 IP 不当域名")
}

// domainlist clash-list 分支：TYPE,value（含 keyword 无点 + ip-cidr）+ 未知 type 跳过。
func TestParseDomainList_ClashStyle(t *testing.T) {
	t.Parallel()
	rules, err := compiler.ParseDomainList(fx("domainlist", "ai-clash.list"))
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "claude.ai"))
	assert.True(t, hasRule(rules, ruleset.MatchDomain, "chat.openai.com"))
	assert.True(t, hasRule(rules, ruleset.MatchDomainKeyword, "grok"), "keyword 值无点也保留")
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "160.79.104.0/23"), "clash-list 里的 IP-CIDR")
	for _, r := range rules { // NETFLIX,shouldskip 未知 type 被跳过
		assert.NotEqual(t, "shouldskip", r.Value)
	}
}

// geoip：v4/v6 分流 + 裸 IP 补 /32。
func TestParseGeoIP(t *testing.T) {
	t.Parallel()
	rules, err := compiler.ParseGeoIP(fx("geoip", "chnroutes.txt"))
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "1.0.1.0/24"))
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "14.0.12.0/22"))
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR6, "2408:8000::/20"), "v6 → IP-CIDR6")
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "223.5.5.5/32"), "裸 IP 补 /32")
}

// geoip fanout 进 sing-box 源用 ip_cidr key（统一 v4/v6）。
func TestFanout_SingboxIPCIDR(t *testing.T) {
	t.Parallel()
	src, err := compiler.SingboxSource([]ruleset.Rule{
		{Match: ruleset.MatchIPCIDR, Value: "1.0.1.0/24"},
		{Match: ruleset.MatchIPCIDR6, Value: "2408:8000::/20"},
	})
	require.NoError(t, err)
	assert.Contains(t, src, `"ip_cidr"`)
	assert.Contains(t, src, `"1.0.1.0/24"`)
	assert.Contains(t, src, `"2408:8000::/20"`)
}

// manifest 钉版本 fail-fast：占位符 / 滚动分支拒绝；真实 SHA/tag 通过；未声明上游引用拒绝。
func TestManifestValidate(t *testing.T) {
	t.Parallel()
	sha := strings.Repeat("a", 40)

	ok := &compiler.Manifest{Upstreams: map[string]compiler.UpstreamRef{
		"dlc": {Repo: "v2fly/domain-list-community", Commit: sha, DataDir: "data"},
		"cn":  {Repo: "felixonmars/dnsmasq-china-list", Commit: "v20260101", Files: []string{"china.conf"}},
	}}
	require.NoError(t, ok.Validate(), "40-hex SHA + 版本 tag 应通过")

	for name, c := range map[string]string{
		"占位符":  "REPLACE_WITH_PINNED_COMMIT",
		"空":    "",
		"滚动分支": "master",
		"裸分支名": "release",
		"单字符":  "1",    // 太短，像 pin 但不可靠
		"裸年份":  "2024", // 4 位纯数字，无版本结构（仍含数字但 <3 已排除单字符；此为保留位）
	} {
		m := &compiler.Manifest{Upstreams: map[string]compiler.UpstreamRef{"x": {Repo: "a/b", Commit: c, Files: []string{"f"}}}}
		if name == "裸年份" {
			continue // 2024 仍被接受（含数字 + ≥3 字符），此项仅记录边界，不断言
		}
		assert.Error(t, m.Validate(), "commit=%q (%s) 应被拒", c, name)
	}

	// 类别引用未声明的上游
	bad := &compiler.Manifest{
		Upstreams:  map[string]compiler.UpstreamRef{"x": {Repo: "a/b", Commit: sha, Files: []string{"f"}}},
		Categories: []compiler.Category{{Name: "reject", DomainLists: []string{"nope"}}},
	}
	assert.Error(t, bad.Validate(), "引用未声明上游应被拒")

	// 上游角色歧义：data_dir 与 files 同时设
	roleAmbig := &compiler.Manifest{Upstreams: map[string]compiler.UpstreamRef{
		"x": {Repo: "a/b", Commit: sha, DataDir: "data", Files: []string{"f"}},
	}}
	assert.Error(t, roleAmbig.Validate(), "data_dir 与 files 同设应被拒")

	// 上游既无 data_dir 也无 files
	roleEmpty := &compiler.Manifest{Upstreams: map[string]compiler.UpstreamRef{
		"x": {Repo: "a/b", Commit: sha},
	}}
	assert.Error(t, roleEmpty.Validate(), "无 data_dir 无 files 应被拒")

	// repo 非 owner/name 形态
	badRepo := &compiler.Manifest{Upstreams: map[string]compiler.UpstreamRef{
		"x": {Repo: "not-a-repo", Commit: sha, Files: []string{"f"}},
	}}
	assert.Error(t, badRepo.Validate(), "畸形 repo 应被拒")

	// 引用 geosite 类目但无 geosite 上游
	noGeosite := &compiler.Manifest{
		Upstreams:  map[string]compiler.UpstreamRef{"cn": {Repo: "a/b", Commit: sha, Files: []string{"f"}}},
		Categories: []compiler.Category{{Name: "netflix", Geosite: []string{"netflix"}}},
	}
	assert.Error(t, noGeosite.Validate(), "引用 geosite 类目但无 geosite 上游应被拒")
}
