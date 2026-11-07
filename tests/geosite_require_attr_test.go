// 通天河 编译器：geosite_require_attr 测试（只筛 geosite 这一路，每个叶子都必须带该属性）。
package compiler_test

import (
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func leafAttrManifest() *compiler.Manifest {
	return &compiler.Manifest{
		Upstreams: map[string]compiler.UpstreamRef{
			"geosite":   {Repo: "v2fly/domain-list-community", Commit: "76199d8e61d348c58433e3c36d9d4f85131d80e7", DataDir: "."},
			"chinalist": {Repo: "felixonmars/dnsmasq-china-list", Commit: "d87a5bc1d76b87a6729e9e5355c35be0fc2ff6cd", Files: []string{"china.conf"}},
		},
	}
}

// 带该属性的叶子留下，不带的丢弃，包括经过子 include 的。
func TestGeositeRequireAttr_KeepsOnlyTaggedLeaves(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(leafAttrManifest(), "fixtures", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:               "cn",
		Geosite:            []string{"leafattr"},
		GeositeRequireAttr: "cn",
	})
	require.NoError(t, err)

	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "cdn.example"), "顶层 @cn")
	assert.True(t, hasRule(rules, ruleset.MatchDomain, "exact-cdn.example"), "full: 也算叶子，匹配类型不变")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "tagged.example"), "带多个属性，含 @cn 即可")
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "plain-cdn.example"), "穿透无属性的子 include")

	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "store.example"), "顶层无属性")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "plain-store.example"), "子 include 里无属性")
}

// 回归：子 include 自带 @!cn 时，只把属性传给现有 requireAttr 会被覆盖，商店域名混进 cn。
func TestGeositeRequireAttr_NotOverriddenByChildInclude(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(leafAttrManifest(), "fixtures", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:               "cn",
		Geosite:            []string{"leafattr"},
		GeositeRequireAttr: "cn",
	})
	require.NoError(t, err)
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "negated-store.example"), "@!cn 必须被拒收")
}

// include 自身的 @attr 过滤继续生效，和 geosite_require_attr 是“且”的关系。
func TestGeositeRequireAttr_AndWithIncludeAttr(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(leafAttrManifest(), "fixtures", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:               "cn",
		Geosite:            []string{"leafattr"},
		GeositeRequireAttr: "cn",
	})
	require.NoError(t, err)
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "both.example"), "@cdn 且 @cn")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "onlycdn.example"), "只有 @cdn")
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "onlycn.example"), "include 要求 @cdn，只有 @cn 不算")
}

// 只筛 geosite 这一路：chinalist 和 localdomains 的条目没有属性，也必须全部留下。
func TestGeositeRequireAttr_LeavesOtherSourcesAlone(t *testing.T) {
	t.Parallel()
	m := leafAttrManifest()
	r := compiler.NewResolver(m, "fixtures", "fixtures")

	base, err := r.ResolveCategory(compiler.Category{
		Name:         "cn",
		ChinaList:    "chinalist",
		LocalDomains: "localdomains/cn-extra.txt",
	})
	require.NoError(t, err)
	require.NotEmpty(t, base)

	merged, err := r.ResolveCategory(compiler.Category{
		Name:               "cn",
		ChinaList:          "chinalist",
		LocalDomains:       "localdomains/cn-extra.txt",
		Geosite:            []string{"leafattr"},
		GeositeRequireAttr: "cn",
	})
	require.NoError(t, err)

	for _, rr := range base {
		assert.True(t, hasRule(merged, rr.Match, rr.Value), "原有条目 %s 不能少", rr.Value)
	}
	assert.True(t, hasRule(merged, ruleset.MatchDomainSuffix, "cdn.example"), "geosite 这一路并入")
	assert.False(t, hasRule(merged, ruleset.MatchDomainSuffix, "store.example"), "geosite 这一路被筛")
}

// 不配时行为不变，也不影响 exclude_attrs。
func TestGeositeRequireAttr_UnsetKeepsExistingBehavior(t *testing.T) {
	t.Parallel()
	r := compiler.NewResolver(leafAttrManifest(), "fixtures", "fixtures")
	rules, err := r.ResolveCategory(compiler.Category{
		Name:         "steam",
		Geosite:      []string{"leafattr"},
		ExcludeAttrs: []string{"cn"},
	})
	require.NoError(t, err)
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "store.example"))
	assert.True(t, hasRule(rules, ruleset.MatchDomainSuffix, "plain-store.example"))
	assert.False(t, hasRule(rules, ruleset.MatchDomainSuffix, "cdn.example"), "exclude_attrs 丢 @cn")
}

// 配了属性但没有 geosite：字段不起作用，属于写错的 manifest，校验期就要报错。
func TestValidateGeositeRequireAttr_NeedsGeosite(t *testing.T) {
	t.Parallel()
	m := leafAttrManifest()
	m.Categories = []compiler.Category{{Name: "cn", ChinaList: "chinalist", GeositeRequireAttr: "cn"}}
	err := m.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "geosite_require_attr")

	m.Categories = []compiler.Category{{Name: "cn", Geosite: []string{"leafattr"}, GeositeRequireAttr: "cn"}}
	require.NoError(t, m.Validate())
}
