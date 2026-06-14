// Package compiler 是 通天河 的自编译 CI 编译器：拉上游 → 解析各上游格式 →
// 统一规则模型（ruleset.Rule）→ 扇出各客户端 rule-set。
//
// 上游解析器（见 README「选定的源头」）：
//   - geosite：v2fly/domain-list-community 的 DSL（include / @attr / full,keyword,regexp）
//   - chinalist：felixonmars/dnsmasq-china-list（`server=/domain/dns` 行 → CN 域名）
//   - gfwlist：gfwlist/gfwlist（base64 AutoProxy 语法 → 被墙域名）
//   - domainlist：anti-AD / hagezi / AI 长尾等纯域名清单
//   - geoip：misakaio/chnroutes2（IP-CIDR 列表 → GEOIP-CN）
//
// 铁律：只取源数据，分类边界 / 格式 / 托管全由本编译器掌控；上游钉 commit（manifest）。
package compiler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/river2heaven/tongtian/ruleset"
)

// rawEntry 是 geosite 文件解析出的一行（include 未展开）。
type rawEntry struct {
	kind        string   // domain | full | keyword | regexp | include
	value       string   // 域名 / 关键字 / 正则 / 被 include 的类别名
	attrs       []string // @ads @cn → ["ads","cn"]
	includeAttr string   // include:cat @attr → 仅取该类别中带此 attr 的条目
}

// Geosite 是 domain-list-community data 目录的内存表示。
type Geosite struct {
	categories map[string][]rawEntry
}

// LoadGeosite 解析 data 目录下所有类别文件（文件名即类别名）。
func LoadGeosite(dir string) (*Geosite, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取 geosite 目录 %s: %w", dir, err)
	}
	g := &Geosite{categories: make(map[string][]rawEntry)}
	for _, de := range entries {
		if de.IsDir() {
			continue
		}
		rows, err := parseGeositeFile(filepath.Join(dir, de.Name()))
		if err != nil {
			return nil, err
		}
		g.categories[de.Name()] = rows
	}
	return g, nil
}

func parseGeositeFile(path string) ([]rawEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开 %s: %w", path, err)
	}
	defer f.Close()

	var out []rawEntry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		if e, ok := parseGeositeLine(sc.Text()); ok {
			out = append(out, e)
		}
	}
	return out, sc.Err()
}

// parseGeositeLine 解析一行 geosite DSL（翻译 regexp,keyword,full）。
func parseGeositeLine(line string) (rawEntry, bool) {
	// 去行内注释
	if i := strings.IndexByte(line, '#'); i >= 0 {
		line = line[:i]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return rawEntry{}, false
	}

	fields := strings.Fields(line)
	head := fields[0]

	var attrs []string
	for _, a := range fields[1:] {
		if strings.HasPrefix(a, "@") {
			attrs = append(attrs, strings.TrimPrefix(a, "@"))
		}
	}

	if rest, ok := strings.CutPrefix(head, "include:"); ok {
		ia := ""
		if len(attrs) > 0 {
			ia = attrs[0]
		}
		return rawEntry{kind: "include", value: rest, includeAttr: ia}, true
	}

	kind, value := "domain", head
	if i := strings.IndexByte(head, ':'); i >= 0 {
		switch head[:i] {
		case "full", "keyword", "regexp", "domain":
			kind, value = head[:i], head[i+1:]
		}
	}
	return rawEntry{kind: kind, value: value, attrs: attrs}, true
}

// Resolve 把一个类别展开成统一规则模型（递归 include + @attr 过滤 + 去重）。
//
//	excludeAttrs: manifest 级全局排除（如某流媒体类排掉 @ads 条目）。
//	所有产物无策略 —— rule-set 文件只含 matcher，策略在订阅引用 RULE-SET 时施加。
func (g *Geosite) Resolve(category string, excludeAttrs []string) ([]ruleset.Rule, error) {
	if _, ok := g.categories[category]; !ok {
		return nil, fmt.Errorf("geosite 类别不存在: %s", category)
	}
	seen := make(map[string]bool)
	var out []ruleset.Rule
	g.resolve(category, "", excludeAttrs, map[string]bool{}, seen, &out)
	// 稳定排序，golden 可比对
	sort.Slice(out, func(i, j int) bool {
		if out[i].Match != out[j].Match {
			return out[i].Match < out[j].Match
		}
		return out[i].Value < out[j].Value
	})
	return out, nil
}

func (g *Geosite) resolve(category, requireAttr string, excludeAttrs []string, visiting, seen map[string]bool, out *[]ruleset.Rule) {
	if visiting[category] {
		return // 环保护
	}
	visiting[category] = true
	defer delete(visiting, category)

	for _, e := range g.categories[category] {
		if e.kind == "include" {
			// include 的 @attr 过滤需穿透嵌套 include：子 include 无自带 attr 时继承父 requireAttr。
			childRequire := e.includeAttr
			if childRequire == "" {
				childRequire = requireAttr
			}
			g.resolve(e.value, childRequire, excludeAttrs, visiting, seen, out)
			continue
		}
		if requireAttr != "" && !hasAttr(e.attrs, requireAttr) {
			continue
		}
		if hasAnyAttr(e.attrs, excludeAttrs) {
			continue
		}
		r, ok := entryToRule(e)
		if !ok {
			continue
		}
		key := string(r.Match) + "|" + r.Value
		if seen[key] {
			continue
		}
		seen[key] = true
		*out = append(*out, r)
	}
}

func entryToRule(e rawEntry) (ruleset.Rule, bool) {
	switch e.kind {
	case "domain":
		return ruleset.Rule{Match: ruleset.MatchDomainSuffix, Value: e.value}, true
	case "full":
		return ruleset.Rule{Match: ruleset.MatchDomain, Value: e.value}, true
	case "keyword":
		return ruleset.Rule{Match: ruleset.MatchDomainKeyword, Value: e.value}, true
	case "regexp":
		return ruleset.Rule{Match: ruleset.MatchDomainRegex, Value: e.value}, true
	}
	return ruleset.Rule{}, false
}

func hasAttr(attrs []string, a string) bool {
	for _, x := range attrs {
		if x == a {
			return true
		}
	}
	return false
}

func hasAnyAttr(attrs, set []string) bool {
	for _, a := range set {
		if hasAttr(attrs, a) {
			return true
		}
	}
	return false
}

// ParseChinaList 解析 dnsmasq-china-list 的 `server=/domain/dns` 行 → CN DOMAIN-SUFFIX 集合
// （CN 域名表，DNS 分流核心）。同时容忍 `ipset=/domain/...` 行。
func ParseChinaList(paths ...string) ([]ruleset.Rule, error) {
	seen := make(map[string]bool)
	var out []ruleset.Rule
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("打开 china-list %s: %w", p, err)
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			// 只接受 dnsmasq-china-list 真正使用的两种指令：
			//   server=/baidu.com/114.114.114.114  或  ipset=/baidu.com/china
			// 收紧 key（此前接受任意 key=/.../，坏上游可借 address=/domain=/ 等行混入域名）。
			eq := strings.IndexByte(line, '=')
			if eq < 0 {
				continue
			}
			key := line[:eq]
			if key != "server" && key != "ipset" {
				continue
			}
			rest := line[eq+1:]
			if !strings.HasPrefix(rest, "/") {
				continue
			}
			parts := strings.Split(rest, "/") // ["", "baidu.com", "114..."]
			if len(parts) < 2 || parts[1] == "" {
				continue
			}
			domain := parts[1]
			// 域名形态校验：拒注入字符 + 非域名行（与扇出 SafeValue 一致的最后一道防线）。
			if !ruleset.SafeValue(ruleset.MatchDomainSuffix, domain) {
				continue
			}
			if seen[domain] {
				continue
			}
			seen[domain] = true
			out = append(out, ruleset.Rule{Match: ruleset.MatchDomainSuffix, Value: domain})
		}
		f.Close()
		if err := sc.Err(); err != nil {
			return nil, err
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Value < out[j].Value })
	return out, nil
}
