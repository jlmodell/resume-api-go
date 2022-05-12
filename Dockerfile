
# Building the binary of the App
FROM golang:1.18 AS build

WORKDIR /go/src/resume-api-go

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

# moving the binary to the root of the image
FROM alpine:latest

WORKDIR /app

COPY --from=build /go/src/resume-api-go/app .

EXPOSE 8001

CMD ["./app"]