#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
SOURCE_DIR="${ROOT_DIR}/compose"
TARGET_DIR="${1:-${ROOT_DIR}/deploy}"

mkdir -p "${TARGET_DIR}"

cp "${SOURCE_DIR}/compose.yaml" "${TARGET_DIR}/compose.yaml"
cp "${SOURCE_DIR}/Caddyfile" "${TARGET_DIR}/Caddyfile"
cp "${SOURCE_DIR}/.env.example" "${TARGET_DIR}/.env.example"

ENV_FILE="${TARGET_DIR}/.env"
EXAMPLE_FILE="${SOURCE_DIR}/.env.example"

printf 'Deploy files copied to: %s\n' "${TARGET_DIR}"
printf 'Fill the variables below. Press Enter to keep the default from .env.example.\n\n'

: > "${ENV_FILE}"

while IFS= read -r line || [ -n "${line}" ]; do
  if [[ -z "${line}" || "${line}" =~ ^[[:space:]]*# ]]; then
    printf '%s\n' "${line}" >> "${ENV_FILE}"
    continue
  fi

  key="${line%%=*}"
  default_value="${line#*=}"

  read -r -p "${key} [${default_value}]: " user_value
  value="${user_value:-${default_value}}"

  printf '%s=%s\n' "${key}" "${value}" >> "${ENV_FILE}"
done < "${EXAMPLE_FILE}"

printf '\nGenerated %s\n' "${ENV_FILE}"
printf 'Next step:\n'
printf '  cd %s && docker compose up -d\n' "${TARGET_DIR}"
