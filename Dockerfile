FROM golang:1.21.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /kuberpc ./cmd

EXPOSE 8080

CMD ["/kuberpc"]
