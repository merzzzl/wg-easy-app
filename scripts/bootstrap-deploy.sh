#!/bin/sh

set -eu

RAW_BASE_URL="${RAW_BASE_URL:-https://raw.githubusercontent.com/merzzzl/wg-easy-app/main}"
TARGET_DIR="$(pwd)"
ENV_FILE="${TARGET_DIR}/.env"
EXAMPLE_FILE="${TARGET_DIR}/.env.example"
TTY_AVAILABLE=0

if [ "$#" -ne 0 ]; then
  printf 'This script does not accept arguments. Run it from the target directory.\n' >&2
  exit 1
fi

if [ -c /dev/tty ] && [ -r /dev/tty ] && [ -w /dev/tty ] && { : >/dev/tty; } 2>/dev/null; then
  TTY_AVAILABLE=1
fi

download_file() {
  source_path="$1"
  target_path="$2"
  source_url="${RAW_BASE_URL}/${source_path}"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$source_url" -o "$target_path"
    return
  fi

  if command -v wget >/dev/null 2>&1; then
    wget -qO "$target_path" "$source_url"
    return
  fi

  printf 'curl or wget is required to download deploy files\n' >&2
  exit 1
}

prompt_value() {
  key="$1"
  default_value="$2"

  if [ "$TTY_AVAILABLE" -eq 1 ]; then
    printf '%s [%s]: ' "$key" "$default_value" > /dev/tty
    IFS= read -r user_value < /dev/tty || user_value=""
  else
    user_value=""
  fi

  if [ -n "$user_value" ]; then
    printf '%s' "$user_value"
    return
  fi

  printf '%s' "$default_value"
}

mkdir -p "$TARGET_DIR"

download_file "compose/compose.yaml" "$TARGET_DIR/compose.yaml"
download_file "compose/Caddyfile" "$TARGET_DIR/Caddyfile"
download_file "compose/.env.example" "$EXAMPLE_FILE"

printf 'Deploy files copied to: %s\n' "$TARGET_DIR"

if [ "$TTY_AVAILABLE" -eq 1 ]; then
  printf 'Fill the variables below. Press Enter to keep the default from .env.example.\n\n' > /dev/tty
else
  printf 'No interactive terminal detected. Using defaults from .env.example.\n\n'
fi

: > "$ENV_FILE"

while IFS= read -r line || [ -n "$line" ]; do
  case "$line" in
    ''|'#'*)
      printf '%s\n' "$line" >> "$ENV_FILE"
      continue
      ;;
  esac

  key=${line%%=*}
  default_value=${line#*=}
  value=$(prompt_value "$key" "$default_value")

  printf '%s=%s\n' "$key" "$value" >> "$ENV_FILE"
done < "$EXAMPLE_FILE"

printf '\nGenerated %s\n' "$ENV_FILE"
printf 'Next step:\n'
printf '  cd %s && docker compose up -d\n' "$TARGET_DIR"
