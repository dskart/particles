FROM golang:1.24-alpine AS build

WORKDIR /go/src/github.com/dskart/particles
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go generate ./...
RUN CGO_ENABLED=0 go build -o particles .


FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /usr/bin

COPY --from=build /go/src/github.com/dskart/particles/particles .
RUN ./particles --help > /dev/null

ENV PARTICLES__API__MAX_NUM_SESSIONS="10"

# 2222: SSH, 8080: health checks + SSH over websocket (/ssh)
EXPOSE 2222 8080

ENTRYPOINT ["/usr/bin/particles"]
CMD ["serve", "--port", "2222", "--http-port", "8080"]
