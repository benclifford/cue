FROM ubuntu:16.04

RUN apt-get update && apt-get install -y sudo
RUN apt-get update && apt-get install -y curl

RUN curl -L -O https://github.com/ethereum/mist/releases/download/v0.9.0/Ethereum-Wallet-linux64-0-9-0.deb

# skip exit code because we expect it to fail (due to missing
# dependencies which we will fix up later)
RUN dpkg --install Ethereum-Wallet-linux64-0-9-0.deb || true

# fix up the dpkg specified dependencies for what we just installed
RUN apt-get -y -f install

# fix up the dependencies apparently needed by hand
RUN apt-get -y install libxss1 libasound2
