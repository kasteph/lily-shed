FROM golang:1.18.3-alpine AS build

WORKDIR /go/src/github.com/kasteph/lily-shed

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o .

FROM gcr.io/distroless/base-debian11

COPY --from=build /go/src/github.com/kasteph/lily-shed/lily-shed /usr/bin/

ENTRYPOINT ["/usr/bin/lily-shed"]
CMD ["--help"]