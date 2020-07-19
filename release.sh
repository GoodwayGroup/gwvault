#!/bin/bash

set +e
set -o noglob

#
# Set Colors
#

bold="\e[1m"
dim="\e[2m"
underline="\e[4m"
blink="\e[5m"
reset="\e[0m"
red="\e[31m"
green="\e[32m"
blue="\e[34m"

#
# Common Output Styles
#

h1() {
  printf "\n${bold}${underline}%s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
h2() {
  printf "\n${bold}%s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
info() {
  printf "${dim}➜ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
success() {
  printf "${green}✔ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
error() {
  printf "${red}${bold}✖ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
warnError() {
  printf "${red}✖ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
warnNotice() {
  printf "${blue}✖ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
note() {
  printf "\n${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}

typeExists() {
  if [ $(type -P $1) ]; then
    return 0
  fi
  return 1
}

if ! typeExists "git-chglog"; then
  error "git-chglog is not installed"
  note "To install run: go get -u github.com/git-chglog/git-chglog/cmd/git-chglog"
  exit 1
fi

VERSION=${1}

if [ "x${VERSION}x" = "xx" ]; then
  error "Must supply version number as first argument"
  exit 1
fi

h1 "Preparing release of $VERSION"

h2 "Tagging version: $VERSION"
git tag $VERSION

h2 "Updating CHANGELOG.md"
git-chglog -o CHANGELOG.md && git add CHANGELOG.md
git commit -m "feat(release): $VERSION"

note "Building assets to be uploaded"
make ci

echo ""
note "What you still need to do:"
info "1. Push the tag: git push origin $VERSION"
info "2. Update the release in github with compiled assets."
echo ""
success "Done!"
