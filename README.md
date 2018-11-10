# gwvault - GoodwayGroup Ansible Vault

> `ansible-vault` CLI reimplemented in go

`ansible-vault` is a very powerful tool and we wanted to simplifying the install and management of the tool as a standalone, cross platform tool.

## Basic Usage

Use in place of `ansible-vault`. All commands are reimplemented. The tool will default to asking for your Vault password.

```
$ gwvault -h
NAME:
   gwvault - encryption/decryption utility for Ansible data files

USAGE:
   gwvault [global options] command [command options] [arguments...]

VERSION:
   1.2.0

COMMANDS:
     encrypt         encrypt file
     decrypt         decrypt file
     edit            edit file and re-encrypt
     rekey           alter encryption password and re-encrypt
     create          create a new encrypted file
     view            view contents of encrypted file
     encrypt_string  encrypt provided string
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --vault-password-file VAULT_PASSWORD_FILE  vault password file VAULT_PASSWORD_FILE
   --help, -h                                 show help
   --version, -v                              print the version
```

## Installation

```
$ curl https://i.jpillora.com/GoodwayGroup/gwvault! | bash
```

## Benchmarks

Benchmarking done using [`bench`](https://github.com/Gabriel439/bench). Execute the `benchmark/run.sh` script to generate a new benchmark.

As compared to `ansible-vault`, typical actions take a 1/10th the time to complete.

|Action|`ansible-vault`|`gwvault`|
|------|---------------|---------|
| encrypt | 633 ms | **76 ms** |
| decrypt | 639 ms | **72 ms** |
| encrypt_string | 536 ms | **42 ms** |
| encrypt + decrypt | 1,160 ms | **105 ms** |

See [`./benchmark/results.html`](./benchmark/results.html) for a detailed breakdown of the results.

## Built With

* go v1.10+
* make
* [github.com/mitchellh/gox](https://github.com/mitchellh/gox)

## Deployment

Run `./release.sh $VERSION`

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
