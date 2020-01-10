# tinyfpga bx
# based on instructions at
# https://tinyfpga.com/bx/guide.html

FROM debian:bullseye

RUN apt-get update && apt-get install -y sudo

# had some problems with python3 but I think it might have
# been modem manager running on host
# RUN apt-get update && apt-get install -y python3 python3-pip

RUN apt-get update && apt-get install -y python python-pip

# i also installed the following as my own user later on,
# as part of fiddling, but I think it might not be necessary
# and I was experiencing a different problem
RUN pip install apio==0.4.0b5 tinyprog
RUN apio install system scons icestorm iverilog
RUN apio drivers --serial-enable

RUN apt-get update && apt-get install -y yosys arachne-pnr


# this bit installs the riscv5 version of gcc etc
# following the instructions at https://github.com/cliffordwolf/picorv32/blob/master/README.md#building-a-pure-rv32i-toolchain

RUN apt-get update && apt-get install -y autoconf automake autotools-dev curl libmpc-dev libmpfr-dev libgmp-dev gawk build-essential bison flex texinfo gperf libtool patchutils bc zlib1g-dev git libexpat1-dev
RUN mkdir /opt/riscv32i
RUN mkdir /build
WORKDIR /build

# do this all in one go because it uses 9gb in /build that doesn't need to be encoded into an image
RUN git clone https://github.com/riscv/riscv-gnu-toolchain riscv-gnu-toolchain-rv32i && cd /build/riscv-gnu-toolchain-rv32i && git checkout 411d134 && git submodule update --init --recursive && mkdir build && cd /build/riscv-gnu-toolchain-rv32i/build && ../configure --with-arch=rv32i --prefix=/opt/riscv32i && make && cd / && rm -rf /build

# this puts the executables in /opt/riscv32i/bin/ which is not on the PATH

# get rid of ModemManager on host which opens the serial port
# and so confuses things

# run docker like this: cue -D "--device /dev/ttyACM0:/dev/ttyACM0" tinyfpga-bx 
