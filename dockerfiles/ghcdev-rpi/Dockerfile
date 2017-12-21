FROM resin/rpi-raspbian:stretch

RUN apt-get update && apt-get install -y sudo

# need /etc/protocol for stack so that getProtocolByName
# works
# and enough route CAs to talk to github/aws
# and to decompress stuff
# we need a C compiler, apparently
# the haskell stack installer does some dependency stuff
# that is ubuntu aware, though, so maybe see how much of this
# we don't need?
# curl we need to download the installer
RUN apt-get update && apt-get install -y netbase
RUN apt-get update && apt-get install -y ca-certificates
RUN apt-get update && apt-get install -y zlib1g-dev
RUN apt-get update && apt-get install -y xz-utils
RUN apt-get update && apt-get install -y gcc
RUN apt-get update && apt-get install -y curl

# this installs distro stack, but that is pretty
# outdated: 0.1.10.0 when at time of writing,
# latest from 'stack upgrade' is 1.4.0
# RUN apt-get update && apt-get install -y haskell-stack

# or haskell stack websites recommended generic linux
# method:
# RUN curl -sSL https://get.haskellstack.org/ | sh
# 
# or stack is too annoying so install cabal?

RUN apt-get update && apt-get install -y cabal-install
