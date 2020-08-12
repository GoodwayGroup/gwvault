# gwvault - GoodwayGroup Ansible Vault

> `ansible-vault` CLI reimplemented in go

`ansible-vault` is a very powerful tool and we wanted to simplifying the install and management of the tool as a standalone, cross platform tool.

## Basic Usage

Please see [the docs for details on the commands.](./docs/gwvault.md)

Use in place of `ansible-vault`. All commands are reimplemented. The tool will default to asking for your Vault password.

```
$ gwvault -h
NAME:
   gwvault - encryption/decryption utility for Ansible data files

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   encrypt                            encrypt file
   decrypt                            decrypt file
   edit                               edit file and re-encrypt
   rekey                              alter encryption password and re-encrypt
   create                             create a new encrypted file
   view                               view inputs of encrypted file
   encrypt_string, av_encrypt_string  encrypt provided string, output in ansible-vault format
   install-manpage                    Generate and install man page
   version, v                         Print version info
   help, h                            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --vault-password-file VAULT_PASSWORD_FILE          vault password file VAULT_PASSWORD_FILE
   --new-vault-password-file NEW_VAULT_PASSWORD_FILE  new vault password file for rekey NEW_VAULT_PASSWORD_FILE
   --help, -h                                         show help (default: false)
```

## Installation

```
$ curl https://i.jpillora.com/GoodwayGroup/gwvault! | bash
```

## Benchmarks

Benchmarking done using [`bench`](https://github.com/Gabriel439/bench). Execute the `benchmark/run.sh` script to generate a new benchmark.

As compared to `ansible-vault` (v2.9.11 on python v3.8.5), typical actions take a 80% less time to complete.

![image](https://user-images.githubusercontent.com/1429775/89319494-cb4b0400-d645-11ea-9f6a-592900130125.png)

|Action|`ansible-vault`|`gwvault`|
|------|---------------|---------|
| encrypt | 482 ms | **94 ms** |
| decrypt | 499 ms | **96 ms** |
| rekey | 650 ms | **162 ms** |
| encrypt_string | 429 ms | **64 ms** |
| encrypt + decrypt | 1,087 ms | **168 ms** |

See [`./benchmark/results.html`](./benchmark/results.html) for a detailed breakdown of the results after running the benchmark.

## Built With

* go v1.14+
* make
* [github.com/mitchellh/gox](https://github.com/mitchellh/gox)

## Deployment

Run `./release.sh $VERSION`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We employ [git-chglog](https://github.com/git-chglog/git-chglog) to manage the [CHANGELOG.md](CHANGELOG.md). For the versions available, see the [tags on this repository](https://github.com/GoodwayGroup/gwvault/tags).

## Authors

* **Derek Smith** - [@clok](https://github.com/clok)
* **Paulo Black** - [@paulojblack](https://github.com/paulojblack)

See also the list of [contributors](https://github.com/GoodwayGroup/gwvault/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Thank you to [@pbthorste](https://github.com/pbthorste) for doing the heavy lifting on [avtool](https://github.com/pbthorste/avtool)

## Sponsors

[![goodwaygroup][goodwaygroup]](https://goodwaygroup.com)

[goodwaygroup]: https://s3.amazonaws.com/gw-crs-assets/goodwaygroup/logos/ggLogo_sm.png "Goodway Group"
