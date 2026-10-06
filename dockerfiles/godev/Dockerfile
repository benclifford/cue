# 'go' build environment

FROM debian:trixie

RUN apt-get update && apt-get install -y sudo curl 

RUN apt-get update && apt-get install -y git

WORKDIR /

RUN curl -O https://dl.google.com/go/go1.27.1.linux-amd64.tar.gz
RUN tar -C /usr/local -x -f go1.27.1.linux-amd64.tar.gz

RUN echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh

