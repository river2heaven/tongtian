// 通天河 编译器：localfile 类别源（仓内 curated IP-CIDR，如 china-company-ip）解析 + 越权防护测试。
package compiler_test

import (
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// localfile：复用 ParseGeoIP —— 跳 # 注释、剥尾注、裸 IP 补位、v4/v6 分流。
func TestResolveLocalFile(t *testing.T) {
	t.Parallel()
	// repoRoot = "fixtures"，故 LocalFile 相对它解析。
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:      "china-company-ip",
		LocalFile: "localfile/china-company-ip.cidr",
	})
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "43.128.0.0/10"), "腾讯海外段")
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "119.29.0.0/16"), "腾讯国内段")
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR6, "2402:4e00::/32"), "v6 分流到 IPCIDR6")
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "1.2.3.0/24"), "尾注被剥离")
	assert.True(t, hasRule(rules, ruleset.MatchIPCIDR, "203.0.113.7/32"), "裸 IP 补 /32")
	// 注释行不产规则
	for _, rr := range rules {
		assert.NotContains(t, rr.Value, "#")
	}
}

// 越权防护：localfile 含 .. 应报错（防坏 manifest 读仓外任意文件）。
func TestResolveLocalFile_Traversal(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "fixtures")
	_, err := r.ResolveCategory(compiler.Category{Name: "evil", LocalFile: "../../../etc/passwd"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "越出仓库根")
}

// 未配 repoRoot 时 localfile 类别应报错（而非静默读到 cwd 相对路径）。
func TestResolveLocalFile_NoRepoRoot(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(&compiler.Manifest{}, "", "")
	_, err := r.ResolveCategory(compiler.Category{Name: "cci", LocalFile: "data/x.cidr"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "repoRoot")
}

// Manifest.Validate 也应在校验期挡住越权 localfile（fail-fast）。
func TestValidateLocalFileTraversal(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Upstreams: map[string]compiler.UpstreamRef{
			"dnsmasq-china-list": {Repo: "felixonmars/dnsmasq-china-list", Commit: "d87a5bc1d76b87a6729e9e5355c35be0fc2ff6cd", Files: []string{"accelerated-domains.china.conf"}},
		},
		Categories: []compiler.Category{{Name: "evil", LocalFile: "/etc/passwd"}},
	}
	err := m.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "越出仓库根")
}
