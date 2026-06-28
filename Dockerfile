FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /kuberpc ./cmd

EXPOSE 8080
EXPOSE 9443

CMD ["/kuberpc"]
