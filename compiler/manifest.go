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
	hasGeositeUpstream := false
	for name, u := range m.Upstreams {
		if u.Repo == "" {
			return fmt.Errorf("上游 %q 缺 repo", name)
		}
		if err := validateRepo(name, u.Repo); err != nil {
			return err
		}
		if err := validateCommit(name, u.Commit); err != nil {
			return err
		}
		// 上游角色自洽：geosite 上游（data_dir 非空）与文件型上游（files 非空）二选一，
		// 既不能两者皆空（无从读取），也不该同时设（角色歧义）。
		switch {
		case u.DataDir != "" && len(u.Files) > 0:
			return fmt.Errorf("上游 %q 同时设了 data_dir 与 files（角色歧义）", name)
		case u.DataDir != "":
			hasGeositeUpstream = true
		case len(u.Files) == 0:
			return fmt.Errorf("上游 %q 既无 data_dir 也无 files（无从读取）", name)
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
		// 引用 geosite 类目的类别，必须有一个声明了 data_dir 的 geosite 上游兜底，
		// 否则 ResolveCategory 会在运行期才炸（fail-fast 提前到校验期）。
		if len(cat.Geosite) > 0 && !hasGeositeUpstream {
			return fmt.Errorf("类别 %q 引用 geosite 类目，但 manifest 无 geosite 上游（无上游设 data_dir）", cat.Name)
		}
	}
	return nil
}

// validateRepo 校验上游 repo 形如 `<owner>/<name>`（防注入进 git clone URL 的畸形值）。
func validateRepo(name, repo string) error {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("上游 %s 的 repo=%q 不是 <owner>/<name> 形态", name, repo)
	}
	for _, r := range repo {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9',
			r == '/', r == '-', r == '_', r == '.':
		default:
			return fmt.Errorf("上游 %s 的 repo=%q 含非法字符 %q", name, repo, string(r))
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
//
// 加固：要求 tag 至少 3 字符且含数字，拒掉裸单字符 / 短 hex 前缀这类「像 pin 但不可靠」的值
// （git 解析短 SHA 前缀是有歧义的，不是不可变引用）；字符集限定 [A-Za-z0-9._-]，
// 防把畸形 / 注入串当成合法 pin。
func looksLikeTag(c string) bool {
	switch strings.ToLower(c) {
	case "master", "main", "head", "latest", "release", "dev", "develop", "trunk", "stable", "nightly":
		return false
	}
	if len(c) < 3 {
		return false
	}
	hasDigit := false
	for _, r := range c {
		switch {
		case r >= '0' && r <= '9':
			hasDigit = true
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '.', r == '-', r == '_':
		default:
			return false // 含非 tag 合法字符 → 拒
		}
	}
	return hasDigit
}
