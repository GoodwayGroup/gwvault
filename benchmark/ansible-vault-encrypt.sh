#!/bin/bash -e

cp bm.enc bm.enc.back
ansible-vault encrypt --vault-password-file vault.pw bm.enc
cp bm.enc.back bm.enc
