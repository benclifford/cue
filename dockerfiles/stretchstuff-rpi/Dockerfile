FROM resin/rpi-raspbian:stretch

RUN apt-get update && apt-get install -y sudo curl

WORKDIR /

RUN apt-get update && apt-get install -y dos2unix

# rust
RUN apt-get update && apt-get install -y rustc
RUN apt-get update && apt-get install -y cargo

# erlang

RUN apt-get update && apt-get install -y erlang
RUN apt-get update && apt-get install -y erlang-jiffy

