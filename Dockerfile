FROM golang:1.24-alpine AS build

WORKDIR /go/src/github.com/dskart/particles
COPY . .

RUN go generate ./...
RUN go build .


FROM golang:1.24-alpine

WORKDIR /usr/bin

COPY --from=build /go/src/github.com/dskart/particles/particles .
RUN ./particles --help > /dev/null

ENV PARTICLES__API__MAX_NUM_SESSIONS="10"

ENTRYPOINT ["/usr/bin/particles"]