// 通天河 编译器：geosite.dat 扇出测试（issue #3）。
//
// geosite.dat 是 xray/v2ray 原生格式：单个 protobuf 文件，内含全部类目 tag，
// 消费端用 `geosite:<tag>` 引用。schema：
//
//	GeoSiteList { repeated GeoSite entry = 1 }
//	GeoSite     { string country_code = 1; repeated Domain domain = 2 }
//	Domain      { Type type = 1; string value = 2 }  // Type: Plain=0 Regex=1 Domain=2 Full=3
//
// 映射（统一模型 → Domain.Type）：
//
//	DOMAIN(exact)  → Full(3)
//	DOMAIN-SUFFIX  → Domain(2)   // geosite Domain type 即 root-domain/suffix 语义
//	DOMAIN-KEYWORD → Plain(0)
//	DOMAIN-REGEX   → Regex(1)
//	IP-CIDR / IP-CIDR6 → 跳过（IP 属 geoip.dat）
//
// 测试用自写 protowire 反解验证：类目 tag（大写）下域名集与统一模型一致 + 各 Type 映射正确。
package compiler_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

// decodedDomain 是反解出的一条 Domain（type + value）。
type decodedDomain struct {
	Type  int32
	Value string
}

// decodedSite 是反解出的一个类目（country_code + 域名集）。
type decodedSite struct {
	CountryCode string
	Domains     []decodedDomain
}

// parseDomain 反解 Domain message（field 1 = type varint, field 2 = value string）。
func parseDomain(t *testing.T, b []byte) decodedDomain {
	t.Helper()
	var d decodedDomain
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		require.GreaterOrEqual(t, n, 0, "ConsumeTag")
		b = b[n:]
		switch {
		case num == 1 && typ == protowire.VarintType:
			v, n := protowire.ConsumeVarint(b)
			require.GreaterOrEqual(t, n, 0)
			d.Type = int32(v)
			b = b[n:]
		case num == 2 && typ == protowire.BytesType:
			v, n := protowire.ConsumeBytes(b)
			require.GreaterOrEqual(t, n, 0)
			d.Value = string(v)
			b = b[n:]
		default:
			n := protowire.ConsumeFieldValue(num, typ, b)
			require.GreaterOrEqual(t, n, 0)
			b = b[n:]
		}
	}
	return d
}

// parseSite 反解 GeoSite message（field 1 = country_code string, field 2 = repeated Domain）。
func parseSite(t *testing.T, b []byte) decodedSite {
	t.Helper()
	var s decodedSite
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		require.GreaterOrEqual(t, n, 0)
		b = b[n:]
		switch {
		case num == 1 && typ == protowire.BytesType:
			v, n := protowire.ConsumeBytes(b)
			require.GreaterOrEqual(t, n, 0)
			s.CountryCode = string(v)
			b = b[n:]
		case num == 2 && typ == protowire.BytesType:
			v, n := protowire.ConsumeBytes(b)
			require.GreaterOrEqual(t, n, 0)
			s.Domains = append(s.Domains, parseDomain(t, v))
			b = b[n:]
		default:
			n := protowire.ConsumeFieldValue(num, typ, b)
			require.GreaterOrEqual(t, n, 0)
			b = b[n:]
		}
	}
	return s
}

// parseGeositeDat 反解 GeoSiteList（field 1 = repeated GeoSite），按 country_code 索引。
func parseGeositeDat(t *testing.T, b []byte) map[string]decodedSite {
	t.Helper()
	out := map[string]decodedSite{}
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		require.GreaterOrEqual(t, n, 0)
		b = b[n:]
		if num == 1 && typ == protowire.BytesType {
			v, n := protowire.ConsumeBytes(b)
			require.GreaterOrEqual(t, n, 0)
			s := parseSite(t, v)
			out[s.CountryCode] = s
			b = b[n:]
			continue
		}
		n2 := protowire.ConsumeFieldValue(num, typ, b)
		require.GreaterOrEqual(t, n2, 0)
		b = b[n2:]
	}
	return out
}

// 各 MatchType → Domain.Type 映射正确
func TestGeositeDat_TypeMapping(t *testing.T) {
	t.Parallel()
	cats := []compiler.GeositeCategory{{
		Name: "netflix",
		Rules: []ruleset.Rule{
			{Match: ruleset.MatchDomain, Value: "www.netflix.com"},
			{Match: ruleset.MatchDomainSuffix, Value: "netflix.com"},
			{Match: ruleset.MatchDomainKeyword, Value: "netflix"},
			{Match: ruleset.MatchDomainRegex, Value: `.*\.nflxext\.com`},
		},
	}}
	b, err := compiler.GeositeDat(cats)
	require.NoError(t, err)

	sites := parseGeositeDat(t, b)
	site, ok := sites["NETFLIX"]
	require.True(t, ok, "country_code 为类目名大写")

	got := map[string]int32{}
	for _, d := range site.Domains {
		got[d.Value] = d.Type
	}
	assert.Equal(t, int32(3), got["www.netflix.com"], "DOMAIN → Full(3)")
	assert.Equal(t, int32(2), got["netflix.com"], "DOMAIN-SUFFIX → Domain(2)")
	assert.Equal(t, int32(0), got["netflix"], "DOMAIN-KEYWORD → Plain(0)")
	assert.Equal(t, int32(1), got[`.*\.nflxext\.com`], "DOMAIN-REGEX → Regex(1)")
	assert.Len(t, site.Domains, 4, "四条 domain 全进")
}

// IP-CIDR / IP-CIDR6 跳过（不进 geosite.dat）
func TestGeositeDat_SkipIP(t *testing.T) {
	t.Parallel()
	cats := []compiler.GeositeCategory{{
		Name: "ai",
		Rules: []ruleset.Rule{
			{Match: ruleset.MatchDomainSuffix, Value: "claude.com"},
			{Match: ruleset.MatchIPCIDR, Value: "10.0.0.0/8"},
			{Match: ruleset.MatchIPCIDR6, Value: "2001:db8::/32"},
		},
	}}
	b, err := compiler.GeositeDat(cats)
	require.NoError(t, err)

	site := parseGeositeDat(t, b)["AI"]
	require.Len(t, site.Domains, 1, "仅 domain 进，两条 IP 跳过")
	assert.Equal(t, "claude.com", site.Domains[0].Value)
	for _, d := range site.Domains {
		assert.NotContains(t, d.Value, "/", "无 CIDR 残留")
	}
}

// 多类目合一：单文件含全部 tag，各 tag 独立索引
func TestGeositeDat_MultiCategory(t *testing.T) {
	t.Parallel()
	cats := []compiler.GeositeCategory{
		{Name: "ai", Rules: []ruleset.Rule{{Match: ruleset.MatchDomainSuffix, Value: "claude.com"}}},
		{Name: "netflix", Rules: []ruleset.Rule{{Match: ruleset.MatchDomainSuffix, Value: "netflix.com"}}},
		{Name: "gfw", Rules: []ruleset.Rule{{Match: ruleset.MatchDomainSuffix, Value: "google.com"}}},
	}
	b, err := compiler.GeositeDat(cats)
	require.NoError(t, err)

	sites := parseGeositeDat(t, b)
	require.Len(t, sites, 3)
	assert.Equal(t, "claude.com", sites["AI"].Domains[0].Value)
	assert.Equal(t, "netflix.com", sites["NETFLIX"].Domains[0].Value)
	assert.Equal(t, "google.com", sites["GFW"].Domains[0].Value)
}

// 域名集与统一模型一致（含顺序无关比对）
func TestGeositeDat_DomainSetMatchesModel(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "a.com"},
		{Match: ruleset.MatchDomainSuffix, Value: "b.com"},
		{Match: ruleset.MatchDomain, Value: "full.c.com"},
		{Match: ruleset.MatchIPCIDR, Value: "1.1.1.0/24"}, // 应被剔除
	}
	b, err := compiler.GeositeDat([]compiler.GeositeCategory{{Name: "streaming", Rules: rules}})
	require.NoError(t, err)

	site := parseGeositeDat(t, b)["STREAMING"]
	var got []string
	for _, d := range site.Domains {
		got = append(got, d.Value)
	}
	sort.Strings(got)
	assert.Equal(t, []string{"a.com", "b.com", "full.c.com"}, got)
}

// 空类目（全 IP 或无规则）仍产出 GeoSite 条目（合法 country_code，0 域名）
func TestGeositeDat_EmptyCategory(t *testing.T) {
	t.Parallel()
	cats := []compiler.GeositeCategory{
		{Name: "geoip-cn", Rules: []ruleset.Rule{{Match: ruleset.MatchIPCIDR, Value: "1.0.1.0/24"}}},
		{Name: "empty", Rules: nil},
	}
	b, err := compiler.GeositeDat(cats)
	require.NoError(t, err)

	sites := parseGeositeDat(t, b)
	require.Contains(t, sites, "GEOIP-CN")
	assert.Empty(t, sites["GEOIP-CN"].Domains, "纯 IP 类目无 domain")
	require.Contains(t, sites, "EMPTY")
	assert.Empty(t, sites["EMPTY"].Domains)
}

// WriteGeositeDat 在 outDir 写出单文件 geosite.dat，内容可被反解
func TestWriteGeositeDat(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cats := []compiler.GeositeCategory{
		{Name: "ai", Rules: []ruleset.Rule{{Match: ruleset.MatchDomainSuffix, Value: "claude.com"}}},
		{Name: "netflix", Rules: []ruleset.Rule{{Match: ruleset.MatchDomain, Value: "www.netflix.com"}}},
	}
	require.NoError(t, compiler.WriteGeositeDat(dir, cats))

	path := filepath.Join(dir, "geosite.dat")
	require.FileExists(t, path)
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	sites := parseGeositeDat(t, b)
	assert.Equal(t, "claude.com", sites["AI"].Domains[0].Value)
	assert.Equal(t, int32(3), sites["NETFLIX"].Domains[0].Type, "DOMAIN → Full(3)")
}
