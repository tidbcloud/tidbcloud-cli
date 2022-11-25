#!/bin/sh

version="0.1.0-rc2"

repo='https://github.com/tidbcloud/tidbcloud-cli/releases/download'
if [ -n "$TICLOUD_MIRRORS" ]; then
    repo=$TICLOUD_MIRRORS
fi

case $(uname -s) in
    Linux|linux) os=linux ;;
    Darwin|darwin) os=macos ;;
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

install_binary() {
    curl -L "$repo/v${version}/ticloud_${version}_${os}_$arch.tar.gz" -o "/tmp/ticloud_${version}_${os}_$arch.tar.gz" || return 1
    tar -zxf "/tmp/ticloud_${version}_${os}_$arch.tar.gz" -C "$bin_dir" || return 1
    rm "/tmp/ticloud_${version}_${os}_$arch.tar.gz"
    return 0
}

check_depends() {
    pass=0
    command -v curl >/dev/null || {
        echo "Dependency check failed: please install 'curl' before proceeding."
        pass=1
    }
    command -v tar >/dev/null || {
        echo "Dependency check failed: please install 'tar' before proceeding."
        pass=1
    }
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