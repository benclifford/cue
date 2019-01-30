FROM debian:stretch

RUN apt-get update && apt-get install -y sudo curl

RUN apt-get update && apt-get install -y python

RUN curl -O https://downloads.globus.org/globus-connect-personal/linux/stable/globusconnectpersonal-2-latest.tgz
RUN tar xzvf globusconnectpersonal-2-latest.tgz
RUN rm globusconnectpersonal-2-latest.tgz
RUN echo export PATH=$(ls -d /globusconnectpersonal*):'$PATH' > /etc/profile.d/go.sh

