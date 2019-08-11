FROM debian:buster

RUN apt-get update && apt-get install -y sudo curl

RUN apt-get update && apt-get install -y cargo

# needed for thrussh-libsodium crate?
RUN apt-get update && apt-get install -y libsodium-dev libsodium23
RUN apt-get update && apt-get install -y pkg-config libssl-dev

# needed for nettle-sys crate
RUN apt-get update && apt-get install -y nettle-dev libclang-dev clang

RUN cargo install --root /pijul pijul

RUN echo 'export PATH=/pijul/bin:$PATH' > /etc/profile.d/pijul.sh

