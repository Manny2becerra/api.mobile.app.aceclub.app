#!/usr/bin/env bash
set -euo pipefail

ENV="${1:-prod}"                       # ./fetch-env prod  (default) or ./fetch-env stage
SSM_PATH="/${ENV}/api-embed/"  # note the trailing slash

# 1. Truncate—or create—.env
> .env

# 2. Pull every parameter under the path (auto-paginates, decrypts SecureStrings)
aws ssm get-parameters-by-path \
    --path "$SSM_PATH"                \
    --recursive                       \
    --with-decryption                 \
    --query 'Parameters[].[Name,Value]' \
    --output text |

# 3. Convert full SSM names to KEY="value" lines
while read -r full value; do
  echo "$(basename "$full")=\"${value}\"" >> .env
done

# 4. Lock down the file
chmod 600 .env
echo ".env written for ${ENV}"