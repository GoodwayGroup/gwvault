#!/bin/bash -e

gwvaultPath=''
{
  gw="$(ls ../build/gwvault*darwin*/gwvault)"
  gwvaultPath=${gw::${#gw}-8}
} || {
  echo "Please run 'make release' first."
  exit 1
}

export GWVAULT_PATH="${gwvaultPath}"

echo 'aStr0ngP455w0rd!' >./vault.pw
echo 'aN0therStr0ngP455w0rd!' >./new_vault.pw

curl https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.18.md --output bm.md --silent

cp bm.md bm.enc
cp bm.md bm.dec
ansible-vault encrypt --vault-password-file vault.pw bm.dec

bench \
  './ansible-vault-encrypt.sh' \
  './gwvault-encrypt.sh' \
  './ansible-vault-decrypt.sh' \
  './gwvault-decrypt.sh' \
  './ansible-vault-string.sh' \
  './gwvault-string.sh' \
  './ansible-vault-file.sh' \
  './gwvault-file.sh' \
  './ansible-vault-rekey.sh' \
  './gwvault-rekey.sh' \
  --output results.html -L 15

open results.html

rm ./vault.pw ./new_vault.pw ./bm.* || true
