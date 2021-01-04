FROM debian:buster

RUN apt-get update && apt-get upgrade -y

RUN apt-get update && apt-get install -y sudo

RUN apt-get update && apt-get install -y curl less vim python3-pip
RUN apt-get update && apt-get install -y python3-virtualenv virtualenv

# git is needed for parsl to figure out it's own repo-specific
# version string
RUN apt-get update && apt-get install -y git stgit

# useful stuff to have around
RUN apt-get update && apt-get install -y procps

# for building documentation
RUN apt-get update && apt-get install -y pandoc

# for monitoring visualization
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

# now you can
# virtualenv --python=python3 .
# pip3 install jupyter
# 
ADD cue.shell /
RUN chmod a+x /cue.shell
