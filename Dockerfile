
FROM golang:1.19.3-alpine AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go

FROM gcr.io/distroless/base-debian10

COPY --from=build /src/main /usr/bin/babel
ENTRYPOINT ["/usr/bin/babel"]