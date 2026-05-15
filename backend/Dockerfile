FROM golang:1.26.2-alpine3.23
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ft_transcendence_backend
EXPOSE 7772
CMD ["./ft_transcendence_backend"]
