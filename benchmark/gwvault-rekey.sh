#!/bin/bash -e

if [ "x${GWVAULT_PATH}x" = "xx" ]; then
  echo "Please provide GWVAULT_PATH"
  exit 1
fi

cp bm.dec bm.dec.back
$GWVAULT_PATH/gwvault rekey --vault-password-file vault.pw --new-vault-password-file new_vault.pw bm.dec
cp bm.dec.back bm.dec
