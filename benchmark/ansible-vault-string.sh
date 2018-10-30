#!/bin/bash -e

ansible-vault encrypt_string --vault-password-file vault.pw aStringToEncrypt
