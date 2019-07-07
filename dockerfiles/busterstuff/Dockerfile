# debian buster general work environment

FROM debian:buster

RUN apt-get update && apt-get install -y sudo

# misc small tools
RUN apt-get update && apt-get install -y less vim git sqlite3 dnsutils whois alpine curl irssi

# sound and vision
RUN apt-get update && apt-get install -y youtube-dl sox libsox-fmt-all

# radio
RUN apt-get update && apt-get install -y direwolf gqrx-sdr

# pics
RUN apt-get update && apt-get install -y less gthumb

# office
RUN apt-get update && apt-get install -y less libreoffice

# tex
RUN apt-get update && apt-get install -y texlive-latex-base texlive-latex-recommended texlive-latex-extra xpdf


