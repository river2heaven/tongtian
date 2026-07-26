package compiler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/river2heaven/tongtian/ruleset"
)

// Resolver 按 manifest + 各上游本地 checkout 目录解析类别成统一规则。
//
// 约定：每个上游 checkout 在 upstreamsDir 下的子目录里，子目录名 = manifest 的 upstream key；
// UpstreamRef.Files / DataDir 相对该子目录。geosite 上游 = DataDir 非空那个。
type Resolver struct {
	m            *Manifest
	upstreamsDir string
	repoRoot     string   // 仓库根(manifest.yaml 所在目录); Category.LocalFile 相对此解析
	geosite      *Geosite // lazy
}

// NewResolver 构造解析器。repoRoot 是 manifest.yaml 所在目录，用于解析 Category.LocalFile
// （仓内 curated 数据文件，如 data/china-company-ip.cidr）。传空串则 localfile 类别会报错。
func NewResolver(m *Manifest, upstreamsDir, repoRoot string) *Resolver {
	return &Resolver{m: m, upstreamsDir: upstreamsDir, repoRoot: repoRoot}
}

// localFilePath 把 Category.LocalFile 解析为仓内绝对路径，并挡住 `..` 越权（防坏 manifest
// 读到仓外任意文件）。
func (r *Resolver) localFilePath(rel string) (string, error) {
	if r.repoRoot == "" {
		return "", fmt.Errorf("localfile %q 需要 repoRoot，但 Resolver 未配置", rel)
	}
	clean := filepath.Clean(rel)
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("localfile %q 越出仓库根（禁绝对路径 / ..）", rel)
	}
	return filepath.Join(r.repoRoot, clean), nil
}

func (r *Resolver) upstreamFiles(key string) ([]string, error) {
	u, ok := r.m.Upstreams[key]
	if !ok {
		return nil, fmt.Errorf("manifest 未声明上游 %q", key)
	}
	if len(u.Files) == 0 {
		return nil, fmt.Errorf("上游 %q 未配 files", key)
	}
	var fs []string
	for _, f := range u.Files {
		fs = append(fs, filepath.Join(r.upstreamsDir, key, f))
	}
	return fs, nil
}

func (r *Resolver) loadGeosite() (*Geosite, error) {
	if r.geosite != nil {
		return r.geosite, nil
	}
	key, dataDir := "", ""
	for k, u := range r.m.Upstreams {
		if u.DataDir != "" {
			key, dataDir = k, u.DataDir
			break
		}
	}
	if key == "" {
		return nil, fmt.Errorf("manifest 无 geosite 上游（没有上游设了 data_dir）")
	}
	g, err := LoadGeosite(filepath.Join(r.upstreamsDir, key, dataDir))
	if err != nil {
		return nil, err
	}
	r.geosite = g
	return g, nil
}

// ResolveCategory 取一个类别所有输入源的并集 + 去重 + 稳定排序。
func (r *Resolver) ResolveCategory(cat Category) ([]ruleset.Rule, error) {
	var all []ruleset.Rule

	if len(cat.Geosite) > 0 {
		g, err := r.loadGeosite()
		if err != nil {
			return nil, err
		}
		for _, gc := range cat.Geosite {
			rs, err := g.Resolve(gc, cat.ExcludeAttrs)
			if err != nil {
				return nil, err
			}
			all = append(all, rs...)
		}
	}

	type src struct {
		key   string
		parse func(...string) ([]ruleset.Rule, error)
	}
	for _, s := range []src{
		{cat.ChinaList, ParseChinaList},
		{cat.GFWList, ParseGFWList},
		{cat.GeoIP, ParseGeoIP},
	} {
		if s.key == "" {
			continue
		}
		fs, err := r.upstreamFiles(s.key)
		if err != nil {
			return nil, err
		}
		rs, err := s.parse(fs...)
		if err != nil {
			return nil, err
		}
		all = append(all, rs...)
	}

	for _, key := range cat.DomainLists {
		fs, err := r.upstreamFiles(key)
		if err != nil {
			return nil, err
		}
		rs, err := ParseDomainList(fs...)
		if err != nil {
			return nil, err
		}
		all = append(all, rs...)
	}

	// localfile：仓内 curated IP-CIDR 文件（如 data/china-company-ip.cidr，由
	// scripts/gen-china-company-ip.sh 从 ipverse/asn-ip 钉定 SHA 生成）。复用 ParseGeoIP
	// （跳 # 注释、v4/v6 分流），不占外部上游、CI 无需 clone 大仓。
	if cat.LocalFile != "" {
		p, err := r.localFilePath(cat.LocalFile)
		if err != nil {
			return nil, err
		}
		rs, err := ParseGeoIP(p)
		if err != nil {
			return nil, err
		}
		all = append(all, rs...)
	}

	// localdomains：仓内 curated 域名种子（如 data/cn-extra.txt，补上游清单未收录的域名）。
	// 复用 ParseDomainList（裸域 → DOMAIN-SUFFIX，容忍 clash 规则行与 #/! 注释）。
	if cat.LocalDomains != "" {
		p, err := r.localFilePath(cat.LocalDomains)
		if err != nil {
			return nil, err
		}
		rs, err := ParseDomainList(p)
		if err != nil {
			return nil, err
		}
		all = append(all, rs...)
	}

	// exclude_domains：仓内 curated 白名单（如 data/reject-allowlist.txt）。**减法**，
	// 必须在所有输入源并集之后执行——否则先并后加的源会把已剔除的域名带回来。
	// 复用 ParseDomainList 读白名单（同样容忍裸域 / clash 规则行 / #! 注释）。
	if cat.ExcludeDomains != "" {
		p, err := r.localFilePath(cat.ExcludeDomains)
		if err != nil {
			return nil, err
		}
		rs, err := ParseDomainList(p)
		if err != nil {
			return nil, err
		}
		allow := make([]string, 0, len(rs))
		for _, x := range rs {
			allow = append(allow, x.Value)
		}
		all = excludeDomainTree(all, allow)
	}

	return dedupSortRules(all), nil
}
