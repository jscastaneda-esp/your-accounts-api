FROM golang:1.20

RUN go install -v golang.org/x/tools/gopls@latest
RUN apt update -y && apt install nano
RUN git config --global core.editor "nano"

WORKDIR /code

COPY ["go.mod", "go.sum", "./"]

RUN go mod download