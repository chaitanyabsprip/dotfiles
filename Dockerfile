FROM golang:latest

WORKDIR /root/dot
COPY go.mod go.sum /root/dot/
RUN go mod download
RUN apt update -y && apt install zsh tmux -y
COPY . /root/dot/
RUN go install ./cmd/dot
RUN dot alacritty setup; \
        dot bat setup; \
        dot bin setup; \
        dot gh setup; \
        dot git setup; \
        dot kitty setup; \
        dot tmx setup;
