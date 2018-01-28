FROM ubuntu:latest

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
RUN apt-get update && apt-get install -y netbase ca-certificates zlib1g-dev xz-utils gcc curl

# this installs distro stack, but that is pretty
# outdated: 0.1.10.0 when at time of writing,
# latest from 'stack upgrade' is 1.4.0
# RUN apt-get update && apt-get install -y haskell-stack

# or haskell stack websites recommended generic linux
# method:
RUN curl -sSL https://get.haskellstack.org/ | sh

RUN apt-get update && apt-get install -y postgresql-9.5 postgresql-server-dev-9.5

RUN echo "local   all             postgres                                trust" | sudo tee /etc/postgresql/9.5/main/pg_hba.conf 

RUN apt-get update && apt-get install -y texlive-latex-base texlive-latex-recommended
