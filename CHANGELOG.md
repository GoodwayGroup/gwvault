<a name="unreleased"></a>
## [Unreleased]


<a name="v2.1.5"></a>
## [v2.1.5] - 2021-06-17
### Pull Requests
- chore(deps): update alpine docker tag to v3.14.0 ([#43](https://github.com/GoodwayGroup/gwvault/issues/43))
- fix(deps): update golang.org/x/term commit hash to 6886f2d ([#42](https://github.com/GoodwayGroup/gwvault/issues/42))
- ci: update goreleaser token ([#41](https://github.com/GoodwayGroup/gwvault/issues/41))


<a name="v2.1.4"></a>
## [v2.1.4] - 2021-05-13
### Bug Fixes
- **deps:** update module github.com/clok/avtool to v2.1.1

### Chore
- **release:** 2.1.4

### Pull Requests
- fix(deps): update golang.org/x/term commit hash to a79de54 ([#39](https://github.com/GoodwayGroup/gwvault/issues/39))
- fix(deps): update module github.com/clok/cdocs to v1.2.0 ([#40](https://github.com/GoodwayGroup/gwvault/issues/40))
- fix(deps): update golang.org/x/term commit hash to f5beecf ([#37](https://github.com/GoodwayGroup/gwvault/issues/37))
- fix(deps): update module github.com/alecaivazis/survey/v2 to v2.2.12 ([#38](https://github.com/GoodwayGroup/gwvault/issues/38))
- chore(deps): update alpine docker tag to v3.13.5 ([#36](https://github.com/GoodwayGroup/gwvault/issues/36))
- fix(deps): update golang.org/x/term commit hash to 72f3dc4 ([#35](https://github.com/GoodwayGroup/gwvault/issues/35))


<a name="v2.1.3"></a>
## [v2.1.3] - 2021-04-05
### Chore
- **release:** v2.1.3

### Pull Requests
- fix(deps): update all non-major dependencies ([#34](https://github.com/GoodwayGroup/gwvault/issues/34))
- fix(deps): update golang.org/x/term commit hash to de623e6 ([#33](https://github.com/GoodwayGroup/gwvault/issues/33))


<a name="v2.1.2"></a>
## [v2.1.2] - 2021-04-01
### Chore
- update README with new install methods
- **release:** v2.1.2

### Ci
- **docker:** add goreleaser docker GitHub Container Registry support


<a name="v2.1.1"></a>
## [v2.1.1] - 2021-03-16
### Chore
- remove deprecated github action
- **goreleaser:** remove workflow trigger
- **release:** update homebrew push

### Features
- **release:** v2.1.1


<a name="v2.1.0"></a>
## [v2.1.0] - 2021-03-15
### Bug Fixes
- **build:** correct naming of build artifact

### Chore
- **ci:** remove deprecated release commands
- **ci:** port to using golangci-lint and goreleaser github action
- **go.mod:** bump to go 1.16

### Features
- **release:** v2.1.0


<a name="v2.0.1"></a>
## [v2.0.1] - 2021-03-03
### Bug Fixes
- **v2:** bump go.mod to v2

### Chore
- **github actions:** add go proxy warming

### Features
- **release:** v2.0.1


<a name="v2.0.0"></a>
## [v2.0.0] - 2021-03-02
### Chore
- **renovate:** add extension for group:allNonMajor

### Features
- **avtool:** port to clok/avtool/v2 module
- **release:** v2.0.0


<a name="v1.8.3"></a>
## [v1.8.3] - 2021-02-26
### Bug Fixes
- **deps:** update module github.com/alecaivazis/survey/v2 to v2.2.8
- **deps:** update golang.org/x/term commit hash to 6a3ed07
- **deps:** update golang.org/x/crypto commit hash to 5ea612d

### Chore
- **renovate:** add gomodTidy option

### Features
- **release:** v1.8.3

### Pull Requests
- chore(deps): update golang.org/x/term commit hash to 2321bbc ([#25](https://github.com/GoodwayGroup/gwvault/issues/25))


<a name="v1.8.2"></a>
## [v1.8.2] - 2020-12-29
### Chore
- **deps:** port to using golang.org/x/term
- **deps:** update module urfave/cli/v2 to v2.3.0
- **deps:** update module clok/cdocs tp v0.2.3
- **docs:** updated README with asdf plugin installation
- **workflow:** add github workflow for linting

### Features
- **release:** v1.8.2

### Pull Requests
- chore(deps): update golang.org/x/crypto commit hash to eec23a3 ([#24](https://github.com/GoodwayGroup/gwvault/issues/24))
- chore(deps): update module alecaivazis/survey/v2 to v2.2.7 ([#23](https://github.com/GoodwayGroup/gwvault/issues/23))


<a name="v1.8.1"></a>
## [v1.8.1] - 2020-10-28
### Bug Fixes
- **file rename:** deprecate use of os.Rename in favor of ioutil.WriteFile

### Features
- **release:** v1.8.1

### Pull Requests
- fix(file rename): Overwrite original file rather than use Rename ([#21](https://github.com/GoodwayGroup/gwvault/issues/21))
- chore(deps): update golang.org/x/crypto commit hash to 9e8e0b3 ([#19](https://github.com/GoodwayGroup/gwvault/issues/19))


<a name="v1.8.0"></a>
## [v1.8.0] - 2020-10-13
### Features
- **release:** v1.8.0

### Pull Requests
- chore(deps): update golang.org/x/crypto commit hash to 84dcc77 ([#17](https://github.com/GoodwayGroup/gwvault/issues/17))
- fix(windows) Fix Windows specific issue: Can't rename open files ([#18](https://github.com/GoodwayGroup/gwvault/issues/18))


<a name="v1.7.2"></a>
## [v1.7.2] - 2020-08-25
### Chore
- updated release script to include publish to github
- **make:** don't update go.mod with gox

### Features
- **release:** v1.7.2

### Pull Requests
- chore(deps): update module clok/kemba to v0.7.1 ([#16](https://github.com/GoodwayGroup/gwvault/issues/16))
- chore(deps): update golang.org/x/crypto commit hash to 5c72a88 ([#13](https://github.com/GoodwayGroup/gwvault/issues/13))


<a name="v1.7.1"></a>
## [v1.7.1] - 2020-08-21
### Chore
- **deps:** update clok/kemba and clok/cdocs
- **renovate:** clean up dupe config
- **renovate:** add renovate.json

### Features
- **release:** v1.7.1

### Pull Requests
- chore(deps): update module alecaivazis/survey/v2 to v2.1.1 ([#12](https://github.com/GoodwayGroup/gwvault/issues/12))
- chore(deps): add renovate.json ([#11](https://github.com/GoodwayGroup/gwvault/issues/11))


<a name="v1.7.0"></a>
## [v1.7.0] - 2020-08-13
### Chore
- update README.md
- **docs:** updating docs for version v1.7.0

### Features
- **release:** v1.7.0

### Fest
- **cdocs:** integrate cdocs library


<a name="v1.6.0"></a>
## [v1.6.0] - 2020-08-12
### Chore
- updated make and release process
- **docs:** initial generation of docs

### Features
- **docs:** updated cli to v2, added in docs generation
- **release:** v1.6.0


<a name="v1.5.0"></a>
## [v1.5.0] - 2020-08-04
### Chore
- remove benchmark results html to decrease bloat

### Features
- **release:** v1.5.0
- **release:** 1.4.0

### Pull Requests
- feat(modernize): support go.mod, cleaned up code, brought in line with current ansible-vault ([#10](https://github.com/GoodwayGroup/gwvault/issues/10))


<a name="1.4.0"></a>
## [1.4.0] - 2020-07-19
### Chore
- update .gitignore to exclude bin dir

### DevOps
- updated release process and changelog to git-chglog

### Pull Requests
- feat: Read ANSIBLE_VAULT_PASSWORD_FILE env variable if no password provided ([#9](https://github.com/GoodwayGroup/gwvault/issues/9))


<a name="1.3.0"></a>
## [1.3.0] - 2019-08-19
### Pull Requests
- Support check for TTY terminal when using `view` ([#7](https://github.com/GoodwayGroup/gwvault/issues/7))


<a name="1.2.1"></a>
## [1.2.1] - 2019-08-18
### Pull Requests
- Patch: Use `cat` instead of `more` ([#6](https://github.com/GoodwayGroup/gwvault/issues/6))
- Support rekey method ([#4](https://github.com/GoodwayGroup/gwvault/issues/4))
- Merge pull request [#3](https://github.com/GoodwayGroup/gwvault/issues/3) from GoodwayGroup/release/v1.2.0


<a name="1.2.0"></a>
## [1.2.0] - 2018-10-30
### Pull Requests
- Support encrypt_string method ([#2](https://github.com/GoodwayGroup/gwvault/issues/2))


<a name="1.1.0"></a>
## [1.1.0] - 2018-08-20
### Pull Requests
- Added support for file globs ([#1](https://github.com/GoodwayGroup/gwvault/issues/1))


<a name="1.0.1"></a>
## [1.0.1] - 2018-08-17

<a name="1.0.0"></a>
## 1.0.0 - 2018-08-17

[Unreleased]: https://github.com/GoodwayGroup/gwvault/compare/v2.1.5...HEAD
[v2.1.5]: https://github.com/GoodwayGroup/gwvault/compare/v2.1.4...v2.1.5
[v2.1.4]: https://github.com/GoodwayGroup/gwvault/compare/v2.1.3...v2.1.4
[v2.1.3]: https://github.com/GoodwayGroup/gwvault/compare/v2.1.2...v2.1.3
[v2.1.2]: https://github.com/GoodwayGroup/gwvault/compare/v2.1.1...v2.1.2
[v2.1.1]: https://github.com/GoodwayGroup/gwvault/compare/v2.1.0...v2.1.1
[v2.1.0]: https://github.com/GoodwayGroup/gwvault/compare/v2.0.1...v2.1.0
[v2.0.1]: https://github.com/GoodwayGroup/gwvault/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/GoodwayGroup/gwvault/compare/v1.8.3...v2.0.0
[v1.8.3]: https://github.com/GoodwayGroup/gwvault/compare/v1.8.2...v1.8.3
[v1.8.2]: https://github.com/GoodwayGroup/gwvault/compare/v1.8.1...v1.8.2
[v1.8.1]: https://github.com/GoodwayGroup/gwvault/compare/v1.8.0...v1.8.1
[v1.8.0]: https://github.com/GoodwayGroup/gwvault/compare/v1.7.2...v1.8.0
[v1.7.2]: https://github.com/GoodwayGroup/gwvault/compare/v1.7.1...v1.7.2
[v1.7.1]: https://github.com/GoodwayGroup/gwvault/compare/v1.7.0...v1.7.1
[v1.7.0]: https://github.com/GoodwayGroup/gwvault/compare/v1.6.0...v1.7.0
[v1.6.0]: https://github.com/GoodwayGroup/gwvault/compare/v1.5.0...v1.6.0
[v1.5.0]: https://github.com/GoodwayGroup/gwvault/compare/1.4.0...v1.5.0
[1.4.0]: https://github.com/GoodwayGroup/gwvault/compare/1.3.0...1.4.0
[1.3.0]: https://github.com/GoodwayGroup/gwvault/compare/1.2.1...1.3.0
[1.2.1]: https://github.com/GoodwayGroup/gwvault/compare/1.2.0...1.2.1
[1.2.0]: https://github.com/GoodwayGroup/gwvault/compare/1.1.0...1.2.0
[1.1.0]: https://github.com/GoodwayGroup/gwvault/compare/1.0.1...1.1.0
[1.0.1]: https://github.com/GoodwayGroup/gwvault/compare/1.0.0...1.0.1
