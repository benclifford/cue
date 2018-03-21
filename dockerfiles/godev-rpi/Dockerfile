# 'go' build environment

FROM resin/rpi-raspbian:stretch

RUN apt-get update && apt-get install -y sudo curl

RUN apt-get update && apt-get install -y git

WORKDIR /

RUN curl -O https://storage.googleapis.com/golang/go1.8.linux-armv6l.tar.gz

RUN tar -C /usr/local -x -f go1.8.linux-armv6l.tar.gz

RUN echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh

