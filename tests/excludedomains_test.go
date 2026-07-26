// 通天河 编译器：exclude_domains 类别源（仓内 curated 白名单，减法）测试。
// 语义：并集算完后，按域名树剔除命中白名单的 DOMAIN / DOMAIN-SUFFIX 规则。
package compiler_test

import (
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 域名树语义：白名单写根域即剔除其本身 + 全部子域；兄弟域与「后缀像但非子域」的必须留下。
func TestResolveExcludeDomains_Tree(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:           "reject",
		LocalDomains:   "excludedomains/blocklist.txt",
		ExcludeDomains: "excludedomains/allowlist.txt",
	})
	require.NoError(t, err)

	// 命中白名单的：自身 + 各级子域全部剔除
	for _, v := range []string{
		"example.com", "ads.example.com", "deep.sub.example.com",
		"media.test", "img.media.test",
	} {
		assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, v), "%s 应被剔除", v)
	}

	// 未命中的必须保留 —— 尤其 example.com.evil.net：它以 example.com 开头但
	// 并非其子域，若用朴素 strings.Contains/HasPrefix 会误杀。
	for _, v := range []string{"notexample.com", "example.com.evil.net", "other.test"} {
		assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, v), "%s 不该被剔除", v)
	}
}

// 减法必须发生在并集之后：白名单要能剔掉来自 domainlists 上游的条目，
// 而不只是剔掉 localdomains 里的（否则上游误杀根本修不掉）。
func TestResolveExcludeDomains_AppliesToUpstreamUnion(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Upstreams: map[string]compiler.UpstreamRef{
			"domainlist": {Repo: "privacy-protection-tools/anti-AD", Commit: "a8f65d8043940b20ba7ee75d402c60b885e547e6", Files: []string{"plain.txt"}},
		},
	}
	r := compiler.NewResolver(m, "fixtures", "fixtures")

	// 先不加白名单：确认上游 fixture 里确有 example.com（否则下面的断言是假绿）。
	base, err := r.ResolveCategory(compiler.Category{Name: "reject", DomainLists: []string{"domainlist"}})
	require.NoError(t, err)
	require.True(t, hasRule(base, ruleset.MatchDomainSuffix, "example.com"), "前置条件：上游本来收录了 example.com")

	// 加白名单后，来自**上游**的这条必须被剔除 —— 这正是修上游误杀的关键路径。
	withAllow, err := r.ResolveCategory(compiler.Category{
		Name:           "reject",
		DomainLists:    []string{"domainlist"},
		ExcludeDomains: "excludedomains/allowlist.txt",
	})
	require.NoError(t, err)
	assert.False(t, hasRule(withAllow, ruleset.MatchDomainSuffix, "example.com"), "上游条目应被白名单剔除")
	assert.Less(t, len(withAllow), len(base), "确实少了条目")
	// 未命中白名单的上游条目不受牵连
	assert.True(t, hasRule(withAllow, ruleset.MatchDomainSuffix, "adguard.com"), "无关条目保留")
}

// 只作用于域名规则：IP-CIDR 与 keyword 不受白名单影响。
func TestResolveExcludeDomains_LeavesNonDomainRules(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:           "mixed",
		LocalDomains:   "localdomains/cn-extra.txt", // 含 DOMAIN-KEYWORD 规则行
		ExcludeDomains: "excludedomains/allowlist.txt",
	})
	require.NoError(t, err)
	assert.True(t, hasRule(rules, ruleset.MatchDomainKeyword, "examplekw"), "keyword 规则不受影响")
}

// 空 / 未配置 exclude_domains 时结果与不配一致（不引入行为变化）。
func TestResolveExcludeDomains_NoopWhenUnset(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	with, err := r.ResolveCategory(compiler.Category{Name: "cn", LocalDomains: "excludedomains/blocklist.txt"})
	require.NoError(t, err)
	assert.True(t, hasRule(with, ruleset.MatchDomainSuffix, "example.com"), "未配白名单则原样保留")
}

// 越权防护：exclude_domains 含 .. / 绝对路径应在解析期与校验期都被挡下。
func TestResolveExcludeDomains_Traversal(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	_, err := r.ResolveCategory(compiler.Category{Name: "evil", ExcludeDomains: "../../../etc/passwd"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "越出仓库根")
}

func TestValidateExcludeDomainsTraversal(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Upstreams: map[string]compiler.UpstreamRef{
			"dnsmasq-china-list": {Repo: "felixonmars/dnsmasq-china-list", Commit: "d87a5bc1d76b87a6729e9e5355c35be0fc2ff6cd", Files: []string{"accelerated-domains.china.conf"}},
		},
		Categories: []compiler.Category{{Name: "evil", ExcludeDomains: "/etc/passwd"}},
	}
	err := m.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "越出仓库根")
}
