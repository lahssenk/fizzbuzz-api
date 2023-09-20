FROM golang:1.21 as build

WORKDIR /go/src/app
COPY . .

# RUN go mod download
RUN go vet ./...
RUN go test -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/...

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/app /
EXPOSE 8080
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV ADMIN_PORT=8081
ENTRYPOINT ["/app"]
