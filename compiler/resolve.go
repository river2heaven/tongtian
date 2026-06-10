package compiler

import (
	"fmt"
	"path/filepath"

	"github.com/river2heaven/tongtian/ruleset"
)

// Resolver 按 manifest + 各上游本地 checkout 目录解析类别成统一规则。
//
// 约定：每个上游 checkout 在 upstreamsDir 下的子目录里，子目录名 = manifest 的 upstream key；
// UpstreamRef.Files / DataDir 相对该子目录。geosite 上游 = DataDir 非空那个。
type Resolver struct {
	m            *Manifest
	upstreamsDir string
	geosite      *Geosite // lazy
}

// NewResolver 构造解析器。
func NewResolver(m *Manifest, upstreamsDir string) *Resolver {
	return &Resolver{m: m, upstreamsDir: upstreamsDir}
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

	return dedupSortRules(all), nil
}
