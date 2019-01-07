FROM i386/ubuntu:trusty
# this wants to run in i386 mode, not x86, so we need to use an i386
# distro.

RUN apt-get update && apt-get install -y sudo

# for downloading, if this is how you want to download
RUN apt-get update && apt-get install -y curl


# assorted X libraries needed to run. I figured out this list by 
# running repeatedly, fixing library error messages until it worked.
RUN apt-get update && apt-get install -y libxext6 libxext-dev libfreetype6 libxrender1 libfontconfig1


# Inside a container it seems like QT gets upset running X windows without
# being able to do shared memory - this causes it to not try to do so.
# The symptoms for me were a window appearing, but blank.
# This line will set the variable appropriately if running inside `cue`,
# but set it by hand in your shell if you're running some other way.
RUN echo 'export QT_X11_NO_MITSHM=1' > /etc/profile.d/hmrc.sh


# you can install PAYE tools by downloading the zip file, and running
# the only executable inside. It will install itself by default in ~/HMRC
# so make sure you have # that mounted - if you're using `cue` that will
# happen for you but # maybe you've come across this Dockerfile elsehow.

# curl -L -O https://www.gov.uk/government/uploads/uploaded/hmrc/payetools-rti-18.1.18115.285-linux.zip


# inside this container:
# cd ~/HMRC/payetools-rti
# ./rti.linux
