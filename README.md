# 通天河 (Tongtian)

> 一份**自己掌控**的代理分流规则供应链：选定可信上游 → 统一规则模型 → 自编译 → 扇出到每个客户端的原生格式。
>
> 名字取自长江上游的**通天河**——河流天然就是这套系统的隐喻：源头（上游数据源）→ 汇流（统一规则模型）→ 分汊扇出（各客户端格式）。

---

## 0. 这是什么 / 为什么存在

`通天河` 是一份**自托管的代理分流规则供应链**。它面向一个具体问题：

> 一个**多客户端**代理订阅（clash/mihomo、sing-box、surge、shadowrocket、egern、quantumult-x、v2ray），如果每个客户端的规则各自外包给不同发行方，就会出现**同一个 Netflix / Claude，换个客户端路由结果不同**；而上游发行方又各有冻结 / 信任 / 格式覆盖缺口。

`通天河` 的主张：**规则供应链应当自己掌控一份上游，自编译，扇出所有客户端格式**。这样跨客户端一致、可钉版本、可回退、可审计。

> 上游选型的完整调查见 [`docs/upstream-audit.md`](docs/upstream-audit.md)。

---

## 1. 现状与思考来源（先读这个）

在选定上游之前，我们做了一次**全面的规则供应链审计**：对照已知的主流规则库，又用多角度网络搜索扇出、对抗式核实了 **80 个公开规则库**（覆盖 clash / sing-box / surge / 广告 reject / GEOIP / AI / 区域专项 / 聚合器）。

完整调查（调查了哪些 repo、分别是什么、为什么选 / 为什么否决）见 **[`docs/upstream-audit.md`](docs/upstream-audit.md)**。一句话结论：

> 这个生态**绝大多数规则库同源**——血统都能回溯到 `v2fly/domain-list-community`（多经 Loyalsoldier 工具链）或 `blackmatrix7`，且很多已冻结 / 停更 / 无 license / 单格式 / 类目错配。**真正给我们带来"新血统"的只有 4 类源**（GFW 一手源、GEOIP-CN 一手源、更强的广告 reject、AI 长尾补全）。

---

## 2. 选定的源头（最终上游）

> **铁律**：上游只取**源数据**，分类边界 / 格式 / 托管全部由 `通天河` 的编译器掌控。**钉版本（commit/tag）而非滚动分支**，防上游某次坏 PR 直接进生产。

### 🥇 Tier 1 — 骨架（一手源，构成规则的主干）

| 用途 | 上游 | License | 为什么是它 |
|------|------|---------|-----------|
| **geosite 骨架**（流媒体 / 大厂 / 社交 / AI 等分类） | [`v2fly/domain-list-community`](https://github.com/v2fly/domain-list-community) | MIT | geosite 事实标准源，全生态（v2ray/Xray/sing-box/Clash）共用，PR 驱动 + 海量消费 = 坏改动暴露快。provenance 比个人维护、已冻结的 blackmatrix7 更硬 |
| **CN 域名表**（域名级 DNS 分流核心） | [`felixonmars/dnsmasq-china-list`](https://github.com/felixonmars/dnsmasq-china-list) | 宽松（待最终核对条款） | 11 万+ 条穷举 CN 域名，**域名级 DNS 分流的唯一正确工具**。也是 Loyalsoldier `direct.txt` 的真实上游 |
| **GFW 语义一手源** ⭐新 | [`gfwlist/gfwlist`](https://github.com/gfwlist/gfwlist) | LGPL-2.1 | **整条 anticensor 谱系的祖先**——Loyalsoldier 的 gfw.dat 反而是从它抽取。直接消费可绕过中间层的停更 / 决策风险。自写 AutoProxy parser 转中性域名集 |
| **GEOIP-CN 一手源**（IP 侧分流） ⭐新 | [`misakaio/chnroutes2`](https://github.com/misakaio/chnroutes2) | CC-BY-SA-4.0 | 一手 BGP route-collector 聚合、逐小时刷新。补上 **IP 侧上游**（多数现成方案只有域名级）。也是客户端 DNS 防污染 geoip 裁决的数据来源 |

### 🥈 Tier 2 — 分类富化（择优合并，编译期去重）

| 用途 | 上游 | License | 备注 |
|------|------|---------|------|
| **广告 / 隐私 reject** | [`privacy-protection-tools/anti-AD`](https://github.com/privacy-protection-tools/anti-AD)（CN 广告权威，MIT）+ [`hagezi/dns-blocklists`](https://github.com/hagezi/dns-blocklists)（分级最细，GPL-3.0） | MIT / GPL-3.0 | 比只用 Sukka/Loyalsoldier 更专、更分级。接入需叠加 AI/API 白名单防误杀。[`217heidai/adblockfilters`](https://github.com/217heidai/adblockfilters) 已含 .mrs/.srs 编译态，可作格式参照 |
| **AI / LLM** | `domain-list-community` AI 类 + Sukka `ai`（基底）；**cross-check 长尾**：[`VPSDance/ai-proxy-rules`](https://github.com/VPSDance/ai-proxy-rules)（MIT，厂商 taxonomy 最全）、[`fmz200/wool_scripts`](https://github.com/fmz200/wool_scripts)（GPL-3.0，含 Anthropic IP-CIDR）、[`dler-io/Rules`](https://github.com/dler-io/Rules)（无 license，**仅抽域名**） | 混合 | 直接补 AI 覆盖短板：claude.com / grok / perplexity / cursor / windsurf。无 license 源只抽域名当情报，不 vendor 文件 |

### 🔧 Tier 3 — 工程参照（**不是数据源**，是编译/扇出管线的开源蓝本）

| 项目 | License | 借鉴点 |
|------|---------|--------|
| [`xkww3n/Rules`](https://github.com/xkww3n/Rules) | MIT | 单源 → clash `.mrs` + singbox `.srs` + surge + egern… 全客户端扇出的**可运行 MIT 管线**；dedup / CIDR 合并 |
| [`DustinWin/domain-list-custom`](https://github.com/DustinWin/domain-list-custom) | MIT | Go trie 去重 + 多源合并 + 扇出**编译器参考** |
| [`QuixoticHeart/rule-set`](https://github.com/QuixoticHeart/rule-set) · [`FuGfConfig`](https://github.com/Elysian-Realme/FuGfConfig) | GPL-3.0 / MIT | 「单源 → 全客户端格式」的 working blueprint（CI 内跑真二进制编译 `.mrs`/`.srs`） |

### 🔬 仅审计 / cross-check（**不进 pipeline**）

学术自动化测量项目——经核实**不适合做 pipeline 上游**：

| 项目 | 实情（已核实） | 角色 |
|------|---------------|------|
| **OONI** | 数据**活**（每小时上传 S3，AWS Open Data），**但是原始逐次测量 JSONL，不是黑名单**；提取干净列表需自建检测管线 + 处理高假阳性 | 偶尔审计：抽样核对，不消费 |
| **GFWatch** (USENIX'21) | "日更"出自 2021 论文（2020 数据）；`gfwatch.org` 现为 JS dashboard，**未确认仍有可下载的新鲜批量 feed**（学术平台普遍发完即冻结） | 一次性快照交叉核对 |
| **GFWeb** (USENIX'24) | 覆盖 **SNI/HTTP 层**（人工 gfwlist 可能滞后的层），但是**周期性数据集发布、非 live feed** | 周期性 diff 审计「SNI 层我们漏了啥」 |

> **为什么只是审计而非源**：① 新鲜度 / 可下载性多数不可靠（论文型，易冻结）；② 格式是测量数据非规则集，OONI 尤其要先建检测管线；③ 对**路由**而言边际价值低——它们的强项是穷举长尾（用户多数不访问），还会把 **GFW 自己的过封 bug**（GFWatch 自报 4.1 万正则连坐误封）抄进来。真正有用的窄场景：拿 GFWeb 的 SNI 层发现做**人工审计输入**，过「先测试再入库」的闸门，不自动消费。

---

## 3. 统一规则模型（汇流：唯一真理源）

所有来源（内置 + 未来的用户自定义）都收口到**同一种内部表示**，编译器不区分来源：

```yaml
# 类别集合（set）— 由编译器从上游编出，命名的 matcher 集合
netflix: [ "+.netflix.com", "+.nflxvideo.net", ... ]

# 路由规则（rule）— matcher + value + policy
- match: DOMAIN-SUFFIX        # 见 §4 能力表
  value: example.com
  policy: PROXY               # PROXY | DIRECT | REJECT | <命名策略组>
- set: netflix                # 引用类别集合
  policy: STREAMING
```

- `policy: PROXY` 渲染时映射到各客户端声明的主代理组。
- 命名策略组（如「② Netflix」）由我们生成的模板保证**跨客户端一致命名**。

---

## 4. 编译 → 扇出（怎么编成每个客户端的规则集）

### 4.1 编译流水（CI）

```
拉上游(钉 commit/tag)
   │  ① 解析各上游 DSL：展开 include、过滤 @attr、翻译 regexp/keyword/full
   ▼
统一规则模型（§3）：去重 + CIDR 合并 + 分类边界归一
   │  ② 按各客户端能力扇出（不支持的 matcher 按目标降级 + 告警）
   ▼
   ┌────────┬─────────┬──────────────┬─────────┬──────────┬──────────────┐
 clash/    surge    shadowrocket    egern    sing-box   v2ray
 mihomo    .list                  (复用     .srs       geosite.dat
 .list+.mrs                        clash)                (客户端 bundled)
   │  ③ 托管 + 日更 + 版本钉 / 回退（§5）
   ▼
   ruleset CDN（自有 repo @ jsdelivr） + 可选自有镜像
```

### 4.2 扇出目标与工具

| 客户端 | 输出 | 工具 |
|--------|------|------|
| clash / mihomo | `.list`（domain behavior）+ `.mrs`（二进制） | `mihomo convert-ruleset` |
| sing-box | `.srs`（二进制） | `sing-box rule-set compile` |
| surge / shadowrocket | surge 风格 `.list` | 文本模板 |
| egern | 复用 clash `.list`（egern 兼容） | — |
| v2ray | `geosite.dat`（客户端侧 bundled，按需产出） | v2fly generator |

### 4.3 能力天花板（按目标降级）

统一格式能表达的 = **各客户端能力的交集**；超出的按目标降级（支持的正常编、不支持的跳过 + 编译告警）。

| matcher | clash | sing-box | surge/sr | 说明 |
|---------|:---:|:---:|:---:|------|
| DOMAIN / -SUFFIX / -KEYWORD | ✅ | ✅ | ✅ | 交集 |
| IP-CIDR / IP-CIDR6 | ✅ | ✅ | ✅ | 交集 |
| DOMAIN-REGEX | ✅ | ✅ | ⚠️ | 内置层可用，按目标降级 |
| GEOIP / GEOSITE-ref | ✅ | ✅ | ⚠️ | 内置层用 |

### 4.4 产物布局（示意）

```
release/                       # CI 产出的 release 分支，jsdelivr 直接指向
├── clash/{netflix,ai,cn,reject,...}.{list,mrs}
├── singbox/{netflix,ai,cn,reject,...}.srs
├── surge/{netflix,ai,cn,reject,...}.list
├── geoip/cn.{list,mrs,srs}          # 来自 chnroutes2
└── version.json                     # 上游 commit/tag 钉点 + 构建时间戳
```

---

## 5. 托管 / 更新 / 回退 / 非致命投递

- **编译运行**：CI（GitHub Action，可复现）拉上游 → 编译 → 产物 commit 到 release 分支。
- **分发**：主走 `自有 repo @ jsdelivr`（免费 CDN、零自有服务器负载）；自有域名 `ruleset.<域名>` 作可选镜像。
- **更新频率**：日更（cron）+ 上游变更触发。
- **版本钉 + 回退** ⭐：上游钉到 commit/tag；release 分支保留历史 tag，订阅模板可指向 `@<tag>` 而非滚动分支；坏批次一键回退到上一 tag。
- **非致命投递** ⭐：clash/singbox/surge 一律用 **per-category rule-provider**（`RULE-SET` 远程拉取，**失败 = 降级，非致命**），**绝不**用会致命 bootstrap 的整包 geodata。

---

## 6. License 分层（开源 / 再分发合规）

新增源 license 跨度大，**vendor / 再分发前必须过一遍**：

| 档 | 协议 | 我们怎么用 |
|----|------|-----------|
| 友好 | MIT / BSD / LGPL（domain-list-community、anti-AD、xkww3n、VPSDance…） | 可 vendor 产物，保留署名 |
| copyleft | GPL-3.0 / CC-BY-SA（chnroutes2、hagezi、217heidai…） | 自编译再分发需**保留署名 + 同协议**，与你的分发 / 商业性质**单独评估** |
| 无 license | dler-io 等 | **只抽域名当情报，不 vendor 文件** |

---

## 7. 如何使用

各客户端在自己的配置 / 订阅里以 per-category `RULE-SET` 远程指向 `通天河` 的产物（产物布局见 §4.4，jsdelivr 地址见 §5），失败即降级、非致命（§5）。

**一个额外用法**：`chnroutes2` 编出的 GEOIP-CN 数据除了用作 IP-CIDR 路由规则，也可直接喂给 clash/mihomo 的 DNS `fallback-filter`（`geoip: true, geoip-code: CN`）做防污染裁决——同一份数据双用。

---

## 8. 致谢上游

`通天河` 站在巨人的肩膀上。所有数据来自上述开源上游的辛勤维护者——尤其是 `v2fly/domain-list-community` 的贡献者社区、`gfwlist` 十余年的人工核验者、`felixonmars` 与 `misakaio` 的基础数据维护者。我们承诺：经多点复测确认的长尾域名，将**回馈上游**。
