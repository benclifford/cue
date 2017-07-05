# 'racket' build environment

FROM ubuntu:16.04

RUN apt-get update && apt-get install -y sudo

# needed for developer
RUN apt-get update && apt-get install -y emacs-nox

# prereqs for installing idris:
# 'make' is needed for idris build
# 'zlib1g-dev' is needed for hackage zlib build
# 'netbase' is needed to deal with this message:
#   idris: getProtocolByName: does not exist (no such protocol name: tcp)
# that appears when starting idris interactively - I think it is
# looking in /etc/protocols


RUN apt-get update && apt-get install -y cabal-install zlib1g-dev make netbase

# 'pkg-config' for building idris with ffi
RUN apt-get update && apt-get install -y pkg-config

RUN cabal update && cabal install --global -f FFI -f GMP idris

