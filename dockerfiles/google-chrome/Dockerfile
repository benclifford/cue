# google chrome, from a google repo

FROM debian:stretch

RUN apt-get update && apt-get install -y sudo curl
RUN apt-get update && apt-get install -y wget gnupg2

# key installing instructions from google:

RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add -

RUN curl -O https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb

# dpkg --install is expected to return a failure, and apt-get -f
# should not return a failure. so ; to sequence the operations.
# this probably means that some errors (eg a corrupt .deb) will not
# cause a build failure here - as this line does not distinguish between
# dpkg expected dependency failure and some other failure.

RUN dpkg --install google-chrome-stable_current_amd64.deb ; apt-get -y -f install

