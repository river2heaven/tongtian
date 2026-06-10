// Command rulesetc 是 通天河 的自编译 CI 编译器入口。
//
// 拉好的上游（每个 git checkout 到 --upstreams-dir/<key>）由 CI
// (.github/workflows/ruleset-build.yml) 准备，本程序读 manifest → 校验钉版本 →
// 解析各上游 → 扇出各客户端 rule-set 到 --out。
//
// 用法：
//
//	rulesetc --manifest manifest.yaml \
//	         --upstreams-dir <work>            # 子目录名 = manifest 的 upstream key \
//	         --out dist --sing-box $(which sing-box) --tag ruleset-20260610
//
// 不带 --sing-box 也能跑：产出 .list + .singbox.json 文本产物（.srs 跳过）。
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/river2heaven/tongtian/compiler"
)

func main() {
	manifestPath := flag.String("manifest", "manifest.yaml", "manifest 路径")
	upstreamsDir := flag.String("upstreams-dir", "", "各上游 checkout 根目录（子目录名 = manifest 的 upstream key）")
	outDir := flag.String("out", "dist", "产物输出目录")
	singboxBin := flag.String("sing-box", "", "sing-box 可执行路径（产 .srs；空则跳过二进制）")
	mihomoBin := flag.String("mihomo", "", "mihomo 可执行路径（为纯 domain/纯 ip-cidr 类别产 .mrs；空则跳过）")
	buildTag := flag.String("tag", "", "构建 tag（写入 version.json，回退用）")
	check := flag.Bool("check", false, "只校验 manifest（钉版本 / 上游引用）然后退出")
	flag.Parse()

	m, err := compiler.LoadManifest(*manifestPath)
	if err != nil {
		log.Fatalf("manifest: %v", err)
	}
	if err := m.Validate(); err != nil {
		log.Fatalf("manifest 校验失败（钉版本 / 上游引用）: %v", err)
	}
	if *check {
		log.Printf("✓ manifest 校验通过（%d 上游 / %d 类别）", len(m.Upstreams), len(m.Categories))
		return
	}
	if *upstreamsDir == "" {
		log.Fatalf("必须提供 --upstreams-dir")
	}

	res := compiler.NewResolver(m, *upstreamsDir)
	tools := compiler.Tools{SingBox: *singboxBin, Mihomo: *mihomoBin}
	built := map[string]int{}

	for _, cat := range m.Categories {
		rules, err := res.ResolveCategory(cat)
		if err != nil {
			log.Fatalf("类别 %s: %v", cat.Name, err)
		}
		if err := compiler.WriteCategory(*outDir, cat.Name, rules, tools); err != nil {
			log.Fatalf("写 %s: %v", cat.Name, err)
		}
		built[cat.Name] = len(rules)
		log.Printf("✓ %s: %d 条规则", cat.Name, len(rules))
	}

	writeVersion(*outDir, *buildTag, m, built)
	log.Printf("✓ 编译完成 → %s（%d 类别 / %d 上游）", *outDir, len(m.Categories), len(m.Upstreams))
}

func writeVersion(outDir, tag string, m *compiler.Manifest, built map[string]int) {
	pins := map[string]string{}
	for name, u := range m.Upstreams {
		pins[name] = u.Commit
	}
	version := map[string]any{
		"tag":        tag,
		"upstreams":  pins,
		"categories": built,
	}
	vb, _ := json.MarshalIndent(version, "", "  ")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "version.json"), append(vb, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
}
