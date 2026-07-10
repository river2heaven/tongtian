#!/usr/bin/env bash
# 生成 data/china-company-ip.cidr —— 大厂 ASN 全 IP（含海外接入段）的聚合 CIDR。
#
# 数据源：ipverse/asn-ip（各 ASN 的 announced CIDR 聚合，as/<asn>/ipv{4,6}-aggregated.txt）。
# 钉 commit SHA 复现：改 IPVERSE_SHA 重跑即刷新（大厂 ASN 的 BGP 宣告会变）。
#
# ASN 清单在下方 ASNS 数组，每个都经 ipverse header 的 org 名逐个核对（排除 ISP 骨干/无关公司）。
# 增删公司：改 ASNS + 注释，重跑，提交 data/china-company-ip.cidr。
#
# 用法：bash scripts/gen-china-company-ip.sh
set -euo pipefail

# —— 钉定 ipverse/asn-ip 版本（bump：换成新 master SHA 重跑）——
IPVERSE_SHA="b5889fcdbfc99d5a73f691474842d0d35a85a3c0"
RAW="https://raw.githubusercontent.com/ipverse/asn-ip/${IPVERSE_SHA}/as"
OUT="$(cd "$(dirname "$0")/.." && pwd)/data/china-company-ip.cidr"

# ASN → 公司（header org 名逐个核对；ISP 骨干如 China Unicom/Telecom/Mobile 一律不收——过宽）。
ASNS=(
  "132203 Tencent-global"           # 微信 mmtls 43.128/10 海外接入点在此
  "45090 Tencent-Shenzhen"
  "133478 Tencent-Cloud-Beijing"
  "137876 Tencent-Thailand"
  "45102 Alibaba-China"
  "37963 Alibaba-Hangzhou"
  "134963 AlibabaCloud-Singapore"
  "24429 Alibaba-Taobao"
  "396986 ByteDance-US"
  "138699 ByteDance-TikTok"
  "150436 ByteDance-Byteplus"
  "137718 ByteDance-VolcanoEngine"
  "38365 Baidu"
  "55967 Baidu"
  "55990 Huawei-Cloud"
  "136907 Huawei-Intl-Singapore"
  "149167 Huawei-Intl"
  "45062 NetEase"
  "55992 Qihoo360"
  "131486 JD"
)

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

{
  echo "# china-company-ip — 大厂 ASN 全 IP(含海外接入段)→ 供分流直连"
  echo "# 生成源: ipverse/asn-ip @ ${IPVERSE_SHA}"
  echo "# 生成脚本: scripts/gen-china-company-ip.sh (可复现: 换 pinned SHA 重跑)"
  echo "# ASN 清单(header 逐个核对 org 名, 排除 ISP 骨干/无关公司):"
  for e in "${ASNS[@]}"; do echo "#   AS$(echo "$e" | awk '{print $1}') ($(echo "$e" | awk '{print $2}'))"; done
  echo "#"
} > "$OUT"

for e in "${ASNS[@]}"; do
  asn="$(echo "$e" | awk '{print $1}')"
  for v in ipv4 ipv6; do
    curl -fsSL -m 20 "${RAW}/${asn}/${v}-aggregated.txt" 2>/dev/null >> "$tmp/all.txt" || true
  done
done

# 跳注释/空行/404，取 CIDR，去重 + 版本排序
grep -hE '^[0-9a-fA-F]' "$tmp/all.txt" | grep -vE '404' | sort -u -V >> "$OUT"

n="$(grep -cE '^[0-9a-fA-F]' "$OUT")"
echo "✓ 写出 $OUT （${n} CIDR，源 ipverse@${IPVERSE_SHA}）"
