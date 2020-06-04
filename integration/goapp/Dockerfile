ARG GOVERSION=1.14
FROM golang:${GOVERSION} as builder
COPY main.go .
RUN go build -gcflags="all=-N -l" -o /app main.go

FROM gcr.io/distroless/base
CMD ["./app"]
ENV GOTRACEBACK=single
COPY --from=builder /app .
