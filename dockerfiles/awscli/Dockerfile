FROM debian:bullseye

RUN apt-get update && apt-get install -y sudo

RUN apt-get update && apt-get install -y locales
ENV LANG en_US.UTF-8
RUN echo en_US.UTF-8 UTF-8 >> /etc/locale.gen
RUN locale-gen

RUN apt-get update && apt-get install -y curl unzip groff

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN sudo ./aws/install
