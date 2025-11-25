FROM debian:trixie

RUN apt-get update && apt-get upgrade -y

RUN apt-get update && apt-get install -y sudo openssh-server

RUN apt-get update && apt-get install -y curl less vim

RUN apt-get update && apt-get install -y git stgit

# useful stuff to have around
RUN apt-get update && apt-get install -y procps

# for building documentation
RUN apt-get update && apt-get install -y pandoc

# for Parsl monitoring visualization
RUN apt-get update && apt-get install -y graphviz wget

# for commandline access to monitoring database
RUN apt-get update && apt-get install -y sqlite3

# flake8 needs locales to be installed, otherwise
# it encounters some kind of decoding error reading
# the config file
RUN apt-get update && apt-get install -y locales
ENV LANG en_US.UTF-8
RUN echo en_US.UTF-8 UTF-8 >> /etc/locale.gen
RUN locale-gen

RUN apt-get update && apt-get install -y python3 python3-virtualenv

RUN useradd lolhu -m -s /bin/bash

ADD cue.shell /
RUN chmod a+x /cue.shell
