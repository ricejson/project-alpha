#!/usr/bin/env bash
#
# One-command public deployment for project-alpha.
#
#   - Frontend  -> Cloudflare Pages  (https://project-alpha-shilianjie.pages.dev)
#   - Backend   -> Fly.io            (https://project-alpha-shilianjie.fly.dev)
#   - Database  -> Neon free tier     (you provide DATABASE_URL once)
#
# Usage:
#   DATABASE_URL="postgres://user:pass@ep-xxx.region.aws.neon.tech/neondb?sslmode=require" \
#       ./scripts/deploy.sh
#
# Or just run it and paste the URL when prompted.

set -euo pipefail

# ---------- config ----------
APP_NAME="project-alpha-shilianjie"
PRIMARY_REGION="hkg"
PAGES_PROJECT_NAME="project-alpha-shilianjie"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# ---------- pretty output ----------
RED=$'\033[0;31m'; GRN=$'\033[0;32m'; YEL=$'\033[0;33m'; BLU=$'\033[0;34m'; NC=$'\033[0m'
info()  { printf "%b\n" "${BLU}▸${NC} $*"; }
ok()    { printf "%b\n" "${GRN}✓${NC} $*"; }
warn()  { printf "%b\n" "${YEL}!${NC} $*"; }
die()   { printf "%b\n" "${RED}✗${NC} $*" >&2; exit 1; }

# ---------- prerequisites ----------
require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "'$1' is required. Run scripts/install-tools.sh first."
}

require_cmd fly
require_cmd wrangler
require_cmd node
require_cmd npm

# ---------- auth ----------
info "Checking Fly.io auth..."
if ! fly auth whoami >/dev/null 2>&1; then
  die "Not logged in to Fly.io. Run: fly auth signup (GitHub login is fastest)"
fi
ok "Fly.io: $(fly auth whoami)"

info "Checking Cloudflare auth..."
if ! wrangler whoami >/dev/null 2>&1; then
  die "Not logged in to Cloudflare. Run: wrangler login"
fi
ok "Cloudflare: $(wrangler whoami | head -1)"

# ---------- DATABASE_URL ----------
if [[ -z "${DATABASE_URL:-}" ]]; then
  warn "DATABASE_URL not set."
  echo "  Create a free Neon project: https://neon.tech → New Project → copy the connection string."
  echo "  The string looks like: postgres://user:pass@ep-xxx.region.aws.neon.tech/neondb?sslmode=require"
  read -r -p "  Paste DATABASE_URL here: " DATABASE_URL
  [[ -n "${DATABASE_URL}" ]] || die "DATABASE_URL is required"
  export DATABASE_URL
fi
[[ "${DATABASE_URL}" == *"sslmode="* ]] || warn "DATABASE_URL has no sslmode=... — Neon requires SSL"

# ---------- BACKEND on Fly.io ----------
info "Ensuring Fly.io app '${APP_NAME}' exists in region ${PRIMARY_REGION}..."
if ! fly status --app "${APP_NAME}" >/dev/null 2>&1; then
  info "App does not exist yet — creating it..."
  cd "${REPO_ROOT}"
  fly launch --no-deploy \
      --name "${APP_NAME}" \
      --region "${PRIMARY_REGION}" \
      --copy-config \
      --ha=false
  cd "${REPO_ROOT}"
  ok "App created"
else
  ok "App already exists"
fi

info "Setting DATABASE_URL secret on Fly.io..."
fly secrets set DATABASE_URL="${DATABASE_URL}" --app "${APP_NAME}"

info "Deploying backend to Fly.io..."
fly deploy --remote-only --app "${APP_NAME}" "${REPO_ROOT}"

BACKEND_URL="https://${APP_NAME}.fly.dev"
ok "Backend live: ${BACKEND_URL}"

# ---------- wait for health ----------
info "Waiting for /healthz..."
for i in {1..30}; do
  if curl -fsS "${BACKEND_URL}/healthz" >/dev/null 2>&1; then
    ok "Backend healthy"
    break
  fi
  sleep 2
done
curl -fsS "${BACKEND_URL}/healthz" >/dev/null || die "Backend never became healthy — run 'fly logs --app ${APP_NAME}'"

# ---------- seed ----------
info "Seeding production database..."
cd "${REPO_ROOT}/backend"
DATABASE_URL="${DATABASE_URL}" go run ./cmd/seed
cd "${REPO_ROOT}"
ok "Database seeded"

# ---------- FRONTEND on Cloudflare Pages ----------
info "Ensuring Cloudflare Pages project '${PAGES_PROJECT_NAME}' exists..."
if ! wrangler pages project list 2>/dev/null | grep -q "${PAGES_PROJECT_NAME}"; then
  info "Creating Pages project..."
  wrangler pages project create "${PAGES_PROJECT_NAME}" --production-branch=main --compatibility-date=2025-01-01
  ok "Pages project created"
else
  ok "Pages project already exists"
fi

info "Building frontend with production API URL..."
cd "${REPO_ROOT}/frontend"
VITE_API_BASE_URL="${BACKEND_URL}/api/v1" npm run build
ok "Frontend built"

info "Deploying frontend to Cloudflare Pages..."
wrangler pages deploy dist --project-name "${PAGES_PROJECT_NAME}" --branch=main --commit-dirty=true
cd "${REPO_ROOT}"

FRONTEND_URL="https://${PAGES_PROJECT_NAME}.pages.dev"
ok "Frontend live: ${FRONTEND_URL}"

# ---------- summary ----------
cat <<EOF

${GRN}========================================${NC}
${GRN}  Deployment complete${NC}
${GRN}========================================${NC}

  Frontend : ${FRONTEND_URL}
  Backend  : ${BACKEND_URL}
  API base : ${BACKEND_URL}/api/v1

  Health   : ${BACKEND_URL}/healthz

EOF
