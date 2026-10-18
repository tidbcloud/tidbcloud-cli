#!/bin/sh
# Copyright 2026 PingCAP, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

version="1.0.0-beta.11"
if [ -n "$TICLOUD_VERSION" ]; then
    version=${TICLOUD_VERSION#v}
fi

release_repo='https://github.com/tidbcloud/tidbcloud-cli/releases/download'
repo=$release_repo
if [ -n "$TICLOUD_MIRRORS" ]; then
    repo=${TICLOUD_MIRRORS%/}
    case "$repo" in
        https://?*) ;;
        *) echo "TICLOUD_MIRRORS must use HTTPS." >&2; exit 1 ;;
    esac
    mirror_host=${repo#https://}
    mirror_host=${mirror_host%%/*}
    case "$mirror_host" in
        ''|*@*) echo "TICLOUD_MIRRORS must have a host without credentials." >&2; exit 1 ;;
    esac
    case "$repo" in
        *'?'*|*'#'*|*'\'*|*[[:space:]]*)
            echo "TICLOUD_MIRRORS must not contain a query, fragment, backslash, or whitespace." >&2
            exit 1
            ;;
    esac
fi

case $(uname -s) in
    Linux|linux) os=linux ;;
    Darwin|darwin) os=darwin ;;
    *) os= ;;
esac

if [ -z "$os" ]; then
    echo "OS $(uname -s) not supported." >&2
    exit 1
fi

case $(uname -m) in
    amd64|x86_64) arch=amd64 ;;
    arm64|aarch64) arch=arm64 ;;
    *) arch= ;;
esac

if [ -z "$arch" ]; then
    echo "Architecture  $(uname -m) not supported." >&2
    exit 1
fi

if [ -z "$TICLOUD_HOME" ]; then
    TICLOUD_HOME=$HOME/.ticloud
fi
bin_dir=$TICLOUD_HOME/bin
mkdir -p "$bin_dir"

install_binary() (
    tmp_dir=$(mktemp -d) || exit 1
    trap 'rm -rf "$tmp_dir"' 0
    trap 'exit 1' HUP INT TERM
    archive="ticloud_${version}_${os}_$arch.tar.gz"

    curl -q --proto '=https' --proto-redir '=https' -fsSL "$repo/v${version}/$archive" -o "$tmp_dir/$archive" || exit 1
    curl -q --proto '=https' --proto-redir '=https' -fsSL "$release_repo/v${version}/checksums.txt" -o "$tmp_dir/checksums.txt" || exit 1

    expected_checksum=$(awk -v name="$archive" '
        $2 == name || $2 == "*" name { checksum = $1; matches++ }
        END { if (matches != 1) exit 1; print checksum }
    ' "$tmp_dir/checksums.txt") || {
        echo "Expected exactly one checksum for $archive." >&2
        exit 1
    }

    if command -v sha256sum >/dev/null; then
        actual_checksum=$(sha256sum "$tmp_dir/$archive") || exit 1
    else
        actual_checksum=$(shasum -a 256 "$tmp_dir/$archive") || exit 1
    fi
    actual_checksum=${actual_checksum%% *}
    if [ "$actual_checksum" != "$expected_checksum" ]; then
        echo "Checksum verification failed for $archive." >&2
        exit 1
    fi

    tar -zxf "$tmp_dir/$archive" -C "$bin_dir" || exit 1
)

check_depends() {
    pass=0
    for dependency in curl tar awk mktemp; do
        command -v "$dependency" >/dev/null || {
            echo "Dependency check failed: please install '$dependency' before proceeding."
            pass=1
        }
    done
    if ! command -v sha256sum >/dev/null && ! command -v shasum >/dev/null; then
        echo "Dependency check failed: please install 'sha256sum' or 'shasum' before proceeding."
        pass=1
    fi
    return $pass
}

if ! check_depends; then
    exit 1
fi

if ! install_binary; then
    echo "Failed to download and/or extract ticloud archive."
    exit 1
fi

chmod 755 "$bin_dir/ticloud"

bold=$(tput bold 2>/dev/null)
sgr0=$(tput sgr0 2>/dev/null)

# Refrence: https://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux-unix
shell=$(echo $SHELL | awk 'BEGIN {FS="/";} { print $NF }')
echo "Detected shell: ${bold}$shell${sgr0}"
if [ -f "${HOME}/.${shell}_profile" ]; then
    PROFILE=${HOME}/.${shell}_profile
elif [ -f "${HOME}/.${shell}_login" ]; then
    PROFILE=${HOME}/.${shell}_login
elif [ -f "${HOME}/.${shell}rc" ]; then
    PROFILE=${HOME}/.${shell}rc
else
    PROFILE=${HOME}/.profile
fi
echo "Shell profile:  ${bold}$PROFILE${sgr0}"

case :$PATH: in
    *:$bin_dir:*) : "PATH already contains $bin_dir" ;;
    *) printf '\nexport PATH=%s:$PATH\n' "$bin_dir" >> "$PROFILE"
        echo "$PROFILE has been modified to add ticloud to PATH"
        echo "open a new terminal or ${bold}source ${PROFILE}${sgr0} to use it"
        ;;
esac

echo "Installed path: ${bold}$bin_dir/ticloud${sgr0}"
echo "==============================================="
echo "Have a try:     ${bold}ticloud${sgr0}"
echo "==============================================="
