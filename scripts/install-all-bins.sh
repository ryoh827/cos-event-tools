#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: install-all-bins.sh [-d DIR] [-n] [-v]

Build and install all Go CLI binaries found under ./cmd/... (package main).

Options:
  -d DIR  Install destination directory (default: ~/bin)
  -n      Dry-run (print actions without building)
  -v      Verbose output
  -h      Show this help
USAGE
}

install_dir="${HOME}/bin"
dry_run=false
verbose=false

while getopts ":d:nvh" opt; do
  case "${opt}" in
    d)
      install_dir="${OPTARG}"
      ;;
    n)
      dry_run=true
      ;;
    v)
      verbose=true
      ;;
    h)
      usage
      exit 0
      ;;
    :) 
      echo "Error: option -${OPTARG} requires an argument" >&2
      usage >&2
      exit 2
      ;;
    \?)
      echo "Error: invalid option -${OPTARG}" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if ! command -v go >/dev/null 2>&1; then
  echo "Error: 'go' command not found in PATH" >&2
  exit 127
fi

if [[ ! -f go.mod ]]; then
  if [[ "${dry_run}" == true ]]; then
    echo "[dry-run] Error: go.mod not found. Please run this script from the module root."
    echo "Dry-run completed."
    exit 0
  fi
  echo "Error: go.mod not found. Please run this script from the module root." >&2
  exit 1
fi

if [[ "${install_dir}" == ~* ]]; then
  install_dir="${install_dir/#\~/${HOME}}"
fi

if [[ "${dry_run}" == false ]]; then
  mkdir -p "${install_dir}"
else
  echo "[dry-run] mkdir -p ${install_dir}"
fi

packages=()
while IFS= read -r pkg; do
  [[ -n "${pkg}" ]] && packages+=("${pkg}")
done < <(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./cmd/...)

if ((${#packages[@]} == 0)); then
  echo "No main packages found under ./cmd/..." >&2
  exit 1
fi

if [[ "${verbose}" == true || "${dry_run}" == true ]]; then
  echo "Detected main packages (${#packages[@]}):"
  for pkg in "${packages[@]}"; do
    echo "  - ${pkg}"
  done
fi

for pkg in "${packages[@]}"; do
  bin_name="${pkg##*/}"
  out_path="${install_dir}/${bin_name}"

  if [[ "${dry_run}" == true ]]; then
    echo "[dry-run] go build -o ${out_path} ${pkg}"
    continue
  fi

  if [[ "${verbose}" == true ]]; then
    echo "Building ${pkg} -> ${out_path}"
  fi

  go build -o "${out_path}" "${pkg}"
done

if [[ "${dry_run}" == true ]]; then
  echo "Dry-run completed."
else
  echo "Installed ${#packages[@]} binaries to ${install_dir}"
fi
