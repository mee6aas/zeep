FROM golang:1.12.9

WORKDIR /zeep

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install ./cmd/zeep-agent

EXPOSE 5122

ENTRYPOINT ["zeep-agent"]
