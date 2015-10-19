FROM golang:1.4.3-wheezy
MAINTAINER Qlean

# Install system dependencies
RUN apt-get update -qq && \
    apt-get install -qq -y pkg-config build-essential

RUN mkdir -p /app
WORKDIR /app
COPY . /app/
ENV GOPATH /go/
ENV _ENV production
RUN go get -d -v
RUN go build -o go-sms
EXPOSE 8080

CMD ["./go-sms"]