#!/bin/bash -e

VERSION=${1}

if [ "x${VERSION}x" = "xx" ]; then
    echo "Must supply version number as first argument"
    exit 1
fi

echo "Updating CHANGELOG.md"
auto-changelog -v $VERSION -l false --template keepachangelog && git add CHANGELOG.md
git commit -m "$VERSION"

echo "Tagging version: $VERSION"
git tag $VERSION

echo "Building assets to be uploaded (currently manual process)"
make ci

echo ""
echo "What you still need to do:"
echo "1. Push the tag: git push origin $VERSION"
echo "2. Update the release in github"
echo "3. Add the assets to the release."
echo ""
echo "Done!"
