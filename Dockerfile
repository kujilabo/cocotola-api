FROM golang:1.16-alpine as builder

RUN apk add --no-cache build-base

WORKDIR /go/src/app
ADD . .
# COPY main.go .

RUN go build -o cocotola

# Application image.
FROM alpine:latest

RUN apk --no-cache add tzdata

COPY --from=builder /go/src/app/cocotola .
COPY --from=builder /go/src/app/configs ./configs
COPY --from=builder /go/src/app/sqls ./sqls

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

EXPOSE 8080

CMD ["./cocotola"]
