package compiler

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/netip"
	"os"
	"sort"
	"strings"

	"github.com/river2heaven/tongtian/ruleset"
)

// 本文件是 geosite/chinalist 之外的「新血统」上游解析器（见 README「选定的源头」）：
//   - gfwlist：base64 AutoProxy → 被墙域名（GFW 一手源头）
//   - domainlist：纯域名 / AdGuard / clash-list 行（anti-AD / hagezi / AI 长尾，去广告 + AI 补全）
//   - geoip：IP-CIDR 列表（misakaio/chnroutes2 → GEOIP-CN，补 IP 侧分流）
//
// 所有解析器都产 *headless* 规则（无策略），并在出口去重 + 稳定排序（golden 可比对）。

// dedupSortDomains 按 Match|Value 去重，按 (Match,Value) 稳定排序。
func dedupSortRules(in []ruleset.Rule) []ruleset.Rule {
	seen := make(map[string]bool, len(in))
	out := in[:0:0]
	for _, r := range in {
		if r.Value == "" {
			continue
		}
		k := string(r.Match) + "|" + r.Value
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Match != out[j].Match {
			return out[i].Match < out[j].Match
		}
		return out[i].Value < out[j].Value
	})
	return out
}

// extractDomain 从一行（可能带 ||、|scheme://、前导点、路径、^ 等）里抽出裸域名并小写化。
// 返回 ok=false 表示这行不含可用域名（注释 / 正则 / IP / 畸形）。
func extractDomain(s string) (string, bool) {
	s = strings.TrimSpace(s)
	// 去 AutoProxy / AdGuard 前缀
	s = strings.TrimPrefix(s, "||")
	s = strings.TrimPrefix(s, "|")
	s = strings.TrimPrefix(s, "*.")
	s = strings.TrimPrefix(s, ".")
	// 去 scheme
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	// 截到第一个分隔符（路径 / AdGuard ^ / 通配 / 端口 / 选项）
	if i := strings.IndexAny(s, "/^*?:= \t"); i >= 0 {
		s = s[:i]
	}
	s = strings.ToLower(strings.TrimSpace(s))
	if !isPlausibleDomain(s) {
		return "", false
	}
	return s, true
}

// isPlausibleDomain 结构性校验：含点、1+ 标签、每标签 [a-z0-9_-] 1-63、总长 ≤253、非纯数字（防 IP）。
func isPlausibleDomain(s string) bool {
	if s == "" || len(s) > 253 || !strings.Contains(s, ".") {
		return false
	}
	allDigit := true
	for _, label := range strings.Split(s, ".") {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		for _, c := range label {
			switch {
			case c >= 'a' && c <= 'z', c >= '0' && c <= '9', c == '-', c == '_':
			default:
				return false
			}
			if c < '0' || c > '9' {
				allDigit = false
			}
		}
	}
	return !allDigit // 全数字 = IP，不当域名
}

// scanLines 逐行读文件，跳过空行，交给 fn（fn 返回的 rule 累加）。
func scanLines(path string, fn func(line string) (ruleset.Rule, bool)) ([]ruleset.Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开 %s: %w", path, err)
	}
	defer f.Close()
	var out []ruleset.Rule
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if r, ok := fn(line); ok {
			out = append(out, r)
		}
	}
	return out, sc.Err()
}

// ParseGFWList 解析 gfwlist（base64 编码的 AutoProxy 规则）→ 被墙域名 DOMAIN-SUFFIX 集合。
//
// AutoProxy 语法：`||domain` / `.domain` / `|http://host/...` / 裸域 → 代理；
// `!`/`[` 注释、`@@` 白名单、`/regex/` 正则 → 跳过（v1 不消费白名单与正则，避免误差）。
func ParseGFWList(paths ...string) ([]ruleset.Rule, error) {
	var all []ruleset.Rule
	for _, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("读取 gfwlist %s: %w", p, err)
		}
		text := decodeMaybeBase64(raw)
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "[") {
				continue // 注释 / 头
			}
			if strings.HasPrefix(line, "@@") {
				continue // 白名单（v1 不消费）
			}
			if strings.HasPrefix(line, "/") && strings.HasSuffix(line, "/") {
				continue // 正则（v1 不消费）
			}
			if d, ok := extractDomain(line); ok {
				all = append(all, ruleset.Rule{Match: ruleset.MatchDomainSuffix, Value: d})
			}
		}
	}
	return dedupSortRules(all), nil
}

// decodeMaybeBase64 容忍两种 gfwlist 镜像：base64 编码 或 已解码明文。
func decodeMaybeBase64(raw []byte) string {
	compact := strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == ' ' || r == '\t' {
			return -1
		}
		return r
	}, string(raw))
	if dec, err := base64.StdEncoding.DecodeString(compact); err == nil && len(dec) > 0 {
		return string(dec)
	}
	return string(raw)
}

// ParseDomainList 容忍式解析「纯域名 / AdGuard / clash-list」清单
// （anti-AD、hagezi onlydomains、VPSDance/dler/fmz200 的 AI 列表都能吃）：
//   - `# ! comment` → 跳过
//   - 含逗号 `TYPE,value[,policy]`（clash/surge 规则行）→ 反查 TYPE → 对应 matcher
//   - 否则视为裸域（容忍 `*.` / `.` / `||...^` 前缀）→ DOMAIN-SUFFIX
func ParseDomainList(paths ...string) ([]ruleset.Rule, error) {
	var all []ruleset.Rule
	for _, p := range paths {
		rs, err := scanLines(p, parseDomainLine)
		if err != nil {
			return nil, err
		}
		all = append(all, rs...)
	}
	return dedupSortRules(all), nil
}

func parseDomainLine(line string) (ruleset.Rule, bool) {
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
		return ruleset.Rule{}, false
	}
	if strings.ContainsRune(line, ',') { // clash/surge 规则行 TYPE,value[,policy]
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			return ruleset.Rule{}, false
		}
		m, ok := clashTypeToMatch(strings.ToUpper(strings.TrimSpace(parts[0])))
		if !ok {
			return ruleset.Rule{}, false
		}
		v := strings.TrimSpace(parts[1])
		switch m {
		case ruleset.MatchIPCIDR, ruleset.MatchIPCIDR6:
			if _, err := netip.ParsePrefix(v); err != nil {
				return ruleset.Rule{}, false
			}
			return ruleset.Rule{Match: m, Value: v}, true
		case ruleset.MatchDomainKeyword: // 关键字值无需含点
			v = strings.ToLower(v)
			if !isKeywordSafe(v) {
				return ruleset.Rule{}, false
			}
			return ruleset.Rule{Match: m, Value: v}, true
		default: // DOMAIN / DOMAIN-SUFFIX
			v = strings.ToLower(strings.TrimPrefix(strings.TrimPrefix(v, "*."), "."))
			if !isPlausibleDomain(v) {
				return ruleset.Rule{}, false
			}
			return ruleset.Rule{Match: m, Value: v}, true
		}
	}
	if d, ok := extractDomain(line); ok {
		return ruleset.Rule{Match: ruleset.MatchDomainSuffix, Value: d}, true
	}
	return ruleset.Rule{}, false
}

// isKeywordSafe 校验 DOMAIN-KEYWORD 值（子串匹配，可无点）：[a-z0-9._-]，≤100。
func isKeywordSafe(s string) bool {
	if s == "" || len(s) > 100 {
		return false
	}
	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9', c == '-', c == '_', c == '.':
		default:
			return false
		}
	}
	return true
}

// clashTypeToMatch 把 clash/surge 规则关键字反查成 MatchType（仅域名 + IP 类）。
func clashTypeToMatch(t string) (ruleset.MatchType, bool) {
	switch t {
	case "DOMAIN":
		return ruleset.MatchDomain, true
	case "DOMAIN-SUFFIX":
		return ruleset.MatchDomainSuffix, true
	case "DOMAIN-KEYWORD":
		return ruleset.MatchDomainKeyword, true
	case "IP-CIDR":
		return ruleset.MatchIPCIDR, true
	case "IP-CIDR6":
		return ruleset.MatchIPCIDR6, true
	}
	return "", false
}

// ParseGeoIP 解析 IP-CIDR 列表（misakaio/chnroutes2 等）→ IP-CIDR / IP-CIDR6 集合。
// 容忍裸 IP（补 /32 或 /128）；按 v4/v6 分流到对应 matcher。
func ParseGeoIP(paths ...string) ([]ruleset.Rule, error) {
	var all []ruleset.Rule
	for _, p := range paths {
		rs, err := scanLines(p, parseGeoIPLine)
		if err != nil {
			return nil, err
		}
		all = append(all, rs...)
	}
	return dedupSortRules(all), nil
}

func parseGeoIPLine(line string) (ruleset.Rule, bool) {
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
		return ruleset.Rule{}, false
	}
	// 容忍尾部注释
	if i := strings.IndexAny(line, " \t#"); i >= 0 {
		line = strings.TrimSpace(line[:i])
	}
	var p netip.Prefix
	if pp, err := netip.ParsePrefix(line); err == nil {
		p = pp
	} else if a, err2 := netip.ParseAddr(line); err2 == nil {
		bits := 32
		if a.Is6() && !a.Is4In6() {
			bits = 128
		}
		p = netip.PrefixFrom(a, bits)
	} else {
		return ruleset.Rule{}, false
	}
	if p.Addr().Is4() {
		return ruleset.Rule{Match: ruleset.MatchIPCIDR, Value: p.String()}, true
	}
	return ruleset.Rule{Match: ruleset.MatchIPCIDR6, Value: p.String()}, true
}
