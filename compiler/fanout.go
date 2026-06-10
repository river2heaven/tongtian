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
	Mihomo  string // mihomo 可执行；非空时为「纯 domain / 纯 ip-cidr」类别编译 .mrs
}

// WriteCategory 把一个类别的规则扇出到 outDir：
//
//	<name>.list          clash classical + surge 共用（remote rule-provider, format: text）
//	<name>.singbox.json  sing-box rule-set 源
//	<name>.srs           sing-box 二进制（仅 tools.SingBox 非空时）
//	<name>.mrs           mihomo 二进制（仅 tools.Mihomo 非空 且 类别为纯 domain / 纯 ip-cidr）
//
// 注：mihomo .mrs 只支持 domain / ipcidr behavior。含 DOMAIN-KEYWORD / DOMAIN-REGEX 或
// domain+ip 混合的类别（netflix / ai / disney）无 .mrs，clash 端对它们继续用 .list（format: text）。
// 纯域名大表（cn / reject 等）有 .mrs 时体积 3-5× 缩小 + 客户端加载/匹配更快。
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
	if tools.Mihomo != "" {
		if err := writeMRS(tools.Mihomo, outDir, name, rules); err != nil {
			return err
		}
	}
	return nil
}

// mrsBehavior 判定类别的 .mrs 资格：全 domain → "domain"；全 ip-cidr → "ipcidr"；
// 含 keyword/regex 或 domain+ip 混合 → ""（不产 .mrs）。
func mrsBehavior(rules []ruleset.Rule) string {
	if len(rules) == 0 {
		return ""
	}
	hasDomain, hasIP := false, false
	for _, r := range rules {
		switch r.Match {
		case ruleset.MatchDomain, ruleset.MatchDomainSuffix:
			hasDomain = true
		case ruleset.MatchIPCIDR, ruleset.MatchIPCIDR6:
			hasIP = true
		default: // keyword / regex / 其它 → mrs 不支持
			return ""
		}
	}
	switch {
	case hasDomain && !hasIP:
		return "domain"
	case hasIP && !hasDomain:
		return "ipcidr"
	default:
		return "" // domain + ip 混合（如 ai）
	}
}

// writeMRS 为符合条件的类别产 <name>.mrs：把规则转成 mihomo domain/ipcidr behavior payload，
// 写临时文件喂 `mihomo convert-ruleset`。不符合（含 keyword/regex/混合）则跳过、维持 .list。
func writeMRS(mihomoBin, outDir, name string, rules []ruleset.Rule) error {
	behavior := mrsBehavior(rules)
	if behavior == "" {
		return nil
	}
	tmp, err := os.CreateTemp("", name+"-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(mihomoPayload(behavior, rules)); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	mrs := filepath.Join(outDir, name+".mrs")
	cmd := exec.Command(mihomoBin, "convert-ruleset", behavior, "text", tmp.Name(), mrs)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mihomo convert-ruleset %s (%s): %w", name, behavior, err)
	}
	return nil
}

// mihomoPayload 渲染 mihomo domain/ipcidr behavior 的 text payload：
// DOMAIN-SUFFIX,x → "+.x"；DOMAIN,x → "x"；IP-CIDR,x → "x"。
func mihomoPayload(behavior string, rules []ruleset.Rule) string {
	var b strings.Builder
	for _, r := range rules {
		switch {
		case behavior == "domain" && r.Match == ruleset.MatchDomainSuffix:
			b.WriteString("+.")
			b.WriteString(r.Value)
			b.WriteByte('\n')
		case behavior == "domain" && r.Match == ruleset.MatchDomain:
			b.WriteString(r.Value)
			b.WriteByte('\n')
		case behavior == "ipcidr" && (r.Match == ruleset.MatchIPCIDR || r.Match == ruleset.MatchIPCIDR6):
			b.WriteString(r.Value)
			b.WriteByte('\n')
		}
	}
	return b.String()
}
