#!/usr/bin/env bash
set -euo pipefail

OLD_MODULE="project-layout"
YES=false
FORCE=false

usage() {
  cat <<EOF
Usage: $(basename "$0") <new-module-path> [--yes] [--force]

Example:
  $(basename "$0") github.com/antyiwei/my-api
  $(basename "$0") github.com/antyiwei/my-api --yes

Replaces module path "${OLD_MODULE}" with your new module path,
then runs go mod tidy and go build ./...
EOF
}

for arg in "$@"; do
  case "$arg" in
    --yes) YES=true ;;
    --force) FORCE=true ;;
    -h|--help)
      usage
      exit 0
      ;;
  esac
done

NEW_MODULE="${1:-}"
if [[ -z "$NEW_MODULE" || "$NEW_MODULE" == --* ]]; then
  usage
  exit 1
fi

if [[ "$NEW_MODULE" != */* ]]; then
  echo "error: module path must contain a slash, e.g. github.com/user/repo" >&2
  exit 1
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ "$NEW_MODULE" == "$OLD_MODULE" ]]; then
  echo "error: new module path is the same as the template default" >&2
  exit 1
fi

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  remote_url="$(git remote get-url origin 2>/dev/null || true)"
  if [[ "$remote_url" == *"antyiwei/project-layout"* && "$FORCE" != true ]]; then
    echo "error: this looks like the template source repo." >&2
    echo "Run from a project created via 'Use this template', or pass --force." >&2
    exit 1
  fi
fi

echo "Will replace:"
echo "  ${OLD_MODULE} -> ${NEW_MODULE}"
echo
echo "Files to update (excluding docs/):"
grep -rl "$OLD_MODULE" . \
  --exclude-dir=.git \
  --exclude-dir=docs \
  --exclude="init-project.sh" \
  2>/dev/null || true
echo

if [[ "$YES" != true ]]; then
  read -r -p "Continue? [y/N] " reply
  if [[ ! "$reply" =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
  fi
fi

if [[ "$(uname)" == "Darwin" ]]; then
  SED_INPLACE=(-i '')
else
  SED_INPLACE=(-i)
fi

while IFS= read -r -d '' file; do
  sed "${SED_INPLACE[@]}" "s#${OLD_MODULE}#${NEW_MODULE}#g" "$file"
done < <(grep -rl "$OLD_MODULE" . \
  --exclude-dir=.git \
  --exclude-dir=docs \
  --exclude="init-project.sh" \
  -print0 2>/dev/null || true)

sed "${SED_INPLACE[@]}" "s#^module ${OLD_MODULE}#module ${NEW_MODULE}#" go.mod

echo "Running go mod tidy..."
go mod tidy

echo "Running go build ./..."
go build ./...

cat <<EOF

Done! Module path updated to: ${NEW_MODULE}

Next steps:
  make build
  make run
  git add -A && git commit -m "chore: init module path"
EOF
