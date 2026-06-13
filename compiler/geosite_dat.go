package compiler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/river2heaven/tongtian/ruleset"
	"google.golang.org/protobuf/encoding/protowire"
)

// geosite.dat 是 xray/v2ray 原生 domain 数据库格式：**单个** protobuf 文件，内含全部
// 类目 tag，消费端用 `geosite:<tag>` 引用（tag 大小写不敏感，本编译器统一产大写）。
//
// protobuf schema（与 v2fly/domain-list-community generator 产物等价）：
//
//	message GeoSiteList { repeated GeoSite entry = 1; }
//	message GeoSite     { string country_code = 1; repeated Domain domain = 2; }
//	message Domain      { Type type = 1; string value = 2; }  // 余字段（attribute 等）本期不产
//	enum   Type         { Plain = 0; Regex = 1; Domain = 2; Full = 3; }
//
// 本文件直接拼 protobuf wire（protowire），零 codegen / 零外部 CLI，与现有纯 Go
// 文本/二进制渲染风格一致；wire 字段号与 v2ray schema 对齐，可被 xray/v2ray 直接加载。

// geosite Domain.Type 枚举值（与上方 schema 注释一致）。
const (
	geositeTypePlain  = 0 // 关键字子串匹配
	geositeTypeRegex  = 1 // 正则匹配
	geositeTypeDomain = 2 // root-domain / suffix 语义（含自身与子域）
	geositeTypeFull   = 3 // 精确域名
)

// geosite wire 字段号。
const (
	fieldEntry       = 1 // GeoSiteList.entry
	fieldCountryCode = 1 // GeoSite.country_code
	fieldDomain      = 2 // GeoSite.domain
	fieldDomainType  = 1 // Domain.type
	fieldDomainValue = 2 // Domain.value
)

// domainTypeFor 把统一模型 MatchType 映射成 geosite Domain.Type；IP 侧返回 ok=false（跳过）。
func domainTypeFor(m ruleset.MatchType) (uint64, bool) {
	switch m {
	case ruleset.MatchDomain:
		return geositeTypeFull, true
	case ruleset.MatchDomainSuffix:
		return geositeTypeDomain, true
	case ruleset.MatchDomainKeyword:
		return geositeTypePlain, true
	case ruleset.MatchDomainRegex:
		return geositeTypeRegex, true
	default: // IP-CIDR / IP-CIDR6 → 属 geoip.dat，本 target 只做 domain 侧
		return 0, false
	}
}

// GeositeCategory 是一个类目的 geosite.dat 输入：tag 名 + 该类目规则。
// 编译器遍历完所有类目后，把每个类目收集成切片，一次性产单文件 geosite.dat。
type GeositeCategory struct {
	Name  string
	Rules []ruleset.Rule
}

// GeositeDat 把多个类目渲染成单个 geosite.dat 的 protobuf 字节流（GeoSiteList）。
//
// 每个类目产一个 GeoSite entry，country_code = 类目名大写（xray tag 惯例不分大小写）；
// 仅 domain 侧规则（DOMAIN / -SUFFIX / -KEYWORD / -REGEX）进 entry，IP-CIDR 跳过。
// 即便某类目无 domain 规则（纯 IP / 空），仍产出一个 0 域名的 GeoSite entry，
// 保证 tag 存在（消费端 `geosite:<tag>` 引用不至于解析失败）。
func GeositeDat(cats []GeositeCategory) ([]byte, error) {
	var out []byte
	for _, cat := range cats {
		site := encodeSite(cat)
		out = protowire.AppendTag(out, fieldEntry, protowire.BytesType)
		out = protowire.AppendBytes(out, site)
	}
	return out, nil
}

// encodeSite 渲染一个 GeoSite message（country_code + repeated Domain）。
func encodeSite(cat GeositeCategory) []byte {
	var site []byte
	site = protowire.AppendTag(site, fieldCountryCode, protowire.BytesType)
	site = protowire.AppendString(site, strings.ToUpper(cat.Name))
	for _, r := range cat.Rules {
		typ, ok := domainTypeFor(r.Match)
		if !ok {
			continue
		}
		dom := encodeDomain(typ, r.Value)
		site = protowire.AppendTag(site, fieldDomain, protowire.BytesType)
		site = protowire.AppendBytes(site, dom)
	}
	return site
}

// encodeDomain 渲染一个 Domain message（type + value）。
func encodeDomain(typ uint64, value string) []byte {
	var dom []byte
	dom = protowire.AppendTag(dom, fieldDomainType, protowire.VarintType)
	dom = protowire.AppendVarint(dom, typ)
	dom = protowire.AppendTag(dom, fieldDomainValue, protowire.BytesType)
	dom = protowire.AppendString(dom, value)
	return dom
}

// WriteGeositeDat 把全类目渲染成单文件写到 outDir/geosite.dat。
// 与 per-category 的 WriteCategory 不同：geosite.dat 本就是按 tag 索引的单库，
// 故由编译器遍历完所有类目后一次性调用。
func WriteGeositeDat(outDir string, cats []GeositeCategory) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	b, err := GeositeDat(cats)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "geosite.dat"), b, 0o644)
}
