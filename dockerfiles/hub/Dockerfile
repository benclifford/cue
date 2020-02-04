# github commandline tool

FROM debian:stretch

RUN apt-get update && apt-get install -y sudo curl
RUN apt-get update && apt-get install -y git
RUN curl -L -O https://github.com/github/hub/releases/download/v2.8.4/hub-linux-amd64-2.8.4.tgz
RUN tar xf hub-linux-amd64-2.8.4.tgz
RUN echo 'export PATH=/hub-linux-amd64-2.8.4/bin:$PATH' > /etc/profile.d/hub.sh
ADD cue.shell /
RUN chmod a+x /cue.shell
