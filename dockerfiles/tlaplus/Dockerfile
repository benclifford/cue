FROM debian:bullseye

RUN apt-get update && apt-get install -y sudo
RUN apt-get update && apt-get install -y curl

WORKDIR /

RUN curl -O -L https://github.com/tlaplus/tlaplus/releases/download/v1.7.1/TLAToolbox-1.7.1.deb


# Setting up tla+toolbox (1.7.1~20201231-1918) ...
# /var/lib/dpkg/info/tla+toolbox.postinst: 10: /usr/bin/update-mime-database: not found
# dpkg: error processing package tla+toolbox (--install):
# installed tla+toolbox package post-installation script subprocess returned error exit status 127
# ... so:
RUN apt-get update && apt-get install -y shared-mime-info

# eclipse inside TLA+ can't start, with a error in the log file about not being able to load gtk
# ... so:
RUN apt-get update && apt-get install -y libgtk-3-0

RUN echo export PATH=/opt/TLA+Toolbox/:'$PATH' > /etc/profile.d/tlaplus.sh

RUN dpkg --install TLAToolbox-1.7.1.deb

