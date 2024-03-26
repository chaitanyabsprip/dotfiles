# syntax=docker/dockerfile:1

FROM ubuntu:latest
WORKDIR /root/dotfiles
RUN yes | unminimize && \
    apt update && \
    apt-get -y --no-install-recommends upgrade && \
    apt-get -y --no-install-recommends install sudo && \
    apt-get install -y --no-install-recommends \
    build-essential ca-certificates vim curl wget \
    git make ssh man-db sudo
COPY . /root/dotfiles
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /var/log/dmesg.* && \
    cat /dev/null > /var/log/dmesg
RUN apt-get update
ENV PATH=.:$PATH:/root/.local/bin
ENV TERM=xterm-256color
RUN ./dotme serve
# && setup all
