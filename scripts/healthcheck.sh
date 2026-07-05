#!/usr/bin/env bash
# Verifies that the deployed frontend, backend, and database are healthy.

set -euo pipefail

APP_NAME="${APP_NAME:-project-alpha-shilianjie}"
PAGES_NAME="${PAGES_NAME:-project-alpha-shilianjie}"

FRONTEND="https://${PAGES_NAME}.pages.dev"
BACKEND="https://${APP_NAME}.fly.dev"

RED=$'\033[0;31m'; GRN=$'\033[0;32m'; YEL=$'\033[0;33m'; NC=$'\033[0m'

check() {
  local label="$1" url="$2"
  if curl -fsS --max-time 10 "$url" >/dev/null 2>&1; then
    printf "%b %s\n" "${GRN}✓${NC}" "$label  $url"
  else
    printf "%b %s\n" "${RED}✗${NC}" "$label  $url"
  fi
}

check "Frontend (HTTP)"      "${FRONTEND}/"
check "Backend  healthz"     "${BACKEND}/healthz"
check "Backend  API status"  "${BACKEND}/api/v1"

echo
echo "Tag list sample:"
curl -fsS --max-time 10 "${BACKEND}/api/v1/tags" | head -c 400 || true
echo
