FROM gocv:latest

USER root

RUN mkdir -p /imgs && chmod 777 /imgs

RUN mkdir -p /imgs/results && chmod 777 /imgs/results

ENV GOPATH /go

WORKDIR /app

COPY . /app

RUN go build -o bin/main .

CMD ["bin/main"]

