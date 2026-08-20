FROM golang:1.27.0-alpine3.24 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

# Загружаю все необходимые зависимости:
RUN go mod download

COPY ./ ./ 

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-server ./main.go

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /app/auth-server .  
EXPOSE 50051 

CMD [ "./auth-server" ]
