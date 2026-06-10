package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/river2heaven/tongtian/ruleset"
)

// ClashClassicalList 渲染 mihomo classical text rule-provider 内容：每行 `<matcher>,<value>`
// （rule-set 文件只含 matcher，无 policy——策略在订阅引用 RULE-SET 时施加）。
// clash（behavior: classical, format: text）与 surge（RULE-SET text）共用此格式。
func ClashClassicalList(rules []ruleset.Rule) string {
	var b strings.Builder
	for _, r := range rules {
		kw, ok := ruleset.ClashKeyword(r.Match)
		if !ok {
			continue
		}
		b.WriteString(kw)
		b.WriteByte(',')
		b.WriteString(r.Value)
		b.WriteByte('\n')
	}
	return b.String()
}

// SingboxSource 渲染 sing-box rule-set 源 JSON（version 2，支持 domain_regex / ip_cidr），
// 喂 `sing-box rule-set compile` 产 .srs。
func SingboxSource(rules []ruleset.Rule) (string, error) {
	headless := map[string][]string{}
	for _, r := range rules {
		var key string
		switch r.Match {
		case ruleset.MatchDomain:
			key = "domain"
		case ruleset.MatchDomainSuffix:
			key = "domain_suffix"
		case ruleset.MatchDomainKeyword:
			key = "domain_keyword"
		case ruleset.MatchDomainRegex:
			key = "domain_regex"
		case ruleset.MatchIPCIDR, ruleset.MatchIPCIDR6:
			key = "ip_cidr" // sing-box 用 ip_cidr 统一 v4/v6
		default:
			continue
		}
		headless[key] = append(headless[key], r.Value)
	}
	doc := map[string]any{"version": 2, "rules": []any{headless}}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

// Tools 是可选外部 CLI 路径（CI 装对应版本）。
// 空 = 跳过该二进制产物——本地 `go run` 仍产出 .list + .singbox.json 文本产物。
type Tools struct {
	SingBox string // sing-box 可执行；非空时编译 .srs
}

// WriteCategory 把一个类别的规则扇出到 outDir：
//
//	<name>.list          clash classical + surge 共用（remote rule-provider, format: text）
//	<name>.singbox.json  sing-box rule-set 源
//	<name>.srs           sing-box 二进制（仅 tools.SingBox 非空时）
//
// 注：mihomo .mrs 只支持 domain/ipcidr behavior，不支持本类 classical 混合 matcher；
// mihomo 可直接远程加载 .list（format: text），故 .mrs 列为后续性能优化，本期不产出。
func WriteCategory(outDir, name string, rules []ruleset.Rule, tools Tools) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	listPath := filepath.Join(outDir, name+".list")
	if err := os.WriteFile(listPath, []byte(ClashClassicalList(rules)), 0o644); err != nil {
		return err
	}
	src, err := SingboxSource(rules)
	if err != nil {
		return err
	}
	srcPath := filepath.Join(outDir, name+".singbox.json")
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		return err
	}
	if tools.SingBox != "" {
		srs := filepath.Join(outDir, name+".srs")
		cmd := exec.Command(tools.SingBox, "rule-set", "compile", "--output", srs, srcPath)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("sing-box rule-set compile %s: %w", name, err)
		}
	}
	return nil
}
