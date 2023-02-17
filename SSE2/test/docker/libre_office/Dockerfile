FROM debian:buster AS libre_office
# Debian 10 O.S., Linux Kernel 4.19 (LTS).
#TODO: Update to 'debian:bullseye' when it is released.
# Debian 11 O.S., Linux Kernel 5.10 (LTS).

# Show version of important things.
RUN echo "$PATH"
RUN uname --all

# Update the packages and C.A. certificates.
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get upgrade && \
  apt-get install -y ca-certificates
RUN apt-get install -y nano mc htop iotop

# Install LibreOffice.
RUN apt-get install -y libreoffice

# Where are we ?
RUN echo $PWD
RUN ls -l
