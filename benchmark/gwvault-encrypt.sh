#!/bin/bash -e

if [ "x${GWVAULT_PATH}x" = "xx" ]; then
    echo "Please provide GWVAULT_PATH"
    exit 1
fi

cp bm.enc bm.enc.back
$GWVAULT_PATH/gwvault encrypt --vault-password-file vault.pw bm.enc
cp bm.enc.back bm.enc
