FROM golang:1.23.3-alpine

WORKDIR /app

RUN apk add --no-cache curl

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY .air.toml ./

COPY . .

RUN curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 && \
    chmod +x /usr/local/bin/dbmate

EXPOSE 8080

CMD ["sh", "-c", "dbmate up && go run main.go -seed && air -c .air.toml"]
