// 通天河 编译器：Anywhere .arrs 扇出测试（issue #1）。
//
// .arrs 格式（NodePassProject/Anywhere Documentations/Routing.md）：
// 头部 `name = <value>`，规则行 `<type>, <value>`，type 0=IPv4 CIDR / 1=IPv6 CIDR /
// 2=domain-suffix / 3=domain-keyword；单文件硬上限 10,000 条，超限客户端整体拒绝。
package compiler_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/river2heaven/tongtian/compiler"
	"github.com/river2heaven/tongtian/ruleset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 四类直映 + 头部 + 行格式（golden）
func TestArrsList_Mapping(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "example.com"},
		{Match: ruleset.MatchDomainKeyword, Value: "example"},
		{Match: ruleset.MatchIPCIDR, Value: "10.0.0.0/8"},
		{Match: ruleset.MatchIPCIDR6, Value: "2001:db8::/32"},
	}
	content, kept, dropped := compiler.ArrsList("netflix", rules)

	want := "name = netflix\n" +
		"2, example.com\n" +
		"3, example\n" +
		"0, 10.0.0.0/8\n" +
		"1, 2001:db8::/32\n"
	assert.Equal(t, want, content)
	assert.Equal(t, 4, kept)
	assert.Equal(t, 0, dropped)
}

// DOMAIN（exact）降级为 suffix（.arrs 无 exact 类型）
func TestArrsList_DomainDowngradesToSuffix(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomain, Value: "www.netflix.com"},
	}
	content, kept, dropped := compiler.ArrsList("x", rules)

	assert.Equal(t, "name = x\n2, www.netflix.com\n", content)
	assert.Equal(t, 1, kept)
	assert.Equal(t, 0, dropped)
}

// DOMAIN-REGEX 无 .arrs 对应 → 丢弃并计数，不进产物
func TestArrsList_RegexDropped(t *testing.T) {
	t.Parallel()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainRegex, Value: `.*\.nflxext\.com`},
		{Match: ruleset.MatchDomainSuffix, Value: "netflix.com"},
	}
	content, kept, dropped := compiler.ArrsList("x", rules)

	assert.Equal(t, "name = x\n2, netflix.com\n", content)
	assert.Equal(t, 1, kept)
	assert.Equal(t, 1, dropped)
	assert.NotContains(t, content, "nflxext")
}

// 空规则集仍产出合法头部（仅 name 行）
func TestArrsList_Empty(t *testing.T) {
	t.Parallel()
	content, kept, dropped := compiler.ArrsList("empty", nil)

	assert.Equal(t, "name = empty\n", content)
	assert.Equal(t, 0, kept)
	assert.Equal(t, 0, dropped)
}

// WriteCategory 增产 <name>.arrs，与 .list / .singbox.json 并列
func TestWriteCategory_ArrsEmitted(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rules := []ruleset.Rule{
		{Match: ruleset.MatchDomainSuffix, Value: "example.com"},
		{Match: ruleset.MatchIPCIDR, Value: "10.0.0.0/8"},
	}
	require.NoError(t, compiler.WriteCategory(dir, "cat", rules, compiler.Tools{}))

	b, err := os.ReadFile(filepath.Join(dir, "cat.arrs"))
	require.NoError(t, err)
	assert.Equal(t, "name = cat\n2, example.com\n0, 10.0.0.0/8\n", string(b))

	// 既有产物不受影响
	assert.FileExists(t, filepath.Join(dir, "cat.list"))
	assert.FileExists(t, filepath.Join(dir, "cat.singbox.json"))
}

// 映射后超 10k 上限：跳过 .arrs（客户端会整体拒绝），其余产物照常
func TestWriteCategory_ArrsOverCapSkipped(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rules := make([]ruleset.Rule, 0, compiler.ArrsMaxRules+1)
	for i := 0; i <= compiler.ArrsMaxRules; i++ {
		rules = append(rules, ruleset.Rule{
			Match: ruleset.MatchDomainSuffix,
			Value: fmt.Sprintf("d%d.example.com", i),
		})
	}
	require.NoError(t, compiler.WriteCategory(dir, "cn", rules, compiler.Tools{}))

	assert.NoFileExists(t, filepath.Join(dir, "cn.arrs"))
	assert.FileExists(t, filepath.Join(dir, "cn.list"))
	assert.FileExists(t, filepath.Join(dir, "cn.singbox.json"))
}

// 恰好 10k 条：不超限，正常产出（边界）
func TestWriteCategory_ArrsAtCapEmitted(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rules := make([]ruleset.Rule, 0, compiler.ArrsMaxRules)
	for i := 0; i < compiler.ArrsMaxRules; i++ {
		rules = append(rules, ruleset.Rule{
			Match: ruleset.MatchDomainSuffix,
			Value: fmt.Sprintf("d%d.example.com", i),
		})
	}
	require.NoError(t, compiler.WriteCategory(dir, "gfw", rules, compiler.Tools{}))

	assert.FileExists(t, filepath.Join(dir, "gfw.arrs"))
}

// 超限判定按「映射后 kept」算：REGEX 丢弃不计入（10k 条 suffix + 1 条 regex 不超限）
func TestWriteCategory_ArrsCapCountsKeptOnly(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	rules := make([]ruleset.Rule, 0, compiler.ArrsMaxRules+1)
	for i := 0; i < compiler.ArrsMaxRules; i++ {
		rules = append(rules, ruleset.Rule{
			Match: ruleset.MatchDomainSuffix,
			Value: fmt.Sprintf("d%d.example.com", i),
		})
	}
	rules = append(rules, ruleset.Rule{Match: ruleset.MatchDomainRegex, Value: `.*\.x\.com`})
	require.NoError(t, compiler.WriteCategory(dir, "edge", rules, compiler.Tools{}))

	assert.FileExists(t, filepath.Join(dir, "edge.arrs"))
}
