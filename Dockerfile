FROM golang:1.23-alpine

ENV GOTOOLCHAIN=auto

WORKDIR /app
COPY . .

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "main.go"]
