FROM golang:latest

RUN apt update -y && apt install zsh tmux -y
WORKDIR /root/dot

COPY . /root/dot/

RUN go mod tidy

RUN go install ./cmd/dot

RUN dot alacritty setup; \
        dot bat setup; \
        dot bin setup; \
        dot gh setup; \
        dot git setup; \
        dot kitty setup; \
        dot tmx setup;
