// 通天河 编译清单：同一 repo 被多个上游 key 引用时，commit 必须一致。
//
// 动因：apple-cn 类别与 cn 同源 felixonmars/dnsmasq-china-list，但各取该仓不同文件
// （apple.china.conf vs accelerated-domains.china.conf），故声明为两个 key。CI 按 key
// 独立 clone，两 key 的 SHA 一旦漂移，同仓不同文件即来自不同版本 —— 静默不一致，
// 且 bump 上游时极易只改一处。
package compiler_test

import (
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/stretchr/testify/require"
)

const (
	shaA = "d87a5bc1d76b87a6729e9e5355c35be0fc2ff6cd"
	shaB = "0000000000000000000000000000000000000000"
)

func TestValidate_SameRepoSameSHA_OK(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Repo: "river2heaven/tongtian",
		Upstreams: map[string]compiler.UpstreamRef{
			"dnsmasq-china-list": {
				Repo: "felixonmars/dnsmasq-china-list", Commit: shaA,
				Files: []string{"accelerated-domains.china.conf"},
			},
			"dnsmasq-china-list-apple": {
				Repo: "felixonmars/dnsmasq-china-list", Commit: shaA,
				Files: []string{"apple.china.conf"},
			},
		},
		Categories: []compiler.Category{
			{Name: "cn", ChinaList: "dnsmasq-china-list"},
			{Name: "apple-cn", ChinaList: "dnsmasq-china-list-apple"},
		},
	}
	require.NoError(t, m.Validate(), "同 repo 同 SHA 应通过")
}

func TestValidate_SameRepoDifferentSHA_Rejected(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Repo: "river2heaven/tongtian",
		Upstreams: map[string]compiler.UpstreamRef{
			"dnsmasq-china-list": {
				Repo: "felixonmars/dnsmasq-china-list", Commit: shaA,
				Files: []string{"accelerated-domains.china.conf"},
			},
			"dnsmasq-china-list-apple": {
				Repo: "felixonmars/dnsmasq-china-list", Commit: shaB, // 漂移
				Files: []string{"apple.china.conf"},
			},
		},
		Categories: []compiler.Category{
			{Name: "cn", ChinaList: "dnsmasq-china-list"},
			{Name: "apple-cn", ChinaList: "dnsmasq-china-list-apple"},
		},
	}
	err := m.Validate()
	require.Error(t, err, "同 repo 不同 SHA 必须被拒")
	require.Contains(t, err.Error(), "同仓多 key 必须同 SHA")
	// 报错须点名两个 key + repo，便于定位该改哪一处
	require.Contains(t, err.Error(), "dnsmasq-china-list")
	require.Contains(t, err.Error(), "dnsmasq-china-list-apple")
	require.Contains(t, err.Error(), "felixonmars/dnsmasq-china-list")
}

// 不同 repo 各自不同 SHA 是常态，不该误伤。
func TestValidate_DifferentRepos_NotAffected(t *testing.T) {
	t.Parallel()
	m := &compiler.Manifest{
		Repo: "river2heaven/tongtian",
		Upstreams: map[string]compiler.UpstreamRef{
			"dnsmasq-china-list": {
				Repo: "felixonmars/dnsmasq-china-list", Commit: shaA,
				Files: []string{"accelerated-domains.china.conf"},
			},
			"gfwlist": {
				Repo: "gfwlist/gfwlist", Commit: shaB,
				Files: []string{"gfwlist.txt"},
			},
		},
		Categories: []compiler.Category{{Name: "cn", ChinaList: "dnsmasq-china-list"}},
	}
	require.NoError(t, m.Validate(), "不同 repo 不同 SHA 属常态，不该误伤")
}
