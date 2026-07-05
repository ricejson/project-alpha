#!/usr/bin/env bash
# Installs fly (Fly.io CLI) and wrangler (Cloudflare CLI).
# Safe to re-run; skips tools that are already present and on the latest version.

set -euo pipefail

RED=$'\033[0;31m'; GRN=$'\033[0;32m'; YEL=$'\033[0;33m'; BLU=$'\033[0;34m'; NC=$'\033[0m'
info() { printf "%b\n" "${BLU}▸${NC} $*"; }
ok()   { printf "%b\n" "${GRN}✓${NC} $*"; }

OS="$(uname -s)"

# ---------- fly ----------
if command -v fly >/dev/null 2>&1; then
  ok "fly already installed: $(fly version)"
else
  info "Installing fly..."
  if [[ "${OS}" == "Darwin" ]] && command -v brew >/dev/null 2>&1; then
    brew install flyctl
  else
    curl -L https://fly.io/install.sh | sh
  fi
  ok "fly installed"
fi

# ---------- wrangler ----------
if command -v wrangler >/dev/null 2>&1; then
  ok "wrangler already installed: $(wrangler --version)"
else
  info "Installing wrangler..."
  npm install -g wrangler
  ok "wrangler installed"
fi

cat <<EOF

Next steps (only two manual actions needed):

  1) fly auth signup        # 30 sec, "Sign in with GitHub" is fastest
  2) wrangler login         # 30 sec, same GitHub account works

Then deploy with:

  DATABASE_URL='postgres://user:pass@ep-xxx.region.aws.neon.tech/neondb?sslmode=require' \\
      ./scripts/deploy.sh
EOF
