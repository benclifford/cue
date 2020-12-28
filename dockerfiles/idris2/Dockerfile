FROM debian:buster

RUN apt-get update && apt-get install -y sudo

RUN apt-get update && apt-get install -y less vim git

# implicit dependencies on top of base buster
RUN apt-get update && apt-get install -y less build-essential

# idris instructions are to install chezscheme
RUN apt-get update && apt-get install -y chezscheme

WORKDIR /
RUN git clone https://github.com/idris-lang/Idris2
WORKDIR /Idris2
RUN make bootstrap SCHEME=chezscheme

# can't install with a different prefix
# (eg PREFIX=/foo) because the installer fails,
# expecting stuff to have been installed in
# ~/.idris2 (~ is ~root, in this case)
RUN make install

RUN echo 'export PATH=$PATH:/root/.idris2/bin' > /etc/profile.d/idris2.sh

# make stuff accessible to all users; otherwise cannot access
# eg the prelude when running as a normal user
RUN chmod a+rx /root
RUN find /root/.idris2/ -type d -exec chmod -v a+rx {} \;
RUN find /root/.idris2/ -type f -exec chmod -v a+r {} \;


