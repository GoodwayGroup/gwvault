#!/bin/bash -e

cp bm.dec bm.dec.back
ansible-vault rekey --vault-password-file vault.pw --new-vault-password-file new_vault.pw bm.dec
cp bm.dec.back bm.dec
