# gwvault - GoodwayGroup Ansible Vault

> `ansible-vault` CLI reimplemented in go

`ansible-vault` is a very powerful tool and we wanted to simplifying the install and management of the tool as a standalone, cross platform tool.

## Basic Usage

Us in place of `ansible-vault`. All commands are reimplemented except for `encrypt_string` and `rekey` (coming soom!). The tool will default to asking for your Vault password.

```
$ gwvault -h
NAME:
   gwvault - encryption/decryption utility for Ansible data files

USAGE:
   gwvault [global options] command [command options] [arguments...]

COMMANDS:
     encrypt  encrypt file
     decrypt  decrypt file
     edit     edit file and re-encrypt
     create   create a new encypted file
     view     view contents of encrypted file
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --vault-password-file VAULT_PASSWORD_FILE  vault password file VAULT_PASSWORD_FILE
   --help, -h                                 show help
   --version, -v                              print the version
```

## Installation

```
$ curl https://i.jpillora.com/GoodwayGroup/gwvault! | bash
```

## Built With

* go v1.10+
* make
* [github.com/mitchellh/gox](https://github.com/mitchellh/gox)

## Deployment

Run `./release.sh $VERISION`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We employ [auto-changelog](https://www.npmjs.com/package/auto-changelog) to manage the [CHANGELOG.md](CHANGELOG.md). For the versions available, see the [tags on this repository](https://github.com/GoodwayGroup/gwvault/tags).

## Authors

* **Derek Smith** - [@clok](https://github.com/clok)

See also the list of [contributors](https://github.com/GoodwayGroup/gwvault/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Thank you to [@pbthorste](https://github.com/pbthorste) for doing the heavy lifting on [avtool](https://github.com/pbthorste/avtool)

## Sponsors

[![goodwaygroup][goodwaygroup]](https://goodwaygroup.com)

[goodwaygroup]: https://s3.amazonaws.com/gw-crs-assets/goodwaygroup/logos/ggLogo_sm.png "Goodway Group"
