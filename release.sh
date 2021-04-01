#!/bin/bash

set +e

NAME=gwvault

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

if [[ "$(git tag -l | grep -c "$VERSION" 2>/dev/null)" != "0" ]]; then
  error "Tag $VERSION already exists in this repo. Please use a different version."
  exit 1
fi

h1 "Preparing release of $VERSION for $NAME"

h2 "Updating docs"
make docs
if [[ "$(git status -s docs/${NAME}.* 2>/dev/null | wc -l)" == "0" ]]; then
  note "No changes to docs"
else
  note "Committing changes to docs"
  git add docs/${NAME}.*
  git commit -m "chore(docs): updating docs for version $VERSION"
fi

h2 "Updating CHANGELOG.md"
make changelog
git add CHANGELOG.md
git commit -m "chore(release): $VERSION"

h2 "Tagging version: $VERSION"
git tag "$VERSION"

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
note "Pushing branch: git push origin $BRANCH"
git push origin "$BRANCH"

note "Pushing tag: git push origin $VERSION"
git push origin "$VERSION"
