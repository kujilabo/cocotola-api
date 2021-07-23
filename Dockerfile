FROM golang:1.16-alpine as builder

WORKDIR /go/src/app
ADD . .
# COPY main.go .

RUN go build -o cocotola

# Application image.
FROM alpine:latest

RUN apk --no-cache add tzdata

COPY --from=builder /go/src/app/cocotola .
COPY --from=builder /go/src/app/configs ./configs

CMD ["./cocotola", "-env", "production"]