# build go code image
# usr golang version 1.19 image
FROM golang:1.19-alpine as builder
# build namespace
WORKDIR /build

COPY go.mod .
# download go package
RUN go mod download

COPY . .
# build go
# RUN go build -o /main main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

# use alpine:3 image run image
FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]