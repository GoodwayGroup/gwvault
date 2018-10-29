#!/bin/bash -e

if [ "x${GWVAULT_PATH}x" = "xx" ]; then
    echo "Please provide GWVAULT_PATH"
    exit 1
fi

$GWVAULT_PATH/gwvault encrypt_string --vault-password-file vault.pw aStringToEncrypt
