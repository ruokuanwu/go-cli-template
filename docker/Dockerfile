FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/go-cli-template ./cmd/go-cli-template

FROM scratch
COPY --from=builder /out/go-cli-template /go-cli-template
USER 65532:65532
ENTRYPOINT ["/go-cli-template"]
CMD ["--help"]
