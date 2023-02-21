FROM golang:1.20

RUN go install -v golang.org/x/tools/gopls@latest

WORKDIR /code

COPY ["go.mod", "go.sum", "./"]

RUN go mod download