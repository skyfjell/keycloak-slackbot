################# Building Binary #################

FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/bot
RUN apk --no-cache add ca-certificates

################# Production Binary #################
FROM scratch

COPY --from=builder /go/bin/bot /go/bin/bot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 5000

ENTRYPOINT [ "/go/bin/bot" ]