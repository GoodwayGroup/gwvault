#!/bin/bash -e

if [ "x${GWVAULT_PATH}x" = "xx" ]; then
  echo "Please provide GWVAULT_PATH"
  exit 1
fi

$GWVAULT_PATH/gwvault encrypt --vault-password-file vault.pw bm.md
$GWVAULT_PATH/gwvault decrypt --vault-password-file vault.pw bm.md
