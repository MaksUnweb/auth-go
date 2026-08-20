FROM golang:1.27.0-alpine3.24 AS builder

WORKDIR /app 

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o admin-web ./main.go

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /app/admin-web .  
EXPOSE 8080

CMD [ "./admin-web" ]
