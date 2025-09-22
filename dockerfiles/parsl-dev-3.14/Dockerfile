FROM debian:forky

RUN apt-get update && apt-get upgrade -y

RUN apt-get update && apt-get install -y sudo openssh-server

RUN apt-get update && apt-get install -y curl less vim

# git is needed for parsl to figure out it's own repo-specific
# version string
RUN apt-get update && apt-get install -y git

# useful stuff to have around
RUN apt-get update && apt-get install -y procps

# for building documentation
RUN apt-get update && apt-get install -y pandoc

# for monitoring visualization
RUN apt-get update && apt-get install -y graphviz wget

# for commandline access to monitoring database
RUN apt-get update && apt-get install -y sqlite3

RUN apt-get update && apt-get install -y stgit


# flake8 needs locales to be installed, otherwise
# it encounters some kind of decoding error reading
# the config file
RUN apt-get update && apt-get install -y locales
ENV LANG en_US.UTF-8
RUN echo en_US.UTF-8 UTF-8 >> /etc/locale.gen
RUN locale-gen

RUN git clone --depth=1 --branch=3.14 https://github.com/python/cpython
RUN apt-get update && apt-get install -y gcc build-essential make pkg-config
WORKDIR /cpython

# dependencies recommended for all optional features,
# at https://devguide.python.org/getting-started/setup-building/#build-dependencies
RUN apt-get update && apt-get install -y build-essential gdb lcov pkg-config libbz2-dev libffi-dev libgdbm-dev libgdbm-compat-dev liblzma-dev libncurses5-dev libreadline6-dev libsqlite3-dev libssl-dev lzma tk-dev uuid-dev zlib1g-dev
# removed lzma-dev for forky - it isn't there but not sure what the right behaviour is

RUN ./configure --with-ensurepip=install && make

RUN make install

# to build ccutils/work queue
# RUN apt-get update && apt-get install -y swig python3-dev libpython3-all-dev 
# configure WQ with this:
# ./configure --prefix=$(pwd)/I --with-python3-path=$(which python3)

# RUN /usr/bin/python3 -m ensurepip --upgrade
RUN python3 -m ensurepip

RUN pip3 install virtualenv

# RUN apt-get update && apt-get install -y postgresql-all

# now you can
# virtualenv --python=python3 .
# pip3 install jupyter
# 
ADD cue.shell /
RUN chmod a+x /cue.shell
