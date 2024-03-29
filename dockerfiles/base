# syntax=docker/dockerfile:1

FROM ubuntu:latest
WORKDIR /root/dotfiles
ENV TERM=xterm-256color
ENV TZ=Asia/Kolkata
RUN yes | unminimize && \
    apt update && \
    ln -fs /usr/share/zoneinfo/Asia/Kolkata /etc/localtime && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends sudo tzdata
RUN dpkg-reconfigure -f noninteractive tzdata
RUN apt-get install -y --no-install-recommends build-essential ca-certificates \
    vim curl wget git make ssh man-db locales
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en  
ENV LC_ALL en_US.UTF-8  
ENV PATH=.:$PATH:/root/.local/bin
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /var/log/dmesg.* && \
    cat /dev/null > /var/log/dmesg && \
    apt-get update && \
    apt-get -y --no-install-recommends upgrade
COPY . /root/dotfiles
RUN ./dotme serve && setup all
ENTRYPOINT ["tmux", "new", "-As", "home"]

