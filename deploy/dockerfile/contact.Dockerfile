FROM golang:1.18.4-alpine3.15 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --update curl

WORKDIR /app

COPY . .

RUN go mod vendor
RUN go build ./services/contact/cmd/app

FROM scratch
WORKDIR /
COPY --from=builder /app/app .
COPY --from=builder /app/services/contact/internal/repository/storage/postgres/migrations /services/contact/internal/repository/storage/postgres/migrations
COPY --from=builder /app/services/contact/internal/delivery/http/swagger/docs /services/contact/internal/delivery/http/swagger/docs



CMD [ "/app" ]
