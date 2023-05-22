#!/bin/bash
set -e

usage() {
  this=$1
  cat <<EOF
$this: download go binaries for kirychukyurii/wdeploy and install ansible

Usage: $this [-b] bindir [-d] [tag]
  -b sets bindir or installation directory, defaults to $HOME/.local/bin
  -d turns on debug logging

  [tag] is a tag from https://github.com/kirychukyurii/wdeploy/releases
        If tag is missing, then the latest will be used.

EOF
  exit 2
}

parse_args() {
  # BINDIR is ./bin unless set be ENV
  # over-ridden by flag below

  BINDIR=${BINDIR:-$HOME/.local/bin}
  while getopts "b:dh?x" arg; do
    case "$arg" in
      b) BINDIR="$OPTARG" ;;
      d) log_set_priority 10 ;;
      h | \?) usage "$0" ;;
      x) set -x ;;
    esac
  done
  shift $((OPTIND - 1))
  TAG=$1
}

# this function wraps all the destructive operations
# if a curl|bash cuts off the end of the script due to
# network, either nothing will happen or will syntax error
# out preventing half-done work
execute() {
  if [[ ":$PATH:" != *":${BINDIR}:"*  ]] ; then
    log_info "your path is missing ${BINDIR}, you might want to add it."
    log_info "temporary adding ${BINDIR} to PATH environment"
    PATH=$PATH:$HOME/.local/bin
    log_debug "PATH: ${PATH}"
  fi

  tmpdir=$(mktemp -d)
  log_debug "downloading files into ${tmpdir}"
  http_download "${tmpdir}/${TARBALL}" "${TARBALL_URL}"
  http_download "${tmpdir}/${CHECKSUM}" "${CHECKSUM_URL}"
  hash_sha256_verify "${tmpdir}/${TARBALL}" "${tmpdir}/${CHECKSUM}"
  srcdir="${tmpdir}/${NAME}"
  rm -rf "${srcdir}"
  (cd "${tmpdir}" && untar "${TARBALL}")
  test ! -d "${BINDIR}" && install -d "${BINDIR}"
  for binexe in $BINARIES; do
    if [ "$OS" = "windows" ]; then
      binexe="${binexe}.exe"
    fi
    install "${tmpdir}/${binexe}" "${BINDIR}/"
    log_info "installed ${BINDIR}/${binexe}"
  done
  rm -rf "${tmpdir}"

  tmpdir=$(mktemp -d)
  log_debug "checking if Ansible is installed"
  if ! is_command ansible; then
    http_download "${tmpdir}/${PIP_INSTALL_SCRIPT_NAME}" "${PIP_INSTALL_SCRIPT}"
    python3 "${tmpdir}/${PIP_INSTALL_SCRIPT_NAME}" --user 1> /dev/null
    python3 -m pip install --user ansible 1> /dev/null

    if is_command ansible; then
      log_info "installed Ansible for user via python pip"
    fi
  fi
  http_download "${tmpdir}/${ANSIBLE_REQUIREMENTS_NAME}" "${ANSIBLE_REQUIREMENTS}"
  log_info "starting galaxy collection install process"
  ansible-galaxy collection install -r "${tmpdir}/${ANSIBLE_REQUIREMENTS_NAME}" 1> /dev/null
  log_info "all requested collections are installed"
  rm -rf "${tmpdir}"
  
  wdeploy run --help
}

get_binaries() {
  case "$PLATFORM" in
    darwin/amd64) BINARIES="wdeploy" ;;
    darwin/arm64) BINARIES="wdeploy" ;;
    darwin/armv6) BINARIES="wdeploy" ;;
    darwin/armv7) BINARIES="wdeploy" ;;
    darwin/mips64) BINARIES="wdeploy" ;;
    darwin/mips64le) BINARIES="wdeploy" ;;
    darwin/ppc64le) BINARIES="wdeploy" ;;
    darwin/s390x) BINARIES="wdeploy" ;;
    freebsd/386) BINARIES="wdeploy" ;;
    freebsd/amd64) BINARIES="wdeploy" ;;
    freebsd/armv6) BINARIES="wdeploy" ;;
    freebsd/armv7) BINARIES="wdeploy" ;;
    freebsd/mips64) BINARIES="wdeploy" ;;
    freebsd/mips64le) BINARIES="wdeploy" ;;
    freebsd/ppc64le) BINARIES="wdeploy" ;;
    freebsd/s390x) BINARIES="wdeploy" ;;
    linux/386) BINARIES="wdeploy" ;;
    linux/amd64) BINARIES="wdeploy" ;;
    linux/arm64) BINARIES="wdeploy" ;;
    linux/armv6) BINARIES="wdeploy" ;;
    linux/armv7) BINARIES="wdeploy" ;;
    linux/mips64) BINARIES="wdeploy" ;;
    linux/mips64le) BINARIES="wdeploy" ;;
    linux/ppc64le) BINARIES="wdeploy" ;;
    linux/s390x) BINARIES="wdeploy" ;;
    linux/riscv64) BINARIES="wdeploy" ;;
    linux/loong64) BINARIES="wdeploy" ;;
    netbsd/386) BINARIES="wdeploy" ;;
    netbsd/amd64) BINARIES="wdeploy" ;;
    netbsd/armv6) BINARIES="wdeploy" ;;
    netbsd/armv7) BINARIES="wdeploy" ;;
    windows/386) BINARIES="wdeploy" ;;
    windows/amd64) BINARIES="wdeploy" ;;
    windows/arm64) BINARIES="wdeploy" ;;
    windows/armv6) BINARIES="wdeploy" ;;
    windows/armv7) BINARIES="wdeploy" ;;
    windows/mips64) BINARIES="wdeploy" ;;
    windows/mips64le) BINARIES="wdeploy" ;;
    windows/ppc64le) BINARIES="wdeploy" ;;
    windows/s390x) BINARIES="wdeploy" ;;
    *)
      log_crit "platform $PLATFORM is not supported.  Make sure this script is up-to-date and file request at https://github.com/${PREFIX}/issues/new"
      exit 1
      ;;
  esac
}

tag_to_version() {
  if [ -z "${TAG}" ]; then
    log_info "checking GitHub for latest tag"
  else
    log_info "checking GitHub for tag '${TAG}'"
  fi
  REALTAG=$(github_release "$OWNER/$REPO" "${TAG}") && true
  if test -z "$REALTAG"; then
    log_crit "unable to find '${TAG}' - use 'latest' or see https://github.com/${PREFIX}/releases for details"
    exit 1
  fi
  # if version starts with 'v', remove it
  TAG="$REALTAG"
  VERSION=${TAG#v}
}

adjust_format() {
  # change format (tar.gz or zip) based on OS
  case ${OS} in
    windows) FORMAT=zip ;;
  esac
  true
}

adjust_os() {
  # adjust archive name based on OS
  true
}

adjust_arch() {
  # adjust archive name based on ARCH
  true
}

is_command() {
  command -v "$1" >/dev/null
}

echoerr() {
  echo "$@" 1>&2
}

log_prefix() {
  echo "$0"
}

_logp=6
log_set_priority() {
  _logp="$1"
}

log_priority() {
  if test -z "$1"; then
    echo "$_logp"
    return
  fi
  [ "$1" -le "$_logp" ]
}

log_tag() {
  case $1 in
    0) echo "emerg" ;;
    1) echo "alert" ;;
    2) echo "crit" ;;
    3) echo "err" ;;
    4) echo "warning" ;;
    5) echo "notice" ;;
    6) echo "info" ;;
    7) echo "debug" ;;
    *) echo "$1" ;;
  esac
}

log_debug() {
  log_priority 7 || return 0
  echoerr "$(log_prefix)" "$(log_tag 7)" "$@"
}

log_info() {
  log_priority 6 || return 0
  echoerr "$(log_prefix)" "$(log_tag 6)" "$@"
}

log_err() {
  log_priority 3 || return 0
  echoerr "$(log_prefix)" "$(log_tag 3)" "$@"
}

log_crit() {
  log_priority 2 || return 0
  echoerr "$(log_prefix)" "$(log_tag 2)" "$@"
}

uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    msys*) os="windows" ;;
    mingw*) os="windows" ;;
    cygwin*) os="windows" ;;
    win*) os="windows" ;;
  esac
  echo "$os"
}

uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    x86) arch="386" ;;
    i686) arch="386" ;;
    i386) arch="386" ;;
    aarch64) arch="arm64" ;;
    armv5*) arch="armv5" ;;
    armv6*) arch="armv6" ;;
    armv7*) arch="armv7" ;;
    loongarch64) arch="loong64" ;;
  esac
  echo ${arch}
}

uname_os_check() {
  os=$(uname_os)
  case "$os" in
    darwin) return 0 ;;
    dragonfly) return 0 ;;
    freebsd) return 0 ;;
    linux) return 0 ;;
    android) return 0 ;;
    nacl) return 0 ;;
    netbsd) return 0 ;;
    openbsd) return 0 ;;
    plan9) return 0 ;;
    solaris) return 0 ;;
    windows) return 0 ;;
  esac
  log_crit "uname_os_check '$(uname -s)' got converted to '$os' which is not a GOOS value."
  return 1
}

uname_arch_check() {
  arch=$(uname_arch)
  case "$arch" in
    386) return 0 ;;
    amd64) return 0 ;;
    arm64) return 0 ;;
    armv5) return 0 ;;
    armv6) return 0 ;;
    armv7) return 0 ;;
    ppc64) return 0 ;;
    ppc64le) return 0 ;;
    mips) return 0 ;;
    mipsle) return 0 ;;
    mips64) return 0 ;;
    mips64le) return 0 ;;
    s390x) return 0 ;;
    riscv64) return 0 ;;
    amd64p32) return 0 ;;
    loong64) return 0 ;;
  esac
  log_crit "uname_arch_check '$(uname -m)' got converted to '$arch' which is not a GOARCH value."
  return 1
}

untar() {
  tarball=$1
  case "${tarball}" in
    *.tar.gz | *.tgz) tar --no-same-owner -xzf "${tarball}" ;;
    *.tar) tar --no-same-owner -xf "${tarball}" ;;
    *.zip) unzip "${tarball}" ;;
    *)
      log_err "untar unknown archive format for ${tarball}"
      return 1
      ;;
  esac
}

http_download_curl() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    code=$(curl -w '%{http_code}' -sL -o "$local_file" "$source_url")
  else
    code=$(curl -w '%{http_code}' -sL -H "$header" -o "$local_file" "$source_url")
  fi
  if [ "$code" != "200" ]; then
    log_debug "http_download_curl received HTTP status $code"
    return 1
  fi
  return 0
}

http_download_wget() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    wget -q -O "$local_file" "$source_url"
  else
    wget -q --header "$header" -O "$local_file" "$source_url"
  fi
}

http_download() {
  log_debug "http_download $2"
  if is_command curl; then
    http_download_curl "$@"
    return
  elif is_command wget; then
    http_download_wget "$@"
    return
  fi
  log_crit "http_download unable to find wget or curl"
  return 1
}

http_copy() {
  tmp=$(mktemp)
  http_download "${tmp}" "$1" "$2" || return 1
  body=$(cat "$tmp")
  rm -f "${tmp}"
  echo "$body"
}

github_release() {
  owner_repo=$1
  version=$2
  test -z "$version" && version="latest"
  giturl="https://github.com/${owner_repo}/releases/${version}"
  json=$(http_copy "$giturl" "Accept:application/json")
  test -z "$json" && return 1
  version=$(echo "$json" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//')
  test -z "$version" && return 1
  echo "$version"
}

hash_sha256() {
  TARGET=${1:-/dev/stdin}
  if is_command gsha256sum; then
    hash=$(gsha256sum "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command sha256sum; then
    hash=$(sha256sum "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command shasum; then
    hash=$(shasum -a 256 "$TARGET" 2>/dev/null) || return 1
    echo "$hash" | cut -d ' ' -f 1
  elif is_command openssl; then
    hash=$(openssl -dst openssl dgst -sha256 "$TARGET") || return 1
    echo "$hash" | cut -d ' ' -f a
  else
    log_crit "hash_sha256 unable to find command to compute sha-256 hash"
    return 1
  fi
}

hash_sha256_verify() {
  TARGET=$1
  checksums=$2
  if [ -z "$checksums" ]; then
    log_err "hash_sha256_verify checksum file not specified in arg2"
    return 1
  fi
  BASENAME=${TARGET##*/}
  want=$(grep "${BASENAME}" "${checksums}" 2>/dev/null | tr '\t' ' ' | cut -d ' ' -f 1)
  if [ -z "$want" ]; then
    log_err "hash_sha256_verify unable to find checksum for '${TARGET}' in '${checksums}'"
    return 1
  fi
  got=$(hash_sha256 "$TARGET")
  if [ "$want" != "$got" ]; then
    log_err "hash_sha256_verify checksum for '$TARGET' did not verify ${want} vs $got"
    return 1
  fi
}

PROJECT_NAME="wdeploy"
OWNER=kirychukyurii
REPO="wdeploy"
BINARY=wdeploy
FORMAT=tar.gz
OS=$(uname_os)
ARCH=$(uname_arch)
PREFIX="$OWNER/$REPO"

# use in logging routines
log_prefix() {
	echo "$PREFIX"
}

PLATFORM="${OS}/${ARCH}"
GITHUB_DOWNLOAD=https://github.com/${OWNER}/${REPO}/releases/download

uname_os_check "$OS"
uname_arch_check "$ARCH"

parse_args "$@"

get_binaries

tag_to_version

adjust_format

adjust_os

adjust_arch

log_info "found version: ${VERSION} for ${TAG}/${OS}/${ARCH}"

NAME=${BINARY}-${VERSION}-${OS}-${ARCH}
TARBALL=${NAME}.${FORMAT}
TARBALL_URL=${GITHUB_DOWNLOAD}/${TAG}/${TARBALL}
CHECKSUM=${PROJECT_NAME}-${VERSION}-checksums.txt
CHECKSUM_URL=${GITHUB_DOWNLOAD}/${TAG}/${CHECKSUM}

PIP_INSTALL_SCRIPT="https://bootstrap.pypa.io/get-pip.py"
PIP_INSTALL_SCRIPT_NAME="get-pip.py"

ANSIBLE_REQUIREMENTS="https://raw.githubusercontent.com/kirychukyurii/wansible/main/requirements.yml"
ANSIBLE_REQUIREMENTS_NAME="requirements.yml"

execute