FROM golang:1.19 AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/ux-server cmd/main.go

FROM scratch

WORKDIR /

COPY --from=build /server/bin .
COPY --from=build /server/migration ./migration

EXPOSE 8080
ENTRYPOINT ["./ux-server"]
