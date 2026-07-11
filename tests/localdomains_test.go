// 通天河 编译器：localdomains 类别源（仓内 curated 域名种子，如 cn-extra）解析 + 并集 + 越权防护测试。
package compiler_test

import (
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// localdomains：复用 ParseDomainList —— 裸域 → DOMAIN-SUFFIX，容忍 clash 规则行与 #/! 注释。
func TestResolveLocalDomains(t *testing.T) {
	t.Parallel()
	// repoRoot = "fixtures"，故 LocalDomains 相对它解析。
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:         "cn",
		LocalDomains: "localdomains/cn-extra.txt",
	})
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "bambulab.com"), "裸域 → DOMAIN-SUFFIX")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "bambulab.cn"), "裸域 → DOMAIN-SUFFIX")
	assert.True(t, hasRule(rules, ruleset.MatchDomain, "api.example.cn"), "DOMAIN 规则行")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "bblmw.cn"), "DOMAIN-SUFFIX 规则行")
	assert.True(t, hasRule(rules, ruleset.MatchDomainKeyword, "examplekw"), "DOMAIN-KEYWORD 规则行")
	// 注释与畸形行不产规则
	for _, rr := range rules {
		assert.NotContains(t, rr.Value, "#")
		assert.NotContains(t, rr.Value, " ")
	}
}

// localdomains 与 chinalist 同类别并集 + 去重（真实用法：cn = dnsmasq-china-list ∪ cn-extra）。
func TestResolveLocalDomains_UnionWithChinaList(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Upstreams: map[string]compiler.UpstreamRef{
			"chinalist": {Repo: "felixonmars/dnsmasq-china-list", Commit: "d87a5bc1d76b87a6729e9e5355c35be0fc2ff6cd", Files: []string{"china.conf"}},
		},
	}
	r := compiler.NewResolver(m, "fixtures", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:         "cn",
		ChinaList:    "chinalist",
		LocalDomains: "localdomains/cn-extra.txt",
	})
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "baidu.com"), "chinalist 侧")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "bambulab.cn"), "localdomains 侧")
	// 去重：两侧合并后同 (Match,Value) 只留一条
	count := 0
	for _, rr := range rules {
		if rr.Match == ruleset.MatchDomainSuffix && rr.Value == "bambulab.cn" {
			count++
		}
	}
	assert.Equal(t, 1, count, "并集去重")
}

// 越权防护：localdomains 含 .. 应报错（防坏 manifest 读仓外任意文件）。
func TestResolveLocalDomains_Traversal(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	_, err := r.ResolveCategory(compiler.Category{Name: "evil", LocalDomains: "../../../etc/passwd"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "越出仓库根")
}

// 未配 repoRoot 时 localdomains 类别应报错（而非静默读到 cwd 相对路径）。
func TestResolveLocalDomains_NoRepoRoot(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "")
	_, err := r.ResolveCategory(compiler.Category{Name: "cn", LocalDomains: "data/cn-extra.txt"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "repoRoot")
}

// Manifest.Validate 也应在校验期挡住越权 localdomains（fail-fast）。
func TestValidateLocalDomainsTraversal(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Upstreams: map[string]compiler.UpstreamRef{
			"dnsmasq-china-list": {Repo: "felixonmars/dnsmasq-china-list", Commit: "d87a5bc1d76b87a6729e9e5355c35be0fc2ff6cd", Files: []string{"accelerated-domains.china.conf"}},
		},
		Categories: []compiler.Category{{Name: "evil", LocalDomains: "/etc/passwd"}},
	}
	err := m.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "越出仓库根")
}
