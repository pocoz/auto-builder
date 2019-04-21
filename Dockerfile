# build binary
FROM golang:1.12-alpine3.9 AS build
RUN apk add git
ENV GO111MODULE=on
WORKDIR /go/mod/github.com/pocoz/auto-builder
COPY . /go/mod/github.com/pocoz/auto-builder
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/auto-builder github.com/pocoz/auto-builder/cmd/auto-builderd

# copy to alpine image
FROM alpine:3.9 AS prod
WORKDIR /app
COPY --from=build /out/auto-builder /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app/auto-builder"]
