FROM golang:1.20-alpine

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY src/go.mod src/go.sum ./
#RUN go mod download && go mod verify

RUN mkdir -p /etc/app/migrations
RUN mkdir -p /etc/app/configs
COPY migrations /etc/app/migrations
COPY configs /etc/app/configs

WORKDIR /usr/src/app
COPY src .
RUN mkdir -p /usr/local/bin/
RUN go mod tidy
RUN go build -v -o /usr/local/bin/app

ENV ENVIRONMENT=PROD

CMD ["app"]
