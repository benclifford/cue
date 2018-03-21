FROM debian:stretch

RUN apt-get update && apt-get install -y sudo

RUN apt-get update && apt-get install -y wget

# instructions from https://github.com/erlang/rebar3
WORKDIR /bin
RUN wget https://s3.amazonaws.com/rebar3/rebar3 && chmod +x rebar3

