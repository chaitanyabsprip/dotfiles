FROM golang:latest

WORKDIR /root/dot
COPY go.mod go.sum /root/dot/
RUN go mod download
RUN apt update -y && apt install zsh -y
COPY . /root/dot/
RUN go install ./cmd/dot
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
