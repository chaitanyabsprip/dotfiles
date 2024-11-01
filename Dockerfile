FROM golang:latest

RUN apt update -y && apt install zsh tmux -y
WORKDIR /root/dot

COPY . /root/dot/

RUN go mod tidy

RUN go install ./cmd/dot && bash -c 'complete -C dot dot'

RUN dot tmx setup; dot gh setup; dot git setup; dot alacritty setup
