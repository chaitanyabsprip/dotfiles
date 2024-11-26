FROM golang:latest

WORKDIR /root/dot
COPY go.mod go.sum /root/dot/
RUN go mod download
RUN apt update -y && apt install -y locales
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Kolkata
ENV LANG=en_US.UTF-8  
ENV LANGUAGE=en_US:en  
ENV LC_ALL=en_US.UTF-8  
RUN dpkg-reconfigure -f noninteractive tzdata && \
        sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
        dpkg-reconfigure --frontend=noninteractive locales && \
        update-locale LANG=en_US.UTF-8
COPY . /root/dot/
RUN go install -v ./cmd/dot
RUN dot alacritty setup; \
        dot bash setup; \
        dot bat setup; \
        dot bin setup; \
        dot brew setup; \
        dot dirs setup; \
        dot fish setup; \
        dot gh setup; \
        dot git setup; \
        dot gitui setup; \
        dot hypr setup; \
        dot kitty setup; \
        dot lsd setup; \
        dot ohmyposh setup; \
        dot shell setup; \
        dot sqlfluff setup; \
        dot starship setup; \
        dot tmx setup; \
        dot vimium setup; \
        dot waybar setup; \
        dot zsh setup;
