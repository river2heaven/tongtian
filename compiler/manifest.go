package compiler

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Manifest 是编译清单：钉上游版本 + 类别定义 + 自有 repo/base_url。
//
// 铁律：上游钉 commit/tag（非滚动分支），防某次坏 PR 直接进生产。
// 钉版本由 Validate() 在 Go 侧 fail-fast 强制（不只靠 CI shell）。
type Manifest struct {
	Repo       string                 `yaml:"repo"`     // 自有产物 repo（<org>/<repo>），jsdelivr 分发
	BaseURL    string                 `yaml:"base_url"` // 订阅模板引用的基址（cutover 用）
	Upstreams  map[string]UpstreamRef `yaml:"upstreams"`
	Categories []Category             `yaml:"categories"`
}

// UpstreamRef 一个上游的钉版本引用。Files / DataDir 相对该上游 checkout 根。
type UpstreamRef struct {
	Repo    string   `yaml:"repo"`
	Commit  string   `yaml:"commit"`   // 钉 commit/tag
	DataDir string   `yaml:"data_dir"` // geosite: data 子目录（设了即视为 geosite 上游）
	Files   []string `yaml:"files"`    // 其它解析器: 要读的文件
}

// Category 一个输出类别。可同时引用多个输入源，编译器取并集 + 去重。
//
//	geosite      → domain-list-community 类别名（合并）
//	chinalist    → dnsmasq-china-list 类上游 key
//	gfwlist      → gfwlist 类上游 key
//	geoip        → IP-CIDR 类上游 key（chnroutes2）
//	domainlists  → 纯域名 / AdGuard / clash-list 类上游 key（合并：anti-AD + hagezi + AI 等）
type Category struct {
	Name         string   `yaml:"name"`
	Geosite      []string `yaml:"geosite"`
	ExcludeAttrs []string `yaml:"exclude_attrs"`
	ChinaList    string   `yaml:"chinalist"`
	GFWList      string   `yaml:"gfwlist"`
	GeoIP        string   `yaml:"geoip"`
	DomainLists  []string `yaml:"domainlists"`
}

// LoadManifest 读取并解析 manifest.yaml（不校验；调用方再调 Validate）。
func LoadManifest(path string) (*Manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 manifest %s: %w", path, err)
	}
	var m Manifest
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("解析 manifest: %w", err)
	}
	return &m, nil
}

// Validate fail-fast 校验所有上游已钉真实版本（40 位 SHA 或合法 tag），
// 顶掉占位符与滚动分支（master/main/...）。校验在 Go 侧而非只在 CI shell。
func (m *Manifest) Validate() error {
	if len(m.Upstreams) == 0 {
		return fmt.Errorf("manifest 未声明任何上游")
	}
	for name, u := range m.Upstreams {
		if u.Repo == "" {
			return fmt.Errorf("上游 %q 缺 repo", name)
		}
		if err := validateCommit(name, u.Commit); err != nil {
			return err
		}
	}
	// 类别引用的上游 key 必须存在
	for _, cat := range m.Categories {
		for _, key := range append([]string{cat.ChinaList, cat.GFWList, cat.GeoIP}, cat.DomainLists...) {
			if key != "" {
				if _, ok := m.Upstreams[key]; !ok {
					return fmt.Errorf("类别 %q 引用了未声明的上游 %q", cat.Name, key)
				}
			}
		}
	}
	return nil
}

func validateCommit(name, c string) error {
	if c == "" || strings.HasPrefix(c, "REPLACE_") {
		return fmt.Errorf("上游 %s 未钉 commit/tag（commit=%q）——填真实 SHA 或 tag", name, c)
	}
	if isHex40(c) || looksLikeTag(c) {
		return nil
	}
	return fmt.Errorf("上游 %s 的 commit=%q 不是 40 位 SHA 也不像 tag（禁滚动分支 master/main/release/...）", name, c)
}

func isHex40(s string) bool {
	if len(s) != 40 {
		return false
	}
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'f', r >= 'A' && r <= 'F':
		default:
			return false
		}
	}
	return true
}

// looksLikeTag 接受版本 / 日期类 tag（含数字），拒掉无数字的滚动分支名。
func looksLikeTag(c string) bool {
	switch strings.ToLower(c) {
	case "master", "main", "head", "latest", "release", "dev", "develop", "trunk", "stable", "nightly":
		return false
	}
	for _, r := range c {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}
