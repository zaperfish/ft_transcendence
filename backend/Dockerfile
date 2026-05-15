FROM golang:1.26.2-alpine3.23
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o ft_transcendence_backend
EXPOSE 4000
CMD ["./ft_transcendence_backend"]
