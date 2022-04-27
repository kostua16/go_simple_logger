FROM golang:1.18

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./...

ENV EUREKA_URL=http://eureka:8761
ENV PORT=8080
EXPOSE ${PORT}
CMD ["simple_logger", "-eureka", "${EUREKA_URL}", "-port", "${PORT}", "-verbose", "INFO"]