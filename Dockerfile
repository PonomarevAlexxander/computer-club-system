FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod ./

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /computer-club-system ./cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /computer-club-system /computer-club-system

ENTRYPOINT ["/computer-club-system"]
