% gwvault 8
# NAME
gwvault - encryption/decryption utility for Ansible data files
# SYNOPSIS
gwvault

```
[--new-vault-password-file]=[value]
[--vault-password-file]=[value]
```


# COMMAND TREE

- [encrypt](#encrypt)
- [decrypt](#decrypt)
- [edit](#edit)
- [rekey](#rekey)
- [create](#create)
- [view](#view)
- [encrypt_string, av_encrypt_string](#encrypt_string-av_encrypt_string)
- [install-manpage](#install-manpage)
- [version, v](#version-v)

**Usage**:
```
gwvault [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--new-vault-password-file**="": new vault password file for rekey `NEW_VAULT_PASSWORD_FILE`

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`


# COMMANDS

## encrypt

encrypt file

>[options] [vaultfile.yml]

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## decrypt

decrypt file

>[options] [vaultfile.yml]

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## edit

edit file and re-encrypt

>[options] [vaultfile.yml]

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## rekey

alter encryption password and re-encrypt

>[options] [vaultfile.yml] [newvaultfile.yml]

**--new-vault-password-file**="": new vault password file for rekey `NEW_VAULT_PASSWORD_FILE`

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## create

create a new encrypted file

>[options] [vaultfile.yml]

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## view

view inputs of encrypted file

>[options] [vaultfile.yml]

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## encrypt_string, av_encrypt_string

encrypt provided string, output in ansible-vault format

>[options] string_to_encrypt

**--name**="": variable name to encrypt

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

## install-manpage

Generate and install man page

## version, v

Print version info

