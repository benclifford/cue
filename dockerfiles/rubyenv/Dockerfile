FROM ubuntu:latest

RUN apt-get update && apt-get install -y ruby
RUN apt-get update && apt-get install -y ruby-dev


# licensee needs (in its stack somewhere) make

# TODO: merge these lines:
RUN apt-get update && apt-get install -y make
RUN apt-get install -y cmake
RUN apt-get install -y pkg-config
RUN gem install licensee

RUN apt-get install -y sudo
