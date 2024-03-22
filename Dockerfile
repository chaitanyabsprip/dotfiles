# syntax=docker/dockerfile:1

FROM ubuntu:latest
WORKDIR /root/dotfiles
RUN yes | unminimize && \
    apt update && \
    apt-get -y --no-install-recommends upgrade && \
    apt-get install -y --no-install-recommends \
    build-essential ca-certificates vim curl wget \
    git zsh tmux make ssh fd-find bat man-db
ENV PATH=.:$PATH:/root/.local/bin
ENV TERM=xterm-256color
COPY . /root/dotfiles
RUN ./dotme serve && make all
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /var/log/dmesg.* && \
    cat /dev/null > /var/log/dmesg
