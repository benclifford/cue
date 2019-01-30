# 'rust' build environment

FROM ubuntu:18.10

RUN apt-get update && apt-get install -y sudo curl

# RUN apt-get update && apt-get install -y rustc
# try using rustup instead

RUN curl -o /main.sh https://sh.rustup.rs -sSf 
RUN chmod a+x /main.sh
RUN CARGO_HOME=/cargo_home RUSTUP_HOME=/rustup_home /main.sh -y

# for pijul build
RUN apt-get install -y gcc pkg-config libsodium-dev libssl-dev

# for use, set /cargo_home/bin onto the path but do not set CARGO_HOME directory.
# *do* set RUSTUP_HOME
