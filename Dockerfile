FROM alpine:latest

RUN apk --no-cache add tzdata

COPY --from=builder /go/src/app/cocotola .
COPY --from=builder /go/src/app/configs ./configs

CMD ["./cocotola", "-env", "production"]
