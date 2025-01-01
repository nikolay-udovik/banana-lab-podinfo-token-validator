FROM golang:1.23.4 as builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY cmd/main.go cmd/main.go
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o app cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/app .
USER 65533:65533
ENTRYPOINT ["/app"]
