# 'go' build environment

FROM ubuntu:16.04

RUN apt-get update && apt-get install -y sudo curl

WORKDIR /

RUN curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz

RUN tar -C /usr/local -x -f go1.8.linux-amd64.tar.gz

RUN echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh

