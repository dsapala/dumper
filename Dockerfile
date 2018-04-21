FROM golang:1.10-alpine
WORKDIR /go/src/github.com/dsapala/dumper/
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dumper .

FROM alpine:latest
VOLUME "/requests"
WORKDIR /
COPY --from=0 /go/src/github.com/dsapala/dumper/dumper .
ENTRYPOINT ["./dumper"]
CMD ["-addr", "0.0.0.0"]

