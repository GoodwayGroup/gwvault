#!/bin/bash -e

cp bm.dec bm.dec.back
ansible-vault decrypt --vault-password-file vault.pw bm.dec
cp bm.dec.back bm.dec
