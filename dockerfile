# building

FROM golang:1.25-alpine3.22 AS build
WORKDIR /project
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p bin && CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# running

FROM gcr.io/distroless/base-debian12:latest AS runner
WORKDIR /
COPY --from=build /project/main /main
COPY --from=build /project/migrations /migrations
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT [ "./main" ]