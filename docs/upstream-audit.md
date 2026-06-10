# 上游审计：我们调查了哪些规则库，结论是什么

> 这是 `通天河` 选定上游前的完整尽调记录。回答三个问题：**调查了哪些 repo、分别是什么、为什么选 / 为什么否决**。
>
> 方法：先列已知主流库，再用 9 个角度（clash / sing-box / surge 生态 · geosite 血统 · AI+流媒体专项 · anticensor · 聚合器 · 广告 reject · 中文社区）做多角度网络搜索扇出，去重后对 **80 个公开规则库**逐一**对抗式核实**（不轻信 README 自述，读 workflow / build 脚本 / fork graph / commit 历史判定血统与活跃度）。

---

## 1. 已知 / 排除集（调查前就在视野内，不计入 80 个新候选）

| 库 | 是什么 | 我们的处置 |
|----|--------|-----------|
| `v2fly/domain-list-community` | geosite 事实标准源，人工 PR 驱动 | ✅ **选为 geosite 骨架**（见 README §2） |
| `felixonmars/dnsmasq-china-list` | 11 万+ 条穷举 CN 域名 | ✅ **选为 CN 域名表** |
| `Loyalsoldier/v2ray-rules-dat` / `clash-rules` | = domain-list-community 自编译 + gfwlist + felixonmars + 广告，**聚合打包**，从不主动探测 | 不选：我们直接消费它的上游，砍掉中间商 |
| `blackmatrix7/ios_rule_script` | per-app 手工库，最全的中文专项 list | 不选：个人维护、**部分冻结**（Netflix/YouTube/GitHub/Cloudflare 停在 2025-06） |
| `DivineEngine/Profiles` (ConnersHua) | Surge 圈手工策展 | 不选：**已停更 / 原 repo 下架** |
| `Sukka / skk.moe` | reject / ai 质量高 | 备选：AGPL + 请求别直连，编译期择优并入 |
| `MetaCubeX/meta-rules-dat` | clash+singbox `.mrs/.srs` | 不选：= domain-list-community 编译产物，只覆盖 clash+singbox |
| SagerNet `sing-geosite` | = domain-list-community 编出的 `.srs` | 不选：同源下游 |

> **血统关键事实**：`Loyalsoldier` 的源头就是 `domain-list-community + felixonmars`；`DivineEngine` 是和 `blackmatrix7` 平行的另一条 Surge 手工血统、已死。所以「参考三巨头」里唯一活跃且公共 PR 驱动的，其源头就是我们选的 `domain-list-community`。

---

## 2. ✅ 值得纳入视野（80 个里的 23 个）

按角色分组。⭐ = 给上游选型带来真正新血统的源。

### A. 一手源 / 新血统（README Tier 1）

| 库 | 是什么 | License | 结论 |
|----|--------|---------|------|
| ⭐ [`gfwlist/gfwlist`](https://github.com/gfwlist/gfwlist) | anticensor 谱系**祖先**，AutoProxy 规则，人工核验 16 年 | LGPL-2.1 | **选**：唯一的「真上游原料」，绕过 Loyalsoldier 中间层 |
| ⭐ [`misakaio/chnroutes2`](https://github.com/misakaio/chnroutes2) | 一手 BGP route-collector 聚合 CN IP，逐小时刷新 | CC-BY-SA-4.0 | **选**：GEOIP-CN source-of-truth，补 IP 侧空白 |

### B. 广告 reject（README Tier 2，比 Sukka/Loyalsoldier 更强）

| 库 | 是什么 | License | 结论 |
|----|--------|---------|------|
| ⭐ [`privacy-protection-tools/anti-AD`](https://github.com/privacy-protection-tools/anti-AD) | CN 广告 reject 权威，聚合 AdGuard/EasyList/v2fly + 白名单防误杀 | MIT | **选**：reject 主力，只取裸域自编译 |
| ⭐ [`hagezi/dns-blocklists`](https://github.com/hagezi/dns-blocklists) | 业界最强 reject/威胁情报，light→ultimate 分级 | GPL-3.0 | **选**：分级最细，需自建转换 |
| [`Cats-Team/AdRules`](https://github.com/Cats-Team/AdRules) | 聚合 EasyList/AdGuard/uBO，DNS 工具扇出强 | 0BSD/混合 | 候选：产物 license 混合 copyleft |
| [`217heidai/adblockfilters`](https://github.com/217heidai/adblockfilters) | 18 广告源 + live-DNS 去死链，**15 格式全覆盖**（含 .mrs/.srs） | GPL-3.0 | 格式参照：已含编译态 |

### C. AI / LLM 长尾（README Tier 2，补 AI 覆盖短板）

| 库 | 是什么 | License | 结论 |
|----|--------|---------|------|
| ⭐ [`VPSDance/ai-proxy-rules`](https://github.com/VPSDance/ai-proxy-rules) | AI 厂商 taxonomy 最全（82 厂商 + ASN/CIDR，含 Claude MCP 域），8 客户端全格式 | MIT | **cross-check 源**：派生非一手，但 AI 垂直最强 |
| [`fmz200/wool_scripts`](https://github.com/fmz200/wool_scripts) | loon/surge/qx/egern，AI 列表含 Anthropic IP-CIDR/grok/pplx/copilot | GPL-3.0 | AI 候选源 + app 级 MITM 品类 |
| [`dler-io/Rules`](https://github.com/dler-io/Rules) | 机场官方手工，AI 极全（claude.com/grok/cursor/windsurf） | **无 license** | 仅**抽域名当情报**，不 vendor |
| [`Accademia/Additional_Rule_For_Clash`](https://github.com/Accademia/Additional_Rule_For_Clash) | dlc+bm7 再策展 + 手工补（AI/银行/网盘） | MIT | 专项采样：Gemini/Grok/AppleAI 补全 |

### D. 工程参照（README Tier 3，**不是数据源**）

| 库 | License | 借鉴点 |
|----|---------|--------|
| [`xkww3n/Rules`](https://github.com/xkww3n/Rules) | MIT | dlc 派生 + 全客户端扇出 MIT 管线（dedup/CIDR 合并） |
| [`DustinWin/domain-list-custom`](https://github.com/DustinWin/domain-list-custom) | MIT | Go trie 去重 + 多源合并编译器 |
| [`DustinWin/ruleset_geodata`](https://github.com/DustinWin/ruleset_geodata) | GPL-3.0 | mihomo `.mrs` + singbox `.srs` 二进制扇出蓝本（daily 真更新） |
| [`QuixoticHeart/rule-set`](https://github.com/QuixoticHeart/rule-set) | GPL-3.0 | 单源 → 全客户端 working blueprint |
| [`FuGfConfig`](https://github.com/Elysian-Realme/FuGfConfig) | MIT | 11 格式自动扇出 + AI 新且全 |
| [`HenryChiao/MIHOMO_YAMLS`](https://github.com/HenryChiao/MIHOMO_YAMLS) | AGPL-3.0 | 8 客户端扇出矩阵工程参照（AGPL 传染，仅看不抄） |

### E. 专项 / 区域（按需，默认低优）

| 库 | 是什么 | License | 结论 |
|----|--------|---------|------|
| [`TG-Twilight/AWAvenue-Ads-Rule`](https://github.com/TG-Twilight/AWAvenue-Ads-Rule) | 逆向安卓 ad-SDK，**安卓 app 内广告**独有盲区 | GPL-3.0 | 选（reject 补充）：bus-factor=1，需钉版本 |
| [`mnixry/direct-android-ruleset`](https://github.com/mnixry/direct-android-ruleset) | 爬应用商店榜单出 PROCESS-NAME 包名直连 | AGPL-3.0 | 旁路源：Android 维度独有 |
| [`Chocolate4U/Iran-v2ray-rules`](https://github.com/Chocolate4U/Iran-v2ray-rules) · [`Iran-sing-box-rules`](https://github.com/Chocolate4U/Iran-sing-box-rules) | 伊朗区域 + 独立 security/malware feeds | GPL-3.0 | 仅扩区域时启用；security feeds 可单看 |
| [`runetfreedom/russia-blocked-geosite`](https://github.com/runetfreedom/russia-blocked-geosite) · [`1andrevich/Re-filter-lists`](https://github.com/1andrevich/Re-filter-lists) | 俄罗斯 RKN 区域 | GPL-3.0 / MIT | 仅扩区域时启用 |
| [`ACL4SSR/ACL4SSR`](https://github.com/ACL4SSR/ACL4SSR) | Clash 生态事实标准底库，gfwlist 直采 + 手工，血统正交 | CC-BY-SA-4.0 | 可作交叉校验；AI 弱、copyleft |

---

## 3. ❌ 否决（80 个里的 57 个）— 按否决原因分组

> 不逐一展开，按「为什么不要」归类。要点：**绝大多数是同源下游 / 无 license / 已死 / 类目错配**。

| 否决原因 | 库（部分） | 说明 |
|----------|-----------|------|
| **同源下游**（dlc/bm7 二阶派生、单格式 sing-box `.srs` 转码镜像） | `lyc8503/sing-box-rules`、`Toperlock/sing-box-geosite`、`Yuu518/sing-box-rules`、`DDCHlsq/sing-ruleset`、`malikshi/sing-box-geo`、`xmdhs/sing-box-ruleset`、`tangnahuaite/...`、`senshinya/singbox_ruleset`、`KaringX/karing-ruleset`、`Dreista/sing-box-rule-set-cn`、`yuumimi/geosite`、`CloudPassenger/geosite` | 转码这步正是我们要自己做的；上游全在视野内 |
| **无 license 下游壳**（供应链审计红线，一票否决） | `Repcz/Tool`、`jnlaoshu/MySelf`、`cutethotw/ClashRule`、`RealSeek/Clash_Rule_DIY`、`Keviin560/Shunt_Rules`、`xndeye/rule-merger`、`cmontage/proxyrules-cm`、`szkane/ClashRuleSet`、`Z-Siqi/Clash-for-Windows_Rule`、`luestr/ShuntRules`、`qRuWGQ/rules` | 无授权 = 不能合规 vendor |
| **已死 / 停更** | `Hackl0us/SS-Rule-Snippet`（2021）、`kyleduo/Surge-Rule-Snippets`（2017）、`Ekko1048/OpenClashRule`、`Dunamis4tw/generate-geoip-geosite`、`runetfreedom/geodat2srs`、`Dracay/ruleset`、`carsondzh/clash-geosite`、`lateautumn2/ruleset_geodata`、`Aethersailor/Custom_Clash_Rules`、`savely-krasovsky/antizapret-sing-box`、`ClashConnectRules`、`zqzess/rule_for_quantumultX` | 含一例反面教材：`Phoroc/sing-rules` CI 每日跑但源码 2 年没动、上游 404 静默产空集 |
| **类目错配**（纯广告 hosts，非分流规则） | `8680/GOODBYEADS`、`uniartisan(zhiyuan1i)/adblock_list`、`neodevpro/neodevhost`、`lingeringsound/10007_auto`、`limbopro/Adblock4limbo`、`StevenBlack/hosts`、`REIJI007/AdBlock_Rule_For_Clash`+`_Sing-box` | anti-AD / hagezi / Sukka 已是更优上游 |
| **IP 库下游**（追上游即可） | `Hackl0us/GeoIP2-CN`、`soffchen/GeoIP2-CN`、`fernvenue/chn-cidr-list`、`ruijzhan/chnroute` | `misakaio/chnroutes2` 更上游 |
| **bm7/Loyalsoldier 纯下游 / 格式桥** | `GMOogway/shadowrocket-rules`、`NSZA156/surge-geox-rules`、`Aoang/Surge`、`sve1r/Rules-For-Quantumult-X`、`vernette/rulesets`、`Ckrvxr/mihomo_yaml`、`powerfullz/override-rules`、`Aethersailor/Custom_OpenClash_Rules`、`SunsetMkt/anti-ip-attribution` | 追本体即可；部分场景错配（社交 IP 归属地等） |

---

## 4. 按血统归类（一图看清生态）

```
                          ┌─ v2fly/domain-list-community (人工 PR, geosite 事实标准)  ← 我们选它
                          │     └─ Loyalsoldier 工具链 → DustinWin*, xkww3n, Chocolate4U,
   geosite 血统 ──────────┤        + 所有 sing-box .srs 转码镜像 (lyc8503/malikshi/Yuu518...)  ← 同源下游, 不选
                          │
   gfwlist 血统 ──────────┴─ gfwlist/gfwlist (祖先, 人工核验)  ← 我们也直接选它 (绕过中间层)

   blackmatrix7 血统 ───── 个人手工, 部分冻结 → DivineEngine(死) + 大量下游壳  ← 不选

   广告 reject 血统 ────── AdGuard/EasyList/uBO → anti-AD / hagezi / Cats / 217heidai  ← 选 anti-AD+hagezi

   GEOIP 血统 ──────────── BGP route-collector → misakaio/chnroutes2  ← 我们选它 (IP 一手源)
                                                  └─ soffchen/Hackl0us/fernvenue  ← 下游, 不选
```

**净结论**：80 个候选里真正新血统只有 4 类——① `gfwlist`（一手 GFW 源头）② `chnroutes2`（一手 GEOIP-CN）③ `anti-AD`/`hagezi`（更强 reject）④ `VPSDance`/`fmz200`/`dler`（AI 长尾）。其余 50+ 要么同源下游、要么无 license / 已死 / 单格式 / 类目错配。

---

## 5. 学术自动测量项目的真实定位（GFWatch / GFWeb / OONI）

经核实，**不适合做 pipeline 上游，仅作审计 / cross-check**。理由与详情见 [README §2「仅审计 / cross-check」](../README.md#2-选定的源头最终上游)。一句话：新鲜度多不可靠（论文型易冻结）、格式是测量数据非规则集（OONI 尤甚）、且对路由边际价值低还会抄进 GFW 自己的过封 bug。

---

## 附 A：调查方法的可复现性

本审计由多角度搜索 + 逐库对抗式核实产生（80 个候选 × 真实 GitHub 仓库核实：读 build 脚本判血统、查 commit/release 判活跃、区分「仓库 CI 活跃」与「规则真更新」）。每库核实 12 个字段：stars / last_activity / formats / provenance / license / ai_coverage / maintenance_status / worth_tracking / specialty / reason 等。完整原始结果见下方 **附 B**。

## 附 B：80 库原始核实数据（可折叠）

> 下面是 80 个候选库逐一对抗式核实的**原始结构化结果**，每库展开可见全部字段。来源标记 ✅ = 值得纳入视野，❌ = 否决。

### ✅ 值得纳入（23）

<details>
<summary>✅ <b>217heidai/adblockfilters</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/217heidai/adblockfilters
- **stars / 最近 / 维护**: ~7.1k (7083) · 2026-06 · active
- **license**: GPL-3.0
- **formats**: clash.list, mihomo.mrs, mihomo.yaml, singbox.json, singbox.srs, surge.list, shadowrocket.sgmodule, quantumultx.conf, loon.list, adguard, adguardhome, smartdns.conf, dnsmasq, mosdns, hosts
- **provenance**: aggregator-of-others — README + clash file headers explicitly enumerate ~18 ad-block upstreams (AdGuard Base/Chinese/Mobile/DNS, EasyList / EasyList China / EasyPrivacy, OISD Basic, StevenBlack + Pollock hosts, AdRules DNS, CJX Annoyance, AWAvenue, jiekouAD, DNS-Blocklists Light). NOT domain-list-community derived, NOT a proxy geosite. Pipeline = merge + dedup + live DNS resolution to prune dead domains (app/*.py converters per format).
- **ai_coverage**: 无有效 AI/LLM 路由覆盖。全库为拦截清单：扫描命中的 38 个含 openai/anthropic/claude/grok/gemini 关键词的条目全部是广告/诈骗/钓鱼域名被 REJECT (如 analogchatgpt.com、chatgpt-premium.com、gemini-ads-team.com、grokai-tech.org)，不含 claude.com/openai.com/x.ai 等真实分流域名。
- **specialty**: 纯去广告/隐私拦截层 (全部规则为 REJECT-DROP)。差异化价值 = 一个把 16+ 主流广告源合并、live-DNS 去死链、并扇出到我们全部 7 个客户端格式 (含 mihomo .mrs / sing-box .srs 已编译二进制) 的现成 REJECT 层；提供 full + lite(国内) 双版本。它填补的是 blackmatrix7/Loyalsoldier 之外更专更全的广告合并维度，而非路由维度。
- **reason**: 对抗式核实通过。(1) 真实存在且确为规则库，但定位是去广告拦截库而非代理分流/geosite 库。(2) 真活跃——非仅"仓库活"：GitHub API 显示 archived=false，最新 commit 2026-06-09(当天)；Update Filters workflow 每~8h 成功跑一次，规则文件头带 Last modified 2026/06/09 09:04:44 + version 时间戳，规则确在更新。(3) 格式经实测落盘确认 (clash/.list、mihomo .mrs+.yaml、singbox .json+.srs、surge、shadowrocket、quantumultx、loon、adguard(home)、smartdns、dnsmasq、mosdns、hosts，各 full+lite)，与我们 7 客户端高度吻合且已含编译态 .mrs/.srs。(4) 血统=聚合器，明确非 domain-list-community 派生。(5) GPL-3.0 copyleft——若我们重分发编译产物需注意许可证传染。差异化：它与 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 正交——那些是 REJECT+路由综合库，本库是更专更全的纯广告合并层，AI/流媒体路由维度为零。结论：值得纳入 本项目 视野，但仅作为"广告拦截 (REJECT) 上游候选"单独追踪，不作为路由规则库对标对象；采用前需评估 GPL-3.0 对我们自编译扇出产物的影响。

</details>

<details>
<summary>✅ <b>Accademia/Additional_Rule_For_Clash</b> — active · MIT</summary>

- **URL**: https://github.com/Accademia/Additional_Rule_For_Clash
- **stars / 最近 / 维护**: ~252 · 2026-05 · active
- **license**: MIT
- **formats**: clash.yaml (rule-provider, Domain/Classical/No_Resolve variants only)
- **provenance**: domain-list-community-derived (partial) + independent-curated (partial) — 混合血统。依据: README 自述 ChinaMax = blackmatrix7/ios_rule_script/ChinaMax 的「修正版」, GeositeCN = v2fly/geosite:cn 的「修正版」(v2fly geosite 即 domain-list-community 谱系), ChinaMax 数据又聚合 ACL4SSR/LM-Firefly/blackmatrix7。故大宗中国规则是对 v2fly geosite + blackmatrix7 的再策展/纠错层; 而 AI/银行/FakeLocation 等专项是作者手工补充(无披露来源, Grok 辅助编辑)。不是纯独立, 也不是纯派生。
- **ai_coverage**: 较强且新: 独立目录覆盖 Gemini、Grok(xAI, 含 IPv6 特殊处理说明)、Copilot(微软)、AppleAI(Apple Intelligence)。但无 Claude/claude.com、无 Perplexity、无 OpenAI/ChatGPT 专项目录(ChatGPT 推测仍依赖 blackmatrix7 主库)。即覆盖了主库滞后的几家新 LLM, 但不完整。
- **specialty**: 差异化在于 blackmatrix7 缺口补充 + 中国规则纠错: (a) 一批 blackmatrix7/Loyalsoldier/Sukka 都没有或不全的专项 — Gemini/Grok/Copilot/AppleAI(Apple Intelligence)/AppleNews、国际 Bank(10+ 国)/VirtualFinance(PayPal/Wise/Revolut)、FakeLocation(B站/抖音/小红书等国内 App 防 IP 探测)、Alipan/BaiduNetDisk/WeiYun 网盘、HomeIP/UnsupportVPN/HijackingPlus; (b) 声称对 v2fly geosite:cn(6800→4800 条纠错) 和 ChinaMax(称原版 2 万+ 错误)做了精修版。对我们而言, AI/银行/网盘/FakeLocation 这几类是 Loyalsoldier/Sukka/blackmatrix7 主线确实薄弱的领域, 有采样价值。
- **reason**: 值得加入 本项目 视野, 但定位为「专项采样源/校对参照」, 不是核心上游。理由: (1) 真实存在、确为 Clash 分流规则库(252★/MIT/未归档/created 2024-11, pushed 2026-05-14 约一月前), 非无关项目; (2) 维护日历上 active 但信号差 — 近 30 commit 里 23 个 message 仅为 "ok"、7 个改 README, 唯一可见实质提交是 2 月一个 GeoIP 文件, 无法验证规则数据真实更新频率, 判 active 但透明度低, 需人工抽查 diff; (3) 致命短板: 仅出 Clash YAML 单格式, README 明确拒绝 .mrs/sing-box/text("纯属自嗨"), 与我们「上游源→自编译→扇出 7 客户端」的目标正交 —— 对我们只能当域名清单原料, 不能当成品格式来源; (4) 血统非纯独立, 大宗规则是 v2fly geosite + blackmatrix7 的再策展, 与已知库高度重叠, 真正增量是 AI(Gemini/Grok/Copilot/AppleAI)/国际银行/网盘/FakeLocation 这几类 Loyalsoldier/blackmatrix7/Sukka 薄弱领域; (5) 缺 Claude/Perplexity/OpenAI 专项, AI 覆盖不完整。结论: 列为低优先级专项采样源, 抽取其 AI/银行/FakeLocation 域名清单做交叉补全, 但不依赖其格式产物, 也不作为权威上游。

</details>

<details>
<summary>✅ <b>ACL4SSR/ACL4SSR</b> — active · CC-BY-SA-4.0</summary>

- **URL**: https://github.com/ACL4SSR/ACL4SSR
- **stars / 最近 / 维护**: ~6.1k (6134, forks 1985) · 2026-06 · active
- **license**: CC-BY-SA-4.0
- **formats**: clash.list (classical .list), clash provider .yaml (rule-provider), mihomo .mrs (binary ruleset, domain+ip), ssr .acl, ssr .rule (gfwlist-user), clash .ini/.yml subscription config (pref.ini, GeneralClashConfig.yml)
- **provenance**: aggregator-of-others + independent-curated 混合 (NOT domain-list-community 派生). 依据 (对抗式核实, 非 README 自述): (1) GFW 黑名单核心由 scripts/gfwlist_parser.py 每日 fetch gfwlist/gfwlist 的 gfwlist.txt 解析生成 (ProxyGFWlist/fullgfwlist/UnBan) → 这部分是 gfwlist 聚合器, 不是 v2fly/domain-list-community. (2) 按服务 Ruleset (OpenAi/Netflix/Telegram 等 ~140 个 .yaml) 是 ACL4SSR 自有历史手工策展的 DOMAIN-SUFFIX/DOMAIN-KEYWORD 风格 (blackmatrix7 早期即从此血统分叉), 非 v2fly geosite 格式. (3) .mrs 由 CI 下载 MetaCubeX/mihomo 跑 `mihomo convert-ruleset` 从自家 .list/.yaml 机械编译, 非从 meta-rules-dat 拉. 即: 它本身是上游之一, 不派生 domain-list-community.
- **ai_coverage**: 有但偏浅. mrs/yaml 含 AI_domain, OpenAi (16 条: openai.com/chatgpt.com/sora.com + DOMAIN-KEYWORD,openai + auth0/arkoselabs 等支撑域, 质量尚可), Gemini, Claude/ClaudeAI (仅 2 条 anthropic.com + claude.ai, 缺 claude.com), 无 grok/perplexity 专项. 综合: AI 聚合分组存在, 但远不及专做 AI 分流的维护者细致, 不能作为 AI 域名权威源.
- **specialty**: 事实标准的 Clash/SSR 中国分流底库 (ChinaDomain/ChinaIp/ChinaMedia/BanAD/BanProgramAD/ProxyGFWlist), 是国内订阅生成器 (subconverter 等) 默认引用源, 生态绑定深. 相对我们已知库的差异化: vs Loyalsoldier (v2fly geosite 派生, 给 xray/sing-box .dat/.srs) — ACL4SSR 是 Clash 生态原生 + gfwlist 直采, 血统正交; vs blackmatrix7 (ACL4SSR 是其上游血统源, 但 blackmatrix7 维护更勤、服务粒度更细、且原生多格式); vs Sukka (Sukka 强工程化/去重/质量审计, ACL4SSR 偏经典手工列表无严格质检); vs MetaCubeX/meta-rules-dat (ACL4SSR 是 .mrs 的下游消费者而非 geo 数据源). 独特价值主要在"中国 GFW 经典底库 + 与国内订阅器的事实标准兼容", 不在质量或 AI 专项.
- **reason**: 真实存在且确为代理分流规则库 (CC-BY-SA-4.0, 6.1k★), 非无关项目, 未 archived. 维护活跃且"规则真在更新": 默认分支最新 push 2026-06-06, 含 GitHub Actions `update.yml` 每日 cron 自动解析 gfwlist 并重生成 .list/.yaml/.mrs ([AutoUpdate] bot commit 由 BROBIRD 推送, acl4ssr-sub/BROBIRD 为主力 author) — 即 active 且规则数据持续刷新, 不是只有仓库门面活跃. 无 GitHub Release (latest 404), 分发靠 raw 文件直链. 值得加入 本项目 视野的原因: 它是 Clash/SSR 中国分流的事实标准底库, 是国内订阅生成器默认引用源, 供应链审计若忽略它会漏掉整条 Clash 生态最常被引用的上游; 同时它的血统 (gfwlist 直采 + 自有手工策展, 非 domain-list-community) 与我们已知的 Loyalsoldier/Sukka 正交, 提供独立交叉校验视角. 但需标注风险: (a) 无严格质检/去重流程, 经典手工列表易陈旧; (b) AI/流媒体专项覆盖浅 (Claude 仅 2 域), 不宜作为 AI 域名权威源; (c) CC-BY-SA-4.0 为 copyleft, 若我们自编译扇出需保留署名+同协议传染, 比 MIT/CC0 系 (Loyalsoldier) 约束更强 — 这点对 本项目 自掌控供应链有 license 影响, 应在 spec 里显式记一笔.

</details>

<details>
<summary>✅ <b>Cats-Team/AdRules</b> — active · script branch (build tooling) = 0BSD; generated rules on main branch inherit upstream licenses (GPL-3.0 / MIT / CC BY-SA mix)</summary>

- **URL**: https://github.com/Cats-Team/AdRules
- **stars / 最近 / 维护**: ~3.6k (3562) · 2026-06 · active
- **license**: script branch (build tooling) = 0BSD; generated rules on main branch inherit upstream licenses (GPL-3.0 / MIT / CC BY-SA mix)
- **formats**: clash domainset (adrules_domainset.txt), mihomo .mrs (adrules-mihomo.mrs), singbox .srs + .json, surge .conf + surge domainset, loon rule-set, quantumultx (qx.conf), clash/surge .list, smartdns (smart-dns.conf), mosdns (mosdns_adrules.txt), AdBlock/ABP syntax (adblock.txt + lite/plus), hosts / dnsmasq (dns.txt, domain.txt)
- **provenance**: aggregator-of-others — Source.md explicitly enumerates upstream blocklists it merges (EasyList, AdGuard, uBlock Origin, Fanboy, anti-AD-style CN lists by xinggsf/damengzhu/cjx82630, StevenBlack/AdAway/someonewhocares hosts, urlhaus). NOT domain-list-community/geosite-derived, NOT hand-curated original. Pure aggregator of third-party ad/tracker blocklists.
- **ai_coverage**: none — reject-only ad/tracker list, no AI/LLM routing domains (claude.com/grok/perplexity); any AI domain would appear only incidentally as a tracker, not as a routing category. No streaming-media routing category either.
- **specialty**: CN-region ad/tracker/malware/HTTPDNS/PCDN REJECT list with unusually broad DNS-tool fan-out (smartdns + mosdns + dnsmasq alongside the usual clash/singbox/surge/qx). Occupies the ad-reject slot orthogonal to our known routing/geosite libs (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX). Closest overlap is Sukka reject + blackmatrix7 BanAD, but with stronger CN-specific ad coverage and three tiers (lite/normal/plus).
- **reason**: Verified real and genuinely active (NOT just cosmetic): main branch is regenerated DAILY by GitHub Actions — last build 2026-06-09 (today), script-branch human commits within last week, 3.6k stars, not archived. Adversarial check overturns one premise: this is NOT a proxy分流/geosite library — it is a pure CN ad-REJECT / DNS-block list (description: 'List for blocking ads in the Chinese region'). Provenance confirmed as aggregator (Source.md lists EasyList/AdGuard/uBO/anti-AD upstreams), so it is NOT differentiated from us on routing intelligence — but it IS differentiated as a maintained, multi-format, CN-focused REJECT source with strong DNS-tool fan-out (smartdns/mosdns). Worth tracking for 本项目 as a candidate REJECT/ad upstream (fills the ad-block slot alongside Sukka/blackmatrix7-BanAD), but explicitly OUT of scope for routing/geosite/AI fan-out. License is favorable for the tooling (0BSD) but the generated rules carry mixed copyleft (GPL-3.0/CC BY-SA) inherited from upstreams — a redistribution-license caveat to note if we self-compile and re-host.

</details>

<details>
<summary>✅ <b>Chocolate4U/Iran-sing-box-rules</b> — active · GPL-3.0 (上游数据混合 license: MaxMind/IP2Location CC BY-SA 4.0, PersianBlocker AGPL-3.0, urlhaus/malware-filter 等)</summary>

- **URL**: https://github.com/Chocolate4U/Iran-sing-box-rules
- **stars / 最近 / 维护**: ~315 · 2026-06 · active
- **license**: GPL-3.0 (上游数据混合 license: MaxMind/IP2Location CC BY-SA 4.0, PersianBlocker AGPL-3.0, urlhaus/malware-filter 等)
- **formats**: singbox.srs, singbox geosite.db/geoip.db (legacy sing-box v1.7.x)
- **provenance**: domain-list-community-derived — geosite 由姊妹库 Iran-v2ray-rules 跑 v2fly/domain-list-community 生成器 + 自家追加的伊朗/波斯/安全数据产出 geosite.dat, 本库再用其 fork 的 sing-geosite 把 geosite.dat 转成 .srs/.db。geoip 来自 MaxMind GeoLite2 + IP2Location LITE 自行编译。属于「domain-list-community 派生 + 重度伊朗区域独立策展」, 不是单纯转封 Loyalsoldier 成品。
- **ai_coverage**: 仅 OpenAI/ChatGPT 专项 (geosite:openai + geoip:openai)。无 claude.com / grok / perplexity / gemini 等独立类目, 无统一 category-ai 聚合。AI 覆盖窄于我们对 本项目 的目标 (多 LLM 域名扇出)。
- **specialty**: 伊朗/波斯区域专项是核心差异化 (geosite:ir 含非 .ir TLD 的伊朗域名 + 全 .ir bypass, geoip:ir 含 Eitaa/Rubika messenger + ArvanCloud/Derak/IranServer/ParsPack 本土 CDN; PersianBlocker/Iran Hosted Domains 波斯广告拦截)。安全集独立: malware/phishing/cryptominers (urlhaus + malware-filter), geoip 与 geosite 双侧都有。有 lite 版本 (24h 校验剔除死域)。AI: 仅 openai 类 (geosite:openai + geoip:openai, ChatGPT/OpenAI), 无更广的 ai-!cn/LLM 聚合类。Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 都没有伊朗本土 CDN + 波斯广告这层覆盖。
- **reason**: 核实通过且对抗式结论可靠: (1) 真实存在、确为代理分流规则库 (sing-box rule-set), 非无关项目, 有 issue/release 活动佐证。(2) 维护 active 且「规则真的在更新」—— release 用 GitHub Actions cron 每日构建, 最新 tag 202606080914 = 2026-06-08, 距今(2026-06-09)仅 1 天, date-stamp 版本号 + lite 集 24h 死域校验, 不是「仓库挂着但规则不动」的假活跃。(3) 产出仅 sing-box 格式 (.srs + 旧版 .db), clash/surge/v2ray 由姊妹库 Iran-clash-rules / Iran-v2ray-rules 各自产出 —— 对我们「一源扇出多客户端」反而是反面教材 (它是按客户端拆库而非统一源)。(4) 血统 = domain-list-community 派生 + 重度伊朗独立策展, 依据明确 (geosite.dat 由 v2fly 生成器在姊妹 v2ray 库构建)。值得加入 本项目 视野的理由: 它是「区域 + 安全」专项独立策展的优质样本, 伊朗本土 CDN/messenger geoip + PersianBlocker 波斯广告是 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 都缺的覆盖, 可作我们安全集 (urlhaus/malware-filter 同源) 与区域专项的参考输入。但 AI 覆盖仅 OpenAI、且为 sing-box-only 按客户端拆库, 不能直接当我们「统一上游源」, 定位为「专项数据源 + 架构反例参考」而非主干上游。

</details>

<details>
<summary>✅ <b>Chocolate4U/Iran-v2ray-rules</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/Chocolate4U/Iran-v2ray-rules
- **stars / 最近 / 维护**: ~668 (v2ray repo); 家族合计约 1k+ (sing-box ~315, clash ~77) · 2026-06 (release artifacts daily 重建) · active
- **license**: GPL-3.0
- **formats**: geoip.dat, geosite.dat, geoip-lite.dat, geosite-lite.dat, security.dat, security-ip.dat, .mmdb (Country/Security-IP/Services), sing-box .srs (Iran-sing-box-rules), sing-box .db/.metadb (deprecated), clash .txt / .yaml (Iran-clash-rules), clash.meta .metadb geoip
- **provenance**: domain-list-community-derived — README 明确: geosite 用 Domain-list-community 源码生成、geoip 用 Loyalsoldier geoip 源码生成; geosite 含 "all categories from domain-list-community" 透传。但绝非纯派生: 叠加独立手工策展的伊朗本地化层 (geosite:ir 非-.ir TLD 伊朗域名 + .ir bypass、Iran ASN/CIDR、ArvanCloud/Derak/ParsPack 等本地 CDN) + 安全富化聚合 (PersianBlocker、abuse.ch URLhaus、Iran Hosted Domains、MaxMind GeoLite2、IP2Location LITE)。属于 "domain-list-community 派生骨架 + 独立区域/安全策展富化" 的混合体。
- **ai_coverage**: 弱。仅 geoip:openai / geosite:openai (OpenAI/ChatGPT)。无 Anthropic/Claude、无 Google AI/Gemini、无 Perplexity/Grok 专项分类。AI 分流不是其设计目标。
- **specialty**: 区域 (伊朗) + 安全富化专项, 与我们已知库正交: (1) 伊朗本地化 (geosite:ir、.ir bypass、Iran ASN、ArvanCloud/Derak/IranServer/ParsPack 本地 CDN CIDR) — Loyalsoldier/blackmatrix7/Sukka 都不覆盖; (2) 安全分流 (malware/phishing/cryptominers/ads，聚合 abuse.ch URLhaus + PersianBlocker)，且把安全做成独立 security.dat/security-ip.dat/.mmdb 资产，比 Sukka 的反广告更偏 malware/phishing IP 维度; (3) 同时产 .dat + .srs + clash 三栈，但缺 mihomo .mrs 二进制 (clash repo 仅 .txt/.yaml)。流媒体非其重点。
- **reason**: 对抗式核实结论: 真实存在、确为代理分流规则库 (v2ray/sing-box/clash 三 repo 家族), 非无关项目。维护状态需区分两层——manual 源码 commit 在 v2ray .dat repo 约 2025-08 后转稀疏, 但 release 由 github-actions bot 每日时间戳 tag (202606080912 等) 重建, 即"规则确实在每日更新"而非仓库假活, 判定 active 成立。血统已证: geosite 由 domain-list-community 源码生成 + Loyalsoldier geoip 源码, README 自述与 commit/release 证据一致, 非纯聚合器也非纯独立手工, 是"DLC 派生 + 独立伊朗/安全富化"混合。值得加入 本项目 视野的理由: (a) 提供我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 都没有的伊朗区域 + 本地 CDN ASN 数据和独立 security.dat/security-ip 安全资产, 是差异化扇入源; (b) GPL-3.0 允许我们自编译再扇出 (但 GPL 传染性需在 本项目 license 矩阵里标注, 区别于 Loyalsoldier 的更宽松条款); (c) 已天然多格式 (.dat/.srs/clash), 可作为"上游源→自编译→扇出"流程的对照实现。不值得高优追的点: AI/LLM 覆盖几乎为空 (仅 OpenAI), 流媒体非重点, 若我们目标偏 AI/流媒体分流则价值有限; 且 clash 侧不出 .mrs, 对 mihomo 二进制 rule-provider 场景需自转。建议: 列为 本项目 的"区域+安全"专项候选源 (中优), 重点白嫖其 Iran ASN/CDN + security feeds, 不依赖其做 AI 分流。

</details>

<details>
<summary>✅ <b>dler-io/Rules</b> — active · 无 (NO LICENSE 文件, 无 root README; 默认 all-rights-reserved, 法律上不可自由再分发/再编译)</summary>

- **URL**: https://github.com/dler-io/Rules
- **stars / 最近 / 维护**: ~1.4k (1430) · 2026-05 · active
- **license**: 无 (NO LICENSE 文件, 无 root README; 默认 all-rights-reserved, 法律上不可自由再分发/再编译)
- **formats**: clash.yaml (rule-provider + full config Head/Rule), surge.list/.conf (Surge 3 + Surge 4, Groups/MITM/Cert/Prototype templates), singbox.json (route config Head.conf/Rule.conf, versioned 1.11 & 1.12 — source JSON, NOT compiled .srs), quantumult.conf, quantumultx.conf (Head/Tail/Script), shadowrocket.conf, egern.yaml
- **provenance**: independent-curated (airport-official). 依据: (1) dler-io 是 Dler Cloud 机场官方 GitHub org (org 主页 + 社区多处 fork 命名 "Dler Cloud.yaml" 佐证); (2) 规则文件是手工策展的 Surge/Clash 语法 (DOMAIN-SUFFIX 列表, 按 app 分节带 # > 注释, 含 IP-CIDR/voice 段), 属 blackmatrix7 式人工策展血统, 非 Loyalsoldier 的 domain-list-community/geosite 范畴派生; (3) 仓库不是纯 rule-provider 库——它分发整套客户端配置 (policy groups / MITM / Cert / Script / DNS 模板), 是机场自用配置仓, 顺带可被 import。非 aggregator-of-others (没有大规模拉别人 ruleset 再编译的痕迹)。初判 'aggregator' 不准确。
- **ai_coverage**: 极全。claude.com/claude.ai/anthropic.com/claudeusercontent.com 全有; grok.com/x.ai 有; perplexity.ai 有; 另含 ChatGPT API+CDN+Voice IP 段、Gemini 全家桶、Cursor/Windsurf/Zed/Copilot/Groq/Cerebras/OpenRouter 等 30+。是本批候选里 AI 覆盖最强的之一。
- **specialty**: AI/LLM 覆盖是其最强差异化: 单个 "AI Suite" provider 覆盖 30+ 家 (Claude/anthropic, ChatGPT/OpenAI 含 API+CDN+Voice, Grok/x.ai, Perplexity, Gemini/Google AI Studio/DeepMind/NotebookLM, Copilot, Cursor, Windsurf, Zed, TRAE, JetBrains, Groq, Cerebras, OpenRouter, Meta AI, Sora, POE, Jasper, Dify 等), 比 Loyalsoldier/Sukka 的 AI 分类更细更全, 且持续更新 (最近 commit 即 'update openai network recommendations')。流媒体解锁也极全 (Netflix/Disney+/Max/Spotify/YouTube/BBC iPlayer/DAZN/Bahamut/Abema/Hulu JP 等数十家, 按区域细分)。但本质是机场自用配置, 域名口径偏 Dler 解锁视角而非通用 geosite。
- **reason**: 值得加入 本项目 视野, 但定位为「AI/流媒体域名情报源」而非「可直接消费的上游 ruleset」。理由: (1) 真实存在且 active——2022 创建, 2026-05-31 仍有 commit (openai 更新), 非 archived/stale, 规则真在更新而非仅仓库活跃; (2) AI/LLM 覆盖 (含 claude.com) 是它相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 的核心差异化, AI Suite 单文件 30+ 家、含 API/CDN/Voice 细分, 适合作为我们自编译 AI 分流的交叉校验/补全源。风险与限制: (a) 无 license = all-rights-reserved, 不能直接打包再分发, 只能当情报参考 (与我们 本项目 '自己掌控上游' 目标契合——抽域名不抽文件); (b) 不是 geosite 派生、不提供 .srs/.mrs/geosite.dat 等可机读编译产物, sing-box 也只给 route 模板而非 rule-set, 我们要的是它的域名清单不是它的扇出; (c) 域名口径绑定 Dler 机场解锁视角, 通用性弱于 domain-list-community。结论: 作为 AI+流媒体域名补全的次级源跟踪, 不作为主上游。

</details>

<details>
<summary>✅ <b>DustinWin/domain-list-custom</b> — active · MIT (本仓库); 下游 ruleset_geodata 为 GPL-3.0</summary>

- **URL**: https://github.com/DustinWin/domain-list-custom
- **stars / 最近 / 维护**: 28 (此 tooling repo); 下游分发库 DustinWin/ruleset_geodata 约 1264 · 2026-05 (真实人工 commit; pushed_at 2026-06-08 为 fork-sync/CI 产物。下游 ruleset_geodata 规则每日更新, 最新 2026-06-09) · active
- **license**: MIT (本仓库); 下游 ruleset_geodata 为 GPL-3.0
- **formats**: mihomo .list (本仓库直接产出), mihomo .mrs (经下游 ruleset_geodata), geosite.dat, geoip.dat, Country.mmdb/.metadb, singbox .srs, singbox source .json
- **provenance**: domain-list-community-derived — 本仓库是 Loyalsoldier/domain-list-custom 的 fork (sync.yml 自动同步), 其血统重构 v2fly/domain-list-community。README 逐类列出数据源: v2fly 骨架 + blackmatrix7/ios_rule_script + Loyalsoldier/clash-rules + felixonmars/dnsmasq-china-list + 自增补域名; ai 类 = v2fly category-ai-!cn + ACL4SSR/AI.list。属 v2fly 根 + 多源合并增补 + 去重, 非纯镜像也非独立手工策展
- **ai_coverage**: 有独立 ai 类 (下游产出 ai.list/ai.mrs/ai.srs)。来源 = v2fly/domain-list-community category-ai-!cn + ACL4SSR/AI.list 合并; 覆盖 ChatGPT/Claude/Gemini 等主流 AI 服务域名, 但本质是聚合上游 AI 列表, 非自家专项策展
- **specialty**: 差异化价值低: 它是「编译器/数据源中间层」而非分发终点, 只产出中间 .list DOMAIN 数据 (data/ 文件甚至不入库, CI 构建时拉取), 真正扇出客户端格式的是下游 DustinWin/ruleset_geodata。源层面它再合并的正是我们已知的 Loyalsoldier/blackmatrix7/v2fly 同批上游, 新增有限。真正价值在「自编译→扇出」的工程参考: Go 实现 trie 去重 + 把 v2fly/ACL4SSR/dnsmasq-china-list 合并后扇出 mihomo .mrs/.dat + sing-box .srs, 恰好对应 本项目 目标
- **reason**: 值得加入 本项目 视野, 但定位要纠偏: 真正该追踪的是下游分发库 DustinWin/ruleset_geodata (1264 star, GPL-3.0, 规则每日更新, 同时产出 mihomo .mrs/.dat + sing-box .srs/.json + mmdb), domain-list-custom 只是其上游 Go 编译器。对抗式核实结论: 仓库真实存在且确为代理分流规则工程 (含 main.go/trie.go 编译器 + sync 工作流); pushed_at 的「活跃」是 fork-sync CI 假象, 但其驱动的下游管线确为每日真更新, 故整体判 active。血统为 v2fly/domain-list-community 派生 + 多源合并 (Loyalsoldier/blackmatrix7/dnsmasq-china-list/ACL4SSR) + Go trie 去重, 非独立策展也非纯聚合。源数据相对我们已知库新意有限 (再合并同批上游), 但其 Go 去重编译 + 单源扇出 mihomo+sing-box 双内核多格式的工程模式, 正是 本项目 自掌控供应链的高价值参考实现, 建议作为「自编译 fan-out」蓝本追踪 ruleset_geodata 而非本仓库

</details>

<details>
<summary>✅ <b>DustinWin/ruleset_geodata</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/DustinWin/ruleset_geodata
- **stars / 最近 / 维护**: ~1.3k (1264) · 2026-06 · active
- **license**: GPL-3.0
- **formats**: clash/mihomo .list, clash/mihomo .mrs, singbox .srs, singbox .json, geosite .dat, geoip .dat, MaxMind .mmdb, mihomo .metadb, singbox geo .db
- **provenance**: domain-list-community-derived (+ selective aggregator). Adversarially verified via fork graph, not README self-claim: README says data sourced from DustinWin/domain-list-custom + DustinWin/geoip. gh API confirms DustinWin/domain-list-custom is a FORK of Loyalsoldier/domain-list-custom (which compiles from v2fly/domain-list-community), and DustinWin/geoip is a FORK of Loyalsoldier/geoip. So the DOMAIN/IP backbone = domain-list-community via Loyalsoldier tooling + DustinWin custom additions; media/games/ads/AI categories additionally aggregate blackmatrix7/ios_rule_script, privacy-protection-tools/anti-AD, gfwlist, dnsmasq-china-list, ACL4SSR. NOT a pure independent hand-curated list, NOT a pure third-party aggregator.
- **ai_coverage**: Dedicated single 'ai' category present in all 4 formats (ai.list/ai.mrs/ai.srs/ai.json), sourced from v2fly category-ai-!cn + ACL4SSR/AI.list, covering non-CN AI services. Coarse: ~3KB list, single bucket, NO per-vendor split (no separate claude.com / openai / grok / perplexity rulesets — would need to inherit upstream's domain set, no vendor-level control). Daily-rebuilt (asset updated 2026-06-08).
- **specialty**: One-stop pre-compiled fan-out for mihomo + sing-box specifically: ships .mrs (mihomo binary ruleset) AND .srs (sing-box binary) AND human-readable .list/.json AND full geodata (.dat/.mmdb/.metadb/.db) from a single daily pipeline. Offers full/lite/mini/compatible variants and ready jsdelivr CDN URLs. Differentiation vs our known libs: more mihomo-native (.mrs/.metadb) + sing-box-native (.srs) binary-format coverage out of the box than Loyalsoldier (geoip/v2ray-rules-dat skews .dat/srs/mrs but less curated category granularity) and pre-built unlike blackmatrix7 (which is mostly source .list needing self-compile to .srs/.mrs). Strong streaming-unlock granularity (netflix/disney/max/primevideo/appletv/spotify/tiktok/youtube/bilibili each as own ruleset + matching mediaip). NOT a Surge/QuantumultX/Egern/Shadowrocket source — only mihomo+sing-box formats.
- **reason**: Verified real, relevant, and genuinely active — and critically the RULES are actually updating, not just the repo: releases use rolling tags (published 2025-12-01) but all assets were re-uploaded 2026-06-08T21:xx (the day before today 2026-06-09), repo pushed_at 2026-06-08, and both upstream data-source forks (domain-list-custom/geoip) also pushed 2026-06-08 — confirming the claimed daily 3AM-CST rebuild is real. Provenance traced through the actual fork graph (not trusting README): backbone is domain-list-community via Loyalsoldier-tooling forks + custom additions + blackmatrix7/anti-AD/ACL4SSR overlays. Worth adding to 本项目 vision as a REFERENCE/CROSS-CHECK target rather than a primary upstream: (a) it overlaps heavily in lineage with libs we already track (domain-list-community + Loyalsoldier), so it is partly redundant as a source; (b) but its differentiated value is the clean, simultaneously-built mihomo .mrs + sing-box .srs binary fan-out with full/lite/mini/compatible variants and fine-grained streaming buckets — a strong blueprint/benchmark for our own self-compile→fan-out pipeline, and a useful diff target to validate our build output against. Caveats for 本项目: (1) AI coverage is coarse single-bucket, no per-vendor granularity — does NOT satisfy a claude.com/grok/perplexity-level requirement on its own; (2) covers only mihomo + sing-box, NOT Surge/QuantumultX/Egern/Shadowrocket — so it does not solve our multi-client (surge/shadowrocket/egern/quantumult-x) fan-out problem; (3) GPL-3.0 — copyleft, relevant if we redistribute derived data. Net: track it as a derived/aggregator benchmark and binary-format reference, not as a net-new independent source.

</details>

<details>
<summary>✅ <b>fmz200/wool_scripts</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/fmz200/wool_scripts
- **stars / 最近 / 维护**: ~5.2k · 2026-06 · active
- **license**: GPL-3.0
- **formats**: loon.list, loon.plugin/.lpx, surge.sgmodule, quantumultx.snippet, egern.sgmodule, shadowrocket, stash
- **provenance**: aggregator-of-others (leaning independent-curated): 无任何 domain-list-community / Loyalsoldier / blackmatrix7 / geosite 引用或派生痕迹。作者 (奶思/fmz200) 手工策展 app/小程序 MITM reject 规则为主, AI.list 署名 lodepuly+fmz200; 同时聚合社区去广告/脚本作者产物。血统与 geosite 路由世界完全独立, 属 reject/去广告 + iOS 客户端工具生态。依据: README-RULE 无致谢/来源致敬段, 规则文件 header 只写 fmz200 自身 homepage/raw-url/tg-channel。
- **ai_coverage**: 强且差异化: AI.list (#!date=2026-03-31) 覆盖 OpenAI/ChatGPT/Sora、Anthropic Claude (含 IP-CIDR/IP-CIDR6)、Google Gemini/Bard/MakerSuite、Perplexity、xAI/Grok、GitHub Copilot、Meta AI、Mistral、JetBrains AI、Apple Intelligence、Trae。规则类型用到 DOMAIN-SUFFIX/DOMAIN/DOMAIN-KEYWORD/AND/IP-CIDR(6)。与 Sukka AI 列表同档。
- **specialty**: ~538 个 app/小程序去广告 reject 专项 (REJECT/REJECT-200/REJECT-img/REJECT-dict/REJECT-array 细粒度策略) + Apple 系统更新拦截 + 维护中的 AI 分流合集。这是我们已知库 (Loyalsoldier geosite 路由 / blackmatrix7 per-app 路由+rewrite / Sukka 工程化去重) 都不深做的「应用级 MITM 去广告」品类。
- **reason**: 对抗式核实结论: 真实存在且确为代理规则库 (5.2k star, GPL-3.0, 真 proxy 语法, 非营销/无关项目)。维护真活: repo 最新 commit 2026-06-08, 且实测规则 payload header 也新 (rejectAd #!date=2026-04-26, AI #!date=2026-03-31) — 区分了"仓库活跃"vs"规则真更新", 两者都成立; 唯一保留意见是近期 commit 流偏 JS 脚本 (青龙签到/监控/小红书), 规则编辑频次低于脚本。差异化价值清晰: ~538 app/小程序去广告 reject 专项 + 维护中 AI 列表, 是我们 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 都不覆盖的品类。值得加入 本项目 视野, 但定位为「上游源 → 自编译 → 扇出」的去广告/AI 候选源, 而非 drop-in provider: 它只产 Loon/Surge/QX/Egern 格式 (.list/.sgmodule/.snippet), 无 Clash .mrs / sing-box .srs / geosite.dat 原生产物, 接入我们多客户端扇出必须自行转换。

</details>

<details>
<summary>✅ <b>FuGfConfig (Elysian-Realme)</b> — active · MIT</summary>

- **URL**: https://github.com/Elysian-Realme/FuGfConfig
- **stars / 最近 / 维护**: ~253 · 2026-03 · active
- **license**: MIT
- **formats**: surge.list/.conf, singbox.srs, singbox.list, loon (.list/.conf/.plugin), quantumultx, clash DomainSet (.list), shadowrocket, adguardhome (.conf), host, littlesnitch, ublacklist
- **provenance**: aggregator-of-others (with automated multi-format compiler). 依据: README 显式致谢 SukkaW/Surge、blackmatrix7/ios_rule_script、VirgilClyne/GetSomeFries、Hackl0us/SS-Rule-Snippet、ACL4SSR、gfwlist、NextDNS、neohosts 等十余个上游；并非 domain-list-community 派生 (未致谢 v2fly/dlc，也无 geosite category 结构)。最近月度 commit 由 bot 账号 muthur6000bot 提交，单 commit 同步改 24 个文件跨 DataFile/DomainSet/Loon/Surge/sing-box(.srs) 多格式 — 证明有自动化抓取上游+扇出编译流水线，README 自述的"手工策展"被推翻。
- **ai_coverage**: 有专门 AI 分类 (ConfigFile/{DomainSet,Loon,sing-box}/AI/domain.list)，覆盖新且全: claude.com + claude.ai + anthropic.com + claudeusercontent.com、grok.com + x.ai、perplexity.ai、openai.com + chatgpt.com + oaistatic/oaiusercontent、openrouter.ai、poe.com/poecdn、groq.com、together.xyz、githubcopilot.com、generativelanguage.googleapis.com。2026-03 commit 还在加 notebooklm.google / hf.co，说明 AI 域名在持续跟进。
- **specialty**: iOS 客户端导向的多格式分流集 (Loon/QX/Surge/sing-box/Shadowrocket)，强项是 Apple 服务精细切分 (国区/外区 CDN、AppleAPI、AppleUpdate、no-cn-cdn) + 国产流氓软件/HTTPDNS/抖音/Bilibili 治理插件 (FuckRogueSoftware/FuckHTTPDNS/FuckDouyin/DNSMap)。这块"中国本地化反流氓+Apple 分流"维度比 Loyalsoldier/Sukka 更细。原生产出 sing-box .srs 二进制 rule-set，可直接 fan-out。
- **reason**: 真实存在、确为代理分流规则库 (created 2021-06，非归档，MIT)。对抗式核实结论: (1) active — 最近 commit 2026-03-24，且关键点是经验证 commit 真正改动规则数据 (新增 .cloudflare.dev/.notebooklm.google/.hf.co 跨 24 文件)，是"规则真在更新"而非仅仓库 churn；节奏约每月一次 bot 自动更新。(2) 血统为聚合器+自动扇出编译，非 dlc 派生、非独立手工——README"手工策展"不可信，实为多上游抓取后编译成 11 种客户端格式 (含 sing-box .srs / Surge / Loon / QX / Clash DomainSet)。差异化价值: 它把"上游源→自编译→扇出全客户端"这套流程做成了现成范本，正是 本项目 想自建的模式，可作架构参考/对照；同时 AI 分类覆盖 Claude+Grok+Perplexity 且持续更新，Apple 国区分流 + 反流氓软件本地化插件是相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 的补充维度。风险/注意: 自动产物 (.srs) 来自不可控的聚合上游，若纳入需把它当"参考源/对照集"而非可信权威，自己重编译时仍应回溯到一线上游 (Sukka/blackmatrix7) 而非二手聚合。值得加入 本项目 视野。

</details>

<details>
<summary>✅ <b>gfwlist/gfwlist</b> — active · LGPL-2.1</summary>

- **URL**: https://github.com/gfwlist/gfwlist
- **stars / 最近 / 维护**: ~25.4k · 2026-06 · active
- **license**: LGPL-2.1
- **formats**: AutoProxy .txt (base64-encoded gfwlist.txt), list.txt (raw AutoProxy 0.2.9 syntax, plaintext)
- **provenance**: independent-curated (源头, 非派生). 依据: 创建于 2015-03 (早于 v2fly domain-list-community); 文件是人工维护的 AutoProxy 黑名单, 带分节注释 (403/451/503 redirects, ehentai 等手工标注), 逐域名追加 (commit msg 形如 "Add dowjones.io"/"Add mozilla.ai"). 它是上游本体: Loyalsoldier 的 gfw.dat / category-gfw 反而是从本文件抽取生成的, 因此本库是祖先而非 domain-list-community 派生, 也不是聚合器.
- **ai_coverage**: 覆盖良好且新鲜. list.txt 中确认含: anthropic.com, claude.ai, openai.com, chatgpt.com, x.ai, grok.com, grokipedia.com, perplexity.ai, huggingface.co, character.ai, poe.com, githubcopilot.com, copilot.microsoft.com, simianx.ai. 但注意它是"被封锁判定"语义 (该不该走代理), 非"AI 服务商精确归类"语义, 颗粒度是单域名级而非服务分组.
- **specialty**: 它是几乎所有 anticensor 派生表的最上游单一真相源 (GFWList 本体), 是"被 GFW 封锁域名"这一语义的权威清单, 约 4200+ 条规则. 与我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 的根本差异: 那四个都是"消费者/再加工者", 输出多客户端格式 (clash/.mrs/.srs/surge); 本库是"生产者", 只产 AutoProxy base64 一种格式, 不做扇出. 对 本项目 的价值正在于此: 若想自己掌控"上游源→自编译→扇出", GFWList 就是要直接消费的源头之一, 绕过 Loyalsoldier 的中间加工 (可消除其编译期决策/延迟/停更风险).
- **reason**: 值得加入 本项目 视野, 但定位是"上游源 (raw source)"而非"成品规则库"。对抗式核实结论: (1) 真实存在, 确为代理分流黑名单 (AutoProxy GFWList 本体), 非无关项目; (2) active — 关键区分: 不是仓库空转, 而是规则文件 gfwlist.txt 本身在更新, 最近一次 2026-06-05, 5-6 月维持每周多次提交, list.txt→gfwlist.txt 编译流可见; (3) 仅产 AutoProxy base64 .txt 一种格式, 无 clash/singbox/surge 扇出; (4) 独立手工策展, 是 domain-list-community/Loyalsoldier 的祖先而非派生, 非聚合器; (5) LGPL-2.1, AI 域名覆盖新鲜 (含 anthropic/claude.ai/openai/grok/perplexity). 差异化价值: 它是供应链的真正起点, 若 本项目 目标是"自己掌控上游→自编译→扇出", 直接消费 GFWList 可绕开 Loyalsoldier 中间层的停更/决策风险。注意三点局限: 仅黑名单 (无 direct/proxy 之外的策略分组), AutoProxy 语法需自写 parser 转中性域名集, 语义是"是否被墙"而非精细服务归类——把它当"原料"而非"半成品"来纳入即可。

</details>

<details>
<summary>✅ <b>hagezi/dns-blocklists</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/hagezi/dns-blocklists
- **stars / 最近 / 维护**: ~23.5k · 2026-06 · active
- **license**: GPL-3.0
- **formats**: adblock, adguard, controld, dnsmasq, domains, hosts, wildcard, ips, pac, rpz, singbox.srs (community-converted only, NOT official)
- **provenance**: aggregator-of-others。依据：sources.md 聚合 AdGuard/OISD/EasyList/uBO/OpenPhish/URLhaus/ThreatFox/Cisco Talos 等广告+追踪+恶意+钓鱼 feed，再做独立去重/策展。明确不是 domain-list-community/v2fly 派生——sources.md 中零 geosite/v2fly 来源，仓库也无 geosite 维度。属于「聚合 + 二次手工策展」混合体，但本质是 reject 域名聚合器。
- **ai_coverage**: 无 AI/LLM 专项分类。这是 reject 黑名单库，不做「AI 域名 → proxy」分流；OpenAI/Claude/Grok/Perplexity 等不在覆盖目标内，AI 相关域名只会作为追踪/widget(如 landbot)被偶然误杀，社区靠 allowlist(如 issue #4354)放行。对我们「AI 分流」诉求无直接价值。
- **specialty**: 业界最强的 reject / 广告 / 恶意软件 / 威胁情报维度规则源：分级 Light/Normal/Pro/Pro++/Ultimate(95k→592k 条) + TIF 三档威胁情报 + Fake/Pop-Up/NRD-DGA 等专项。深度远超 Loyalsoldier/blackmatrix7/MetaCubeX 内置的那点 reject 列表；与 Sukka 重叠最多(Sukka 的 reject 部分本就 HaGeZi 派生)。它补的是「黑名单/拦截」这一维，不补「按地区分流(direct/proxy/geosite)」这一维——和我们已知库正交互补，不是替代。
- **reason**: 对抗式核实结论：真实存在且确为规则库，但要纠正初判定位——它是纯 **DNS reject 黑名单库**，不是代理 geo-分流/geosite 库。证据：(1) gh api 仓库 git tree 顶层只有 adblock/adguard/controld/dnsmasq/domains/hosts/wildcard/ips/pac/rpz，零 clash/sing-box/surge/geosite/.srs 目录；(2) sources.md 聚合的全是广告+恶意+钓鱼 feed，无任何 v2fly/geosite 来源。活跃度：pushed_at=2026-06-08，每天多次 "release" commit=自动重生成，规则是真在更新(active，非仅仓库活跃)；非 archived，GPL-3.0，23.5k★。格式纠偏：原生不产 .srs/clash/surge——issue #7239「请求官方出 .srs」已 closed/未采纳，.srs/.mrs 全靠第三方(MetaCubeX sing-box-geosite 自动转换 HaGeZi)。差异化价值：补「黑名单/拦截/威胁情报」这一维，分级(Light→Ultimate)+TIF 业界最细，与我们已知的分流型库(Loyalsoldier/blackmatrix7/Sukka/MetaCubeX)正交互补。值得加入 本项目 视野——作为 reject 维度的权威上游源(GPL-3.0，注意 copyleft 合规)，但要明确它解决的是拦截而非 AI/地区分流，且对接需我们自建转换流水线(扇出到 clash.list/.mrs/singbox.srs/surge.list)，不能指望它原生给客户端格式。

</details>

<details>
<summary>✅ <b>HenryChiao/MIHOMO_YAMLS</b> — active · AGPL-3.0</summary>

- **URL**: https://github.com/HenryChiao/MIHOMO_YAMLS
- **stars / 最近 / 维护**: ~2.0k (forks ~211) · 2026-06 · active
- **license**: AGPL-3.0
- **formats**: clash/mihomo .yaml (整套订阅配置), mihomo ruleset .list (classical/domain/ipcidr) + .mrs (meta 目录), singbox .srs + .json (version1-4 多版本), surge .list, quantumultx, shadowrocket, stash, loon, egern, geosite.dat/geoip.dat 派生 (GeoData 分支, 重打包上游)
- **provenance**: aggregator-of-others — 硬证据来自 .github/workflows: Merge_ruleset.yml (cron 30 23 * * *) clone/下载 MetaCubeX/meta-rules-dat(meta 分支)、blackmatrix7/ios_rule_script、Loyalsoldier/geoip、ACL4SSR/ACL4SSR、SunsetMkt/anti-ip-attribution、LM-Firefly/Rules,并叠加 SukkaW/felixonmars/NobyDa;update-geodata.yml (cron 0 20 * * *) gh release download 自 MetaCubeX/meta-rules-dat、DustinWin/ruleset_geodata、Loyalsoldier/v2ray-rules-dat+geoip、NobyDa/geoip、xream/geoip。GeoData 分支按上游分目录归档 (MetaCubeX 1795 条 / Loyalsoldier 1767 / xream 504 / DustinWin 353)。scripts/ruleset_process.sh 只做去重/排序/格式互转,自身无策展数据。结论:纯聚合+重打包,不是 domain-list-community 直接派生,也几乎无独立手工策展。
- **ai_coverage**: 强且新鲜。ruleset 分支有专门 ai 分类:meta/ai.list (~197 行) 与 singbox/ai.json+ai.srs。覆盖 anthropic.com/claude.ai/claude.com/claudeusercontent.com、openai.com/chatgpt.com、gemini.google.com、copilot.microsoft.com + GitHub Copilot、perplexity.ai/.com、grok.com/grok.x.com,含 DOMAIN-SUFFIX/KEYWORD/REGEX 混合。随每日 workflow 刷新。
- **specialty**: 差异化有限:它聚合的正是我们已知的 Loyalsoldier/blackmatrix7/MetaCubeX/Sukka/DustinWin 这几家,等于把这些源 + ACL4SSR + anti-ip-attribution 做了一次"多客户端扇出"(8 端:meta/singbox/surge/quantumultx/shadowrocket/stash/loon/egern,且 singbox 出 .srf/.srs 多版本)。对 本项目 唯一参考价值是它的"扇出矩阵工程"——一个上游源同时产 clash .list/.mrs + singbox .srs(v1-4) + surge/qx/egern 等,正是我们想自建的能力,可当工程参照/竞品基线;规则数据本身无新增源。
- **reason**: 对抗式核实结论:仓库真实存在且确为代理分流规则库(非无关项目),~2k star/211 fork,AGPL-3.0。活跃度=真活:main 2026-06-08、ruleset 分支 2026-06-09、GeoData 2026-06-08,均由 github-actions[bot] 提交,workflow 双 cron(23:30 与 04:00 北京时间)每日跑——"仓库活跃"与"规则真在更新"两者都成立,判 active。初判修正两点:(1) 初判"仅 mihomo YAML"错误——ruleset 分支实际扇出 8 客户端(meta/singbox/surge/quantumultx/shadowrocket/stash/loon/egern),singbox 出 .srs+.json 多版本,meta 出 classical/domain/ipcidr 三态;(2) 初判血统"aggregator"正确并有 workflow 硬证据。值得纳入 本项目 视野,但定位是"竞品/工程参照"而非"新上游源":它聚合的源(Loyalsoldier/blackmatrix7/MetaCubeX/Sukka/DustinWin/ACL4SSR)我们基本已覆盖,无独立策展数据;真正可借鉴的是它"单源→全客户端格式扇出"的自动化流水线,正对应我们 本项目 想自建的"上游源→自编译→扇出所有客户端"目标,可作 reference 实现对照。注意 AGPL-3.0 传染性:若我们直接复用其产物或脚本需评估 license 合规;作为思路参照不受限。

</details>

<details>
<summary>✅ <b>misakaio/chnroutes2</b> — active · CC-BY-SA-4.0</summary>

- **URL**: https://github.com/misakaio/chnroutes2
- **stars / 最近 / 维护**: ~850 · 2026-06 · active
- **license**: CC-BY-SA-4.0
- **formats**: chnroutes.txt (IPv4 CIDR plaintext), chnroutes.mmdb (MaxMind GeoIP2 DB)
- **provenance**: aggregator-of-others — 但聚合的是 BGP route-collector feeds (AS917 Misaka / AS906 DMIT / AS131477 / AS138195), 非 geosite/规则库聚合。依据: README + 文件头 '6686117 routes were dumped from route collector, aggregated to 3908'. 非 domain-list-community 派生 (它根本没有 domain, 纯 IP CIDR)。属基础设施级独立 IP 源。
- **ai_coverage**: none — 纯 IPv4 CIDR/GeoIP 源, 不含任何域名, 故无 AI/LLM/流媒体域名覆盖 (claude.com/grok/perplexity 等全部 N/A, 属错误维度)。
- **specialty**: CN GeoIP/CIDR 的"上游 source-of-truth"型供应链节点, 而非客户端规则库。与我们已知库(Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 主打 domain geosite)正交互补: 它解决的是 GEOIP-CN 这一面。差异化在于数据来自真实 BGP 路由表 + 每小时刷新, 比 APNIC delegated-latest 衍生的传统 chnroutes 更聚合(3908 条)、更实时。soffchen/GeoIP2-CN、GeoIP2-CN-misakaio 等下游都消费它来产 clash/sing-box 用的 .mmdb/.srs。
- **reason**: 对抗式核实通过且证据硬: (1) repo 真实存在、未 archive/disable; (2) 关键的"仓库活跃 vs 规则真更新"双轴都过 — RouteBot 'auto update' commit 逐小时落地 (04:00/03:00/02:00...), 且 chnroutes.txt 文件头自带时间戳 'Tue Jun 9 04:00:03 UTC 2026', 证明数据本身而非仅 metadata 在每小时刷新; (3) 但它不是 domain-based 分流规则库 — 纯 IPv4 CIDR(3908 条聚合)+ .mmdb, 不直接产 clash/sing-box/surge 规则。值得加入 本项目 视野, 但定位是 GEOIP-CN 维度的上游 source-of-truth(可自编译成 geoip.dat/.srs/.mrs), 而非 geosite 域名规则来源; 不与 Loyalsoldier/blackmatrix7/Sukka 竞争而是补它们的 IP 短板。注意 license 是 CC-BY-SA-4.0 (copyleft + 署名), 自编译扇出需保留署名并同协议分发, 比 Loyalsoldier 等的 license 约束更强, 供应链合规需单列。

</details>

<details>
<summary>✅ <b>mnixry/direct-android-ruleset</b> — active · AGPL-3.0</summary>

- **URL**: https://github.com/mnixry/direct-android-ruleset
- **stars / 最近 / 维护**: ~43 · 2026-06 · active
- **license**: AGPL-3.0
- **formats**: clash/mihomo classical .yaml (PROCESS-NAME rule-provider), surgio .tpl template
- **provenance**: independent-curated — 自有 scraper 抓取应用商店榜单 (src/provider/appchina.ts + qqdownloader.ts = AppChina + 腾讯应用宝/QQ下载器), 非 domain-list-community 派生, 也不聚合其他规则库; 规则全为 PROCESS-NAME 包名 (com.tencent.mm 等), 与域名系 geosite 血统正交
- **ai_coverage**: 无 — 纯 Android 国产 App 包名直连规则, 不含任何域名/AI 服务维度 (无 claude/grok/perplexity 等)
- **specialty**: 差异化价值在「规则维度」而非格式: 产出的是 Android 包名 PROCESS-NAME 直连规则 (国产 App 让其直连不走代理), 我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 全是 DOMAIN/IP 系, 没有任何一个覆盖 process-name 包名榜单维度。覆盖约 2000+ 个国产 App/Game 包名, 按应用商店榜单每日自动刷新。仅对 Android 客户端 + 支持 PROCESS-NAME 的内核 (mihomo) 有意义。
- **reason**: 对抗式核实结论: 1) 真实存在且确为代理分流规则库 (Mihomo/Clash classical + Surgio), 非无关项目, AGPL-3.0, 2025-02 创建。2) 关键陷阱已识破——main 分支看似 stale (最后人工 commit 2026-03-30, 无 bot commit), 但产出在独立的 rules 分支, 由 github-actions[bot] 经 daily cron (0 0 * * *) 每日自动 push, 最近一次更新就是今天 2026-06-09 → 判定 active (规则真的在更新, 非仅仓库活跃)。3) 格式窄: 仅 clash/mihomo .yaml + surgio .tpl, 无 sing-box .srs / surge .list / geosite.dat / quantumult-x。4) 血统=独立策展: 自有 scraper 抓 AppChina + 腾讯应用宝榜单生成包名, 非 domain-list-community 派生、非聚合他库。5) 差异化在于规则维度: 输出的是 Android PROCESS-NAME 包名规则 (~2000+ 国产 App), 与我们已知 4 库的 DOMAIN/IP 血统完全正交, 它们无一覆盖此维度; AI/流媒体覆盖=0。值得加入 本项目 视野——不是因为格式或域名覆盖, 而是作为「上游源」候选填补我们缺失的 Android 包名直连维度, 但定位小众: 仅对 Android 端 + 支持 PROCESS-NAME 的 mihomo 内核有用, 且产出已是 PROCESS-NAME 而非域名, 不能进 geosite 自编译主线, 只能作旁路源单独扇出。

</details>

<details>
<summary>✅ <b>privacy-protection-tools/anti-AD</b> — active · MIT</summary>

- **URL**: https://github.com/privacy-protection-tools/anti-AD
- **stars / 最近 / 维护**: ~10.5k (10465) · 2026-06 · active
- **license**: MIT
- **formats**: dnsmasq (adblock-for-dnsmasq.conf), AdGuard (anti-ad-adguard.txt), AdGuardHome/EasyList (anti-ad-easylist.txt), clash.yaml (anti-ad-clash.yaml, domain classical), plain domains/Pi-Hole (anti-ad-domains.txt), quantumult-x (anti-ad-quanx.txt), surge.list (anti-ad-surge.txt + anti-ad-surge2.txt, DOMAIN-SUFFIX), smartdns (anti-ad-smartdns.conf), singbox.srs (off-repo: anti-ad.github.io + anti-ad.net CDN), mihomo.mrs (off-repo: anti-ad.github.io + anti-ad.net CDN)
- **provenance**: aggregator-of-others — 依据: adlist-maker 分支含 scripts/prepare-upstream.sh + scripts/build-list.sh + .github/workflows/partial-update.yml, 抓取并去重多个上游 (AdGuard filters / EasyList / Fanboy annoyance / v2fly domain-list / URLhaus / neohosts / yhosts 等) 后扇出。NOT domain-list-community-derived: 它是 reject 黑名单 (纯域名), 不是 geosite 地理/分流数据集; 与 dlc 血统无关。
- **ai_coverage**: 无 AI/LLM 专项覆盖 (本质是 reject 黑名单, 无任何 AI/流媒体分流分类)。README 未声明任何保护开发者/API/AI 服务 (openai/anthropic/claude/github) 的白名单策略 — 若当作整盘 DNS reject 使用, 对 AI/API 域名存在误杀 (false-positive) 风险, 接入需自带白名单兜底。
- **specialty**: CN 区命中率最高的广告/隐私 REJECT 域名黑名单 (~11.3 万裸域 / ~10.6 万 DOMAIN-SUFFIX), 是单一专项 reject 源里覆盖最广、CN 针对性最强的。与我们已知库正交: Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 提供 geo+分流 ruleset (DIRECT/PROXY 路由), anti-AD 只补 reject 这一格。可作为我们 reject 分类的权威 CN 上游, 与 Sukka reject / blackmatrix7 AdvertisingLite 有重叠但更全。
- **reason**: 对抗式核实结论: (1) 真实存在, MIT, 10.5k star, 未 archived; 内容确为广告/隐私 REJECT 黑名单 (读 raw: anti-ad-domains.txt 11.3 万裸域、anti-ad-surge.txt 10.6 万条 DOMAIN-SUFFIX 且无 REJECT/DIRECT/PROXY policy target —— action 由消费端补), 不是 geosite/proxy 分流库。 (2) 关键 trap 已识别: GitHub Releases 停在 2020 v4.3、README 自述 v4.5.2(2022-12) 都是被废弃的人肉版本号; 但规则本体是活的 —— CI 机器人 'Auto renew anti-AD lists with/without upstream changes' 每周多次 commit, 规则文件内嵌 #VER=20260609 为今日生成, pushed_at=2026-06-09。判 active (仓库活跃 ✓ 规则真更新 ✓)。 (3) 血统=聚合器 (adlist-maker 分支 prepare-upstream.sh/build-list.sh 抓多上游去重), 非 dlc 派生。 (4) README over-claim 已抓: 主仓库 git 树只committed dnsmasq/adguard/easylist/clash.yaml/domains/quanx/surge/surge2/smartdns; 大力宣传的 sing-box .srs + mihomo .mrs 在主仓库 404, 实际由独立的 anti-ad.github.io 仓库 + anti-ad.net CDN 派发 —— 我们若要 srs/mrs 不能从本仓库取。 值得纳入 本项目 视野: 它补齐我们已知库缺失的 'CN 广告/隐私 reject' 这一专项, 但仅作 reject 上游之一, 不替代任何 geo/分流库; 接入时只取裸域/clash.yaml 自编译扇出, 不依赖其 off-repo 二进制, 并需叠加 AI/API 白名单防误杀。

</details>

<details>
<summary>✅ <b>QuixoticHeart/rule-set</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/QuixoticHeart/rule-set
- **stars / 最近 / 维护**: 425 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: clash/mihomo .mrs, clash/mihomo .list (domain/ipcidr/classical), sing-box .srs, sing-box .json (ruleset v5), surge .list, loon .list, stash .list, shadowrocket .list, quantumultx .list, egern .yaml
- **provenance**: aggregator-of-others (依据: 构建工作流 .github/workflows/run.yml 第26-33行 checkout MetaCubeX/meta-rules-dat 的 geo/geosite/classical + geo/geoip/classical 作为基底 geosite——这条线本身是 domain-list-community 派生(由 MetaCubeX 编译);随后 curl 叠加 20+ 上游: ruleset.skk.moe(Sukka,7+ 次拉取,最重)、ACL4SSR、Loyalsoldier/clash-rules + v2ray-rules-dat、felixonmars/dnsmasq-china-list、Cats-Team/AdRules、LM-Firefly/Rules、SunsetMkt/anti-ip-attribution、fmz200/wool_scripts、ConnersHua/RuleGo、NobyDa/geoip 等。本仓库自身的产出是 dedup+格式转换引擎(scripts/ruleset_process.sh)与多源合并配方 + 一个小的 custom/ 覆盖层。即: 基底=MetaCubeX meta-rules-dat(承 domain-list-community 血统), 上面聚合合并去重, 非独立手工策展。)
- **ai_coverage**: AI 覆盖很厚但全是聚合: meta/ai.list 合并 5 个 AI 上游——MetaCubeX geosite 的 category-ai-!cn、ACL4SSR Clash/Ruleset/AI.list、ConnersHua/RuleGo Surge AI.list、Sukka ruleset.skk.moe/List/non_ip/ai.conf, 外加 apple-intelligence(geosite + 本仓 custom/apple-intelligence.list 覆盖)。即对各家 AI 列表取并集而非择一。具体 anthropic/claude.com/openai/grok/perplexity 域名条目未逐条验证, 但因吃了 Sukka + geosite category-ai, 主流 LLM 域名实际覆盖度高(承上游)。我们已直接掌握这些上游, 故它在 AI 维度不提供新数据。
- **specialty**: 差异化价值不在「域名数据」而在「扇出打包」: 单仓库统一产出全部 8 客户端格式(含 Egern .yaml / QuantumultX / Loon / Stash 这些较少被一站式覆盖的目标), 二进制(mrs/srs)+文本(list/json/classical/domain/ipcidr 分类型)并存, 目录布局按客户端清晰分层。对 本项目 真正有用的是它的「单源→全客户端」工程蓝本: run.yml 合并配方 + ruleset_process.sh(clash domain→classical 转换、跨类型 DOMAIN/SUFFIX/WILDCARD/REGEX 去重、用真实 mihomo+sing-box 二进制在 CI 内编译 mrs/srs), 正好是我们要自建的 fan-out pipeline 的 working reference。
- **reason**: 真实存在且确为代理分流规则库(描述: 面向 mihomo/surge/loon/stash/shadowrocket/quantumultx/egern/sing-box 的定制规则集; 425★/33 fork, 未 archived)。活跃度上我做了「仓库活跃 vs 规则真更新」的区分核实: master 分支(源码 scripts/custom/.github, 最近人工 commit 2026-05-29)与 ruleset 分支(生成产物)分离; 「Generate RULE-SET」GitHub Action 今日 2026-06-09 及 06-08/07/06 均 success, force-push 到 ruleset 分支, 其 HEAD = github-actions[bot] 的「Auto Update Ruleset 2026-06-09」——所以规则确实每日重建, 不是绿勾但数据冻结的假活跃, 判 active。值得加入 本项目 视野但定位为「参考/低优先」而非「上游源」: 它是纯聚合器, 基底是 MetaCubeX meta-rules-dat(我们应直接取), 其余输入(Sukka/Loyalsoldier/ACL4SSR/AdRules)也都是我们已知库, 给不了任何我们没有的上游。它对我们唯一的实质价值是 prior-art: 整套 run.yml 合并配方 + ruleset_process.sh 转换/去重 + CI 内用真实 mihomo+sing-box 二进制编译 mrs/srs, 正是我们「单源→扇出全客户端」要解决的同一问题的可运行蓝本。注意 GPL-3.0: 可研究其 pipeline, 但不能把其产出直接 vendor 进非 GPL 产物。

</details>

<details>
<summary>✅ <b>Re:filter (1andrevich/Re-filter-lists)</b> — active · MIT</summary>

- **URL**: https://github.com/1andrevich/Re-filter-lists
- **stars / 最近 / 维护**: ~1.2k (1222) · 2026-06 · active
- **license**: MIT
- **formats**: xray/v2fly geoip.dat, xray/v2fly geosite.dat, sing-box geoip.db, sing-box geosite.db, sing-box .srs (refilter_domains / refilter_ipsum), sing-box rule-set .json, plaintext .lst (domains_all/ipsum/community/community_ips/discord_ips/ooni), BGP feed (AS65412)
- **provenance**: aggregator-of-others — 决定性依据: src/step1-download-and-word-filter.py 直接 `Downloads domains.lst from antifilter.download` (antifilter.download/list/domains.lst, 俄罗斯 RKN 封锁聚合源), 再叠加手工策展的 community.lst (PR 逐条加服务: Notion/Cloudflare WARP/DataDome 等) + OONI 列表 (step5) + Discord IP 子网。pipeline 还做主动 QC: step2 DNS/HTTPS 探活、step3 内容核查、step4 域名解析 + 经 bgpview/RIPE/ipinfo 富化 ASN/CIDR。**非 domain-list-community 派生** — README 里 `LoyalsoldierSite.dat:refilter` 只是 V2RayA 用法示例, 不是数据源; geoip/geosite 用 Dunamis4tw/generate-geoip-geosite 仅作打包工具。即"聚合 antifilter/RKN 源 + 主动验活 + 社区策展", 比纯转储更可信。
- **ai_coverage**: 弱/几乎无显式 AI 类目。README 与列表无 ChatGPT/OpenAI/Claude/Anthropic 专门分组; community.lst 偶有逐条加入"对俄封锁"的 SaaS (如 Notion), 但没有可直接扇出的 AI/LLM geosite 类别。要拿 AI 覆盖仍应靠 blackmatrix7/Sukka。
- **specialty**: 俄罗斯/RKN 封锁 + "对俄 IP 反向地理封锁" (geo-restrict RU) 专项规则, 这是 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 全部缺失的区域维度。独特资产: 经探活/内容核查清洗的 RKN 列表 (beta 称 ~8 万域名)、社区维护的"封锁了俄罗斯用户"服务清单 (community.lst)、Discord 语音 IP 子网、以及公开 BGP feed (AS65412) 可直接喂路由器 bird2。对我们而言是补"俄区/反俄地理封锁"这一格的唯一现实候选。
- **reason**: 真实存在且确为代理分流规则库 (README 给出 Xray routing + sing-box rule_set 用法), 非无关项目。维护"真活": 2026-06-07 有 push, 每 2-4 周由 github-actions[bot] 出 release (最近 31052026 / 2026-05-31), 且关键区分点——规则本身在更新: 多名人类贡献者 (Andrevich/German/Melgibzen) 持续向 community.lst 增删真实服务域名, 不是空转的定时器。产物覆盖 xray(.dat) + sing-box(.db/.srs/.json) + 纯文本 lst + BGP, MIT 许可可自由再编译。差异化价值清晰: 唯一成熟的"俄罗斯 RKN / 反俄地理封锁"区域库, 填我们已知四大库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 没有的空格。值得纳入 本项目 视野——但定位是"区域可选源 (俄区)", 非核心 AI/流媒体源; 接入时按 aggregator 对待 (上游 antifilter.download 不可控, 需在我们自编译层做去重/白名单), AI 覆盖仍需别处补。注意它不产 clash .mrs / surge / quantumult-x, 这些格式需我们自己从 geosite.dat 或 .lst 扇出。

</details>

<details>
<summary>✅ <b>runetfreedom/russia-blocked-geosite</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/runetfreedom/russia-blocked-geosite
- **stars / 最近 / 维护**: ~493 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: geosite.dat, per-category .txt domain lists, sha256sum, singbox.srs (via sibling russia-v2ray-rules-dat + geodat2srs converter, NOT in this repo)
- **provenance**: independent-curated (aggregator of Russian-censorship sources). 依据: 它的核心 blocked 列表是独立解析 antifilter.download / community.antifilter.download / re:filter(1andrevich/Re-filter-lists) / AdGuard DNS Filter / Peter Lowe / WindowsSpyBlocker 生成的, 这些才是 ru-blocked / ru-blocked-all / ru-available-only-inside 的真实来源。v2fly/domain-list-community 只是被打包进来的一个辅助输入 (用于 google/youtube/openai 等通用品类), 不是主血统。因此初判 'domain-list-community-derived' 不成立——本质是聚合俄罗斯审查源的独立策展库 (aggregator-of-others 偏向, 但聚合对象是上游审查数据而非别的规则库)。
- **ai_coverage**: 无专项 AI/LLM 策略。仅从 v2fly/domain-list-community 继承一个透传的 openai 品类; 无 anthropic/claude/grok/perplexity 专项。AI 覆盖不是它的价值点。
- **specialty**: 填补地理空白: 权威 Roskomnadzor(RKN) 封锁域名集 + 独有的 ru-available-only-inside (仅俄罗斯境内可访问) 品类。Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 全是 CN-centric 或全球向, 没有任何一家覆盖俄罗斯审查/反审查数据。是 v2rayN 俄罗斯分支的官方 geo 源, 生态成熟 (geosite + geoip + srs + custom-routing 配套仓库)。
- **reason**: 对抗式核实结论: 真实存在且确为代理分流规则库 (RKN 封锁域名 geosite), 非无关项目。维护是真活的——不只是仓库活跃: 2026-06-08 当天多次发布, 累计 ~2414 release, 6 小时 CI 每轮真正重新生成数据 (release 附带 66MB geosite.dat 全量重算)。判定 active。血统订正: 初判 'domain-list-community-derived' 错误——它是独立解析 antifilter/refilter/AdGuard/WindowsSpyBlocker 等俄罗斯审查源的聚合策展库, v2fly 只是辅助输入。值得加入 本项目 视野的理由: (1) 它覆盖的地理维度 (俄罗斯 RKN 封锁 + 境内可达专项) 是我们已知库 (CN/全球向) 的纯空白, 差异化清晰; (2) 配套生态完整 (geosite/geoip/srs), 对'上游源→自编译→扇出'模型是优质上游候选; (3) GPL-3.0 需注意——若我们自编译再分发其派生数据, copyleft 传染风险要在 本项目 license 兼容性矩阵里单列评估。局限: 本仓库不直接产 clash/.mrs/surge/quantumult-x, 只给 geosite.dat+txt, sing-box srs 要走 sibling 仓库或 geodat2srs 转换; 无 AI/LLM 专项价值。建议: 作为'区域(俄罗斯)专项上游源'纳入跟踪, 但不作为 AI/流媒体覆盖来源。

</details>

<details>
<summary>✅ <b>TG-Twilight/AWAvenue-Ads-Rule (秋风广告规则)</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/TG-Twilight/AWAvenue-Ads-Rule
- **stars / 最近 / 维护**: ~6k · 2026-05 · active
- **license**: GPL-3.0
- **formats**: singbox.json, singbox.srs, clash.yaml, clash-classical.yaml, surge.list, quantumultx.list, adguard.txt, hosts, dnsmasq.conf, mosdns_v5.txt, geosite.txt
- **provenance**: independent-curated — 单作者个人项目 ("个人项目，随缘维护更新")，主攻安卓应用内广告 SDK 逆向收集的 reject 域名。非 domain-list-community 派生 (geosite 只是它的一个输出格式而非血统来源)，也非聚合别人列表 (它自称是众多下游列表的"上游"); commit 是手工逐条增删 (如 "add WakeUp ads"、"remove m1.ad.10010.com"、"fix #209")，符合手工策展特征。
- **ai_coverage**: 无。纯广告/SDK reject 列表，无 claude.com/grok/perplexity 等 AI/LLM 分流域名覆盖，也无流媒体分流 (它只做 reject，不做 routing)。
- **specialty**: 极致体积控制的安卓 SDK/应用内广告 reject 单文件列表，逆向 APK 内嵌广告 SDK 域名是其独有数据源。与我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 不重叠：那些是 geo 分流/综合 reject，本库专注 in-app Android ad-SDK 拦截，是一个纯粹的 reject 补充源而非分流源。
- **reason**: 对抗式核实通过：仓库真实存在 (~6k star, GPL-3.0)，且经 commit 历史交叉验证规则数据本身在持续更新 (最近规则数据 commit 2026-05，非仅 docs/CI 灌水)，判定 active 而非 stale。它原生 fan-out 到 sing-box(.json+.srs)/clash/surge/quantumultx/adguard/hosts/dnsmasq/mosdns/geosite，与我们 本项目 的多客户端扇出目标天然契合，自编译成本低。血统为真正的 independent-curated (逆向安卓 ad-SDK 手工策展)，填补我们已知库的盲区——它是 reject/广告专项补充源，而非又一个 geo 分流库，不与 Loyalsoldier/blackmatrix7/Sukka 重复。值得加入视野。注意点 (不轻信 README)：(1) 它定位是上游被众多下游引用，但本身是单作者"随缘维护"，是 bus-factor=1 的供应链单点，纳入需做版本钉死+自编译校验；(2) 它不覆盖 AI/LLM/流媒体分流，不能替代我们对 geo 分流库的需求，只补 reject 这一格。

</details>

<details>
<summary>✅ <b>VPSDance/ai-proxy-rules</b> — active · MIT</summary>

- **URL**: https://github.com/VPSDance/ai-proxy-rules
- **stars / 最近 / 维护**: ~158 · 2026-06 · active
- **license**: MIT
- **formats**: clash/mihomo (.yaml), sing-box (.json, rule-set source form), surge (.list), shadowrocket (.list), loon (.list), stash (.list), quantumult-x (.list), egern (.yaml)
- **provenance**: aggregator-of-others (partly domain-list-community-derived) + thin manual curation. Hard evidence: data/cache/_index.json enumerates 22 upstream sources synced by scripts/sync — ~15 are v2fly/domain-list-community geosite entries (anthropic, openai, xai, cursor, groq, huggingface, perplexity, meta, jetbrains-ai, github-copilot, windsurf, manus, poe, cerebras, elevenlabs, google-deepmind), plus blackmatrix7/ios_rule_script (OpenAI/Gemini/Copilot), xiaolai/anthropic-claude-surge-rules-set, SkywalkerJi/Clash-Rules (Trae), and one html-extract scrape of ip.net.coffee/claude. The per-provider data/providers/*.yaml then hand-add Third-Party CDN domains, ASN and ipCidr/ipCidr6 on top. So: NOT independent from-scratch; it is an aggregator that re-bundles domain-list-community + blackmatrix7 + a few niche curated lists, with a light manual layer.
- **ai_coverage**: Strong/best-in-class for AI vertical. 82 providers covering all majors: OpenAI/ChatGPT, Anthropic/Claude (incl. claude.ai/claude.com/claude.new/MCP domains + ASN 399358 + 160.79.104.0/21), Google AI/Gemini, xAI/Grok, Meta AI, Mistral, Cohere, GitHub Copilot, Cursor, plus long-tail (Perplexity, Groq, HuggingFace, ElevenLabs, Midjourney/black-forest-labs, Character.AI, Devin, Augment, Kiro, JetBrains AI, Bytedance-AI, Cloudflare-AI, etc.). Per-provider files also carry Third-Party CDN domains and ASN/CIDR, which is richer than a bare domain list.
- **specialty**: Narrow but deep AI/LLM-only vertical: 82 per-provider rule files fanned out to 8 client formats from a single source-of-truth YAML (TypeScript/pnpm generator). Differentiator vs our known libs: Loyalsoldier=geoip/geosite .dat only; blackmatrix7=huge generic catalog (AI is a tiny corner, per-protocol .list, no unified AI bundle); Sukka=quality generic Surge/Mihomo lists (not AI-vertical); MetaCubeX=meta-rules/geo .mrs infra. This repo is the only one that (a) is AI-exclusive, (b) splits by individual AI vendor with an 'all' aggregate, (c) emits ALL 8 client formats incl. egern/stash/quantumult-x from one definition. Value to us = a ready-made cross-client AI-provider taxonomy + curated long-tail vendor domains (Cursor/Devin/Augment/Kiro/Trae/Genspark etc.) that domain-list-community lags on.
- **reason**: Verified real and genuinely a proxy split-routing ruleset (org VPSDance, MIT, ~158 stars, created 2026-04-25, TypeScript+JS generator, topics ai-proxy/clash/sing-box/geosite/...). Adversarial caveats that matter for 本项目: (1) Active-vs-updating split — CI runs a daily cron and committed as recently as 2026-06-08, but the commit log is dominated by 'Update rules: no rule changes'; real rule deltas are sparse (e.g. 'openai +1', 'google-ai +1' over weeks), so freshness is CI-freshness, not high content churn. (2) Provenance is downstream — it re-aggregates v2fly/domain-list-community + blackmatrix7 + a couple niche lists + one HTML scrape, so adding it gives us little NEW upstream signal we don't already reach via domain-list-community; its value is the curation/fan-out layer, not new ground-truth. (3) No releases/tags — output served as raw repo files (jsdelivr), so consuming it means pinning a commit, not a versioned artifact. Net: worth keeping in 本项目 VISIBILITY as a reference for AI-vertical taxonomy and long-tail AI-vendor domains + as a cross-client fan-out design exemplar, but it should be treated as a derived/secondary source, NOT a primary upstream — our own 'upstream→self-compile→fan-out' pipeline would directly consume domain-list-community/blackmatrix7 (its sources) rather than depend on this aggregator. Single-org, ~6-week-old, small (4 forks) = bus-factor/longevity risk; do not hard-depend.

</details>

<details>
<summary>✅ <b>xkww3n/Rules</b> — active · MIT</summary>

- **URL**: https://github.com/xkww3n/Rules
- **stars / 最近 / 维护**: ~173 · 2026-06 · active
- **license**: MIT
- **formats**: surge.list, surge-legacy.list (Loon/Shadowrocket/LanceX), clash.list (text-plus), clash.yaml, clash.mrs (Meta-compatible), singbox.srs, geosite.dat, geoip.dat, geoip.mmdb (MaxMind/IPinfo), quantumultx, egern, stash, surfboard
- **provenance**: domain-list-community-derived — 决定性证据：config.py 中 PATH_SOURCE_GEOSITE = Path("domain-list-community/data/")；CI main.yml 显式 checkout v2fly/domain-list-community 并用 Loyalsoldier/domain-list-custom + Loyalsoldier/geoip 编译；workers/v2fly.py 显式枚举上游 category（anthropic/openai/google/youtube 等）。reject 类聚合自 AdGuard/EasyList/秋风/NoCoin，国内 IP 取 gaoyifan/china-operator-ip。自有 source/ 仅约 25 个手工文件（bilibili/steam/wechat/apple-music/apple-intelligence/cmhk + personal/）作薄覆盖层。故非 independent-curated，而是 v2fly+Loyalsoldier 工具链的派生+裁剪，血统最接近 Loyalsoldier 而非 blackmatrix7。
- **ai_coverage**: 有但全部继承自上游 v2fly domain-list-community：workers/v2fly.py 选入 anthropic / openai / google-deepmind；上游存在 category-ai-chat-!cn。无 grok/perplexity/claude.com 级别的自有独立策展；自有 source/ 仅 apple-intelligence.txt（含 ChatGPT 集成解锁域名），无独立 AI 分类文件。
- **specialty**: 差异化在于「构建管线」而非「数据源」：(1) 体积最小化 + 跨规则集去重 + CIDR 合并（commit 与 utils/ruleset.py 可见）；(2) 单一源扇出极广——Egern/LanceX/Surfboard/QX/sing-box SRS/geosite.dat/mmdb 一次产出，正对 本项目「上游源→自编译→扇出所有客户端」诉求。但底层域名数据与 Loyalsoldier 高度重叠（同 v2fly base），不构成新的独立数据源。CDN: rules.xkww3n.cyou。
- **reason**: 真实存在且活跃（源码手工 commit 至 2026-06-05，CI 每日 12:00 UTC 重建，规则数据确在更新，非僵尸）。值得纳入 本项目 视野，但定位是「管线/扇出参考」而非「新数据源」：它的核心价值是一套成熟的 v2fly→去重→CIDR 合并→多客户端格式（含 Egern/LanceX/QX/sing-box SRS/mmdb）的开源 Python 构建管线（main.py + workers/ + utils/，MIT），正好示范我们要做的「上游源→自编译→扇出」；可直接借鉴其 dedup/minimization 与 fan-out target 设计。反例：底层域名数据派生自 v2fly domain-list-community + Loyalsoldier 工具，与我们已知的 Loyalsoldier 高度同源，AI 覆盖也仅继承上游（无 grok/perplexity 独立策展），故不能当作独立可信的"第二数据源"来交叉校验。初判"independent-curated"经对抗核实证伪，应修正为 domain-list-community-derived。

</details>

### ❌ 否决（57）

<details>
<summary>❌ <b>8680/GOODBYEADS</b> — active · MIT</summary>

- **URL**: https://github.com/8680/GOODBYEADS
- **stars / 最近 / 维护**: ~1.8k · 2026-06 · active
- **license**: MIT
- **formats**: AdGuard adblock (adblock.txt), AdGuard Home DNS (dns.txt), domain-only list (ad-domain.txt), QuantumultX (qx.list), SmartDNS (smartdns.conf + whitelist), allowlist (allow.txt)
- **provenance**: aggregator-of-others — README 自述合并 9 个上游去广告源 (AdGuard规则 / Perflyst SmartTV-AGH(TV) / EasyPrivacy / 乘风(Xinggsf)视频过滤 / NoAppDownload / OISD / AWAvenue秋风广告 / CJX Annoyance List + 自补充), GitHub Actions(actions-user) 每日 2 次自动 merge+去重。无任何 domain-list-community / geosite / GFW 派生痕迹 —— 因为它根本不是分流库。
- **ai_coverage**: 无 AI/LLM 路由覆盖。该库目标是"拦截"而非"分流",不存在 claude.com/openai/gemini 的代理走向规则;即便其黑名单中可能命中某些广告/追踪域名,也与 AI 服务分流无关。流媒体同理 —— 只有"视频广告过滤"(乘风规则),没有 Netflix/Disney 解锁分流。
- **specialty**: 无差异化价值(对 本项目 而言)。它是纯去广告(ad-blocking/拦截)规则聚合器,产出的是 block 域名黑名单 + DNS sinkhole + 去广告过滤语法,不产出代理分流(routing/policy)规则,不区分国内外/流媒体/直连代理走向。与我们已知的 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX(分流定位)正交;若真要用,它对标的是 anti-AD / AdGuardSDNSFilter / 217heidai 这类去广告库,而非分流库。
- **reason**: 核实结论:仓库真实存在、确为去广告规则库,GitHub Actions 每日 2 次自动 commit(最新 2026-06-09),所以"规则真的在更新"而非仅"仓库活跃"——active 判定成立,约 1.8k star,MIT。但对抗式核查推翻初判中"代理分流"的隐含定位:它是 ad-blocking 聚合器(aggregator),不是代理分流/geosite 派生库。产出全部是 block 类格式(AdGuard/AdGuardHome DNS/QuantumultX 去广告 list/SmartDNS),无 Clash/sing-box/Surge 的 routing rule-provider,无国内外/流媒体/AI 分流概念,与 domain-list-community 无血缘。本项目 关注的是"上游分流源→自编译→扇出各客户端分流格式";GOODBYEADS 解决的是正交的"去广告"问题域。除非 本项目 显式扩展到"去广告规则也要自托管扇出",否则不值得纳入分流供应链视野。即便要做去广告,更应优先对标 anti-AD(已支持 Clash/mihomo/sing-box/Surge/SmartDNS 多格式且命中率口碑更高),GOODBYEADS 是其下位替代。

</details>

<details>
<summary>❌ <b>Aethersailor/Custom_Clash_Rules</b> — stale · none (无 LICENSE 文件, API license=null)</summary>

- **URL**: https://github.com/Aethersailor/Custom_Clash_Rules
- **stars / 最近 / 维护**: ~68 · 2026-02 (real content); default branch wiped 2026-03 · stale
- **license**: none (无 LICENSE 文件, API license=null)
- **formats**: clash/mihomo subconverter .ini templates (assembles []GEOSITE/[]GEOIP + remote .list rule-providers)
- **provenance**: independent-curated (but NOT a rule-data repo) — 依据: 它本身不含规则数据, 只有 cfg/*.ini subconverter 模板. 模板 ruleset= 指向作者自己的 Aethersailor/Custom_OpenClash_Rules@main/rule/*.list (Custom_Direct/Custom_Proxy/Steam_CDN) + Loyalsoldier/v2ray-rules-dat 的 GeoSite/GeoIP. 上游 .list 由作者 Telegram 提交机器人 (Rule-Bot) 人工/半自动策展, 非 domain-list-community 派生. 无 blackmatrix7/ACL4SSR 引用.
- **ai_coverage**: 无专项 AI/LLM 域名规则 (README 与模板均未见 Claude/OpenAI/Gemini/Copilot 分组); AI 流量靠 GeoSite 兜底, 无 claude.com/grok/perplexity 细粒度覆盖
- **specialty**: 差异化价值低且已贬值: 卖点是「OpenClash 教程配套 + 无 DNS 泄露的 mihomo 订阅模板」, 即配置工程 (subconverter 模板 + DNS 防泄露布局) 而非规则数据本身. 相对 Loyalsoldier(geosite 数据)/blackmatrix7(多端规则)/Sukka(自编译)/MetaCubeX(官方 geo) 没有新规则资产 — 它消费 Loyalsoldier + 作者自家 .list. 真正有价值的规则资产在其姊妹库 Custom_OpenClash_Rules (6.1k star, rule/ 目录到 2026-06 仍由 Telegram bot 活跃更新, 产出 .list/.yaml/.mrs), 本候选只是它的客户端模板壳.
- **reason**: 对抗式核实结论: (1) 真实存在但不是规则库 — 它是 clash/mihomo 的 subconverter .ini 模板集合 (cfg/Custom_Clash*.ini), 自身零规则数据. (2) 已死/停更: README 明写「2026/2/18 本项目暂不维护」, main 分支真实内容停在 2026-02-25; 更糟的是默认分支被切到名为 rm 的分支, 该分支唯一 commit「rm」于 2026-03-05 删掉全部 715 行内容, README 只剩单词 rm — 作者主动废弃此壳库. 注意区分: 仓库 pushed_at 看似新只因 main 分支残留, 规则其实没在这里更新. (3) 产出格式单一: 仅 clash/mihomo 模板, 无 sing-box .srs / surge / quantumult-x. (4) 血统: 非 domain-list-community 派生; 模板引用 Loyalsoldier/v2ray-rules-dat 的 geo 数据 + 作者自家 Custom_OpenClash_Rules/rule 的人工策展 .list. (5) 无 license, 无 AI 专项覆盖. (6) 不值得加入 本项目 视野 — 它对我们「上游源→自编译→扇出多客户端」目标没有可消费的规则资产, 且已废弃. 真正值得单独评估的是其姊妹库 Aethersailor/Custom_OpenClash_Rules (rule/ 目录活跃产出 .list/.yaml/.mrs, 由 Telegram 提交机器人驱动) — 建议把追踪目标换成它, 而非本候选.

</details>

<details>
<summary>❌ <b>Aethersailor/Custom_OpenClash_Rules</b> — active · CC-BY-SA-4.0 (文件名 LICENCE; README 附加声明: 不鼓励转载、禁止转载到中国大陆平台)</summary>

- **URL**: https://github.com/Aethersailor/Custom_OpenClash_Rules
- **stars / 最近 / 维护**: ~6.1k (6146) · 2026-06 · active
- **license**: CC-BY-SA-4.0 (文件名 LICENCE; README 附加声明: 不鼓励转载、禁止转载到中国大陆平台)
- **formats**: clash/mihomo .list, clash/mihomo .mrs (Mihomo binary), clash rule-provider .yaml (Classical / Domain / IP-CIDR variants), subconverter .ini template (Custom_Clash.ini, drives clash/mihomo + sing-box/surge/etc via subconverter)
- **provenance**: independent-curated (with a twist). 依据: (1) rule/README 自述 "个人维护的轻量级规则碎片"; Custom_Proxy.list 头部 "强制代理列表,来自个人收集"; 由 Telegram Rule-Bot + GitHub Actions 自动收录单个域名 (commit msg "add direct domain X by Telegram Bot")。(2) 自有规则极小 (Custom_Proxy.list 仅 ~18 条, 820B; 主体是 Custom_Direct/Game_Download_CDN/IPTVMainland/Steam_CDN 等大陆直连碎片), 非 domain-list-community 批量派生, 非聚合 blackmatrix7/Loyalsoldier 的成品 list。(3) BUT 它的品类级分流 (AI/Netflix/Telegram/Google...) 全部靠 subconverter 模板里的 []GEOSITE,openai / category-ai-!cn / netflix 等 builtin tag, 即消费 MetaCubeX/domain-list-community 的 geosite.dat —— 自己不重新策展这些品类。所以血统 = 独立手工策展(仅限大陆直连碎片) + 上游 geosite 消费者(品类分流)。
- **ai_coverage**: 弱/间接。仓库自有规则里 0 个 AI/LLM 域名 (grep Custom_Proxy.list 无 openai/claude/anthropic/gemini/perplexity/grok)。模板通过 []GEOSITE,openai + []GEOSITE,category-ai-!cn 把 AI 分流, 但这些 AI 域名数据来自 MetaCubeX geosite, 不是本仓库策展。即: 它给 AI 一个独立策略组, 但 AI 域名清单不归它维护。无 claude.com/grok/perplexity 专项自有清单。
- **specialty**: 差异化价值 = DNS 防泄漏/防污染的 OpenClash 整套配置方法论 + 模板 (Custom_Clash.ini 及 Full/GFW/Lite/Mainland 变体), 而非规则库本身。规则面相对 Loyalsoldier/blackmatrix7/Sukka 几乎无新增覆盖 (品类全靠 MetaCubeX geosite)。唯一独家增量是手工收录的大陆直连冷门域名碎片 (Custom_Direct ~58KB, Game_Download_CDN, IPTVMainland, Steam_CDN 精确匹配, Talkatone) + 自动化 Telegram 提交机器人这套运营模式。对"自建供应链/扇出多客户端格式"的目标, 它本身不是上游源, 反而是下游消费者。
- **reason**: 真实存在且确为代理分流相关仓库 (OpenClash 生态, 6.1k star), 仓库高度活跃 (Telegram Rule-Bot + Actions 每日自动生成, 最近 commit 2026-06-07, 距今 2 天; 规则确实在更新而非僵尸活跃)。注意一个对抗式发现: 默认分支被设为空的 'rm' 分支 (仅 README+.github, 单 commit msg "rm"), 真正内容和活跃度都在 main 分支——只看默认分支会误判为 dead, 必须看 main。 不值得加入 本项目 视野, 原因: (1) 血统上它是 MetaCubeX/domain-list-community geosite 的下游消费者, 不是上游源——品类分流(含 AI/流媒体/Telegram/Google)全靠 []GEOSITE 内置 tag, 自己不重做这些品类; 我们要的"上游源→自编译"它不提供新源。(2) 自有规则仅大陆直连冷门碎片 (Custom_Direct/IPTVMainland/Game CDN/Steam CDN/Talkatone), 与 Loyalsoldier/blackmatrix7/Sukka 重叠且覆盖面更小, 无差异化增量。(3) AI/LLM 自有覆盖为 0。(4) 它的核心价值是 OpenClash DNS 防泄漏配置方法论 + subconverter 模板, 属于"客户端配置范式", 与我们 本项目"规则供应链审计"目标正交。(5) license CC-BY-SA-4.0 且 README 明确不鼓励转载/禁止搬运到大陆平台, 复用其规则碎片有 ShareAlike 传染性约束。结论: 可作为"OpenClash DNS 防泄漏配置参考"收藏, 但不是规则供应链上游, 不纳入 本项目 规则源清单。

</details>

<details>
<summary>❌ <b>Aoang/Surge</b> — active · MIT (主体); China 规则集 CC BY-SA 2.0; Base 去广告列表 GPLv3 — 混合 license, 派生需注意</summary>

- **URL**: https://github.com/Aoang/Surge
- **stars / 最近 / 维护**: ~31 · 2026-05 · active
- **license**: MIT (主体); China 规则集 CC BY-SA 2.0; Base 去广告列表 GPLv3 — 混合 license, 派生需注意
- **formats**: surge.list (DOMAIN-SET/RULE-SET), surge.sgmodule (.sgmodule modules)
- **provenance**: independent-curated (with selective aggregation). README 自述「继神机策略 DivineEngine/Profiles 之后的自用规则」, 大部分 Ruleset(Apple/CDN/Telegram/PayPal/Microsoft) 标注「手动维护」. 验证 src/ 构建逻辑(Deno/TS): 仅 China IP 自动取自 misakaio/chnroutes2, Bogus 取自 felixonmars/dnsmasq-china-list, 去广告取自 AdGuard/EasyList。无任何 domain-list-community / GeoSite 血统, 不是 v2fly 派生, 也不是大聚合器, 而是小型个人手工策展 + 少量上游 IP/广告列表拼装。
- **ai_coverage**: 几乎无。全仓 grep 仅一条 AI/LLM 相关域名: Download.list 中「## Claude Code → downloads.claude.ai」, 这是为下载 Claude Code 二进制的开发者下载源, 非 AI 代理分流。无 api.openai.com / claude.ai 对话 / perplexity / grok / gemini 等任何 AI 服务分流规则, 也无专门 AI 类目。所谓「Add Claude AI」commit 只是往下载镜像表加了一个域名。
- **specialty**: 相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 无差异化价值: 体量极小(20 commits, 95 paths, 约 14 个 .list), 仅覆盖 Apple/CDN/Telegram/PayPal/Microsoft/China/去广告等大库早已完备的常见类目, 且只产出 Surge 单一格式(无 clash/sing-box/.srs/.mrs/geosite)。手工维护意味着覆盖窄、滞后于上游。唯一略特别的是 sgmodule 模块(MitM/FakeIP 绕过/隐藏VPN图标), 但属 Surge 客户端配置而非分流规则, 与 本项目 供应链审计无关。
- **reason**: 真实存在且确为代理分流规则库(Deno/TS 构建 → _site/ 输出 Surge .list + .sgmodule), 非无关项目。仓库活跃(最近 commit 2026-05-27, 非 archived/disabled, created 2023, 持续小步提交), 但「活跃」是指作者偶尔加几个域名, 规则本身大多手工维护、类目窄。结论: 不值得加入 本项目 视野。理由: (1) 单一 Surge 输出格式, 我们要的是「上游源→自编译→扇出多客户端」, 它既非高质量上游源也非多格式产物; (2) 血统是个人手工策展, 覆盖面被 Loyalsoldier/blackmatrix7/Sukka 全面覆盖且后者更全更新更勤; (3) AI/流媒体专项几乎为零(仅一条 Claude Code 下载域名); (4) 混合 license(MIT+CC-BY-SA+GPLv3) 反而给派生再分发增加合规负担。无任何差异化价值, 排除。

</details>

<details>
<summary>❌ <b>carsondzh/clash-geosite</b> — dead · GPL-3.0</summary>

- **URL**: https://github.com/carsondzh/clash-geosite
- **stars / 最近 / 维护**: 1 · 2023-11 · dead
- **license**: GPL-3.0
- **formats**: geosite.dat, geosite.db (sing-box), geoip.dat, Country.mmdb, geoip.metadb, geoip.db, clash user.yaml (fake-ip/redir-host)
- **provenance**: domain-list-community-derived (two hops, frozen). GitHub API confirms fork chain: carsondzh/clash-geosite -> parent DustinWin/ruleset_geodata -> source Loyalsoldier/v2ray-rules-dat. Loyalsoldier is the canonical v2fly/domain-list-community enhanced derivative. README also lists v2fly/domain-list-community, Privacy-protection-tools/anti-AD, blackmatrix7/ios_rule_script, felixonmars/dnsmasq-china-list as upstreams. So lineage is genuine but inherited entirely from DustinWin, not independently curated by carsondzh.
- **ai_coverage**: README 自述含 OpenAI 分类(继承自 DustinWin geosite-all 变体),但本 fork 无任何产物文件/release,无法验证;实际 AI 覆盖只有上游 DustinWin 在持续更新。
- **specialty**: 无差异化价值。它只是 DustinWin/ruleset_geodata 的一个废弃 fork,fork 后从未被 owner 改动。所宣称的 anti-AD/Netflix/Disney+/YouTube/TikTok/OpenAI full/lite 双版本能力全部继承自上游 DustinWin,而 DustinWin 本体(1264 star, 每天真实构建, 最近 push 2026-06-08)才是有价值的源。相对我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 没有任何新增覆盖面。
- **reason**: 对抗式核实结论:这是一个死掉的 fork-of-a-fork,不值得加入 本项目 视野。证据链:(1) 仓库真实存在且属代理分流规则范畴,但内容仅 .github/+LICENSE+README,无任何规则数据;(2) 全部 commit 停在 2023-11-25 单日,且作者均为上游 DustinWin(fork 那刻继承的),pushed_at=2023-11-25,owner carsondzh 自 fork 后从未提交;updated_at=2026-02 只是 star/watch 元数据噪音,不代表规则更新 -> 判定 dead;(3) releases=0 tags=0 branches 仅 master,README 宣称的"每天 3AM 自动构建 + jsDelivr 分发 geosite.dat/.db"在本 fork 完全不成立(fork 的 scheduled Actions 默认不运行,且无任何 release 资产)-> README 自述与事实矛盾,典型的 fork 不轻信案例;(4) 血统为 domain-list-community 派生但两跳且冻结,无独立策展;(5) GPL-3.0,无差异化覆盖。真正该追踪的是其上游 DustinWin/ruleset_geodata(1264 star,每天真实构建,最近 push 2026-06-08,topics 覆盖 clash/mihomo/sing-box/geox,且自身又派生自 Loyalsoldier)——建议把 DustinWin 列入 本项目,而非这个 fork。

</details>

<details>
<summary>❌ <b>Ckrvxr/mihomo_yaml</b> — active · Apache-2.0</summary>

- **URL**: https://github.com/Ckrvxr/mihomo_yaml
- **stars / 最近 / 维护**: ~3 · 2026-05 · active
- **license**: Apache-2.0
- **formats**: clash.yaml (mihomo override JS scripts: PLUS.js / STD.js), mihomo rule-provider .yaml (behavior: domain/classical/ipcidr payloads)
- **provenance**: aggregator-of-others — 依据: 核心 geo 分流直接拉 MetaCubeX/meta-rules-dat 的 geosite.dat/geoip.dat (即 domain-list-community 派生数据,被消费而非自维护); AD/防诈/IP 规则引用 217heidai/adblockfilters、StevenBlack/hosts、TG-Twilight/AWAvenue-Ads-Rule、Firehol level1、zsokami/ACL4SSR。仅 Source/Addition/ 下少量手工 yaml (AntiAntiFraud.yaml 仅 ~11 条 process-name+1 domain-suffix) 算自编,但量极小且无 domain-list-community 自派生证据。本质是把别人的清单编排进 mihomo override 脚本。
- **ai_coverage**: 有 AI 分流组 (OpenAI/Anthropic/X.AI proxy group),但底层匹配靠 geosite:openai 等 MetaCubeX 数据,无任何自维护 AI/LLM 域名清单 — AI 覆盖完全等同上游 meta-rules-dat,无增量。
- **specialty**: 差异化几乎为零: 它不是规则源而是 mihomo override 配置生成器 (FlClash/Sparkle 用 JS 脚本注入)。唯一原创资产是极小的 AntiAntiFraud / AntiPCDN process-name 黑名单 (反国产手机反诈/PCDN 上传, ~10 条 Tencent/Qihoo/Xiaomi 包名), 这块在 Loyalsoldier/blackmatrix7/Sukka 里确实没有对等物 — 但规模太小,且是 PROCESS-NAME 维度,无法 fan-out 到非进程感知客户端 (surge/quantumultx/singbox 多数走 domain),对我们 本项目 的「域名规则供应链」基本无用。
- **reason**: 对抗式核实结论: 仓库真实存在且近期活跃 (最近 commit 2026-05-29, 但总共仅 5 commit / ~3 star / 0 fork, 是 2026-04-25 新建的个人小项目, 注意"仓库活跃"≠"规则在更新"—近期 commit 多为 CDN URL 迁移和 README 文案, 非规则数据更新)。它本质是 mihomo override 脚本聚合器 (aggregator), 不是独立策展的规则源: geo 分流直接消费 MetaCubeX/meta-rules-dat (domain-list-community 派生), AD/防诈引用 217heidai/StevenBlack/AWAvenue/Firehol, 仅有 ~11 条原创 AntiAntiFraud process-name 黑名单。相对我们已知的 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX, 它没有提供任何可被我们自编译/扇出的新域名数据 — AI 覆盖无增量, 原创资产是无法跨客户端 fan-out 的 PROCESS-NAME 维度。不值得加入 本项目 视野; 唯一边角价值是「国产手机反诈/反 PCDN 上传 process-name 清单」这一冷门专项, 但量级和格式都不适合纳入域名规则供应链, 顶多作为 process-name 灵感参考一次性抄走即可, 无需持续追踪。

</details>

<details>
<summary>❌ <b>ClashConnectRules / Self-Configuration</b> — stale · MIT</summary>

- **URL**: https://github.com/ClashConnectRules/Self-Configuration
- **stars / 最近 / 维护**: ~1.3k · 2026-05 · stale
- **license**: MIT
- **formats**: clash.yaml (subscription config template, not a ruleset)
- **provenance**: aggregator-of-others — 依据: 唯一数据文件 Clash.yaml 不含任何 inline 域名规则, 运行时 100% 从外部 rule-provider URL 拉取规则, 主源 dler-io/Rules (经 jsDelivr) + blackmatrix7/ios_rule_script。它本身不编译/不维护任何规则列表; org 下的 'Rules' repo 也只是 fork 自 dler-io/Rules。本质是下游消费者/编排器, 不是上游规则源。
- **ai_coverage**: 名义有 AI 分组 (ChatGPT/Claude/Gemini), 但规则域名全部外链自 dler-io 'AI Suite', 非自维护; 无独立 AI 域名清单。
- **specialty**: 无差异化规则价值。AI(ChatGPT/Claude/Gemini) + 流媒体(Netflix/Disney+/YouTube/TikTok) 的"覆盖"只是指向 dler-io 的 'AI Suite' 与媒体 provider 的分组映射, 域名数据并非自有。它真正的上游 (blackmatrix7, 间接 domain-list-community via blackmatrix7, 以及 dler-io) 已在我们已知集合内。相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX, 它只多了一个推广付费订阅的 config 模板, 无规则增量。
- **reason**: 对抗式核实结论: 真实存在 (~1.3k star, MIT) 但被错误归类——它根本不是规则库/geosite 派生库, 而是单个开箱即用的 Clash.yaml 订阅配置模板 + README 营销。决定性证据: (1) 唯一数据文件 Clash.yaml 自 2025-12-25 initial commit 后从未被修改过一次, 之后 2026-03~05 的所有 commit 全是 README 改动 (加订阅邀请链接、'ZRJ recommendation' 推广)——即"仓库活跃"是假象, 规则从未更新, 运行时新鲜度完全外包给 dler-io/blackmatrix7。(2) 零 inline 规则, 全部 rule-provider 外链他人。(3) org 的 'Rules' 是 dler-io fork。对 本项目 (上游源→自编译→扇出) 而言它是纯下游消费者, 无任何可控上游价值, 不值得纳入视野; 若要追踪应直接看其真上游 dler-io/Rules 与 blackmatrix7。

</details>

<details>
<summary>❌ <b>CloudPassenger/geosite</b> — stale · GPL-3.0</summary>

- **URL**: https://github.com/CloudPassenger/geosite
- **stars / 最近 / 维护**: 5 · 2026-06 (release/build), 2024-09 (source code) · stale
- **license**: GPL-3.0
- **formats**: geosite.dat (v2ray/xray), geosite.db (sing-box binary, via metacubex/geo), dlc.dat, dlc.db, domain text lists (proxy/direct/reject/china/apple-cn/google-cn/gfw/win-*), tld lists
- **provenance**: domain-list-community-derived — 硬证据: GitHub API "fork":true, parent=Loyalsoldier/v2ray-rules-dat。workflow run.yml 实际 clone Loyalsoldier/domain-list-custom + v2fly/domain-list-community 作为输入,用 `go run ./ --datapath=../community/data` 编 geosite.dat。贡献统计 Loyalsoldier 293 commits vs CloudPassenger 仅 12 commits。本质是 Loyalsoldier 数据集,非独立策展。
- **ai_coverage**: 无独立 AI 策展。仅继承 v2fly/domain-list-community 的 data/openai (595 bytes,极小),与所有 Loyalsoldier 派生库拿到的同一份;README 零次提及 AI/LLM/OpenAI/ChatGPT。相比 blackmatrix7 (有独立 OpenAI/Claude/Gemini 分流) 明显落后。
- **specialty**: 唯一差异化: 在 Loyalsoldier 工具链尾部加了一步 `metacubex/geo convert site -i v2ray -o sing` 生成 sing-box 二进制 geosite.db / dlc.db —— 而 upstream Loyalsoldier 最新 release (202606082308) 只发 .dat + txt + rules.zip,不发 .db。即 "Loyalsoldier 数据 + 现成 sing-box .db 格式封装"。对我们价值极低: 我们 本项目 的目标正是自己掌控 上游源→自编译→扇出,这一步格式转换 (geo convert / sing-box geosite compile) 我们本就要自己做,没必要依赖一个 5-star 个人 fork 代劳。
- **reason**: 对抗式核实结论: 真实存在,确为代理分流规则库,但价值不足以进 本项目 视野。(1) 它是 Loyalsoldier/v2ray-rules-dat 的硬 fork (API parent 字段证实),不是独立来源 —— 我们已知清单里的 Loyalsoldier 就是它的上游,纳入它等于重复。(2) "活跃"是假象: 每日 6AM GitHub Actions cron 仍在跑、release 停在 2026-06-08 看似 active,但源码/策展逻辑冻结在 2024-09 (最后一次自有 commit 2024-08 merge upstream)。规则内容的新鲜度完全来自每次 build 重新拉 v2fly community,CloudPassenger 自身零增量策展。所以"仓库活跃 ≠ 规则在被人维护"。(3) 唯一卖点是把 Loyalsoldier 数据用 metacubex/geo 转成 sing-box .db,这恰好是我们自己要做的扇出步骤,不构成上游源价值。(4) 只有 5 star、单人维护、has_issues=false (关了 issue),供应链信任度低。建议: 不追踪此 fork;若需 sing-box geosite 直接用 SagerNet/sing-geosite (官方) 或自己 geo convert,上游数据继续认 v2fly community / Loyalsoldier 本体即可。

</details>

<details>
<summary>❌ <b>cmontage/proxyrules-cm</b> — stale · GPL-3.0</summary>

- **URL**: https://github.com/cmontage/proxyrules-cm
- **stars / 最近 / 维护**: 8 · 2026-06 · stale
- **license**: GPL-3.0
- **formats**: clash.yaml (rule-provider behavior:domain/classical lists, 非 .mrs/.list), quantumultx.conf, quantumultx.yaml, singbox source rule-set .json (未编译 .srs)
- **provenance**: aggregator-of-others — 依据: scripts/sync_rules.py 直接从 https://raw.githubusercontent.com/blackmatrix7/ios_rule_script/master/ 拉取 (Apple/China/Microsoft/Google/Netflix/PayPal + 自定义 AI=OpenAI/Gemini/Claude 合并); README 自述"参考于 ACL4SSR / blackmatrix7/ios_rule_script 并根据自己需求修改"; 连产出文件头都保留 repo: blackmatrix7/ios_rule_script。是 blackmatrix7 的下游薄壳 (blackmatrix7 本身大量派生自 domain-list-community)。初判 independent-curated 不成立。
- **ai_coverage**: 有 OpenAI/ChatGPT、Claude (claude.ai)、Google Gemini/Bard/Colab、Anthropic; AI.yaml 仅 ~50 域名, 是 blackmatrix7 OpenAI/Claude/Gemini 列表的子集。初判提到的 Copilot 实际未单独包含。
- **specialty**: 几乎无差异化价值: 是 blackmatrix7/ios_rule_script 的人工挑选小子集 (AI 类仅 ~50 域名, 严格子集), 三态 DIRECT/PROXY/REJECT 目录结构按个人偏好组织 (Google 走美国/Netflix 走新加坡), 但规则内容零自有情报。相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 没有任何独占覆盖或更优编译产物。
- **reason**: 对抗式核实结论: 真实存在且确为代理分流规则库, 但(1)血统是 blackmatrix7 薄壳聚合器, 非独立策展, 也非 domain-list-community 直接派生(隔了一层 blackmatrix7); (2)"活跃"是假象——97 commit 里 73 个是 github-actions[bot] 每日 'Auto-sync rules from blackmatrix7/ios_rule_script', 人工 commit 自 2025-02 起基本停摆(2026-03-28 那次只是接入自动化), 规则新鲜度完全继承上游而非自有维护; (3)8 star/1 fork, AI 仅 ~50 域名子集, 格式只有 clash.yaml/QX/singbox-source-json, 无 .mrs/.srs/surge/geosite.dat 编译产物。对 本项目(上游源→自编译→扇出) 无任何增量: 要追就直接追 blackmatrix7, 这个下游壳子不值得进视野。

</details>

<details>
<summary>❌ <b>cutethotw/ClashRule</b> — active · 无 LICENSE 文件 (root 目录无 license, README 为订阅转换教程, 定位"自用"个人库)</summary>

- **URL**: https://github.com/cutethotw/ClashRule
- **stars / 最近 / 维护**: ~901 · 2026-06 · active
- **license**: 无 LICENSE 文件 (root 目录无 license, README 为订阅转换教程, 定位"自用"个人库)
- **formats**: clash.list (classical rule-provider), clash.yaml (classical payload), subconverter .ini (config templates)
- **provenance**: aggregator-of-others — 依据: AI.yaml 头部注明来源 github.com/VPSDance/ai-proxy-rules; Claude.list 注释引用 ip.net.coffee/claude/site.html 和 nodeseek.com 论坛帖。手工汇编社区来源, 非 domain-list-community/geosite 派生, 也非完全原创手策。
- **ai_coverage**: 覆盖强: 专门的 Claude(claude.com/anthropic.com/clau.de + 辅助依赖)、ChatGPT、Grok、Perplexity 独立 list, 外加 AI.yaml 聚合 OpenAI/Google/Anthropic/Meta/Mistral/Cohere 等 ~150 域。Claude.list 于 2026-06-05 更新。
- **specialty**: AI 域名专项是其相对亮点: 独立 Claude.list/ChatGPT.list/Grok.list/Perplexity.list + 宽口径 AI.yaml(~150域)。Claude.list 不止裸 suffix, 还含 Auth0/Intercom/Sentry/Statsig/Datadog 等第三方 auth/遥测依赖, 边界比裸列表更细。但仅 Clash 单生态, 规模远小于 blackmatrix7/Sukka。
- **reason**: 真实存在且活跃(最近 commit 2026-06-07, 规则文件本身在更新, 非僵尸库)。但不符合 本项目「掌控上游源→自编译→扇出」的定位: (1) 无 LICENSE, 个人"自用"框架, 法律上不能安全二次分发/扇出; (2) 单一 Clash 生态(.list/.yaml/.ini), 无 geosite.dat/.srs/.mrs, 无 surge/quantumultx, 不利于多客户端扇出; (3) 它本身是下游聚合器, 上游(VPSDance/ai-proxy-rules、ip.net.coffee、nodeseek 帖)才是更值得直接追踪的源。结论: 不纳入 本项目 作为权威上游, 顶多留作 AI 域名 edge-case 交叉校验参考(尤其 Claude 的 auth/遥测第三方依赖那几条), 但它引用的上游更值得直接对接。

</details>

<details>
<summary>❌ <b>DDCHlsq/sing-ruleset</b> — active · 无 (license: null, 仓库无 LICENSE 文件) — 采用层面的硬伤</summary>

- **URL**: https://github.com/DDCHlsq/sing-ruleset
- **stars / 最近 / 维护**: ~29 · 2026-06 · active
- **license**: 无 (license: null, 仓库无 LICENSE 文件) — 采用层面的硬伤
- **formats**: singbox.srs
- **provenance**: domain-list-community-derived — main.py 仅从 raw.githubusercontent.com/Loyalsoldier/surge-rules/release/ruleset/{name}.txt 下载全部 11 条基础规则 (direct/proxy/reject/private/apple/icloud/google/gfw/tld-not-cn/telegramcidr/cncidr), 而 Loyalsoldier/surge-rules 本身派生自 v2fly domain-list-community。故本库是 Loyalsoldier 的"薄壳重编译器", 传递性地属于 domain-list-community 派生。仅 alibabainc/pcdn/httpdns/steamcn/xbox 等几个小 JSON 是独立手工策展。
- **ai_coverage**: 无任何 AI/LLM 覆盖 — README 及全部 17 个 .srs/6 个 JSON 中均无 ChatGPT/OpenAI/Claude/Gemini/AI 分类。
- **specialty**: 相对我们已知库无独立增量。核心规则就是 Loyalsoldier 的 .srs 重编译。唯一差异化是几个 CN niche 手工列表: Xbox 下载 (global/cn)、Steam-CN (lhie1 归档)、PCDN、HTTP-DNS、Alibaba Inc。这些可作为我们 reject/direct 策略的小补丁参考, 但量极小。
- **reason**: 对抗式核实结论: 真实存在、确为 sing-box 代理分流规则库, 且"规则真的在更新"——但更新仅靠下游自动化, 价值有限。关键证据链: (1) 存在性确认, 2024-03-09 创建, Python 100%, ~29 star, 1 fork, 非 archived。(2) 活跃度需拆两层: master 分支人工开发半停滞 (最后人工 commit 2026-04-29 by DDCHlsq/iamddch); 但编译产物在独立 'ruleset' 分支, 由 GitHub Actions (actions-user/"GitHub Action") 每天 force-push (cron '0 0 * * *' = 08:00 GMT+8), API pushed_at=2026-06-09T00:17Z(今天)。所以 .srs 确实每日刷新——但只是因为上游 Loyalsoldier 每日刷新, 本库自身无增量贡献。(3) 格式: 仅 sing-box .srs (17 个, 在 ruleset 分支); master 只有 6 个 .json 源+main.py。无 clash/surge/quantumultx 等多客户端输出。(4) 血统: main.py 实锤——11 条基础规则 100% 拉自 Loyalsoldier/surge-rules, 即 domain-list-community 传递派生, README 自述"基于 Loyalsoldier 10 条基础规则"属实但低估了依赖程度(其实是全部基础规则)。(5) 无 license(采用阻断)、零 AI 覆盖、相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 严格劣势。结论: 不值得纳入 本项目 视野作为上游源——我们已直接掌握 Loyalsoldier 上游且能自行编译 .srs(扇出能力本就是 本项目 目标), 本库既无 license 又单格式又无 AI。唯一可留作 0.5% 参考的是其 CN niche 手工列表 (Xbox/Steam-CN/PCDN/HTTP-DNS/Alibaba)。

</details>

<details>
<summary>❌ <b>Dracay/ruleset</b> — dead · GPL-3.0</summary>

- **URL**: https://github.com/Dracay/ruleset
- **stars / 最近 / 维护**: 1 · 2025-03 · dead
- **license**: GPL-3.0
- **formats**: clash.list, mihomo.mrs, singbox.srs, singbox.json, geosite.dat, geoip.mmdb, geoip.metadb
- **provenance**: aggregator-of-others (借道 DustinWin) — 决定性证据: 该 repo 全部 commit 的 author/committer 均为 "DustinWin" 而非 Dracay; description 与 homepage(proxy-tutorials.dustinwin.top / dustinwin 教程站) 直指 DustinWin; 分支/tag 结构 (mihomo / mihomo-ruleset / sing-box / sing-box-ruleset / sing-box-ruleset-compatible) 与 DustinWin/ruleset_geodata 完全同构。即它是 DustinWin/ruleset_geodata 的 fork/快照副本, 自身无任何独立策展。其上游 DustinWin 本体血统为多源聚合: 构建在 v2fly/domain-list-community + DustinWin/domain-list-custom + blackmatrix7/ios_rule_script + privacy-protection-tools/anti-AD + felixonmars/dnsmasq-china-list + gfwlist 之上 (IP 侧 DustinWin/geoip + gaoyifan/china-operator-ip + GeoLite2)。所以严格说不是纯 domain-list-community 派生, 而是 dlc 派生 + blackmatrix7 + 自维护 custom 的聚合器。
- **ai_coverage**: 有 AI 分类 (上游 README 含 geosite,ai 🤖 人工智能, 覆盖 OpenAI/Copilot/Gemini/Claude), 但因 fork 冻结在 2025-03, 此处 AI 域名表已陈旧近 15 个月, 无法反映 claude.com/grok/perplexity 等新增域名。要 AI 覆盖应直接看活跃上游而非此 fork。
- **specialty**: 无差异化价值。它等于一份过期的 DustinWin/ruleset_geodata。DustinWin 本体确有特色 (mihomo .mrs/.metadb + sing-box 新旧内核双轨 compatible 版 + 教程站生态), 但都是上游 DustinWin 的, 不是 Dracay 的。相对我们已知库: 多源聚合定位与 Loyalsoldier(同为 dlc 派生) 重叠, 流媒体专项不如 blackmatrix7 细, 全格式扇出工程化不如 Sukka/MetaCubeX。Dracay 这个 fork 只是把上游某天产物冻在 2025-03。
- **reason**: 不值得加入 本项目 视野。对抗式核实结论: (1) 仓库真实存在且确为代理分流 ruleset/geodata 库, 非无关项目。(2) 状态=dead: 默认分支 master 最后 commit 2025-03-14, 三条实际规则数据分支 (mihomo-ruleset / sing-box-ruleset / sing-box-ruleset-compatible) 最后构建 2025-03-18/19 后再无更新; repo pushed_at=2025-03-19; 且 GitHub Actions total_count=0 — 继承自上游的 run.yml 工作流在此 fork 从未运行过, 所谓"每天凌晨3点自动构建"完全是抄来的死自述, 规则已停更约 15 个月。这正是"仓库看着像规则库但规则根本没在更新"的典型陷阱。(3) 它是 DustinWin/ruleset_geodata 的陈旧 fork (commit author 全是 DustinWin, homepage 指向 dustinwin 站, 分支/tag 同构), 零独立策展, 1 star, 0 fork。(4) 真身上游 DustinWin/ruleset_geodata 才是活的 (1264 stars, pushed 2026-06-08, 有自动 build release tag) —— 若要追这条血统应直接纳入 DustinWin 本体, 绝不追这个会污染供应链审计的死 fork。建议: 把 DustinWin/ruleset_geodata 记入候选, 丢弃 Dracay/ruleset。

</details>

<details>
<summary>❌ <b>Dreista/sing-box-rule-set-cn</b> — active · none (无 LICENSE 文件, GitHub 检测不到 license → 默认 all-rights-reserved)</summary>

- **URL**: https://github.com/Dreista/sing-box-rule-set-cn
- **stars / 最近 / 维护**: 71 · 2026-06 · active
- **license**: none (无 LICENSE 文件, GitHub 检测不到 license → 默认 all-rights-reserved)
- **formats**: singbox.srs, singbox v2 .json (rule-set source)
- **provenance**: aggregator-of-others — 依据来自 generate_rule_set.py 源码内显式 URL 列表 (非 README 自述): felixonmars/dnsmasq-china-list (accelerated/apple/google.china.conf), misakaio/chnroutes2, APNIC delegated-apnic-latest, Dreamacro/maxmind-geoip Country.mmdb, IPinfo Lite mmdb, AdGuardSDNSFilter, hagezi/dns-blocklists (pro/pro.mini), gfwlist/gfwlist。纯聚合编译, 非 domain-list-community 派生, 无独立手工策展。
- **ai_coverage**: 无。generate_rule_set.py 中无 openai/anthropic/claude/gemini/netflix/youtube/streaming 任何 AI 或流媒体分类; 范围仅限 CN-direct + adblock + gfwlist。
- **specialty**: CN 直连/中国 IP 聚合专项: 一次性把 dnsmasq-china-list 域名 + 多源 CN IP (APNIC delegated + maxmind + IPinfo Lite, IPv4/IPv6 各一份) + gfwlist + AdGuard/hagezi adblock 编译成 sing-box 原生 .srs。相对 Loyalsoldier/blackmatrix7/Sukka 的唯一差异点是把 APNIC+maxmind+IPinfo Lite 三套 CN-IP 源并排打包进单一 sing-box pipeline; 域名侧无任何新增策展, 全是已知上游。
- **reason**: 对抗式核实结论: 真实存在且确为代理分流规则库 (rule-set 分支含合法 sing-box v2 规则: 真实 CN domain_suffix 与 CN ip_cidr)。活跃度需区分两层 — 数据管线活: rule-set 分支最新 commit = 今天 2026-06-09 由 github-actions[bot] 生成, workflow cron '0 0 * * *' 确为每日自动; 但生成器开发停滞: trunk 最后一次人工 commit 为 2026-03-25 (Kumiko as a Service), 约 2.5 个月无新源/功能, 判定 'rules 在更新但 repo 在维护意义上偏 stale'。不值得纳入 本项目 视野, 原因: (1) 单客户端 — 只产 sing-box .srs, 无 clash/.mrs/surge/quantumult-x/geosite, 与我们多客户端扇出诉求不匹配; (2) 无 license, 再分发法律上有风险; (3) 纯聚合器且所有上游 (felixonmars/chnroutes2/APNIC/maxmind/IPinfo Lite/AdGuard/hagezi/gfwlist) 我们都能自行直接拉取 — 这正是 本项目 '自己掌控供应链' 想消灭的不透明中间层, 引入它反而是反模式; (4) 零 AI/LLM、零流媒体覆盖, 与我们重点专项无交集。唯一可借鉴价值: 其 generate_rule_set.py 的 CN-IP 多源清单 (APNIC delegated + maxmind + IPinfo Lite IPv4/IPv6) 可作为我们自建 CN-IP 规则的源参考 checklist, 但无需追踪该 repo 本体。

</details>

<details>
<summary>❌ <b>Dunamis4tw/generate-geoip-geosite</b> — dead · 无 (无 LICENSE 文件, GitHub API license=null → 默认 all-rights-reserved, 法律上不可再分发)</summary>

- **URL**: https://github.com/Dunamis4tw/generate-geoip-geosite
- **stars / 最近 / 维护**: ~62 · 2024-03 · dead
- **license**: 无 (无 LICENSE 文件, GitHub API license=null → 默认 all-rights-reserved, 法律上不可再分发)
- **formats**: singbox.srs, singbox rule-set.json, geosite.dat (geosite.db), geoip.dat (geoip.db)
- **provenance**: aggregator-of-others — 依据: 它本质是一个通用 Go 生成器 (BYO 列表, input dir 里放 include/exclude-{ip/domain}-{category}.lst)，本身不策展任何域名表。仓库自带的 source*.json 全部指向俄罗斯反审查/广告/BT 黑名单 (zapret-info/z-i Antizapret、Roskomsvoboda/rublacklist、Antifilter、AdAway、SM443 Pi-hole Torrent Blocklist)。与 v2fly/Loyalsoldier domain-list-community 完全无血缘，不是其派生。
- **ai_coverage**: 无。无 AI/LLM 域名分类, 也无流媒体分类; 自带 source 仅 RU 反审查黑名单 + 广告 + BT。claude.com/grok/perplexity 等均不覆盖。
- **specialty**: 相对我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 的唯一差异点是: (a) 它是「生成器工具」而非「规则仓库」——把任意 IP/域名表扇出成 sing-box geoip.db/geosite.db/rule-set .json+.srs; (b) 自带 source 模板专攻俄罗斯 RKN/Antizapret 反审查解封场景 (中文圈代理库基本没覆盖)。但它只产 sing-box 一种生态, 不产 clash/surge/quantumultx/egern, 与我们「一源多扇出」目标方向相反。
- **reason**: 对抗式核实结论: 仓库真实存在且确为代理分流相关工具 (Go 源码 generator.go/downloader.go + custom/ 下样例 geoip.db/geosite.db/.srs/.json 实锤), 但不值得纳入 本项目 视野。三条硬伤: (1) 死库 — 最后 push/release 均为 2024-03-02, 距今约 27 个月零提交; API 上 updated_at=2026-05 只是 star/metadata 触碰, 非代码更新。且它是「生成器」非「规则库」, 规则新鲜度取决于用户自己跑工具拉上游, 仓库本身不产出滚动更新的成品规则集, 即便活跃也不等于规则在更新。(2) 无 license — 无 LICENSE 文件、API license=null, 默认 all-rights-reserved, 我们无法合法再分发其产物或代码。(3) 错配 — 只产 sing-box 单一生态 (geoip.db/geosite.db/rule-set json+srs), 不产 clash/surge/quantumultx/egern, 与我们「上游源→自编译→扇出全部客户端格式」的多扇出目标相反; 自带数据全是俄罗斯反审查/广告/BT 黑名单, 无 AI/LLM、无流媒体, 与中文代理分流场景几乎不重叠。血统上是 RU-censorship 列表的聚合器, 非 domain-list-community 派生, 对我们已有的 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 体系无差异化补充价值。若未来要做 RKN/Antizapret 解封专项, 可把它的 source*.json 思路 (z-i CSV / rublacklist JSON 解析) 当参考实现, 但库本身 (死+无 license+单生态) 不应直接追踪或依赖。

</details>

<details>
<summary>❌ <b>Ekko1048/OpenClashRule</b> — dead · GPL-3.0</summary>

- **URL**: https://github.com/Ekko1048/OpenClashRule
- **stars / 最近 / 维护**: 24 · 2025-03 · dead
- **license**: GPL-3.0
- **formats**: clash.yaml (RULE-SET / rule-provider, classical)
- **provenance**: aggregator-of-others (pointer to blackmatrix7, NOT independent-curated). Verified via recursive git tree: repo holds only LICENSE.md + README.md + rules/zhilian.yaml (63B placeholder: example.com / DOMAIN-KEYWORD,test) + empty 1-byte zhilian.yaml. It hosts ZERO real rule files. The 901-line README references blackmatrix7 106x and jsdelivr 102x, pointing users at cdn.jsdelivr.net/gh/blackmatrix7/ios_rule_script@master/rule/Clash/*.yaml (52 distinct rule-sets). Not domain-list-community derived; not self-curated.
- **ai_coverage**: Only via the blackmatrix7 upstream it points to: README references OpenAI (11x) and Copilot (11x) jsdelivr URLs. NO Anthropic/Claude, NO Gemini/Bard, NO Perplexity/Grok mentions. The repo itself hosts no AI rule files — all AI coverage is whatever blackmatrix7 ships at request time.
- **specialty**: None of differentiating value. It is a README/config tutorial telling v2board/airport users to load blackmatrix7 rules over jsdelivr CDN. Compared to our known libs (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) it adds nothing — it is strictly downstream of blackmatrix7, which we already track. Only original artifact is a stub zhilian.yaml with example placeholders.
- **reason**: Adversarial verification debunks the README. (1) Exists and is proxy-routing-related, but is NOT a real ruleset — recursive git tree shows only README/LICENSE + a 63-byte placeholder rule (example.com/test) + an empty file; zero hosted rules. (2) DEAD: last push 2025-03-18 (14+ months stale as of 2026-06); created 2024-01. The README's headline claim of 'daily GitHub Actions build at 6:30 AM Beijing' is FALSE — there is NO .github directory and no CI/bot commits at all (all 12 commits are manual by Ekko1048). 'Rule freshness' is an illusion produced by jsdelivr pulling blackmatrix7 live, not by this repo. (3) Format: only Clash classical .yaml RULE-SET references (no .mrs/.srs/surge/sing-box/geosite). (4) Provenance: pure pointer/aggregator to blackmatrix7 over CDN — initial 'independent-curated' judgment is wrong; no curation exists. (5) GPL-3.0; AI coverage limited to blackmatrix7's OpenAI+Copilot (no Claude/Gemini); zero differentiation vs libs we already track. Conclusion: NOT worth adding to 本项目 — it is strictly downstream of blackmatrix7 (already in scope), supply-chain-fragile (depends on jsdelivr+blackmatrix7), abandoned, and contains no compilable upstream source of its own.

</details>

<details>
<summary>❌ <b>fernvenue/chn-cidr-list</b> — active · BSD-3-Clause</summary>

- **URL**: https://github.com/fernvenue/chn-cidr-list
- **stars / 最近 / 维护**: ~88 · 2026-06 · active
- **license**: BSD-3-Clause
- **formats**: clash.list (.txt 纯 CIDR), clash rule-provider .yaml (payload:), surge/clash .conf (IP-CIDR,前缀), .mmdb (MaxMind geoip DB), IPv4/IPv6/合并三套各自分离
- **provenance**: aggregator-of-others — 经 .gitlab-ci.yml 实证: 数据 100% 来自 gaoyifan/china-operator-ip (BGP/ASN china.txt+china6.txt) + ftp.apnic.net APNIC delegation stats (grep CN), 用 zhanhb/cidr-merger 合并、carrnot/mmdb-go 转 mmdb。无任何手工策展, 无域名。明确 NOT domain-list-community 派生 (domain-list-community 是域名库, 此库纯 IP)。
- **ai_coverage**: none — 该库不含任何域名规则, 仅 IP-CIDR, 因此完全不涉及 AI/LLM/流媒体域名覆盖 (维度不适用)。
- **specialty**: 纯 CN IP-CIDR 维度, 与我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 主打域名 geosite) 正交。差异点: (1) IPv4/IPv6 严格分离的独立文件; (2) 直接产 .mmdb geoip 二进制, 可省去自编译 geoip.dat 的一道工序; (3) 上游链路极简透明 (仅 APNIC+gaoyifan 两源), 供应链审计成本低。但本质只是 gaoyifan/china-operator-ip 的多格式再封装。
- **reason**: 真实存在且确为代理分流规则库 (Clash/Surge IP-CIDR + Clash rule-provider yaml + mmdb), BSD-3-Clause, 维护真实活跃 (GitLab CI 每日 cron, diff 才提交, 最新 commit 今天 2026-06-09; GitHub 仅为镜像)。但血统是聚合器, 数据全部来自 gaoyifan/china-operator-ip + APNIC, 自身无任何手工策展、无域名规则。对 本项目 「上游源→自编译→扇出」的价值有限: 我们若要 CN geoip, 直接消费它的上游 (gaoyifan china-operator-ip / APNIC raw) 比依赖这层封装更可控、少一跳供应链风险; 它唯一便利是预编译 .mmdb, 但我们既然要自掌控编译链, 这层封装反而该绕过。建议: 不纳入追踪库, 但把它的上游 gaoyifan/china-operator-ip 记为 CN-CIDR 候选源。与已知四大域名库正交, 不构成差异化补充。

</details>

<details>
<summary>❌ <b>GMOogway/shadowrocket-rules</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/GMOogway/shadowrocket-rules
- **stars / 最近 / 维护**: ~4.7k (4727) · 2026-06 · active
- **license**: GPL-3.0
- **formats**: shadowrocket .module (sr_direct_list / sr_proxy_list / sr_reject_list)
- **provenance**: aggregator-of-others (部分 domain-list-community-derived)。依据: README「数据来源」段明列上游 = felixonmars/dnsmasq-china-list + v2fly/domain-list-community + gfwlist/gfwlist + Loyalsoldier/cn-blocked-domain + adblockplus easylistchina+easylist + AdGuard DNS filter + pgl.yoyo.org/adservers + someonewhocares.org/hosts + WindowsSpyBlocker。即 DIRECT/PROXY 来自 china-list+DLC+gfwlist, 巨型 REJECT 来自多个 adblock/hosts 源合并。唯一手工部分是 factory/ 下 6 个 append/excludes .txt 覆盖文件。本质是多源聚合器, 且下游消费了我们已追踪的 Loyalsoldier/DLC。
- **ai_coverage**: 几乎无。扁平 DIRECT/PROXY/REJECT 三桶模型, 零 per-service 分类; PROXY 仅显式补充 telegram/gv/gmail/whatsapp 网段, 无 claude/openai/grok/perplexity 等 AI 域名分类条目。
- **specialty**: 超大体量纯 adblock/REJECT 列表 (16.5w 条 REJECT, 11.3w DIRECT) 的中文优化版, 每日自动构建且规则计数真实波动。差异化价值低: 仅 shadowrocket 单一格式, 扁平 3 桶 (DIRECT/PROXY/REJECT) 无任何 per-service 分类, 且上游全是我们已直接消费的源 (DLC/dnsmasq-china-list/Loyalsoldier)。相对 blackmatrix7(分类多格式)/Sukka/MetaCubeX 严格更弱。
- **reason**: 核实结论: 真实存在、确为代理分流规则库、active(pushed_at 2026-06-08, 文件头 Created:06-09 07:46, 每日 bot release 且各桶 line count 真实增减 = 仓库活跃与规则更新同时成立, 非僵尸)。但初判两处需纠正: (1) 初判称「也产 clash/surge 衍生」—— 本仓库递归 tree 搜索仅含 3 个 shadowrocket .module + 1 个无关 docs/.conf, 无任何 clash/surge/sing-box/geosite 产物; .github/workflows 缺失 (构建在仓库外, 此仓库纯发布目标)。(2) 血统更准确是 multi-source aggregator (部分 DLC 派生), 而非纯聚合。不值得纳入 本项目 视野: 仅 shadowrocket 单一非优先客户端格式 (我们自有 pipeline 已能渲染 shadowrocket), 粗粒度 3 桶无 per-service/AI/流媒体分类, 且严格位于我们已直接掌控的上游 (DLC/dnsmasq-china-list/Loyalsoldier) 下游 —— 对「上游源→自编译→扇出」目标无增量。至多作为 REJECT/adblock 桶的低价值交叉校验参考。

</details>

<details>
<summary>❌ <b>Hackl0us/GeoIP2-CN</b> — active · GPL-3.0 (GitHub metadata; README 另含 MaxMind 商标声明)</summary>

- **URL**: https://github.com/Hackl0us/GeoIP2-CN
- **stars / 最近 / 维护**: ~7.4k · 2026-06 · active
- **license**: GPL-3.0 (GitHub metadata; README 另含 MaxMind 商标声明)
- **formats**: geoip2 mmdb (Country.mmdb ~133KB), ip-cidr txt (CN-ip-cidr.txt ~111KB)
- **provenance**: independent-curated / aggregator-of-others (IP-side). 依据: README + 源码显示数据来自 17mon/china_ip_list(ipip.net) + metowolf/iplist(纯真 Chunzhen) 合并去重, 提取 CN 大陆段, 用自研 ip2cidr.go/dedup.c 构建。与 domain-list-community 无关 (整库零域名)。NOT domain-list-community-derived。
- **ai_coverage**: 无 (zero)。整库只有 IP CIDR, 无任何域名, 因此完全不覆盖 AI/LLM/流媒体域名。
- **specialty**: 精简 CN-only GeoIP2: 只产 mmdb + 纯 IP-CIDR txt, 体积 ~133KB (对比 MaxMind 全球库几十 MB)。是 IP 地理定位层, 与我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 均为域名 geosite/规则库) 正交互补, 不重叠也不竞争。可作为 CN 直连 GEOIP 的轻量底座, 但本身不提供任何分流域名规则。
- **reason**: 对抗式核实结论: 真实存在、活跃 (但需区分两条分支)。陷阱在于 master 分支冻结在 2024-05 仅存构建脚本, 而数据 artifact 在独立 release 分支, 其 HEAD 是 github-action[bot] 2026-06-07 (2 天前) 的 "Updated at Sun Jun 7" 提交 —— 规则真在 3 天周期自动更新, 判 active 而非 stale。但它本质是 CN-only GeoIP/IP-CIDR 库, 非域名分流/geosite 规则库: 产出仅 Country.mmdb + CN-ip-cidr.txt 两个文件, 无 clash .mrs / sing-box .srs / surge .list / geosite.dat, 零域名零 AI 覆盖, 血统是 IP 数据聚合器 (17mon + 纯真) 而非 domain-list-community 派生。本项目 的核心目标是「上游域名源→自编译→扇出所有客户端格式」, 域名规则才是供应链审计主体; CN-GeoIP 这类 IP 定位库属于另一条独立、成熟、低风险的依赖线, 可在需要 CN 直连底座时直接消费 mmdb, 无需纳入域名规则供应链视野。故不值得加入 本项目 追踪 (worth_tracking=false), 但可作为 IP 层备选记一笔。

</details>

<details>
<summary>❌ <b>Hackl0us/SS-Rule-Snippet</b> — dead · AGPL-3.0 (repo metadata); README 另注规则内容 CC BY-NC-SA 4.0 — 含 NC 非商用条款, 对订阅服务是潜在合规风险</summary>

- **URL**: https://github.com/Hackl0us/SS-Rule-Snippet
- **stars / 最近 / 维护**: ~11.3k (11,263) · 2023-01 · dead
- **license**: AGPL-3.0 (repo metadata); README 另注规则内容 CC BY-NC-SA 4.0 — 含 NC 非商用条款, 对订阅服务是潜在合规风险
- **formats**: surge.list, clash.list (classical), clash.yaml, quantumultx.list/.conf/.js, shadowrocket.conf, surge.sgmodule (module), full .conf configs
- **provenance**: independent-curated — 依据: README 明确描述手工抓包+搜索引擎逐域名核验 ("善用抓包工具、搜索引擎...反复斟酌"), 全仓递归 tree 无任何 domain-list-community / v2fly / Loyalsoldier / 上游 import 痕迹; 目录按客户端(Surge/Clash/Quantumult)而非按 geosite 类目组织, 典型独立手工策展库, 非派生非聚合
- **ai_coverage**: 无 (零覆盖) — 全仓递归 grep openai/chatgpt/claude/anthropic/gemini/copilot/perplexity 均无命中; 该库活跃期(2016-2021)早于 LLM 路由需求出现
- **specialty**: 理念是"少而精"的手工精简分流 + 明确拒绝去广告; 但相对我们已知库(Loyalsoldier/blackmatrix7/Sukka/MetaCubeX)无差异化价值: 覆盖面更窄(流媒体仅 Netflix/Disney+/YouTube/Spotify/TVB, social 仅 7 个 app), 且内容停更于 2021-2022, 现已严重过时
- **reason**: 真实存在且确为代理分流规则库(Surge/QX/Shadowrocket/Clash 手工精简规则, 11.3k star), 但对抗式核实判定为 DEAD, 不值得加入 本项目 视野。关键证据: (1) 活跃度是假象 — GitHub updated_at=2026-06 仅为 star/watch 元数据变动; 全部分支权威 commit 日期为 main HEAD 2023-01-19(且该 commit 只把 license 改成 AGPL, 非规则更新), 真正的规则内容 commit 停在 2021-12, dev 分支停在 2016 / rm 停在 2018, pushed_at=2024-03 是悬空 ref 推送未推进任何分支, Changelog 仅到 2018 —— 仓库未 archived 但规则事实上已废弃 3-4 年。(2) AI/LLM 零覆盖, 流媒体仅 5 项且陈旧。(3) 格式只输出 legacy .list/.yaml/classical, 无 clash .mrs / 无 sing-box .srs / 无 geosite .dat, 对我们"自编译扇出多客户端"的供应链无源价值。(4) 血统为独立手工策展(非 domain-list-community 派生), 但覆盖面远窄于 blackmatrix7/Loyalsoldier/Sukka, 无任何独家品类。(5) license 含 CC BY-NC-SA 的 NC 非商用条款, 对订阅服务有合规风险。结论: 历史口碑库, 但已死, 无追踪价值。

</details>

<details>
<summary>❌ <b>jnlaoshu/MySelf</b> — active · None (无 LICENSE 文件, GitHub API license=None, /blob/main/LICENSE 返回 404)</summary>

- **URL**: https://github.com/jnlaoshu/MySelf
- **stars / 最近 / 维护**: 565 · 2026-06 · active
- **license**: None (无 LICENSE 文件, GitHub API license=None, /blob/main/LICENSE 返回 404)
- **formats**: egern.yaml, clash.yaml (rule-provider), surge.list/.conf/.sgmodule, stash, loon, quantumultx, shadowrocket
- **provenance**: aggregator-of-others — README 自述 "网上搜集"，致谢 blackmatrix7/Semporia/mieqq/Repcz 等；merge_rules.py 实时拉取 Repcz/Tool、Cats-Team/AdRules、privacy-protection-tools/anti-AD 并合并去重。非 domain-list-community 派生，非独立手工策展，而是二手聚合 + 整套客户端配置打包。AIGC.yaml 只自引用本仓 raw URL（循环），无真实上游署名。
- **ai_coverage**: 有 — Egern/Rule/AIGC.yaml 覆盖 openai/chatgpt、anthropic/claude、google gemini、microsoft copilot、xAI grok、perplexity、midjourney 等 (~60+ 域名 + domain_keyword_set: openai/chatgpt/anthropic/claude/midjourney/perplexity)。无 claude.com 专项独立文件, 统一归入 AIGC.yaml。
- **specialty**: 差异化价值低: 它是 blackmatrix7/Egern 生态 + anti-AD/Cats-Team 的下游聚合, 不是新上游。相对我们已知库无新增格式 (缺 sing-box .srs/.json、缺 Clash Meta .mrs、缺 v2ray geosite.dat)。唯一可用点: Egern/Rule/AIGC.yaml 是一份整合好的 AI 域名清单 (openai/anthropic-claude/gemini/copilot/grok/perplexity/midjourney ~60+ 域名 + keyword set), 可作 AI 分流域名交叉校验参考; 流媒体 (Netflix/Disney/YouTube/HBO/Spotify/TikTok) 与 AdBlock(每日 Actions 自动合并) 覆盖齐全。但本质是整套自用配置 bundle (含 .sgmodule 模块 + JS 脚本 + 完整 profile), 非纯 rule-provider 库。
- **reason**: 对抗式核实结论: (1) 真实存在, 确为代理分流规则+客户端配置聚合库, 非无关项目。(2) 真活跃 — pushed_at 2026-06-08 (今天前一天), 含 actions-user 每日 cron "Auto-update" 提交, 规则本身在更新而非仅仓库挂活; 2209 commits; status=active。(3) 产出 Egern/Clash(rule-provider yaml)/Surge(.list/.conf/.sgmodule)/Stash/Loon/QuantumultX/Shadowrocket; 关键缺口: 无 sing-box(.srs/.json)、无 Clash Meta .mrs、无 v2ray geosite.dat — 我们 本项目 需要的二进制扇出格式它都没有。(4) 血统=聚合器, 非 domain-list-community 派生; merge_rules.py 实时拉 Repcz/Tool+Cats-Team/AdRules+anti-AD, README 致谢 blackmatrix7 等 — 它在我们已知上游 (blackmatrix7) 的下游。(5) 致命: 无 license, 供应链审计角度没有任何再分发授权; 差异化价值低。结论: 不值得加入 本项目 上游源视野 — 二手、无授权、缺关键格式、无血统优势。若需 AI 域名交叉校验, 直接追它的真实上游 (blackmatrix7 / anti-AD / Cats-Team) 而非本 bundle。

</details>

<details>
<summary>❌ <b>KaringX/karing-ruleset</b> — active · CC-BY-SA-4.0 (workflow 分支 LICENCE 文件为 Creative Commons Attribution-ShareAlike 4.0; GitHub API license 字段未识别为 SPDX, 故 repo 页显示 None)</summary>

- **URL**: https://github.com/KaringX/karing-ruleset
- **stars / 最近 / 维护**: ~603 · 2026-06 · active
- **license**: CC-BY-SA-4.0 (workflow 分支 LICENCE 文件为 Creative Commons Attribution-ShareAlike 4.0; GitHub API license 字段未识别为 SPDX, 故 repo 页显示 None)
- **formats**: singbox.srs, singbox-json (sing-box rule-set source JSON), geoip.dat/geosite.dat (geo/*.dat, 2 files), .list (3 AdGuard-style filter lists)
- **provenance**: aggregator-of-others (依据: workflow/run.yml 的 run.yml 明确 checkout 并合并四个上游: (1) ACL4SSR/ACL4SSR 的 Clash 规则 → convert_json.py + convert_srs.sh 用官方 sing-box v1.12.12 二进制编译成 .srs; (2) MetaCubeX/meta-rules-dat 整个 geo/ 目录直接 cp -r 拿来——这是 geoip/geosite/anthropic/openai/category-ai 的真正来源, 而 meta-rules-dat 本身是 domain-list-community + Loyalsoldier 数据的 sing-box 格式编译产物, 所以 geo 部分是 domain-list-community-derived 的二次聚合; (3) Chocolate4U/Iran-sing-box-rules 伊朗规则; (4) runetfreedom/russia-v2ray-rules-dat + antizapret 俄罗斯规则。自身不策展任何域名, 纯做 fetch→编译→扇出。)
- **ai_coverage**: 较好。geo/geosite 下有 anthropic.srs/.json、openai.srs/.json、perplexity.srs/.json, 以及聚合类 category-ai-!cn / category-ai-chat-!cn / category-ai-cn (含 @ads @!cn 变体)。但这些 100% 来自 MetaCubeX/meta-rules-dat (直接 cp -r), 本仓库没有自行维护任何 AI 域名, 也无独立的 claude.com/grok 等手工策展。AI 覆盖等价于 meta-rules-dat 的 category-ai 系列。
- **specialty**: 把 ACL4SSR 体系 (BanAD/ChinaDomain/ProxyGFWlist/Netflix/Telegram 等) 现成编译成 sing-box .srs, 同时把 meta-rules-dat 全量 geosite/geoip 直接捎带, 并叠加伊朗/俄罗斯专项 (Chocolate4U + antizapret + runetfreedom)。差异化主要是: (a) 现成的 ACL4SSR→.srs 编译 (省了我们自己跑 sing-box compile 把 ACL4SSR Clash 规则转 srs 这一步); (b) 多区域 (IR/RU) 反审查规则一站式聚合, 这是 Loyalsoldier/Sukka/blackmatrix7 不专门覆盖的。recommend/ 下还提供 cn/default/ir/ru 四套开箱即用 rule-set 组合 (-set.json)。但 geo 部分对我们零增量——直接就是 MetaCubeX/meta-rules-dat。
- **reason**: 真实存在且确为代理分流规则库 (fork 自 ACL4SSR, 默认服务于 Karing 客户端)。最近一次 push 2026-06-08 (核实当天前一天), 但需区分: sing 分支是 github-actions[bot] 单 commit 的 orphan/产物分支 (force-push 覆盖, 只有 1 个 "Released on" 提交), 真正的"活跃"取决于上游 ACL4SSR/meta-rules-dat 的更新节奏——本仓库自身只是定时跑 CI 重新编译, 无人工策展, 判定 active(CI 层面) 但规则新鲜度完全继承上游。产出 sing-box .srs + 源 JSON + 少量 .dat/.list, 单一面向 sing-box 生态, 不产 clash .list/.mrs / surge / quantumultx, 与我们多客户端扇出目标不匹配。血统上是纯聚合器: geo 部分=MetaCubeX/meta-rules-dat (=domain-list-community 派生), ACL4SSR 部分=别人 Clash 规则的 .srs 重编译, 我们没有任何新上游源可挖。AI 覆盖虽有 anthropic/openai/perplexity 但全部来自 meta-rules-dat, 对已知 MetaCubeX 库零增量。结论: 不值得加入 本项目 视野——它做的正是我们想自己做的事(聚合+编译扇出), 但上游全是我们已知的库(ACL4SSR/MetaCubeX), 唯一新东西是 IR/RU 反审查规则, 而那与 本项目 受众无关。可作为"ACL4SSR→sing-box .srs 编译 + 多源合并 CI"的参考实现 (run.yml 的 convert_json.py/convert_srs.sh 流水线), 但不作为规则数据源跟踪。CC-BY-SA-4.0 的 ShareAlike 若直接分发其产物还会带 copyleft 传染, 是额外减分项。

</details>

<details>
<summary>❌ <b>Keviin560/Shunt_Rules</b> — active · 无 (license: None; /license 端点与 /LICENSE 均 404) — 法律上不可安全再分发/再编译进我们的产物</summary>

- **URL**: https://github.com/Keviin560/Shunt_Rules
- **stars / 最近 / 维护**: 8 · 2026-06 · active
- **license**: 无 (license: None; /license 端点与 /LICENSE 均 404) — 法律上不可安全再分发/再编译进我们的产物
- **formats**: mihomo .mrs (exact+wildcard 双策略, MetaCubeX mihomo convert-ruleset 编译), Loon .lsr (IP 规则带 no-resolve), mihomo override 配置模板 (mihomo-dns.yaml/.js)
- **provenance**: aggregator-of-others — workflow main.yml (cron 0 20 * * *) 实证: git clone blackmatrix7/ios_rule_script --depth1 + 下载 Loyalsoldier/v2ray-rules-dat(geoip.dat) + Loyalsoldier/domain-list-custom(geosite.dat) + v2fly/domain-list-community(master, 取 AI/game 类), 用 urlesistiana/v2dat 解包 + MetaCubeX mihomo 编译成 .mrs。771 个 .mrs 文件名 (115/12306/2KGames/AbemaTV/AcFun...) 与 blackmatrix7 per-app 命名一一对应。data/upstream_all_tags_snapshot.txt 直接是 domain-list-community 的 category 标签快照 (category-ads-all/category-ai-!cn/geolocation-cn)。所以它既派生 domain-list-community, 又叠加 blackmatrix7+Loyalsoldier, 本质是聚合器+重编译器, 非独立策展。
- **ai_coverage**: 覆盖好且粒度细: 独立文件 Anthropic.mrs / Claude.mrs / OpenAI.mrs+OpenAI_IP.mrs / Gemini.mrs / BardAI.mrs / Copilot.mrs+Copilot_IP.mrs / Civitai.mrs / aiXcoder.mrs / Jetbrains.mrs, 外加聚合 AI_Rules.mrs。README 自称含 ChatGPT/Claude/Gemini/Grok/Copilot 等 20+。但这些 AI 域名本身来自 domain-list-community(category-ai) + blackmatrix7, 非自采。
- **specialty**: 差异化价值有限。相对我们已知库基本是子集再包装: 上游就是 blackmatrix7 + Loyalsoldier + domain-list-community(都已在我们视野)。唯一边际价值是工程模式参考——把多个上游 fan-in 后用 MetaCubeX mihomo 官方 binary 编译成 .mrs(exact/wildcard 双策略)并每日 cron auto-sync 提交, 与 本项目「上游源→自编译→扇出」目标同构, 可作为「单人维护的最小聚合-编译流水线」样本。但它不出 sing-box .srs / surge .list / geosite.dat, 客户端覆盖窄(只 mihomo+Loon)。
- **reason**: 真实存在、确为代理分流规则库(topics: clash-meta/loon/mihomo/proxy-rules/ruleset), 非无关项目。维护真活: created 2026-02-11, 477 commits, 每日 20:00 UTC github-actions[bot] auto-sync, 最新 2026-06-08(距今 1 天); 实测最新 commit 真改规则内容(GeoIP_CN.lsr +211/-62、Global_CN/Stripe/STUN 均有真实 diff), 不是空 touch——属"仓库活跃且规则真在更新"。但不值得纳入 本项目 视野, 三条硬伤: (1) 血统是纯聚合器, 上游 blackmatrix7+Loyalsoldier+domain-list-community 全部已在我们清单内, 它给不出我们拿不到的新数据; (2) 无任何 license, 再分发/再编译进 本项目 产物有法律风险, 与 本项目「自己掌控供应链」诉求冲突; (3) 输出面窄(仅 mihomo .mrs + Loon .lsr), 缺 sing-box/surge/geosite, 与我们多客户端扇出目标不匹配, 且 star 仅 8、单人 fork 0、可持续性弱。可作为「最小聚合-编译流水线工程样本」一次性参考其 main.yml(v2dat 解包 + mihomo convert-ruleset 双策略编译思路), 但不必长期追踪。

</details>

<details>
<summary>❌ <b>kyleduo/Surge-Rule-Snippets</b> — dead · GPL-2.0 (fork 自身标注; 注意上游 Hackl0us 现为 AGPL-3.0, license 已不一致, 直接复用有合规风险)</summary>

- **URL**: https://github.com/kyleduo/Surge-Rule-Snippets
- **stars / 最近 / 维护**: 5 · 2017-07 · dead
- **license**: GPL-2.0 (fork 自身标注; 注意上游 Hackl0us 现为 AGPL-3.0, license 已不一致, 直接复用有合规风险)
- **formats**: surge.list (.conf/.list 片段), shadowrocket 可直接导入的 LAZY_RULES, Potatso/Cross 文本规则片段
- **provenance**: independent-curated — 否。GitHub API 实锤这是 Hackl0us/SS-Rule-Snippet 的个人 fork (fork=true, parent/source=Hackl0us/SS-Rule-Snippet)。fork 与 created_at 同日 (2017-07-28), 最新 commit 2017-07-28, 即 fork 后 0 原创提交, 内容全部继承自上游 2017 年历史。上游本身基于 GFWList + Scomper/Super Liu 转换工具手工策展, 与 domain-list-community 无关。
- **ai_coverage**: 无。2017 年的清单, 不含任何 AI/LLM 域名 (OpenAI/Claude/Anthropic/Gemini 当时均不存在或无需分流), README 与 raw 片段均无 AI 类目。
- **specialty**: 无差异化价值。它是 8 年前 (2017) 的死 fork, 既不如上游 Hackl0us/SS-Rule-Snippet (11.3k star, 2024-03 才停更) 新, 更远落后于我们已知库 (Sukka/Loyalsoldier/blackmatrix7 持续维护)。内容是 GFW 清单/懒人配置/少数 app 优化片段, 明确不含去广告, 规模少而精但已严重过时。
- **reason**: 对抗式核实结论: 仓库真实存在且确为代理分流规则库, 但是 Hackl0us/SS-Rule-Snippet 的个人 fork (API fork=true, parent=Hackl0us)。最近一次实质 commit = 2017-07-28, fork 后零原创提交 (updated_at 2026-04 只是 GitHub 元数据 touch, 非内容更新), 判定 dead。初判"independent-curated"被推翻 — 它不是独立策展, 是聚合/复制上游。初判"Surge/ShadowRocket 实用片段、少而精、不含去广告"属实但已 8 年过时。无 AI/流媒体专项, license 与上游已分叉 (GPL-2.0 vs 上游 AGPL-3.0) 存在合规风险。相对 Sukka/Loyalsoldier/blackmatrix7 无任何增量价值。结论: 不值得加入 本项目 视野。若确实想要这一血统, 唯一应追踪的是其上游 Hackl0us/SS-Rule-Snippet (11.3k star, AGPL-3.0, 2024-03 停更, 本身也已 stale), 但该上游同样在我们已知主流库覆盖范围内, 优先级低。

</details>

<details>
<summary>❌ <b>lateautumn2/ruleset_geodata</b> — stale · GPL-3.0</summary>

- **URL**: https://github.com/lateautumn2/ruleset_geodata
- **stars / 最近 / 维护**: ~10 · 2024-08 · stale
- **license**: GPL-3.0
- **formats**: clash .list (classical), clash/mihomo .mrs, clash .yaml ruleset, sing-box .srs, sing-box .json (source), geosite.dat, geoip.dat, Country.mmdb, geoip.metadb, ASN.mmdb
- **provenance**: aggregator-of-others — 依据: (1) README 明确是 DustinWin/ruleset_geodata 的 fork（页面头 "forked from DustinWin/ruleset_geodata"）；(2) 上游 DustinWin 由 GitHub Actions 自动从多个上游聚合: v2fly/domain-list-community + blackmatrix7/ios_rule_script + Loyalsoldier/clash-rules + privacy-protection-tools/anti-AD + XIU2/TrackersListCollection + MaxMind GeoLite2 + gaoyifan/china-operator-ip + gfwlist。非独立手工策展，非纯 domain-list-community 派生，而是多源聚合再编译扇出。
- **ai_coverage**: 有 AI 专项分类 category-ai，覆盖 Claude / ChatGPT(OpenAI) / Gemini / Copilot。但此 fork 的 AI 规则冻结在 2024-08，已严重过时（缺 2024Q4 之后的 claude.ai/anthropic 新域名、grok/x.ai、perplexity 等新增域名）。上游 DustinWin 的 category-ai 是 daily 滚动更新的活规则。
- **specialty**: 差异化价值低且为负: 这是 DustinWin 的过期 fork（最后 commit 2024-08，维护者主动 freeze 以兼容 sing-box <1.10 并停更 DNS 部分）。真正有价值的是上游 DustinWin/ruleset_geodata（1.3k star，daily 3AM 自动 build，今天 2026-06-09 仍在出新规则，已升级到 sing-box .srs v5 / v1.14 系列）。上游相对 Loyalsoldier 的差异: 一次性同时产出 mihomo .mrs + sing-box .srs(v5) + .dat/.mmdb 全套二进制 + compatible 旧版分支，扇出格式比 Loyalsoldier（偏 clash）更全; 流媒体/AI 分类粒度细（category-ai 含 Claude/ChatGPT/Gemini，Netflix/Disney/Max/Prime/AppleTV+/YouTube/TikTok/Bilibili/Spotify 独立分类）。但 fork 本身无任何新增价值。
- **reason**: 对抗式核实结论: 真实存在、确为代理分流规则库，但**候选本身不值得追踪**。它是 DustinWin/ruleset_geodata 的 stale fork（最后 commit 2024-08-21，~10 star，GPL-3.0，维护者主动声明 freeze + 停更 DNS），规则陈旧约 22 个月，AI/流媒体分类已过时，无任何独立策展或新增价值——纯粹是上游某历史版本的快照。关键纠偏: 初判"aggregator"方向正确但漏了它只是 fork; 真正应进入 本项目 视野的是**上游 DustinWin/ruleset_geodata**（1.3k star，GitHub Actions daily 3AM build，规则文件今天 2026-06-09 仍在更新，sing-box .srs 已到 v5/v1.14，一次扇出 mihomo .mrs + sing-box .srs + geosite/geoip .dat + .mmdb 全套）。上游对我们的供应链价值: 是一个"已经在做我们 本项目 想做的事"的参考实现（多上游聚合→自编译→全客户端格式扇出），但血统是聚合器（依赖 domain-list-community + blackmatrix7 + Loyalsoldier 等），不是我们想要的可独立掌控的一手源。建议: 候选 lateautumn2 丢弃; 改把 DustinWin/ruleset_geodata 作为"扇出流水线参考 + .mrs/.srs 编译方式参考"单独建条目追踪，而非作为一手规则源。

</details>

<details>
<summary>❌ <b>limbopro/Adblock4limbo</b> — active · MIT</summary>

- **URL**: https://github.com/limbopro/Adblock4limbo
- **stars / 最近 / 维护**: ~4.4k · 2026-03 · active
- **license**: MIT
- **formats**: quantumultx (.list reject filter), surge (.sgmodule), rewrite (.conf), userscript (.user.js / .js Tampermonkey+AdGuard), .weblist
- **provenance**: independent-curated (mixed-with-aggregator) — 毒奶/limbopro 个人手工 per-site 去广告+重写策展; README 显式致谢并借用 blackmatrix7 / NobyDa / ConnersHua 的分流+去广告规则, reject 段还聚合 EasyList / Peter Lowe 等公开广告表. 明确不是 domain-list-community / geosite 派生 (README 与 .list 内容均无 geosite/dlc 痕迹).
- **ai_coverage**: 无 AI/LLM 域名覆盖 (无 openai/claude/gemini/perplexity 等路由分组); 主题集中在广告/追踪 reject + 成人/影视站重写.
- **specialty**: 手工 per-site 网页/视频广告去除 + URL 重写 + 油猴脚本, 重点覆盖搜索引擎内容农场、成人/ACG/影视站 (Pornhub/Jable/MissAv/禁漫/绅士漫画 等) 的弹窗·视频·Gif 广告. 这是「广告去除/重写」类资产, 与我们已知的「流量分流 geosite」库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX) 不在同一品类——它解决的是 reject 广告+网页改写, 不是 DIRECT/PROXY 国家·服务路由扇出.
- **reason**: 对抗式核实结论: 仓库真实存在 (~4.4k star, MIT) 且活跃 (最近 commit 2026-03-04, 且有月度「系统自动更新」自动提交, 规则确在更新而非仅仓库 churn)——这点 README 自述属实。但初判类目有偏差: 它本质是「毒奶去网页广告计划」(web ad removal / URL rewrite / 油猴脚本), 不是代理分流/geosite 派生规则库。证据: (1) 产出主体是 .conf 重写 + .sgmodule + .user.js 脚本; (2) 其 Adblock4limbo.list 几乎全是 REJECT 广告·追踪域名 (googlesyndication / histats / 51.la, 聚合 EasyList·Peter Lowe), 而非 DIRECT/PROXY 路由; (3) 血统是独立手工策展 + 借用 blackmatrix7/NobyDa 的少量分流块, 与 domain-list-community/geosite 无关。对 本项目『上游源→自编译→扇出多客户端格式』的核心诉求 (国家/服务/AI 分流的 geosite 源) 没有匹配价值: 无 geosite.dat/.srs/.mrs 产物、无 AI/LLM 覆盖、品类是去广告而非路由。结论: 不纳入 本项目 视野。若 本项目 未来单独做「去广告/重写」侧线 (REJECT 列表 + QX/Surge 重写), 它可作为该侧线的活跃中文站点候选源之一, 但那是另一个 spec, 不属当前供应链审计范围。

</details>

<details>
<summary>❌ <b>lingeringsound/10007_auto</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/lingeringsound/10007_auto
- **stars / 最近 / 维护**: unknown · 2026-06 · active
- **license**: GPL-3.0
- **formats**: hosts (all/reward variants), AdGuard/AdBlock Plus (adb.txt), AdAway whitelist (.prop), Adclose (.rule)
- **provenance**: aggregator-of-others — 但聚合的是去广告 hosts 列表, 非 geosite. credits 列出 大圣净化/ad-wars, yhosts, StevenBlack hosts, oisd, neodevhost, 1024_hosts, hblock, 1Host 等; 明确不含 domain-list-community / v2fly / Loyalsoldier. 因此 NOT domain-list-community-derived。
- **ai_coverage**: 无 (N/A) — 仓库不含任何代理路由/geosite 分类, 没有 claude/openai/grok 等 AI 路由分组; 只把广告/追踪域名打到黑洞。
- **specialty**: 无差异化价值 (对 本项目): 它根本不是代理分流规则库, 而是纯去广告 hosts 聚合器. 仅产出 DNS-blackhole/adblock 格式 (hosts/AdGuard/AdBlock Plus/.prop/.rule), 无任何 clash/sing-box/surge/quantumult-x rule-provider 或 geosite 产出. 唯一沾边的方向是"广告拦截 reject 类别"上游, 但该方向 Sukka 已在我们已知库中且更优、多格式。
- **reason**: 对抗式核实结论: 仓库真实存在且高度活跃 (GitHub Actions 自动每日多次 'hosts update' commit, 最近一次 2026-06-09 即今日, 规则确实在更新而非僵尸活跃), 但它**类目错配** — 是 '自动更新去广告 hosts' 项目, 不是代理分流/rule-provider 规则库. 根目录只有 hosts/.txt(AdGuard)/.prop/.rule + shell 脚本, 无 .yaml/.list/.mrs/.srs/.dat, 不针对任何代理客户端格式, 无 geosite/路由概念. 血统是去广告 hosts 聚合器 (StevenBlack/oisd/neodevhost/yhosts...), 与 domain-list-community 供应链无关. 对 本项目 (上游 geosite → 自编译 → 扇出 clash/sing-box/surge 等客户端) 零贡献, 不应纳入视野. 即便未来 本项目 做 reject 广告拦截分类, Sukka (已知库) 也是更好且多格式的上游。

</details>

<details>
<summary>❌ <b>luestr/ShuntRules</b> — stale · none (无 LICENSE 文件, API license=null, GitHub 404)</summary>

- **URL**: https://github.com/luestr/ShuntRules
- **stars / 最近 / 维护**: ~1.5k (1544) · 2026-06 · stale
- **license**: none (无 LICENSE 文件, API license=null, GitHub 404)
- **formats**: clash.list/.yaml, loon.lsr
- **provenance**: domain-list-community-derived 判定错误。实为 blackmatrix7/ios_rule_script 派生：README 第2行明文 "本项目的数据源来自 ios_rule_script 项目"，668 条规则全部是 ios_rule_script 的逐 app 文件经合并/重打包，无 domain-list-community 痕迹。
- **ai_coverage**: 主流齐全但无长尾：含 OpenAI、Claude、Anthropic、Gemini、BardAI、Copilot、Civitai 独立文件；缺 Grok/xAI、Perplexity、Cursor、HuggingFace、Midjourney、Poe、Mistral。无聚合式 "AI" 大类，每个服务一个单文件（继承自 ios_rule_script 的逐 app 切分）。
- **specialty**: 相对我们已知库几乎无差异化价值。血统与 blackmatrix7 完全重叠（同一上游 ios_rule_script），等于 blackmatrix7 的一个"合并+换 Loon/Clash 双格式"再分发版。唯一卖点是把 ios_rule_script 的 .list/_Domain/_Resolve 三件套合并成单文件、并强调 Loon 端 783 万条规则的匹配性能——这是面向 Loon/Clash 终端用户的便利性优化，对我们"自掌控上游→自编译→扇出多客户端"的供应链目标无增量。我们已追踪 blackmatrix7 即覆盖其全部数据。
- **reason**: 对抗式核实结论：1) 真实存在、确为代理分流规则库（topics: clash/loon/rules, 1544 star, 668 条规则）。2) 关键反转——GitHub 仓库是"空壳"：git tree 仅含 README.md 一个文件，全历史仅 1 条 squash 提交 "同步仓库"（size:0，Link 头确认 page=1 即 last），无 LICENSE/无 workflow/无 release/无 tag。真正的规则文件托管在作者自有 Cloudflare CDN rule.kelee.one（/Clash/*.yaml 与 /Loon/*.lsr），不在仓库里。因此 GitHub 的 2026-06-07 提交日期完全不代表规则新鲜度——规则在 CDN 上不透明更新，无法 pin/diff/自编译，从供应链可审计性看判定为 stale（仓库"活"≠规则可追踪）。CDN 对本环境统一 403（Cloudflare 盾），无法直接枚举/验证 CDN 端文件与更新时间。3) 仅 Clash+Loon 两种格式，缺 Surge/sing-box(.srs)/QuantumultX/Egern/mihomo(.mrs)，与我们多客户端扇出目标不匹配。4) 血统为 ios_rule_script(blackmatrix7) 派生，非 domain-list-community（初判错误）。5) 无 license（再分发法律风险），AI 仅主流覆盖。综合：与已追踪的 blackmatrix7 数据完全同源、格式更窄、且仓库不可自编译——不值得加入 本项目 视野。若想要其"合并大类"思路，直接在我们自编译流程里对 blackmatrix7 源做合并即可，无需引入此中间层。

</details>

<details>
<summary>❌ <b>lyc8503/sing-box-rules</b> — stale · none (无 LICENSE 文件; contents 与 license API 均 404)</summary>

- **URL**: https://github.com/lyc8503/sing-box-rules
- **stars / 最近 / 维护**: ~554 · 2026-04 (最后一次成功规则构建/release 2026-04-03; main 上 2026-04-25 仅 "Clean old releases" 维护性 commit, 非规则更新) · stale
- **license**: none (无 LICENSE 文件; contents 与 license API 均 404)
- **formats**: singbox.srs, singbox geosite.db, singbox geoip.db, singbox geoip-cn.db
- **provenance**: domain-list-community-derived — sing-geosite/main.go 硬编码上游为 "Loyalsoldier/v2ray-rules-dat", 拉取其 latest release 的 geosite.dat (校验 geosite.dat.sha256sum), 再用 sing-box 自带 common/geosite + common/srs 库转码成 .db/.srs; geoip 同理。零自有策展, 纯格式移植。而 Loyalsoldier 本身即 DLC 派生, 故此库是 DLC 的二阶派生。
- **ai_coverage**: 继承 Loyalsoldier 的 category-ai-!cn 体系: 存在 geosite-anthropic.srs、geosite-category-ai-!cn.srs、geosite-category-ai-chat-!cn.srs、geosite-bytedance-ai-!cn.srs 等 (含 OpenAI/Claude/Gemini 等)。覆盖尚可但已冻结在 2026-04-03, 之后 AI 域名变动不会反映 (例如新增 claude.ai/grok 子域不会更新)。
- **specialty**: 把 Loyalsoldier v2ray-rules-dat 一键转成 sing-box 原生 .srs/.db 的成品仓 (geosite 分支 1741 个 .srs)。对我们的差异化价值低: 我们 本项目 的目标正是「自己掌控 上游源→自编译→扇出」, 而此库做的恰是「下载 Loyalsoldier 成品 dat → 转码 srs」这一步——这步用 sing-box 官方 common/srs 库我们能自己做, 不必依赖第三方中转 (反而多一层供应链风险 + 无 license)。它唯一省事处是现成的 srs 转码脚本可参考。
- **reason**: 对抗式核实结论: 真实存在、确为代理分流规则库 (sing-box geosite/geoip, 554 star, 2022 建库)。但「仓库活跃」是假象——这是本次核实最关键发现: 每日 cron "Release" workflow 仍在每晚触发, 却自 2026-04-03 最后一次成功后连续每天 FAILURE (今天 2026-06-09 已确认连续 12+ 天失败), 失败根因极可能是 workflow 用了已于 2021 归档的 actions/create-release@v1 + upload-release-asset@v1 (外加 checkout@v2/setup-go@v2)。所以规则成品已冻结约 2 个月并在静默腐烂, 判定 stale/管线实质已死, 不能只看 pushed_at=2026-04-25 (那只是 "Clean old releases" 维护 commit)。血统上它是 Loyalsoldier(DLC 派生) 的纯格式转码, 无任何独立策展; 只出 sing-box 一种格式 (不覆盖我们要扇出的 clash/surge/quantumultx 等); 无 license。综合: 不值得加入 本项目 视野作为上游源——它处在「成品 dat→srs 转码」这一环, 既不是我们想掌控的真·上游 (DLC/各专项源), 转码这步我们用 sing-box 官方库可自做, 引入它反而多一层无 license、已停更的供应链中转。唯一可借鉴的是其 sing-geosite/main.go 转 srs 的实现思路, 可作为我们自编译扇出 sing-box .srs 的代码参考, 但不作为依赖追踪。

</details>

<details>
<summary>❌ <b>malikshi/sing-box-geo</b> — stale · GPL-3.0 (LICENSE files present in sing-geosite/ and sing-geoip/, copyright nekohasekai/SagerNet; GitHub API reports license:None because it is a per-subdir LICENSE not recognized at repo root)</summary>

- **URL**: https://github.com/malikshi/sing-box-geo
- **stars / 最近 / 维护**: ~59 · 2026-02 (release); master code last touched 2025-07 · stale
- **license**: GPL-3.0 (LICENSE files present in sing-geosite/ and sing-geoip/, copyright nekohasekai/SagerNet; GitHub API reports license:None because it is a per-subdir LICENSE not recognized at repo root)
- **formats**: sing-box geosite.db, sing-box geoip.db, sing-box .srs (per-category geosite, 1700+ files), sing-box .srs (per-country geoip, 265 files)
- **provenance**: domain-list-community-derived — verified by reading sing-geosite/main.go: it pulls upstream geosite.dat from a GitHub release of `malikshi/v2ray-rules-dat`, which GitHub API confirms is a fork of `Loyalsoldier/v2ray-rules-dat`, itself built from `v2fly/domain-list-community` + reject/proxy/cn lists. So it is domain-list-community two hops down a personal Loyalsoldier fork. The generator itself is a near-copy of sagernet/sing-geosite's main.go (imports sagernet/sing-box common/geosite + srs and v2fly/v2ray-core router proto to parse geosite.dat → recompile to .db + .srs).
- **ai_coverage**: Strong, contrary to the initial judgement of "no AI". Verified in the rule-set-geosite branch tree: geosite-anthropic.srs, geosite-openai.srs (+@ads), geosite-google-gemini.srs, geosite-perplexity.srs, geosite-github-copilot.srs, geosite-huggingface.srs, geosite-jetbrains-ai.srs, plus aggregate geosite-category-ai-!cn.srs / geosite-category-ai-chat-!cn.srs / geosite-category-ai-cn.srs. Files are non-empty compiled binaries (category-ai-!cn ~1.4KB). This AI coverage is inherited verbatim from the Loyalsoldier/domain-list-community upstream, not independently curated.
- **specialty**: Pre-compiled sing-box .srs fan-out for the ENTIRE Loyalsoldier/domain-list-community taxonomy (1700+ geosite rule-sets including @cn / @ads / @!cn attribute variants, plus 265 per-country geoip .srs). This saves the step of compiling geosite.dat → .srs yourself if you target sing-box. But it is ONLY sing-box format — no clash/.mrs, no surge, no quantumult-x. So for a multi-client fan-out service it covers just one of our output targets.
- **reason**: Real repo, genuinely a sing-box geo/rule-provider library (verified via GitHub API, workflow, generator source, and rule-set branch tree — not a self-report). BUT it fails the value test for 本项目 on two counts. (1) Maintenance is stale, and the distinction "repo active vs rules actually updating" is the killer here: neither this repo's release.yaml NOR its upstream `malikshi/v2ray-rules-dat` fork's run.yml has a `schedule:` cron — both fire only on push/workflow_dispatch. master code is untouched since 2025-07; the lone 2026-02-08 release was a manual trigger. So rule freshness depends on the maintainer remembering to click a button, unlike Loyalsoldier/MetaCubeX which auto-build daily. For a supply-chain-controlled pipeline this is strictly worse than going to the source. (2) Zero differentiation in DATA — it is a downstream recompile of Loyalsoldier (which we already track) with no independent curation; every domain, including its AI categories, comes from domain-list-community. (3) Its only added value is the sing-box .srs compilation step, but 本项目's whole premise is that WE self-compile upstream → fan out to all clients, so we would compile geosite.dat → .srs ourselves anyway and would also need clash/surge/qumtumult-x outputs it does not provide. Net: it is a stale, single-format mirror of a source we already own upstream of. Note it for awareness as a worked example of "Loyalsoldier → sing-box .srs" tooling, but it does not belong in our active tracking set.

</details>

<details>
<summary>❌ <b>neodevpro/neodevhost</b> — active · MIT</summary>

- **URL**: https://github.com/neodevpro/neodevhost
- **stars / 最近 / 维护**: ~1.3k · 2026-06 · active
- **license**: MIT
- **formats**: hosts (pihole/adaway/hblock), uBlock/AdGuard adblocker filter, dnsmasq.conf, smartdns.conf, raw domain list, clash rule-provider (behavior:domain payload, REJECT 用途)
- **provenance**: aggregator-of-others — README + raw 文件均证实: 聚合 anti-AD + 217heidai/adblockfilters 作为屏蔽源, 再用 AnudeepND/oisd/Energized/Ultimate-Hosts-Blacklist/AWAvenue-Ads 等多份 allowlist 去误杀, 外加自维护 dead-domain list。不是 domain-list-community 派生, 也不是代理分流手工策展。
- **ai_coverage**: 无 AI/LLM 专项 (无 claude.com/grok/perplexity 等分流分类)。它根本不做按服务的分流分类, 只有「屏蔽 vs 放行」二元域名集, 谈不上 AI 域名覆盖。
- **specialty**: 对 本项目 (代理分流供应链) 无差异化价值 —— 它是 anti-AD 赛道, 不是分流赛道。其唯一可能复用点是「广告/追踪域名屏蔽集」(116k 屏蔽域 + 11k allowlist + 5.3k dead-domain), 可作为 REJECT 规则源, 但这块上游已有 anti-AD 本体 / AdGuard / 217heidai / Sukka 的 reject 列表覆盖, neodevhost 本身只是这些源的再聚合, 不提供新数据。
- **reason**: 对抗式核实结论: 仓库真实存在且高度活跃 (FusionPlmH 每日 "Auto Update" 自动 commit, 最新 2026-06-08, ~1.3k star, MIT)。但它【不是代理分流/geosite 派生/rule-provider 规则库】, 而是纯 anti-AD/anti-tracking hosts 项目。证据: (1) README 自我定性 "The Powerful Friendly Uptodate AD Blocking Hosts", 明确范围仅广告/追踪屏蔽; (2) 直接拉 raw `clash` 文件验证 —— 它是 `payload:` 下 116566 条裸域名 YAML, 无 REJECT/PROXY/DIRECT 任何 policy, 无 geosite 分类 (cn/google/streaming/ai), 只是给 clash 当一条 behavior:domain 的屏蔽 rule-provider 用; (3) 血统是 aggregator (聚合 anti-AD + 217heidai + 多份 allowlist), 与 domain-list-community/Loyalsoldier/blackmatrix7 那条「分流」血脉无交集。与我们已知库 (Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 全是分流向, Sukka 另含 reject 集) 相比无新增分流数据, AI/流媒体专项为零。结论: 不纳入 本项目 分流供应链视野; 仅当 本项目 未来扩展到「广告屏蔽 reject 子集」时可作边缘候选, 而即便那时它也只是 anti-AD/Sukka-reject 的下游再聚合, 优先级低于直接追上游。

</details>

<details>
<summary>❌ <b>NSZA156/surge-geox-rules</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/NSZA156/surge-geox-rules
- **stars / 最近 / 维护**: 7 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: surge.conf (RULE-SET, DOMAIN/DOMAIN-SUFFIX), surfboard.conf (same RULE-SET), geoip .txt (per-country, Surge IP-CIDR rule-set)
- **provenance**: aggregator-of-others (二次转换, 非一手策展). 依据: README + 字节级核实——GeoSite 上游 = Loyalsoldier/v2ray-rules-dat (该库本身才是 domain-list-community 派生), GeoIP 上游 = Loyalsoldier/GeoIP. 本库只做 geosite.dat/geoip.dat → Surge RULE-SET 的格式转换 + 去 regex + @attribute 拆分, 不自己维护域名清单. 故严格说是 "Loyalsoldier 的 Surge 格式扇出层", 而非 domain-list-community-derived 一级派生 (初判此处需修正).
- **ai_coverage**: 覆盖良好且为 Loyalsoldier 全量透传: 已核实 release 分支存在 anthropic.conf (含 claude.com / claude.ai / anthropic.com / clau.de / claudeusercontent.com)、openai.conf、perplexity.conf、google-gemini.conf、google-deepmind.conf、cursor.conf、github-copilot.conf, 以及 category-ai-!cn / category-ai-chat-!cn 聚合类 (+ @ads/@!cn 拆分). 注意: AI 域名清单本质来自上游 Loyalsoldier, 非本库自有策展, 新鲜度受上游约束.
- **specialty**: 把 Loyalsoldier 全量 geosite/geoip 转成 Surge/Surfboard 原生 RULE-SET (.conf/.txt) 并按 @attribute (如 @ads/@cn/@!cn) 拆成独立文件, Surge 不支持 regex 故剔除 regex: 规则避免误匹配. 1000+ geosite 文件 + 全套 geoip 国家文件. 相对我们已知库: 它不是新数据源, 而是 "Loyalsoldier→Surge" 的格式桥. 我们 本项目 自编译路线本就要自己做这层扇出, 它的价值主要是 (a) 现成的 @attribute 拆分约定可参考, (b) regex 剔除/DOMAIN-WILDCARD 取舍逻辑可借鉴, 而非作为上游源接入.
- **reason**: 对抗式核实结论 (全部经 gh API 字节级验证, 非轻信 README): 1) 真实存在且确为代理分流规则库——geosite .conf 内容确为 Surge RULE-SET (DOMAIN-SUFFIX,anthropic.com 等), geoip 为国家 .txt. 2) 关键区分 "仓库活跃 vs 规则真更新": main 源码分支停在 2025-04-03 看似 stale, 但 release 分支由 github-actions[bot] 每日推送, 且 Actions "Convert Rules for Surge" 工作流连续每日 success 直到 2026-06-08 (今日), 最新发布 "Released on 2026-06-09"——规则是真·每日更新, 判 active. 3) 产出仅 Surge/Surfboard 两种客户端格式 (无 clash/sing-box/quantumultx). 4) 血统: 初判 "domain-list-community-derived" 需修正为 aggregator-of-others——它转换的是 Loyalsoldier/v2ray-rules-dat (一级派生) 与 Loyalsoldier/GeoIP, 本库是二级格式扇出层, 不持有原始域名策展. 5) GPL-3.0; AI 覆盖齐全但全部继承自 Loyalsoldier. 不值得加入 本项目 视野作为"上游源": 我们已直接掌握 Loyalsoldier 作上游, 而本库只是我们计划自建的 "→Surge 扇出" 那一步, 无新数据增量、star 仅 7、单一客户端、单人项目巴士因子高. 唯一可借鉴点是其 @attribute 拆分命名与 regex 剔除策略, 可作为我们 Surge 扇出实现的参考样本, 但不必长期追踪.

</details>

<details>
<summary>❌ <b>Phoroc/sing-rules</b> — stale · none (无 LICENSE 文件, license=null; 转换产物的合法性继承上游 HaGeZi(GPL-3.0 类) / 1Hosts(GPL-3.0))</summary>

- **URL**: https://github.com/Phoroc/sing-rules
- **stars / 最近 / 维护**: 4 · 2026-06 · stale
- **license**: none (无 LICENSE 文件, license=null; 转换产物的合法性继承上游 HaGeZi(GPL-3.0 类) / 1Hosts(GPL-3.0))
- **formats**: singbox.srs, singbox-ruleset.json
- **provenance**: aggregator-of-others — generate_rule_set.py 直接 curl 两个上游 (hagezi/dns-blocklists 的 wildcard/*-onlydomains.txt 6 档 + badmojr/1Hosts 的 domains.wildcards 4 档), 包成 sing-box ruleset json 后 sing-box rule-set compile 成 .srs。零自有策展、零去重/合并逻辑, 纯格式转换搬运。非 domain-list-community 派生。
- **ai_coverage**: 无任何 AI/LLM 专项域名集 (claude/openai/grok/perplexity 等)。它只产 adblock(广告/追踪/恶意软件)reject 类规则, 不涉及流媒体分流、不涉及 AI 解锁分流。AI 覆盖 = 不适用/无。
- **specialty**: 差异化价值几乎为零。它只是把 HaGeZi(light/normal/pro/proplus/ultimate/tif) 与 1Hosts(mini/lite/pro/xtra) 转成 sing-box .srs 单格式。我们已知库里 MetaCubeX/meta-rules-dat 与 Sukka 都已覆盖 reject/广告且质量更高、格式更全; HaGeZi/1Hosts 上游本身在 本项目 里就该作为「源」直接消费, 没必要经这个无名转换器中转。唯一边际信息: 它演示了「上游 onlydomains.txt → sing-box ruleset json → srs compile」这条最小转换流水线 (workflow + python 30 行), 对我们自编译 .srs 的实现有参考意义, 但不构成「值得追踪的规则库」。
- **reason**: 真实存在, 确为代理分流(adblock reject)规则转换库, 非无关项目。但对抗式核实揭穿了「活跃」假象与一处实质缺陷: (1) 活跃度真相: repos.pushed_at=2026-06-08 看似昨天活跃, 但 default 分支 main 最后人工 commit 停在 2024-01-06 (源码 2 年没动)。每日 cron(0 9 * * *) 只是 git init + force-push 到独立的 rule-set 数据分支(--allow-empty-message, GitHub Action bot)。即「CI 在跑」≠「规则在策展」: 维护者两年未碰逻辑, 属 stale(自动机器人续命, 非真维护)。 (2) 已坏且无人发现的核心证据: 脚本对上游 HTTP 404 无任何错误处理(status_code!=200 直接产出空 domain_suffix:[])。实测上游路径 badmojr/1Hosts/master/mini/domains.wildcards 与 /Pro/domains.wildcards 现已 404(1Hosts 改了目录结构), 导致 rule-set 分支里 1hosts-mini.srs=20 字节、1hosts-pro.srs=20 字节(空 ruleset, json 内 domain_suffix:[]), 但 README 仍照常推荐这两个 URL 给用户。10 档里有 2 档静默失效且每天被 force-push 覆盖, 这正是「仓库活跃 vs 规则真更新」必须分开判的反面教材。HaGeZi 6 档 + 1Hosts Lite/Xtra 2 档(共 8 档)上游路径仍 200, 这部分仍真实更新。 (3) 另一隐患: workflow 用 sing-box v1.8.0-rc.5 编译 .srs, 但 README 宣称 1.9.0+, .srs 二进制格式版本可能与新版 sing-box 不匹配。 (4) 血统: 纯聚合搬运(aggregator), 无自有域名、无去重合并, 非 domain-list-community 派生。 (5) 相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 无差异化价值: 单格式(.srs only)、单领域(adblock)、上游(HaGeZi/1Hosts)本就该在 本项目 直接作为源消费, 无须经此中转。无 AI/流媒体覆盖, 无 license。 结论: 不值得加入 本项目 视野作为「规则源」。唯一可借鉴的是它那条 30 行的 onlydomains.txt→singbox-json→srs-compile 最小流水线写法(可作我们自编译 .srs 的参考样例), 以及它作为「转换器静默 404 产空集」的反面教训, 提醒我们自己的扇出管线必须对上游非 200 fail-closed。

</details>

<details>
<summary>❌ <b>powerfullz/override-rules</b> — active · MIT</summary>

- **URL**: https://github.com/powerfullz/override-rules
- **stars / 最近 / 维护**: ~441 · 2026-06 · active
- **license**: MIT
- **formats**: clash/mihomo (生成的 override .yaml 配置, 非独立 rule-provider), JS/TS 动态覆写脚本 (convert.min.js), 自维护少量 clash .list (ruleset/*.list)
- **provenance**: aggregator-of-others — 依据: 路由主体委托给 mihomo 内置 GEOSITE/GEOIP(背后是 Loyalsoldier/MetaCubeX geo 数据) + 外部 provider(SukkaW/skk.moe, 217heidai/adblockfilters, Loyalsoldier/clash-rules gfw.txt); 自维护仅 src/rule_providers.ts 里指向 ruleset/ 下 ~9 个小众 .list (TikTok/Crypto/EHentai/TruthSocial/Weibo/SteamFix/FCM/2 个 CDN/Filter). 非 domain-list-community 自派生, 也非大规模独立策展, 而是"覆写脚本 + 薄手工层 + 大量上游聚合"
- **ai_coverage**: 弱且为委托式: 仅一条 GEOSITE,category-ai-!cn → AI_SERVICE 分组规则, 依赖上游 geosite category-ai-!cn(Loyalsoldier/MetaCubeX 维护)。无自维护 OpenAI/Anthropic/Claude/Gemini/xAI 域名清单, 无独立 AI 数据可供我们 vendoring
- **specialty**: 差异化价值低: 它是 Mihomo/Sub-Store 的 override 配置生成器(动态按节点国家/地区生成 proxy-group + DNS + TUN), 而不是可被我们消费的上游域名规则源。其规则覆盖几乎全部转嫁给 Loyalsoldier/MetaCubeX/SukkaW(我们已知库)。唯一自有增量是少量小众 .list(TruthSocial/EHentai/Crypto/SteamFix/Weibo/TikTok), 体量小、价值边际
- **reason**: 真实存在且活跃(最新 commit 2026-06-03, v2.4.3 2026-05-23, ~441 stars, MIT), 但对抗式核实后判定不属于 本项目 视野的"上游规则源"。它本质是 Mihomo/Sub-Store 的 override 配置生成脚本(TS→clash YAML), 输出 mihomo-only(README 明确不建议 Stash, 无 sing-box/surge/quantumult-x), 与我们"上游源→自编译→扇出多客户端"的目标方向相反——它是消费端整合器而非源。血统为聚合器: 路由主体委托 GEOSITE/GEOIP(Loyalsoldier/MetaCubeX)+SukkaW+217heidai, 自维护仅 ~9 个小众 .list。AI 覆盖为委托式单条 category-ai-!cn, 无独立 AI 域名资产。"活跃"主要是 TS/依赖/GeoIP 链接 churn, 规则内容自更新有限。相对我们已知库(Loyalsoldier/blackmatrix7/Sukka/MetaCubeX)无可 vendoring 的差异化上游数据, 故不值得追踪; 可作为"竞品 override 脚本设计参考"一次性了解即可

</details>

<details>
<summary>❌ <b>qRuWGQ/rules</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/qRuWGQ/rules
- **stars / 最近 / 维护**: 6 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: clash.mrs (mihomo rule-set), geosite.dat (v2ray), geoip.dat (v2ray), geoip.mmdb (maxmind)
- **provenance**: aggregator-of-others — 依据见 conf/*.yaml + .github/workflows/main.yml: 全部规则数据来自下游已知库——blackmatrix7/ios_rule_script (CN域名+CN IP)、Loyalsoldier/clash-rules (gfw/private/direct)、Loyalsoldier/geoip (private ipcidr + geoip CLI 工具)、217heidai/adblockfilters (adblock)。它 clone v2fly/domain-list-community 仅作为 geosite.dat 的编译工具链 (go run ./ -outputname=geosite.dat)，并 rm -fr ./data/* 后把上游下载的清单灌进去——domain-list-community 是编译器而非数据源。因此初判 'domain-list-community-derived' 不准确：它既非 DLC 派生，也非独立策展，而是 Loyalsoldier+blackmatrix7+217heidai 的薄聚合/重打包器。
- **ai_coverage**: 无。分类只有: CN域名、CN IP、GFW、private、direct、adblock。完全没有 AI/LLM (claude/grok/perplexity)、流媒体、地区分流等专项类别。
- **specialty**: 几乎无差异化价值：唯一可能有用的是它把 mihomo convert-ruleset + Loyalsoldier/geoip CLI + domain-list-community 编译流串成了一个干净的 GitHub Actions 样例 (cron 每12h、yaml 驱动 source 列表、git-auto-commit)。这对我们 本项目「自编译扇出」的 CI pipeline 有参考价值——但仅是工程范式参考，不是规则源。它本身只产出 mihomo/.mrs + v2ray .dat + mmdb，不覆盖 sing-box .srs / surge / quantumult-x。
- **reason**: 对抗式核实结论: 仓库真实存在且确为代理分流规则库 (mihomo .mrs + geoip/geosite.dat + mmdb), GPL-3.0, CI 每12h 自动重编 (cron 0 */12 * * *)，最后一次 'auto update' commit = 2026-06-09 当天，所以「仓库活跃」=true。但「规则真的在更新」并非它的功劳——它零原创策展，所有数据 100% 来自我们已追踪的下游库 (Loyalsoldier×2 + blackmatrix7 + 217heidai)，新鲜度完全由上游决定。它只是个 6 star / 0 fork 的单人重打包器，无 AI/流媒体/地区覆盖，格式也只到 mihomo+v2ray，不扩到 sing-box/surge/qx。差异化价值≈0。结论: 不值得作为「规则源」加入 本项目 视野；唯一可借鉴的是它那套 mihomo convert-ruleset + Loyalsoldier/geoip CLI + domain-list-community 编译器 的 Actions 工作流，可作为我们自编译 pipeline 的工程参考样例 (但应直接对接上游 Loyalsoldier/blackmatrix7，绕过这个中间层)。

</details>

<details>
<summary>❌ <b>RealSeek/Clash_Rule_DIY</b> — active · None (仓库无 LICENSE 文件; README 仅声明'仅供爱好者学习使用', 法律上等同 all-rights-reserved, 不可直接二次分发)</summary>

- **URL**: https://github.com/RealSeek/Clash_Rule_DIY
- **stars / 最近 / 维护**: 148 · 2026-05 · active
- **license**: None (仓库无 LICENSE 文件; README 仅声明'仅供爱好者学习使用', 法律上等同 all-rights-reserved, 不可直接二次分发)
- **formats**: clash.yaml (rule-provider, payload: 格式), clash.mrs (mihomo 编译二进制 rule-set), mihomo overwrite .js 覆写脚本, mihomo config .yaml (整config, 不常更新)
- **provenance**: independent-curated (Sukka 结构血统, 非 domain-list-community 派生). 依据: (1) 目录用 Sukka SukkaW/Surge 的标志性词汇 ip/no_ip/domainset, CDN_domainset.yaml 里是 Sukka 风格的 `+.` domain-suffix 语法 + .mrs 编译产物, 明显照搬 Sukka 的 ruleset 组织法; (2) 但内容是手工再策展 — AI_no_ip.yaml/DNS/各 feature 都是 RealSeek 本人 commit 手改 (commit 全是 'feat: 添加xx规则' 个人编辑, 无 bot/sync 自动化); (3) 完全没有 domain-list-community 的 geosite category 痕迹, 也没出 geosite.dat/.srs/.list, 不是 Loyalsoldier 那条线. 结论: 个人自用规则, 骨架抄 Sukka, 数据自己维护, 非聚合器也非 dlc 派生
- **ai_coverage**: 覆盖且非常新. AI_no_ip.yaml 含 claude.ai/claude.com/clau.de/claudemcpclient.com/claudeusercontent.com/anthropic.com (+auth0/cloudflare cdn 细分), chatgpt.com/chat.com/oaistatic/oaiusercontent + DOMAIN-KEYWORD,openai, x.ai/grok.com + PROCESS-NAME/DOMAIN-KEYWORD grok, perplexity.ai/poe.com/meta.ai/dify.ai/jasper.ai, gemini.google.com/aistudio/makersuite/notebooklm/generativelanguage.googleapis.com/aiplatform, api.jetbrains.ai. AI_ip.yaml 含 Anthropic IP-CIDR 160.79.104.0/21 + IPv6 2607:6BC0::/32 + IP-ASN 399358. 含 2025+ 新域名(clau.de/claude.com/jetbrains.ai), 鲜度优于多数老库
- **specialty**: AI 分流是其相对亮点: 独立的 PROXY/ip/AI_ip.yaml + PROXY/no_ip/AI_no_ip.yaml, 且 IP 段维护到 Anthropic IP-CIDR/IP-CIDR6/IP-ASN 399358 这种细粒度 (比纯域名列表更前沿). 另有 Emby/Stream/Telegram/CDN_domainset/Download_domainset 等专项. 但整体规模小 (mihomo 分支 51 文件 ~4.8MB, 仅 41 yaml), 是'够个人用'量级, 远不及 Sukka/blackmatrix7 的体量与品类广度. 差异化主要在: 一份维护很新的中文个人 AI 域名清单可作为 cross-check 参考源
- **reason**: 真实存在且确为 Clash/mihomo 代理分流规则库 (default branch=mihomo, 2022-09 创建, 最近 push 2026-05-23, 提交是手工 feat/fix 个人编辑 → 判 active 且'规则真的在更新'). 但不建议加入 本项目 视野, 原因: (1) 格式面太窄 — 只产 mihomo 专属的 .yaml(payload:)/.mrs/.js, 不出 .list/.srs/geosite.dat, 对我们'一份源扇出 7 客户端'的供应链目标几乎无复用价值; (2) 血统是抄 Sukka 骨架的个人自用库, 上游事实源仍是 Sukka, 我们直接跟踪 Sukka(已知)即可, 这个是下游再加工; (3) license=None + '仅供学习'声明 → 法律上不可安全二次分发, 对要自编译再扇出的供应链是硬伤; (4) 规模仅个人量级(51 文件). 唯一可留的价值: 它那份维护很新、含 Anthropic IP-ASN 的中文 AI 域名清单, 可作为我们 AI 规则的 cross-check 旁证源(非主源), 但不值得列为 本项目 正式追踪对象

</details>

<details>
<summary>❌ <b>REIJI007/AdBlock_Rule_For_Clash</b> — active · NOASSERTION (双协议: GPL-3.0 + CC BY-NC-SA 4.0,后者 NC 非商业条款,GitHub 无法归一化故显示 NOASSERTION)</summary>

- **URL**: https://github.com/REIJI007/AdBlock_Rule_For_Clash
- **stars / 最近 / 维护**: ~345 · 2026-06 · active
- **license**: NOASSERTION (双协议: GPL-3.0 + CC BY-NC-SA 4.0,后者 NC 非商业条款,GitHub 无法归一化故显示 NOASSERTION)
- **formats**: clash.yaml (rule-provider domain set), clash.mrs (mihomo binary ruleset), txt (plain domain list)
- **provenance**: aggregator-of-others — 依据: 仓库 `Referencing rule sources.txt` 列出 ~300 个上游 URL,host 指纹为 EasyList(125) + AdGuard FiltersRegistry/AdguardFilters(123+10) + uBlockOrigin uAssets(22) + malware-filter(10) + urlhaus.abuse.ch。全部是 EasyList/ABP/DNS 语法的 adblock 过滤器,经 PowerShell 脚本(adblock_rule_generator_*.ps1)规范化去重后编译。与 domain-list-community/geosite 体系零交集,不是其派生,也非独立手工策展,而是纯 adblock 上游聚合器。
- **ai_coverage**: 无 AI/LLM 路由域名覆盖。这是纯 reject 黑名单,不含 claude.com/openai/grok/perplexity 等正向分流域。源列表里唯一 "ai" 命中是 uBlockOrigin uAssets 的 easylist-ai.txt,属于"屏蔽 AI 网站上广告"的 cosmetic 过滤器,与我们关心的 AI 站点分流/解锁完全无关。
- **specialty**: 单一垂直: 广告/隐私/恶意域名 REJECT 拦截。聚合 EasyList+AdGuard+uBO+malware-filter 全家桶,体量巨大(yaml/txt 各 ~15MB,数十万域名),覆盖深度超过 blackmatrix7/Loyalsoldier 里的 reject 子集。但只产出 Clash/mihomo 三种格式,且仅 reject 单类目——没有 geosite 式的国别/服务/流媒体分流,没有 sing-box .srs / surge / quantumultx 等多客户端扇出。
- **reason**: 对抗式核实结论: 仓库真实存在且确为代理(Clash/mihomo)规则库,非无关项目。维护 active 且为真更新——commit 历史全部是 github-actions[bot] 自动提交,实际 .yaml/.mrs/.txt 文件每 30-90 分钟滚动变更(非仅 timestamp 摆拍),README 的"20 分钟更新"基本属实(由 GitHub Actions 定时拉上游重编译)。但对 本项目 价值有限,理由: (1) 血统是纯 adblock 上游聚合器(EasyList+AdGuard+uBO+malware-filter),我们若要自掌控供应链,直接订阅这些一线上游 + 自编译即可,无需经它中转,徒增一层不可控的第三方聚合节点;(2) 单类目(仅 reject 广告/隐私/恶意),无国别/服务/流媒体/AI 分流,与我们已有的 blackmatrix7/Loyalsoldier/Sukka/MetaCubeX 体系正交但价值低——它们的 reject 子集已够用,且 Sukka(sukkaw/ruleset)的 reject 列表同源 + 工程化更强;(3) 只产 Clash 三格式,不做多客户端扇出,正是我们 本项目 要自己解决的环节,它帮不上;(4) 供应链审计红线: 仓库把一个 49MB 的 mihomo.exe 二进制 blob 提交进 tree 用于编译 .mrs,且许可证含 CC BY-NC-SA 的 NC 非商业条款——对商业订阅服务有法务风险。综合: 可作为"adblock 上游清单参考"瞄一眼(它整理的 300 个 URL 是不错的 reject 源 checklist),但不值得作为规则库纳入 本项目 直接依赖。

</details>

<details>
<summary>❌ <b>REIJI007/AdBlock_Rule_For_Sing-box</b> — active · GPL-3.0 + CC BY-NC-SA 4.0 双许可 (GitHub 识别为 "Other"); 含 NC 非商业条款 — 对付费订阅服务是合规风险点</summary>

- **URL**: https://github.com/REIJI007/AdBlock_Rule_For_Sing-box
- **stars / 最近 / 维护**: ~139 (sing-box repo); sibling Clash repo ~345, Collection ~103 · 2026-06 · active
- **license**: GPL-3.0 + CC BY-NC-SA 4.0 双许可 (GitHub 识别为 "Other"); 含 NC 非商业条款 — 对付费订阅服务是合规风险点
- **formats**: singbox.srs, singbox json rule-set, plain domain txt, (sibling repo) clash .yaml/.mrs/.txt, (sibling repo) raw adblock-filter ABP/uBO/AdGuard/hosts/dnsmasq syntax
- **provenance**: aggregator-of-others — 读 build 脚本 adblock_rule_generator_json.ps1 确认: 抓取约330个上游 URL (AdGuard filter registry / EasyList+EasyPrivacy / uBlock Origin assets / AdAway / Peter Lowe / URLhaus+phishing+malware-filter), 解析 ABP ||domain^ / @@ 白名单 / hosts / dnsmasq / 裸域名, PSL 校验 + HashSet 去重 + 5级白/黑名单优先级, 输出 sing-box domain/domain_suffix。脚本中**完全没有** geosite / domain-list-community / Loyalsoldier 引用, 非 DLC 派生, 是对西方 adblock 过滤器列表的独立聚合。
- **ai_coverage**: 无 AI/LLM 分类覆盖 (纯 reject, 无 claude.com/openai/grok/perplexity 路由域名表)
- **specialty**: 广告/追踪/恶意域名 reject 专项 (DNS 级)。差异化价值低: Sukka(sukkaw/Surge) 已提供更强、更透明、多客户端的 reject 集; Loyalsoldier/blackmatrix7/MetaCubeX 覆盖 geosite + 分类路由本仓库完全没有。唯一亮点是 20 分钟自动重建 + sing-box srs 原生产出, 但属 cosmetic churn。
- **reason**: 真实存在且活跃 (cron */20 验证, 2026-06-09 两次 release 仅隔37分钟, 非"仓库被碰一下"假活跃)。但对 本项目 价值低: (1) 纯 adblock reject 专项, 非 geosite/分类路由, 无 AI/流媒体专项, 与我们关心的 AI 路由维度零贡献; (2) 血统是西方 adblock 过滤器聚合器, 与 Sukka 高度重叠且不如其透明/可审计; (3) workflow 周期性 Clear Git History + Delete All Releases, 强制 force-push, 无干净历史审计轨迹 — 对供应链审计是减分项; (4) CC BY-NC-SA 的 NC 条款对付费订阅服务有合规风险。结论: 不值得加入 本项目 视野; 若未来确需托管式 reject 源, 优先 Sukka, 此库仅作"已知存在但不追踪"备注即可。同 org 的 AdBlock_Rule_For_Clash / Adblock-Rule-Collection 是同一聚合管线的不同格式扇出, 价值判断相同。

</details>

<details>
<summary>❌ <b>Repcz/Tool</b> — active · none (license: null; README 含"禁止转载/发布至国内平台"+反 fork 免责声明 — 无明确开源授权,默认 all-rights-reserved,二次分发受限)</summary>

- **URL**: https://github.com/Repcz/Tool
- **stars / 最近 / 维护**: ~976 · 2026-06 · active
- **license**: none (license: null; README 含"禁止转载/发布至国内平台"+反 fork 免责声明 — 无明确开源授权,默认 all-rights-reserved,二次分发受限)
- **formats**: surge.list, clash.list, mihomo.list/.yaml (no .mrs), singbox.json, singbox.srs (compiled), loon.list, quantumultx.list, shadowrocket.list, stash, egern.yaml, surfboard, lancex, geoip.mmdb (Loyalsoldier passthrough), js scripts/rewrites/MitM, client config profiles
- **provenance**: aggregator-of-others — 决定性证据来自 .github/workflows/Build.yml: 用 curl download_and_merge 从固定上游列表拉取再扇出。上游包括 Sukka/skk.moe(reject/cdn/ai/apple_cn/domestic/stream)、blackmatrix7/ios_rule_script(Apple*/HBO/Disney/Spotify/PayPal/Steam)、ACL4SSR(OpenAI/Claude/Netflix/YouTube/Google)、Loyalsoldier/geoip+surge-rules/cncidr、ConnersHua/RuleGo、limbopro、TG-Twilight/AWAvenue、NobyDa、dler-io、VirgilClyne。原创层仅 Surge/Custom/*.list 约 20 个手工覆盖(AI/xAI/DeepSeek/Emby/Crypto/Telegram/Porn/Talkatone/TronLink)。非 domain-list-community 派生,非独立手工策展。
- **ai_coverage**: 较强:专设 AI/OpenAI/Claude/xAI(覆盖 Grok)规则集 + Custom/DeepSeek.list;AI 数据主要来自 Sukka ai.conf + ACL4SSR + 自家 Custom 覆盖。未单列 Perplexity(疑并入 AI)。
- **specialty**: 真正差异化点不是规则数据(数据本身就是 Sukka/blackmatrix7/ACL4SSR 的二次聚合,我们已知),而是它的"打包/扇出工程":一套 GitHub Actions 把多上游合并→去噪(删 skk.moe DOMAIN 占位)→排序→同时输出 8+ 客户端格式,并 sing-box rule-set compile 出 .srs 二进制。可作为 本项目「自编译→多客户端扇出」流水线的参考实现/对照样本,而非作为我们追踪的上游数据源。另含 Egern 这种小众客户端格式映射可参考。
- **reason**: 真实存在、确为代理分流规则库、维护非常活跃(每日 2 次自动构建,今日仍在更新,非僵尸 repo)。但对抗式核实后:它是纯聚合器,数据源 100% 是我们已经覆盖的库(Sukka/blackmatrix7/ACL4SSR/Loyalsoldier),不提供任何新的上游血统或独家域名数据,因此作为「上游数据源」零增量。两个硬性减分:(1) license=null + README 明令禁止转载/发布到国内平台 + 反 fork,法律上不适合纳入我们自编译再分发的供应链;(2) 它做的正是我们 本项目 想自建的"合并→去噪→扇出 8 客户端→编译 .srs"那一层,属于竞品/参考实现而非数据源。结论:不值得作为追踪源加入 本项目 视野;但其 Build.yml 流水线(尤其 sing-box .srs 编译 + Egern/多客户端格式映射)值得作为工程参考样本归档一次,无需持续跟踪。

</details>

<details>
<summary>❌ <b>ruijzhan/chnroute</b> — active · none (无 LICENSE 文件, 无 SPDX, GitHub license=null → 默认 all-rights-reserved)</summary>

- **URL**: https://github.com/ruijzhan/chnroute
- **stars / 最近 / 维护**: 290 · 2026-06 · active
- **license**: none (无 LICENSE 文件, 无 SPDX, GitHub license=null → 默认 all-rights-reserved)
- **formats**: routeros.rsc (CN.rsc/CN_mem.rsc/LAN.rsc/gfwlist_v7.rsc), dnsmasq.conf (03-gfwlist.conf), gfwlist.txt (plain domain list)
- **provenance**: aggregator-of-others — 直接拉两个上游: CN IP 来自 felixonmars/dnsmasq-china-list (generate_cn.sh 内 wget accelerated-domains.china.conf, README 自称 iwik.org 是错的), 域名来自 gfwlist/gfwlist (经 gfwlist2dnsmasq.sh 转换)。非 domain-list-community 派生, 非独立手工策展。
- **ai_coverage**: gfwlist.txt 含约 10 个 AI/LLM 域名 (anthropic.com, claude.ai, claude.com, openai.com, x.ai, grok.com, perplexity.ai, huggingface.co, githubcopilot.com, copilot.microsoft.com) — 但仅因上游 gfwlist 恰好包含, 非独立策展的 AI 分类。流媒体同理 (disneyplus/hbomax/spotify 等零散出现)。
- **specialty**: CN IPv4 CIDR + gfwlist 域名双输出, 但只产 RouterOS/dnsmasq 路由器固件格式, 与我们 7 个代理客户端格式 (clash/sing-box/surge/shadowrocket/egern/quantumultx/v2ray) 完全不交集。相对已知库无差异化策展价值: 两个上游 (dnsmasq-china-list、gfwlist) 我们本就可直接消费。
- **reason**: 真实存在且确为代理分流规则库, 仓库 active (每日 GH Actions cron 0 21 * * *, 最近 commit 2026-06-06), 且经 commit diffstat 核实规则内容真在变 (CN.rsc/gfwlist.txt 有真实增删, 非 no-op)。但对 本项目 无价值, 三点硬伤: (1) 产出全是 RouterOS .rsc / dnsmasq .conf 路由器固件格式, 零代理客户端格式 (无 clash/.mrs、sing-box .srs、surge.list、geosite.dat、quantumultx), 与我们 7 客户端扇出目标完全错位; (2) 无 license (默认保留所有权利), 再分发有法律风险; (3) 血统是纯聚合器, 上游 felixonmars/dnsmasq-china-list + gfwlist/gfwlist 我们可直接抓, 它不增加任何相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 的差异化策展。AI 覆盖只是上游 gfwlist 顺带, 非卖点。建议不加入视野。

</details>

<details>
<summary>❌ <b>runetfreedom/geodat2srs</b> — dead · GPL-3.0</summary>

- **URL**: https://github.com/runetfreedom/geodat2srs
- **stars / 最近 / 维护**: 41 · 2024-10 · dead
- **license**: GPL-3.0
- **formats**: singbox.srs (output only), input: v2ray geoip.dat, input: v2ray geosite.dat
- **provenance**: independent-curated (实为工具, 非规则库): repo tree 只有 Go 源码 (main.go/geoip.go/geosite.go/lib.go + proto), 零 .dat/.srs/domain 数据文件, 不携带任何域名清单, 因此既非 domain-list-community 派生也非聚合器, 而是 owner runetfreedom 独立写的转换 CLI。同 owner 的 russia-v2ray-rules-dat / russia-blocked-geosite/geoip 才是真规则库(俄罗斯 RKN 封锁清单), 但本 repo 只是把任意 .dat 转 .srs 的管道工具。
- **ai_coverage**: none — 本身不含任何域名数据 (纯转换工具, repo 内无 list/dat/srs 数据文件), 无 AI/LLM/流媒体专项覆盖可言
- **specialty**: 几乎没有差异化价值: 它不产出/不携带任何规则数据, 只做 geoip.dat+geosite.dat → sing-box .srs 单向转换。功能与 sing-box 官方 `sing-box rule-set compile/convert`、MetaCubeX geo 工具链重叠, 且更老更停滞。对我们 本项目「上游源→自编译→扇出多客户端」的扇出环节无帮助 (只产 .srs, 不产 clash/surge/qx/egern)。
- **reason**: 对抗式核实结论: (1) 真实存在但确认为「工具非规则库」—— GitHub API + git tree 显示仓库只含 Go 源码与 proto, 无任何 .dat/.srs/域名数据, README 仅 222 字节 usage。初判「工具(非规则库)」正确。(2) 活跃度=dead: created_at 与 pushed_at 均为 2024-10-14, 全部 3 个 commit 都是同日的 "Initial commit", 0 release/0 tag; 仓库页 updated_at=2026-05-22 是 star/watch 等元数据事件而非代码推送, 易误判为活跃 —— 实际代码 ~20 个月未动。且因它本身不产规则, 不存在「规则在更新」一说。(3) 格式: 仅 v2ray .dat → sing-box .srs 单向单目标, 无多客户端扇出。(4) 血统: 独立工具, 非 dlc 派生/非聚合; owner 是俄罗斯反审查生态作者, 其真规则在 russia-* 系列(那些才活跃, pushed 2026-06-08)。(5) GPL-3.0; 无 AI/流媒体覆盖; 相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 在数据轴零价值, 在工具轴与 sing-box 官方 rule-set compile 重叠且更停滞。(6) 不值得加入 本项目 视野: 既无规则数据可作上游源, 又无多客户端扇出能力, 工具功能已被 sing-box 官方/MetaCubeX 更好覆盖。如果将来要关注 owner, 应直接看其 russia-blocked-geosite/geoip 真规则库 (区域专项: 俄罗斯 RKN 封锁), 而非这个转换器。

</details>

<details>
<summary>❌ <b>savely-krasovsky/antizapret-sing-box</b> — stale · MIT</summary>

- **URL**: https://github.com/savely-krasovsky/antizapret-sing-box
- **stars / 最近 / 维护**: ~349 · 2026-03 · stale
- **license**: MIT
- **formats**: singbox.srs, singbox-ruleset.json, geosite.db (sing-box legacy), geoip.db (sing-box legacy)
- **provenance**: independent-curated (machine-derived generator). Go tool that pulls Roskomnadzor blocklist CSV dumps from zapret-info/z-i (+ Antizapret lists) and compiles them into sing-box rulesets. NOT domain-list-community / Loyalsoldier / v2fly derived (README + go.mod show no v2fly dependency, no category taxonomy). Not an aggregator of other rule repos — it consumes a single upstream government blocklist source. So "independent" in the sense of its own pipeline, but the underlying data is a flat RKN dump, not hand-curated categories.
- **ai_coverage**: No dedicated AI/LLM category. It does not curate by service. AI/LLM domains (claude.com, openai, etc.) would appear only incidentally if RKN happens to block them; there is no claude.com/grok/perplexity-specific maintained coverage. Effectively unusable as an AI-routing source.
- **specialty**: Single-purpose Russia censorship-bypass ruleset (Roskomnadzor blocked domains/IPs) for sing-box only. Its only real differentiator vs our known libs is the RU-blocklist coverage — a geographic niche none of Loyalsoldier/blackmatrix7/Sukka/MetaCubeX target. But it has NO category structure (no streaming/ads/AI separation): it is one undifferentiated "blocked-in-Russia" list.
- **reason**: Adversarial verdict: real project (Go, MIT, ~349 stars), but FAILS our 本项目 needs on three axes. (1) Activity trap: the repo looked "active in March 2026", but the single most-recent commit (2026-03-26) is literally `feat: disable scheduled job` which comments out the daily `cron: 0 0 * * *`. The last auto-release `20260326002123` is therefore the last data refresh — rules have NOT updated since 2026-03-26 (~2.5 months stale as of 2026-06). Classic 'repo activity != rules updating'. (2) Format mismatch: sing-box ONLY (.srs/.json/.db). No clash/.mrs, no surge, no quantumult-x, no v2ray geosite.dat — we'd get one of seven client formats. (3) Provenance/scope mismatch: it is a flat Russia/Roskomnadzor blocklist machine-derived from zapret-info/z-i CSV, with zero category taxonomy and zero AI/streaming curation, so no differentiated value over Loyalsoldier/blackmatrix7/Sukka/MetaCubeX for our multi-client categorized fan-out. Only conceivable value is RU-geo censorship niche, which is irrelevant to our use case. Do not add to 本项目 vision.

</details>

<details>
<summary>❌ <b>senshinya/singbox_ruleset</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/senshinya/singbox_ruleset
- **stars / 最近 / 维护**: ~469 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: singbox.json (source rule-set), singbox.srs (compiled binary)
- **provenance**: aggregator-of-others — main.py 下载 blackmatrix7/ios_rule_script/archive/refs/heads/master.zip 后机械转码; 非 domain-list-community 派生, 非独立手工策展; 唯一附加处理是拉 MaxMind ASN CSV 做 IP-ASN 展开
- **ai_coverage**: 继承自 blackmatrix7, 无独立 AI 策展。已确认存在 Anthropic 文件夹 (Anthropic.json + Anthropic.srs); blackmatrix7 上游同时含 OpenAI/Claude/Gemini/Copilot 等, 故覆盖等同 blackmatrix7
- **specialty**: 相对已知库差异化价值低: 它只是 blackmatrix7 的单目标(仅 sing-box).srs/.json 转码器, 而我们已知的 MetaCubeX/meta-rules-dat 同样源自 blackmatrix7 且同样产出 sing-box .srs(并额外产 clash/mihomo YAML), 覆盖面更广。它唯一边际优势是 rule/ 下扁平 per-service 文件夹(115/Amazon/Anthropic/AppleMusic...) 1:1 镜像 blackmatrix7, 便于挑单个服务; 每个文件夹自带 jsDelivr CDN + GitHub raw 链接
- **reason**: 对抗式核实结论: 仓库真实存在且确为代理分流规则库; 维护是真活跃(最近 commit 2026-06-08 由 github-actions[bot], Build RuleSet 工作流 660 次运行近期全绿 55-76s, 每日 04:00 同步真在跑成功——是 bot 驱动的真实新鲜度而非僵尸仓库); 产出仅 sing-box(.json 源 + .srs 编译, main.py 产 JSON、compile.sh 跑 sing-box rule-set compile 出 .srs), 无 clash/.list/.mrs/surge/quantumultx。血统是纯聚合器(blackmatrix7 下游转码器), 非 dlc 派生、非手工。不值得加入 本项目 视野: (1) 它是 blackmatrix7 的下游单目标转码器, 我们要掌控上游应直接 ingest blackmatrix7 本体; (2) 同源同目标的 MetaCubeX/meta-rules-dat 已在我们已知清单且更全(多 clash/mihomo YAML)。仅可作为 'blackmatrix7→.srs 编译流水线' 的参考实现留底, 不作为源追踪。

</details>

<details>
<summary>❌ <b>soffchen/GeoIP2-CN</b> — active · GPL-3.0 (上游 chnroutes2 数据 CC-BY-SA 4.0)</summary>

- **URL**: https://github.com/soffchen/GeoIP2-CN
- **stars / 最近 / 维护**: ~247 · 2026-06 · active
- **license**: GPL-3.0 (上游 chnroutes2 数据 CC-BY-SA 4.0)
- **formats**: Country.mmdb (MaxMind GeoIP2, 供 Surge/Shadowrocket/QuantumultX/Clash 走 GEOIP 匹配), clash-ruleset.list (Clash IP-CIDR behavior=ipcidr), clash-rule-provider.yml (Clash rule-provider payload), surge-ruleset.list (Surge IP-CIDR list), CN-ip-cidr.txt (纯 CIDR 列表, 供 iptables/ipset/squid/gost), chnroute.ipset (ipset restore 格式)
- **provenance**: independent-curated (实为 BGP/路由聚合, 非域名策展)。证据链: 仓库 fork 自 Hackl0us/GeoIP2-CN, master 分支 main.go 仅做 chnroutes2 -> mmdb/CIDR 转换; 上游唯一数据源 = misakaio/chnroutes2, 而 chnroutes2 是从 AS917/AS906/AS131477/AS138195 等 BGP route collector 聚合的中国大陆 IP 段 (CC-BY-SA 4.0), 与 domain-list-community / geosite 域名生态完全无关。所有产物文件抽样均为 `IP-CIDR,<cidr>` 或裸 CIDR, 无任何 domain 规则。属于 GeoIP(IP 维度) 而非 geosite(域名维度) 派生。
- **ai_coverage**: 无 (N/A)。该库零域名内容, 不含 OpenAI/Claude/Grok/Perplexity 等任何 AI/LLM 域名, 也不含流媒体域名分类 — 它只输出中国大陆 IP CIDR。
- **specialty**: 差异化价值低且品类不同: 这是一个纯「中国大陆 IP CIDR / GeoIP2 mmdb」库, 解决 GEOIP,CN 这一格子, 与我们已知的 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 这些「域名分流规则(geosite/AI/流媒体)」库不是同一品类, 不构成域名规则的替代或补充。唯一卖点是 IP 段「更小巧 + 每小时实时(基于 BGP chnroutes2)」, 比 MaxMind 官方 GeoLite2 CN 段更新更勤、体积更小。但同类 GeoIP-CN 镜像极多(pmkol/mac-zhou/apadpro/mangoclover 等), 可替代性强。
- **reason**: 真实存在且确为代理分流相关库, 但定位是「IP 维度的 GeoIP-CN」而非 本项目 关心的「域名分流/geosite 派生库」。对抗式核实结论: (1) 存在性确认 — fork 自 Hackl0us/GeoIP2-CN, 247 stars, GPL-3.0, 未 archived。 (2) 活跃度真实 active 而非"仓库假活" — release 分支最新 commit 即今天 2026-06-09 03:45 UTC 的自动 "Updated at..." bot 提交, Actions run #46814 在跑, workflow cron `30 * * * *` 每小时, 上游 chnroutes2 也是每小时 BGP 更新, 数据确实在动。 (3) 格式: mmdb + Clash/Surge IP-CIDR list + rule-provider.yml + ipset + 裸 CIDR。 (4) 血统: 非 domain-list-community 派生, 而是 misakaio/chnroutes2 (BGP route collector 聚合 CC-BY-SA 4.0) 的下游转换器 — 抽样所有产物均为 IP-CIDR, 印证 README 自述无虚, 但本质是 IP 聚合非域名策展。 (5) 无 AI/LLM/流媒体域名覆盖, 与我们四大已知域名库无重叠也无互补。 不值得纳入 本项目「上游域名源 -> 自编译 -> 扇出客户端」的视野: 它不提供域名规则; 我们若需要 GEOIP,CN 段, 上游应直接锁定 misakaio/chnroutes2 (一手源, 避免这一层 fork 噪音), 且同类 GeoIP-CN 镜像泛滥, 单独追踪此 fork 无收益。

</details>

<details>
<summary>❌ <b>StevenBlack/hosts</b> — active · MIT</summary>

- **URL**: https://github.com/StevenBlack/hosts
- **stars / 最近 / 维护**: 30.5k · 2026-06 · active
- **license**: MIT
- **formats**: /etc/hosts (plaintext 0.0.0.0 domain) ONLY — no proxy client format
- **provenance**: aggregator-of-others — 依据: README + commit log 明确 consolidate ~15 个第三方 hosts 源 (AdAway, MVPS, URLHaus, someonewhocares.org, KADhosts, hostsVN, yoyo.org 等), 每次 release commit msg 写明 "Updates from URLHaus, someonewhocares.org, and KADhosts"。与 V2Ray domain-list-community / geosite 生态无任何血缘关系, 不是其派生
- **ai_coverage**: 无 AI/LLM 路由覆盖 (claude.com/grok/perplexity 等)。本质是黑名单, 不含"路由到代理"的 AI 域名集; 这类项目的目标是封堵而非分流
- **specialty**: 对 本项目 无差异化价值。它是 DNS-sinkhole 广告/恶意域名黑名单 (可选 porn/social/gambling/fakenews 扩展, 31 variant), 不是代理分流/路由策略规则库。语义是"封堵这些域名", 不是"走代理/直连"路由决策, 方向相反 (它会 block social, 而我们要 route)。我们已知的 Sukka (sukkaw/surge) 已经在更上游消费同类 ad/malware hosts 源并直接产出真正的 clash/surge/sing-box 格式, StevenBlack 提供的合并域名集已被覆盖, 且自身不产出任何我们需要的客户端格式
- **reason**: 对抗式核实结论: 仓库真实存在且高度活跃 (GitHub API 确认: 未 archived, 最近 push 2026-06-07, latest release 3.16.88 @ 2026-06-07, 855+ release 自动化每 2-5 天重生, 规则确实在更新而非仅仓库活跃)。WebFetch release 页一度返回 "2024" 系小模型对相对日期误读, 已用 GitHub API 纠正为 2026。但**它不是代理分流规则库**: 唯一产出 /etc/hosts 纯文本格式, 零客户端 (clash/.mrs/sing-box .srs/surge/geosite.dat/quantumult-x) 输出; 语义是 DNS sinkhole 封堵 (ads/malware + 可选 porn/social/gambling/fakenews), 与我们"上游源→自编译→扇出多客户端分流格式"的目标方向相反。血统为纯聚合器, 非 domain-list-community 派生。相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 无差异化价值 — Sukka 已在更上游消费同类 ad/malware hosts 源并产出真格式。故不值得加入 本项目 视野。唯一可考虑的极弱用法: 若未来 本项目 需要一个"广告/恶意域名"专项 reject 集, 可把 StevenBlack 当**纯域名输入源之一**自行编译成各客户端格式 — 但这属于"输入源"而非"规则库", 且 Sukka 已覆盖, 优先级很低。

</details>

<details>
<summary>❌ <b>SunsetMkt/anti-ip-attribution</b> — stale (active-leaning: 仓库未归档、CI 在跑, 但手工策展规则节奏为数月一次, 由 issue 驱动) · GPL-3.0 (copyleft, 若再分发编译产物需注意)</summary>

- **URL**: https://github.com/SunsetMkt/anti-ip-attribution
- **stars / 最近 / 维护**: ~1.1k (1055) · 2026-05 (repo push: 依赖/CI 自动提交); 源规则 rules.yaml 实质更新最近为 2026-01 (斗鱼直播流) · stale (active-leaning: 仓库未归档、CI 在跑, 但手工策展规则节奏为数月一次, 由 issue 驱动)
- **license**: GPL-3.0 (copyleft, 若再分发编译产物需注意)
- **formats**: clash rule-provider.yaml (含 direct/proxy/reject 拆分), clash-for-windows parser.yaml (已废弃), surge.list, quantumultx.list (含 domesticsocial 变体)
- **provenance**: independent-curated — 依据: 单一真相源 rules.yaml + rules/ 下按平台拆分的 YAML, 每条目内联引用 GitHub issue, generate.py 从该单源扇出; 全仓无任何上游列表 import, 既非 domain-list-community 派生也非聚合器
- **ai_coverage**: 无 (零 AI/LLM 域名: 无 openai/claude/anthropic/gemini/perplexity; 也无 Netflix/Disney 等西方流媒体, 仅中国直播平台且仅为 IP 归属地角度)
- **specialty**: 极小众独有专项: 不是"隐藏自己落地 IP 地区"(初判错误), 而是反向覆盖"中国社交/内容平台对外公开显示用户 IP归属地"这一 2022 后强制功能 (微博/B站/知乎/小红书/抖音/斗鱼直播/网易云/贴吧 等), 把这些平台域名定向分流使显示的省份一致。Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 均不覆盖此场景
- **reason**: 真实存在、确为代理分流规则库, 但与 本项目 目标正交。它解决的是"中国境内社交 App 对外显示 IP 归属地省份一致"这一国内表层体验问题, 而 本项目 是面向海外访问的订阅服务, 用不到这套域名; 且与已知库 0 重叠价值(无 AI/LLM、无西方流媒体、无 sing-box .srs / .mrs / geosite.dat, 不进我们的全客户端扇出)。维护上需区分"仓库活跃"(CI/依赖自动提交到 2026-05)与"规则真更新"(实质手工策展最近 2026-01, 数月一次), 属 stale-leaning。结论: 不值得纳入视野; 仅作为"独立手工策展 + 单源扇出"这一架构样本可一瞥, 但其数据与血统对我们无增量。GPL-3.0 也是潜在再分发约束。

</details>

<details>
<summary>❌ <b>sve1r/Rules-For-Quantumult-X</b> — active · MIT</summary>

- **URL**: https://github.com/sve1r/Rules-For-Quantumult-X
- **stars / 最近 / 维护**: ~3.8k (API: 3766) · 2026-02 · active
- **license**: MIT
- **formats**: quantumultx (.list, QX-native host-suffix/host-wildcard syntax), quantumultx rewrite (.conf/.adblock MitM), quantumultx scripts (.js)
- **provenance**: independent-curated — 依据: 无任何 geosite/domain-list 引用; README 自述「所有内容源自互联网,仅作收集整理」; 规则文件为手工分组 (OpenAI.list 用 `# >>` 手动分 ChatGPT/SSO/CDN/Static); 引用 QX 社区作者 (lhie1/NobyDa rewrite). 是 QX 生态内容的 aggregator/手工策展, 非 domain-list-community 派生
- **ai_coverage**: 薄弱: 仅单个 Rules/Services/OpenAI.list (OpenAI/ChatGPT/poe), 无 Claude/Anthropic、Gemini、Grok、Perplexity、Copilot 等; 远不如 Sukka/MetaCubeX 的 AI 集合
- **specialty**: 相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 在「分流域名」维度几乎零差异化, 且只产 QX 格式. 唯一独有资产是 QX rewrite/MitM 去广告脚本 (Bilibili/Weibo/Zhihu/小红书/Apple TestFlight 等), 但属 app-adblock-rewrite 领域, 与 本项目「域名规则→多客户端扇出」正交
- **reason**: 对抗式核实结论: 仓库真实存在且确为代理分流规则库 (2019 建, 662 commit, ~3.8k star, 非 archived, API pushed_at 2026-02-09 有真实 commit). 维护判定 active 但减速且 release 滞后: commit 节奏不均且变稀 (2025-03=20 → 2025-10=2 → 2025-11=1 → 2026-02 突击 4 条后 3-6 月空窗), 且 README 自报的 v25.02.12 release tag 较实际 commit 滞后约 16 个月 —— 典型「仓库活跃但 release 不切、且更新偏 cn-app 去广告 churn 而非核心 geo 路由」。不值得纳入 本项目 视野的硬理由: (1) 仅产 Quantumult X 格式, 无 clash .mrs / singbox .srs / surge / geosite.dat; (2) 其 .list 把策略目标写进第三列 (host-suffix,openai.com,OpenAI), 与 QX 路由强耦合, 不是可复用的纯域名源, 扇出前需剥离改写; (3) 血统为独立手工策展, 无 domain-list-community 可供我们 harvest 的上游; (4) AI 覆盖薄 (仅 OpenAI); (5) 唯一独有价值 (QX rewrite 脚本) 与我们的规则供应链目标正交。建议不追踪。

</details>

<details>
<summary>❌ <b>szkane/ClashRuleSet</b> — active · none (license=null, 无 LICENSE 文件, README 无授权声明 — 默认 all-rights-reserved, 再分发有法律风险)</summary>

- **URL**: https://github.com/szkane/ClashRuleSet
- **stars / 最近 / 维护**: 272 · 2026-05 · active
- **license**: none (license=null, 无 LICENSE 文件, README 无授权声明 — 默认 all-rights-reserved, 再分发有法律风险)
- **formats**: clash.list
- **provenance**: independent-curated (ACL4SSR-fork base + hand-curated personal additions) — NOT domain-list-community-derived. 依据: (1) README 明写"本项目基于 ACL4SSR 项目进行修改"; (2) 规则文件含 PROCESS-NAME 条目(macOS 进程名 "Claude"/"Code Helper"/"Grok"), 这是 Clash/Surge 专属、geosite/domain-list-community 纯域名格式无法表达的, 直接证伪 dlc 血统; (3) commit 是手工逐条加域名(add vooks.com / add files.pythonhosted.org), 非聚合器批量同步。初判 domain-list-community-derived 错误。
- **ai_coverage**: 较好且新鲜: AiDomain.list(~300 行)覆盖 openai/chatgpt、anthropic(claude.ai+claude.com+claudeusercontent.com)、x.ai/grok、perplexity、gemini/aistudio/generativelanguage、copilot/bing、poe、recraft; 另有独立 CiciAi.list 覆盖字节海外 AI。含 PROCESS-NAME 进程级 AI 规则。
- **specialty**: 差异化在专项品类而非规模/格式: (1) 字节海外 AI 单列 CiciAi.list(冷门, 主流库罕见); (2) Developer 品类把 docker 镜像 + huggingface 大模型下载当作"大流量低安全"独立分流策略(原创视角); (3) 美区 AI 强意见分组 + Web3 自建美区节点策略。规模远小于 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX, 仅作品类 taxonomy 灵感来源, 不是可编译的上游 feed。
- **reason**: 核实结论: 真实存在、确为 Clash 分流规则库、active(pushed_at 2026-05-27, updated 2026-06, 规则真在逐条更新, 非僵尸)。但不值得纳入 本项目 作为上游编译源, 原因三条: (1) 仅产出单客户端 Clash .list, 无 .mrs/.srs/surge/geosite, 与我们"上游源→自编译→扇出多客户端"目标不匹配, 还得我们反向解析 PROCESS-NAME(扇出到 sing-box/geosite 会丢信息); (2) license=null, 法律上不可再分发, 供应链审计直接红线; (3) 单人手工维护、ACL4SSR 衍生、覆盖面远窄于已知四大库。可作为"AI/字节海外AI/开发者大下载"品类 taxonomy 的灵感参考(idea source), 但不作 feed 追踪。初判血统 domain-list-community-derived 已证伪, 实为 ACL4SSR-fork 上的独立策展。

</details>

<details>
<summary>❌ <b>tangnahuaite/sing-box_Route-rules</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/tangnahuaite/sing-box_Route-rules
- **stars / 最近 / 维护**: ~101 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: singbox.srs, singbox.json
- **provenance**: aggregator-of-others — source.txt 列出 ~31 个上游 URL, 绝大多数是 blackmatrix7/ios_rule_script 的 .list 文件 (Apple/Google/Netflix/Disney/OpenAI/Telegram 等), 外加维护者自己的少量补充列表 (AD 屏蔽/AD-patch/GoogleFCM/GoogleVoice/Wi-Fi Calling DE-UK)。本质是 blackmatrix7 的下游 re-packager, 非 domain-list-community 派生。
- **ai_coverage**: 极弱: 仅 OpenAI.srs (继承自 blackmatrix7)。无 Claude/Anthropic、Gemini、Copilot、Perplexity 等任何其他 AI/LLM 域名列表。
- **specialty**: 几乎没有差异化价值: 它只是把 blackmatrix7 的 Surge .list 通过 GitHub Actions 转成 sing-box .srs/.json。我们已知 blackmatrix7 (源头) + MetaCubeX (官方 sing-box geosite .srs) 已完全覆盖且更权威。唯一私货是维护者手搓的 AD/Wi-Fi-Calling 列表, 价值边缘。它是单向 (其他格式→sing-box), 不做多客户端扇出, 与我们 本项目 "自编译扇出全客户端" 目标方向相反。
- **reason**: 对抗式核实结论: 仓库真实存在、确为 sing-box 分流规则库, GPL-3.0, ~101 star。"活跃"是假象——每日 commit 全部由 actions-user 机器人自动产生 (929 commits, 消息千篇一律 "Update rules"), 实质只是定时重新拉取上游并转格式, 维护者人工策展极少。血统为聚合器: source.txt 证明它 99% 内容来自 blackmatrix7/ios_rule_script, 我们已直接追踪源头 blackmatrix7, 它属冗余下游。产出仅 sing-box .srs/.json 单一目标且方向单向 (吃别人格式→吐 sing-box), 与我们 本项目 "上游源→自编译→扇出所有客户端" 的诉求正好相反, 无法复用其转换逻辑。AI 覆盖仅 OpenAI 一项, 无 Claude/Gemini。综合: 不值得加入 本项目 视野——既无独立血统也无差异化格式/专项覆盖, 唯一原创内容 (AD/Wi-Fi-Calling 列表) 价值边缘且与代理分流主线无关。

</details>

<details>
<summary>❌ <b>Toperlock/sing-box-geosite</b> — stale · none (GitHub 检测无 SPDX license; api license.spdx_id=null, 默认 all-rights-reserved, 再分发法律上灰色)</summary>

- **URL**: https://github.com/Toperlock/sing-box-geosite
- **stars / 最近 / 维护**: ~403 · 2026-06 · stale
- **license**: none (GitHub 检测无 SPDX license; api license.spdx_id=null, 默认 all-rights-reserved, 再分发法律上灰色)
- **formats**: singbox.srs, singbox source-format .json
- **provenance**: aggregator-of-others — links.txt 拉 blackmatrix7 (Microsoft/Telegram/Netflix/TikTok/OpenAI/AppleMusic/YouTube)、NobyDa (Bilibili/WeChat)、Loyalsoldier (GFW)、adrules.top (广告) 的 qx/surge/loon/clash 列表, 用 main.py + GitHub Actions 重新编译成 sing-box 规则集。完全是下游再编译, 不碰 v2fly domain-list-community, 非 geosite 派生。
- **ai_coverage**: 浅且仅 OpenAI。OpenAI.json/.srs 覆盖 openai.com/chatgpt.com/ai.com (~40 条, 含 Azure/Auth0/Stripe 支撑域), 但无 anthropic/claude.com、无 x.ai/grok、无 perplexity、无 gemini、无 copilot。且该 OpenAI 规则本质就是 blackmatrix7 OpenAI 列表的重编译, 非自研。
- **specialty**: 差异化价值几乎为零: 它的所有上游 (blackmatrix7/NobyDa/Loyalsoldier) 都已在我们已知清单内, 只是把它们重编成 sing-box .srs。唯一可借鉴的是 "links.txt 驱动 → 编译 .srs" 这套 Python pipeline 本身 (正是 本项目 想自建的 fan-out 编译器的一个参考实现), 但它只 fan-out 到 sing-box 单一格式, 且数据源是我们已有源的子集。
- **reason**: 对抗式核实结论: 真实存在且确为代理分流规则库 (Python, links.txt 驱动, GitHub Actions 每日 actions-user "Update rules" cron 提交, 输出 sing-box .json + 编译 .srs)。区分"仓库活跃 vs 规则真更新": cron 每天跑 (pushed_at 2026-06-08), 所以 .srs 数据确在每日重生成; 但最后一次人工提交是 2024-10-16 (Toperlock "upgrade version to 2"), 维护者已约 20 个月零人工策展 → 判 stale。血统为纯聚合器 (aggregator-of-others), 非 domain-list-community 派生; 上游 = blackmatrix7 + NobyDa + Loyalsoldier + adrules, 全部已在我们视野内。无 license。AI 覆盖仅 OpenAI 且系 blackmatrix7 重编。综合: 相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 无新增源、无差异化策展、单一 sing-box 格式、无 license, 不值得作为 本项目 的上游源追踪。唯一可顺手参考的是它 links.txt→.srs 的编译 pipeline 实现思路, 但无需纳入追踪清单。

</details>

<details>
<summary>❌ <b>uniartisan/adblock_list (now zhiyuan1i/adblock_list)</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/uniartisan/adblock_list
- **stars / 最近 / 维护**: ~724 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: AdGuard/EasyList/ABP filter syntax (.txt) — adblock.txt (basic), adblock_lite.txt (mobile, cosmetic/path rules stripped), adblock_plus.txt (~10MB full), adblock_privacy.txt (tracking/spyware)
- **provenance**: aggregator-of-others — 决定性证据在 scripts/lib/data_record.txt + metadata.txt + scripts/tools/raw_download.py：每日 GitHub Action 下载 EasyList / EasyListChina / EasyPrivacy (adblockplus.org)、cjx82630/cjxlist (cjxlist+cjx-annoyance)、AdGuard filters 2/3/14/17/224，再 merge+dedup，加极小本地补充 (raw.txt 248B / mobile.txt 171B / whitelist.txt 100B)。与 domain-list-community / Loyalsoldier 血统零关系。
- **ai_coverage**: 无专门 AI/LLM 域名分流覆盖（claude.com/grok/perplexity 等）。作为广告拦截聚合，其域名集为广告/追踪/隐私域，不含 AI 服务正向分流意图。
- **specialty**: 中文区广告/隐私拦截的 AdGuard 格式聚合，lite 版主动剥离 ## 元素隐藏/路径/cosmetic 规则降低移动端开销。但这是「广告拦截器规则」(AdGuard/uBO/AdGuard Home)，不是代理分流 rule-provider —— 不产出 clash/.mrs/sing-box.srs/surge/geosite 任何一种，无 GEOSITE/策略组分流语义。相对 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 在「分流」维度无差异化价值；它对标的是 anti-AD / AdGuardSDNSFilter 这类，而非我们的供应链对象。
- **reason**: 对抗式核实结论：(1) 真实存在，uniartisan 已改名 zhiyuan1i，仓库自动重定向，非无关项目。(2) 维护「活跃但需澄清」——表面 git 历史只有 2 条今日 commit 且都是 github-actions[bot]，因为 .github/workflows/new.yml 每天北京时间 1 点 `git checkout --orphan new_branch` + `push -f` 抹掉历史；产物 .txt 的 `! Version: 202606090450` 与 `Expires: 1 days` 证明规则确实每日真更新，不是僵尸仓库；2023 的 1.0.0 release/tag 是陈旧残留，别被它误导判 stale。根目录残留 `14.txttyp7aim_.tmp` (3.2MB) 是 wget 临时文件泄漏，属卫生瑕疵但不影响产物。(3) 格式只有 AdGuard/ABP filter 语法（##元素隐藏 + ||domain^ + $domain= + scriptlet），无任何代理客户端分流格式。(4) 血统 = 纯聚合器（EasyList系 + cjxlist + AdGuard 2/3/14/17/224），非 domain-list-community 派生、非独立手工策展。(5) GPL-3.0；无 AI/LLM 分流覆盖。结论：品类错位——它是「广告拦截规则聚合」而非「代理分流 rule-provider」，本项目 关注的是 geosite/分流供应链，二者不重叠；不值得纳入 本项目 视野。若未来 本项目 想加「订阅级广告拦截 rule-set（reject 策略组）」可作候选，但需自行转格式，且 anti-AD/AdGuardSDNSFilter/217heidai 等同类更主流。

</details>

<details>
<summary>❌ <b>vernette/rulesets</b> — active · NONE (无 LICENSE 文件, GitHub API license=null = all-rights-reserved, 默认不可再分发)</summary>

- **URL**: https://github.com/vernette/rulesets
- **stars / 最近 / 维护**: ~28 (forks 5) · 2026-05 · active
- **license**: NONE (无 LICENSE 文件, GitHub API license=null = all-rights-reserved, 默认不可再分发)
- **formats**: singbox.srs, singbox-json (source rule v2), raw-domain-txt
- **provenance**: independent-curated — 经读源码确认, 非 domain-list-community 派生, 非聚合器. json/*.json 是手写的 sing-box rule 文件 (claude.json = 3 个 domain_suffix: claude.ai/claude.com/anthropic.com; grok.json = grok.com/x.ai/grok.x.com/api.x.com), convert_json_to_txt.py 把 json 展平成 raw/*.txt, compile-srs.yml 用 `sing-box rule-set compile` 产 .srs. 唯一外部 feed 是 unavailable-in-russia 自动同步 dartraiden/no-russia-hosts 并 merge 自家 AI/Netflix 列表. 无任何 geosite/domain-list-community 痕迹.
- **ai_coverage**: Claude (claude.ai/claude.com/anthropic.com)、OpenAI、Gemini、Grok (grok.com/x.ai)、Copilot 各为独立 ruleset; 覆盖名录全但每项仅手写少量 apex 域名, 无 IP/CIDR, 颗粒度浅
- **specialty**: 差异化极有限: 纯 sing-box 单格式 (json/srs/raw txt), 不产 clash/surge/qx/geosite. AI 维度名义覆盖全 (Claude/OpenAI/Gemini/Grok/Copilot 各独立 ruleset) + 俄区 RKN/unavailable-in-russia 专项是相对 Loyalsoldier/blackmatrix7/Sukka 唯一稍有特色处, 但每个 ruleset 内容很浅 (仅顶级域 apex, 无 CIDR/IP 段, 无子 API 拆分), 深度与维护强度均不及我们已知三大库.
- **reason**: 真实存在且为代理分流规则库, 活跃 (最近 commit 2026-05-29, 但多为 bot 自动编译/俄区同步, 人工策展规则改动稀疏). 对抗式核实结论: (1) 仅 sing-box 单格式 (json/srs/raw txt), 不覆盖我们多客户端扇出需要的 clash/mihomo/surge/shadowrocket/egern/qx/geosite, 与 本项目「自编译扇出全客户端」诉求正交; (2) 致命问题: 无 license = all-rights-reserved, 不能合法纳入我们供应链; (3) 血统虽确为 independent-curated (非 dlc 派生/非聚合器), 但每个 ruleset 内容浅 (AI 项仅 3-4 个 apex 域、无 CIDR), 深度与维护频度全面劣于已知 Loyalsoldier/blackmatrix7/Sukka; (4) 仅 28 star, 无 release/tag, 仅 git-raw 分发. 俄区 RKN/unavailable-in-russia 是唯一微弱差异点, 但我们用户群非俄区核心. 综合: 不值得加入 本项目 视野, license 缺失单条即可一票否决.

</details>

<details>
<summary>❌ <b>xmdhs/sing-box-ruleset</b> — active · GPL-3.0</summary>

- **URL**: https://github.com/xmdhs/sing-box-ruleset
- **stars / 最近 / 维护**: ~20 · 2026-06 · active
- **license**: GPL-3.0
- **formats**: singbox.srs
- **provenance**: independent-curated — 实为单上游转码器, 非策展也非聚合器. build-adguard.sh 直接 download https://adguardteam.github.io/AdGuardSDNSFilter/Filters/filter.txt 再跑 `sing-box rule-set convert --type adguard`. 全程无 domain-list-community / geosite 引用, 不是 dlc 派生; 只搬 AdGuard 一家上游, 也不聚合 blackmatrix7/Loyalsoldier 等. 若严格分类应归 'transcoder-of-single-upstream', 在给定枚举里最接近 independent-curated.
- **ai_coverage**: 无 AI/LLM 专项域名覆盖. 仓库只做 AdGuard DNS 广告/追踪拦截, 不含 claude.com/grok/perplexity 等分流类目, 也无流媒体/geosite 类目.
- **specialty**: 唯一边际差异: 用 sing-box 原生 `rule-set convert --type adguard` 把 AdGuard SDNSFilter 的 AdGuard 语法(含 regex/例外规则)直接编进 sing-box 二进制 srs, 保留 AdGuard 拦截语义而非降维成纯域名列表. 产出 3 个变体: AdGuardSDNSFilter.srs(标准) / AdGuardSDNSFilter-NoRegex.srs(去正则, 兼容老内核) / AdGuardSDNSFilterSingBox.srs(1.10+ 优化版). 相对 Sukka/blackmatrix7 的广告规则覆盖面更窄, 仅此一类一上游一格式.
- **reason**: 对抗式核实结论: (1) 真实存在且确为 sing-box 广告拦截规则库, 非无关项目. (2) 关键反直觉点已澄清——master 源码分支 stale(末次 commit 2025-10-13 手动), 但 rule-set 输出分支由 github-actions[bot] 按 cron `0 3 * * 1` 每周一自动重建, 末次推送 2026-06-08(即昨天), 故规则确在更新, 判 active(非 stale); 维护活跃度不靠 commit churn 而靠 CI 心跳. (3) 仅产 sing-box .srs(3 个变体), 不扇出 clash/surge/quantumult-x/geosite.dat. (4) 血统: 非 dlc 派生、非聚合器, 是 AdGuard SDNSFilter 单上游的原生 srs 转码器(build-adguard.sh + sing-box rule-set convert --type adguard 为铁证). (5) GPL-3.0, ~20 star, 无 AI/流媒体类目. 不值得纳入 本项目 视野: 我们要做的是「上游源→自编译→扇出全客户端格式」, 它在 本项目 链路里只对应「单上游 AdGuard 广告 → 单格式 srs」这一极窄切片, 且广告拦截我们已有 Sukka(更全)/blackmatrix7 覆盖, 多客户端扇出它也帮不上(只出 srs). 唯一可借鉴的不是它的规则数据, 而是它的方法论——直接用 sing-box 原生 `rule-set convert --type adguard` 保留 AdGuard 语义编译 srs, 这点可作为我们 reject 类目 sing-box 输出的一个实现参考, 但不需要追踪该仓库本身.

</details>

<details>
<summary>❌ <b>xndeye/rule-merger</b> — active · none (无 LICENSE 文件, GitHub API license=null)</summary>

- **URL**: https://github.com/xndeye/rule-merger
- **stars / 最近 / 维护**: 9 · 2026-06 · active
- **license**: none (无 LICENSE 文件, GitHub API license=null)
- **formats**: clash.yaml (mihomo), clash.mrs (mihomo rule-set), text/.txt
- **provenance**: aggregator-of-others — config.yaml 明确以 HTTP 拉取并 merge 第三方上游: Sukka/skk.moe ruleset.skk.moe (主力), MetaCubeX/meta-rules-dat (geoip/geosite gfw), DustinWin/domain-list-custom (ai/apple-cn/microsoft-cn/google-cn), ACL4SSR, blackmatrix7/ios_rule_script (SteamCN), Accademia/Additional_Rule_For_Clash (GeositeCN), 加 xndeye 自己的 adblock_list + local/*.yaml 手工补丁。本身不维护 domain 源,纯聚合器+格式转换器。非 domain-list-community 直系派生(仅间接经 DustinWin/MetaCubeX 引入 geosite 数据)。
- **ai_coverage**: 有专门 ai 类 (output/ai.mrs + ai.yaml, release 分支当日 194 条)。覆盖 OpenAI/Anthropic(claude)/Google Gemini+AIStudio+NotebookLM+Antigravity/Groq/OpenRouter/Perplexity/together.xyz/JetBrains AI/GitHub Copilot 等, 用 DOMAIN-KEYWORD anthropic/claude/openai 兜底。来源 = skk.moe ai.txt + DustinWin ai.list + ACL4SSR AI.list + 本地 local/ai.yaml 手工补 (如 vercel.ai)。覆盖尚可但全部来自已知上游, 无独家域名。无 grok/x.ai 专项 (仅靠 keyword 可能漏)。
- **specialty**: 差异化价值很低: 上游几乎全是我们已知库 (Sukka/skk.moe + MetaCubeX + DustinWin/ACL4SSR/blackmatrix7)。它的唯一增量是"把这些已知源按个人口味 merge 成 12 个粗粒度策略组分类 (reject/direct/proxy/ai/apple@cn/microsoft@cn/steam@cn/lan/fakeip-filter) 并互转 domain↔ipcidr↔classical→.mrs"。属于"成品消费层"而非"源策展层"。对我们想自掌控的"上游源→自编译"链条没有新源贡献——我们若要这种聚合能力,自己的 本项目 编译器就是干这个的,直接接它的上游即可。
- **reason**: 真实存在, 确为 Mihomo 代理分流规则库 (非无关项目), 当前 active: master 配置最近 commit 2026-05-14, CI 输出 release 分支每日自动 commit (最新 2026-06-09 "🚀 Update", ai.yaml 头注释更新时间戳与之吻合) — 规则数据是真在每日刷新的, 不是僵尸仓库。注意 GitHub API 的 pushed_at=2026-06-09 来自 release 分支的机器人提交, 人工配置维护频率其实是月级。产出仅 mihomo 系 (.yaml/.mrs/.txt), 无 sing-box .srs / surge / quantumultx, 与我们多客户端扇出诉求不符。血统为纯聚合器, 上游 (Sukka/MetaCubeX/DustinWin/ACL4SSR/blackmatrix7) 全在我们已知清单内, 无独家源。无 license (法律上不可放心 vendor/再分发)。结论: 不值得加入 本项目 视野 —— 它是"别人成品的二次聚合消费者", 和我们要做的"上游源→自编译→扇出"是同一层的竞品而非可复用的源, 既无新上游也无新格式, 加上零 license 风险, 仅可作为"个人 merge 配置写法"的参考样例 (config.yaml 的 domain/ipcidr/classical 互转 + behavior 声明模式), 不作为追踪对象。

</details>

<details>
<summary>❌ <b>Yuu518/sing-box-rules</b> — active · none (license=null, 无 LICENSE 文件)</summary>

- **URL**: https://github.com/Yuu518/sing-box-rules
- **stars / 最近 / 维护**: ~60 · 2026-06 · active
- **license**: none (license=null, 无 LICENSE 文件)
- **formats**: singbox.srs, singbox.json, txt lists (proxy/direct/reject on release branch), hosts, adguard .srs
- **provenance**: domain-list-community-derived — Build.yml 是铁证: 它 checkout v2fly/domain-list-community, 把多源合并结果写回 community/data/*, 再用 Yuu518/rules-generate 全量编译整个 v2fly 类别集 (go run ./ domain -d ../community/data -f singbox), 故输出 1000+ 文件且带 @attr 命名 (alibaba@!cn / acfun@ads) = v2fly 属性体系签名. 在 v2fly 基底上叠加聚合层: 广告(Cats-Team/AdRules + TG-Twilight/AWAvenue), CN域名(pmkol/easymosdns + felixonmars/dnsmasq-china-list), GFW(gfwlist2dnsmasq), 代理(Loyalsoldier/domain-list-custom), 私有 R2 bucket 自定义清单. IP 侧 = Yuu518/geoip (Loyalsoldier/geoip 的 fork). 所以是 "v2fly 派生 + 聚合器" 混合, 非纯手工策展, 也非纯聚合.
- **ai_coverage**: good/current — 提供 category-ai-!cn / category-ai-chat-!cn / category-ai-cn (按 CN/非CN 切分). 实测 category-ai-!cn.json 含 claude.ai/claude.com/anthropic.com/openai.com/chatgpt.com/grok.com/x.ai/perplexity.ai/pplx.ai/mistral.ai/meta.ai/deepmind/gemini(aistudio,makersuite,bard)/trae.ai/coderabbit.ai/antigravity.google 等, 含很新域名. 但此覆盖来自上游 v2fly category-ai-!cn, 非 Yuu518 自有策展.
- **specialty**: 相对我们已知库差异化很弱: 域名基底就是 v2fly domain-list-community (我们能直接拿上游), IP 是 Loyalsoldier/geoip 的二手 fork, 广告源是 Cats-Team/AWAvenue (可直接接). 唯一"它有我们没有"的是: (a) stream-global / pcdn-cn 两个自定义聚合类别; (b) 把以上全部预编译成 sing-box .srs+.json 的现成产物. 但这正是 本项目 要自建的扇出环节, 没必要消费别人的成品.
- **reason**: 对抗式核实结论: 仓库真实存在且确为 sing-box 代理分流规则库; 活跃度判定 active 且"规则真在更新"——CI cron 每日 23:00 跑, release/rule_set 输出分支显示 Released on 2026-06-09 (核实当日), AI 清单含 antigravity.google/grok.com 等最新域名, 证明是新鲜重编非陈旧快照 (注意 master 源分支停在 2026-04, 别被误判为 stale). 但不值得加入 本项目 视野, 三条硬伤: (1) 血统下游——它 = v2fly domain-list-community (域名) + Loyalsoldier/geoip fork (IP) + 几个公开广告源的预编译聚合, 我们已直接掌控这些上游, 消费它等于多一层不可控中间商, 与 本项目 "自己掌控上游→自编译→扇出" 的目标背道而驰; (2) 无 license, 供应链合规直接红线; (3) 格式单一 (仅 sing-box .srs/.json), 对我们 7 客户端扇出无增量, 我们要的恰是自己做编译/扇出这一段. 可作为"自定义聚合类别 stream-global/pcdn-cn 的 idea 参考", 但不作为依赖源跟踪.

</details>

<details>
<summary>❌ <b>yuumimi/geosite</b> — active · MIT</summary>

- **URL**: https://github.com/yuumimi/geosite
- **stars / 最近 / 维护**: 0 · 2026-06 · active
- **license**: MIT
- **formats**: geosite.dat, geosite.dat.sha256sum
- **provenance**: domain-list-community-derived — GitHub API confirms fork:true, parent=Loyalsoldier/domain-list-custom (the build tool). 工作流 build.yml 在构建时 checkout v2fly/domain-list-community 的 data 目录, 跑 `go run ./ --datapath=./domain-list-community/data` 生成 geosite.dat。本仓库不持有任何手工策展的域名清单, 99% 内容直接来自上游 domain-list-community。本 fork 唯一真实定制 = 在 category-vpnservices 末尾 append 一串硬编码诈骗/VPN 域名 (auvpn.net, ausu.*/ausososo.* 系列, splashtop.com @cn) + 两行 include。
- **ai_coverage**: 间接覆盖 — 仅继承上游 v2fly/domain-list-community 的 AI 相关分类 (如 category-ai-!cn, 含 openai/anthropic/google 等), 本 fork 无任何 AI/LLM 专项增量或自有 claude.com/grok/perplexity 策展。
- **specialty**: 几乎没有差异化价值。它产出的 geosite.dat 在内容上约等于「v2fly/domain-list-community 官方 dat + 一小撮反诈/VPN 域名补丁」, 与我们已知的 Loyalsoldier/v2ray-rules-dat 高度重叠且功能更弱 (Loyalsoldier 还做 geoip、cn-max 等增强, 且有社区维护)。AI/流媒体覆盖完全继承自 domain-list-community 的 category-ai-!cn / geolocation-!cn 等上游分类, 本仓库未做任何 AI/流媒体专项增量。初判中「去 @ads/@!cn 噪音、保护在华海外公司不误代理」属误读: 那段是从 Loyalsoldier 上游 README 抄来的样板文案, 实际 workflow 并未执行该过滤 (只 append 黑名单, 未做 @ads/@!cn 删除)。
- **reason**: 对抗式核实结论: 仓库真实存在且确为代理分流规则库, release 管线每天 21:30 UTC 跑、最新 release 2026-06-08 (active)。但「活跃」仅指自动化打包管线 — 仓库源码/工作流上次人工改动停留在 2022-09, 实质是 Loyalsoldier/domain-list-custom 的一个个人 fork, 构建时 live 拉取 v2fly/domain-list-community 数据。血统判定为 domain-list-community-derived 成立 (GitHub API fork=true + parent + workflow checkout 三重证据)。但它对 本项目 几乎无独立价值: (1) 仅产出单一 geosite.dat, 无 geoip、无 .srs/.mrs/clash/surge 等多客户端格式, 与我们要做的「自编译扇出全客户端」诉求正交; (2) 内容 = 上游 domain-list-community 官方 dat + 一小撮硬编码反诈/VPN 域名, 与 Loyalsoldier/blackmatrix7/Sukka/MetaCubeX 高度重叠且弱; (3) 0 star、0 fork、单人无社区, 维护中断风险高; (4) README 的「去噪/保护在华海外公司」专项是抄上游样板文案, 实际 workflow 未执行, 初判专项不成立。建议: 若 本项目 需要 domain-list-community 这条上游, 应直接追踪 v2fly/domain-list-community 与 Loyalsoldier/v2ray-rules-dat 这两个源头, 没必要纳入这个派生 fork。不值得加入视野。

</details>

<details>
<summary>❌ <b>Z-Siqi/Clash-for-Windows_Rule</b> — stale · none (license: null, /LICENSE 返回 404 → 默认 all-rights-reserved, 无授权再分发)</summary>

- **URL**: https://github.com/Z-Siqi/Clash-for-Windows_Rule
- **stars / 最近 / 维护**: 659 · 2025-12 · stale
- **license**: none (license: null, /LICENSE 返回 404 → 默认 all-rights-reserved, 无授权再分发)
- **formats**: clash-classical-inline-rules (DOMAIN-SUFFIX/DOMAIN-KEYWORD, no extension, policy-group baked into 3rd field)
- **provenance**: independent-curated — 依据: 无任何 domain-list-community 引用/无 attribution 注释; README 自述规则"可能不全面或有误"并征集社区 PR; Netflix 等文件为人工维护的扁平 DOMAIN-SUFFIX 列表 (含 netflixdnstest0-9 探测域等手填条目), 非生成物; 无聚合脚本/无 CI 拉取上游
- **ai_coverage**: 弱。仅 OpenAI 一个文件 (~19 条, ChatGPT 为主, 含 azure/cloudflare/auth0/sentry 等附属域)。无 Claude/Anthropic, 无 Gemini/Bard, 无 Copilot。AI 维度落后于 Sukka / blackmatrix7。
- **specialty**: 差异化几乎为零, 且整体弱于我们已知库。它的卖点(流媒体/游戏分流按国家落地)不是靠数据而是靠 template/ 里的策略组→国家服务器映射实现, 规则文件本身是单一扁平列表 (Netflix 无 US/JP 拆分)。覆盖 62 个服务 (Netflix/Disney+/Spotify/Steam/Epic/Xbox/PSN 等) 但每个域名集都比 blackmatrix7/Loyalsoldier 小且更新慢。唯一可借鉴点: 它把"按国家落地"的策略组模板做得直白 (适合人类抄配置), 但对我们"自编译扇出"无数据价值。
- **reason**: 对抗式核实结论: 真实存在的 Clash 分流规则库 (659★/89 fork, 未 archived), 但不值得加入 本项目 视野。理由: (1) 法律硬伤——无 license (null + /LICENSE 404), 默认保留所有权利, 无法合规拉取再自编译扇出, 这对"掌控供应链"目标是 deal-breaker; (2) 格式单一——只有 Clash 经典内联规则 (无扩展名, policy-group 写死在第3字段), 无 .list/.mrs/.srs/.dat, 且 .github/workflows 404 即零自动编译, 与我们"上游源→自编译→扇出多客户端"模型正交; (3) 维护偏 stale——最后一次代码 push 2025-12-25 (updated_at 2026-06 仅元数据触碰), 提交稀疏且多为零星 Direct/Apple 微调, 远非高频跟踪库; (4) 数据弱于已知库——AI 仅 OpenAI 无 Claude/Gemini, 流媒体按国家落地靠 template 而非数据, 规则集体量小于 blackmatrix7/Loyalsoldier/Sukka; (5) 血统为独立人工策展但价值低。纠正初判: "per-country 规则文件"判断有误 (国家落地在模板层); 血统 independent-curated 判断成立。综合: 既无授权又无格式管线又无数据增量, pass。

</details>

<details>
<summary>❌ <b>zqzess/rule_for_quantumultX</b> — stale · none (无 LICENSE 文件,license=null;README 仅致谢 blackmatrix7/NobyDa/chavyleung 等,法律上 all-rights-reserved)</summary>

- **URL**: https://github.com/zqzess/rule_for_quantumultX
- **stars / 最近 / 维护**: ~2043 · 2026-02 (master pushed_at 2026-02-21; latest visible default-branch commit 2026-01-04, manual hand-edits of personal myProxy/myDirect files) · stale
- **license**: none (无 LICENSE 文件,license=null;README 仅致谢 blackmatrix7/NobyDa/chavyleung 等,法律上 all-rights-reserved)
- **formats**: surge.list (.list + .sgmodule), quantumultx (filter list + rewrite), clash.list / clash .yaml, loon plugin, shadowrocket module, stash, adguardhome, scriptable/js scripts
- **provenance**: aggregator-of-others。依据：factory/*.py 爬虫脚本的上游全部是别人的库——代理分流类(Netflix/YouTube/Mainland 等)直接 raw.githubusercontent 拉 blackmatrix7/ios_rule_script + GeQ1an/Rules;去广告类(ad.py)聚合 anti-AD、AdGuardSDNSFilter、EasyList(China)、adbyby、xinggsf。没有任何 domain-list-community / Loyalsoldier 血统,也不是独立手工策展原始域名;只有维护者自己的 myProxy/myDirect/Mirror 少量私用 list 是手写。
- **ai_coverage**: 无 AI/LLM 专项。爬虫覆盖面是 Netflix/YouTube/Google/Microsoft/Apple/Mainland + 去广告,无 OpenAI/Claude/Grok/Perplexity 等 AI 域名 category(这正是它落后于 blackmatrix7/Sukka 的地方)。
- **specialty**: 相对我们已知库无差异化价值:它本身是 blackmatrix7 的下游聚合+二次封装(去广告侧再聚合 anti-AD/AdGuard/EasyList)。唯一独有内容是维护者私用配置(myProxy.yaml/myDirect/Mirror、PlayBoy/FanQieNovel 去广告、签到/比价 JS),对一个多客户端订阅服务的供应链没有可复用的上游分流数据。
- **reason**: 真实存在、确为代理分流+去广告规则库(2020 建,~2k star),但对抗式核查后不值得纳入 本项目 视野。三点否决:(1) 血统是纯下游聚合器——分流规则爬虫上游就是我们已追踪的 blackmatrix7(外加 GeQ1an),去广告侧聚合 anti-AD/AdGuard/EasyList,无任何 domain-list-community 派生或一手策展,我们追上游即可,追它等于追二道贩子。(2) "爬虫每周自动更新"是过期自述:CI(.github/workflows/main.yml cron 每周一)产生的 "程序自动更新规则" 自动提交在 master 上最后一次是 2024-11-18,此后近 15 个月全是维护者手动改自己私用的 myProxy/myDirect 文件——"仓库 pushed_at 看着新"≠"规则在更新",真正的规则编译管线已死。(3) 无 license(all-rights-reserved 法律风险)、无 AI/LLM 专项、无任何编译产物(geosite.dat/.srs/.mrs),全是明文 list,对我们"自编译扇出多客户端"目标零增量。如果要这条线的数据,直接追它的真上游 blackmatrix7/ios_rule_script 即可。

</details>
