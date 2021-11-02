#!/bin/sh

#
# USAGE:
#   ./semver.sh
#

set -eu

# Computes a semantic version from the latest git tag
# Read more: https://semver.org
semver() {
  git_status=$(git status --porcelain)
	git_describe=$(git describe --tags 2> /dev/null || true)

  if [ -z "$git_describe" ]; then  # No git tag and no previous semantic version --> using the default initial semantic version
    commit_count=$(git rev-list --count HEAD)
		git_sha=$(git rev-parse --short HEAD)

    if [ -z "$git_status" ]; then
      version="0.1.0-${commit_count}.dev"
    else
      version="0.1.0-${commit_count}.${git_sha}"
    fi
  else  # there is at least one git --> interpret the tag as a semantic version
	  # The tag refers to HEAD commit --> current semantic version
	  # Example: v0.2.7
	  release=$(echo "$git_describe" | grep -E -o '^v?[0-9]+\.[0-9]+\.[0-9]+$' || true)
    if [ -n "$release" ]; then
      semver=$(echo "$release" | grep -E -o '[0-9]+\.[0-9]+\.[0-9]+')
		  major=$(echo "$semver" | cut -d '.' -f 1)
		  minor=$(echo "$semver" | cut -d '.' -f 2)
		  patch=$(echo "$semver" | cut -d '.' -f 3)

      next_patch=$(( patch + 1 ))

      if [ -z "$git_status" ]; then
			  version="${major}.${minor}.${patch}"
		  else
		    version="${major}.${minor}.${next_patch}-0.dev"
		  fi
    fi

	  # The tag refers to a previous commit --> next semantic version + pre-release version
	  # Example: v0.2.7-10-gabcdef
	  prerelease=$(echo "$git_describe" | grep -E -o '^v?[0-9]+\.[0-9]+\.[0-9]+-[0-9]+-g[0-9a-f]+$' || true)
    if [ -n "$prerelease" ]; then
      semver=$(echo "$prerelease" | grep -E -o '[0-9]+\.[0-9]+\.[0-9]+')
		  major=$(echo "$semver" | cut -d '.' -f 1)
		  minor=$(echo "$semver" | cut -d '.' -f 2)
		  patch=$(echo "$semver" | cut -d '.' -f 3)

      next_patch=$(( patch + 1 ))
      commit_count=$(echo "$prerelease" | cut -d '-' -f 2)
		  git_sha=$(git rev-parse --short HEAD)

      if [ -z "$git_status" ]; then
			  version="${major}.${minor}.${next_patch}-${commit_count}.${git_sha}"
		  else
			  version="${major}.${minor}.${next_patch}-${commit_count}.dev"
		  fi
    fi
  fi
}


semver
printf "$version"
