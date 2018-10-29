#!/bin/bash -e

ansible-vault encrypt --vault-password-file vault.pw bm.md
ansible-vault decrypt --vault-password-file vault.pw bm.md
