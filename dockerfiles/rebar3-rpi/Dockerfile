FROM resin/rpi-raspbian:stretch

RUN apt-get update && apt-get install -y sudo

# install-time deps
RUN apt-get update && apt-get install -y wget 

# runtime deps for rebar3
RUN apt-get update && apt-get install -y erlang git

# runtime deps for specific erlang packages that I want to compile
# using rebar3:

# jiffy
RUN apt-get update && apt-get install -y gcc g++

# instructions from https://github.com/erlang/rebar3
WORKDIR /bin
RUN wget https://s3.amazonaws.com/rebar3/rebar3 && chmod +x rebar3

