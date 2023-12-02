#!/bin/sh

log() { echo "$@" 1>&2; }

platform=""
case $(uname) in
    Darwin) platform="darwin" ;;
    Linux)  platform="linux" ;;
esac

architecture=""
case $(uname -m) in
    i386)   architecture="386" ;;
    i686)   architecture="386" ;;
    x86_64) architecture="amd64" ;;
    arm64)  architecture="arm64" ;;
esac

# Log to stderr
log "Architecture: $architecture"
log "Platform: $platform"

# Stdout
echo "device_architecture=$architecture"
echo "device_platform=$platform"

protoc_platform=""
case $(uname) in
    Darwin) protoc_platform="osx" ;;
    Linux)  protoc_platform="linux" ;;
esac

protoc_architecture=""
case $(uname -m) in
    i386)   protoc_architecture="x86_32" ;;
    i686)   protoc_architecture="x86_32" ;;
    x86_64) protoc_architecture="x86_64" ;;
    arm64)  protoc_architecture="aarch_64" ;;
esac

# Log to stderr
log "Protoc platform: $protoc_platform"
log "Protoc architecture: $protoc_architecture"

# Stdout
echo "protoc_platform=$protoc_platform"
echo "protoc_architecture=$protoc_architecture"
