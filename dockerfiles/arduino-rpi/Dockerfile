FROM resin/rpi-raspbian:stretch

RUN apt-get update && apt-get install -y sudo curl

WORKDIR /

RUN curl -O https://downloads.arduino.cc/arduino-1.8.5-linuxarm.tar.xz

# to uncompress arduino
RUN apt-get update && apt-get install -y xz-utils

# for arduino to run
RUN apt-get update && apt-get install -y libx11-6 libxext6 libxrender1 libxtst6 libgtk2.0-0

# to unzip neopixel library
RUN apt-get update && apt-get install -y unzip

RUN tar xvf arduino-1.8.5-linuxarm.tar.xz



RUN mkdir /lib-thirdparty


# adafruit neopixel library
WORKDIR /
RUN curl -L -O https://github.com/adafruit/Adafruit_NeoPixel/archive/master.zip
WORKDIR /lib-thirdparty
RUN unzip /master.zip
RUN rm /master.zip


# adafruit boards
WORKDIR /
RUN curl -L -O https://github.com/adafruit/Adafruit_Arduino_Boards/archive/master.zip
RUN mkdir /hardware-thirdparty
RUN mkdir /hardware-thirdparty/adafruit
RUN mkdir /hardware-thirdparty/adafruit/avr
WORKDIR /hardware-thirdparty/adafruit/avr
RUN unzip /master.zip
RUN rm /master.zip



RUN echo 'export PATH=$PATH:/arduino-1.8.5' > /etc/profile.d/arduino.sh
RUN echo 'export _JAVA_AWT_WM_NONREPARENTING=1' >> /etc/profile.d/arduino.sh

# now build stuff like this:
# arduino-builder -hardware /arduino-1.8.5/hardware -tools /arduino-1.8.5/tools -tools /arduino-1.8.5/hardware/tools/ -fqbn adafruit:avr:trinket3 -libraries /lib-thirdparty/ -tools /arduino-1.8.5/tools-builder/ -hardware /hardware-thirdparty/ ./goggles1.ino

# to help look at what is plugged into the usb ports:
RUN apt-get update && apt-get install -y usbutils


# this is for the gertboard using raspberry pi GPIO so don't do
# it in the desktop version of this Dockerfile
RUN curl -O http://project-downloads.drogon.net/gertboard/avrdude_5.10-4_armhf.deb
# dependencies for avrdude .deb
RUN apt-get install -y libusb-0.1-4 libreadline6
RUN dpkg --install avrdude_5.10-4_armhf.deb

RUN mkdir -p /hardware-thirdparty/gertboard/avr/hardware-thirdparty/gertboard/avr

# these are based on a fiddle with the stuff from drogon.net
# and pulling some stuff from the adafruit trinket board files.
ADD gertboard/boards.txt /hardware-thirdparty/gertboard/avr/
ADD gertboard/platform.txt /hardware-thirdparty/gertboard/avr/
ADD gertboard/programmers.txt /hardware-thirdparty/gertboard/avr/

