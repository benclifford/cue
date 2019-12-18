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

# get rid of ModemManager on host which opens the serial port
# and so confuses things

# run docker like this: cue -D "--device /dev/ttyACM0:/dev/ttyACM0" tinyfpga-bx 
