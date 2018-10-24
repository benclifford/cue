FROM debian:stretch

RUN apt-get update && apt-get install -y sudo

RUN apt-get update && apt-get install -y python3-pip

RUN apt-get update && apt-get install -y locales
ENV LANG en_US.UTF-8
RUN echo en_US.UTF-8 UTF-8 >> /etc/locale.gen
RUN locale-gen

RUN pip3 install --system awscli

