FROM golang:1.20
WORKDIR /app
COPY . /app
RUN GIN_MODE=release CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o main .
EXPOSE 6000
CMD ["./main"]