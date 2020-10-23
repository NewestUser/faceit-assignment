FROM golang:1.14

MAINTAINER Marian Zlatev

# create a build dir
RUN mkdir /build
WORKDIR /build

# it seems gcc is not available on plain alpine
ENV CGO_ENABLED=0

# copy module dependencies
COPY go.mod .
COPY go.sum .

# copy the project itself
COPY . .

RUN go mod download && cd cmd && go build -o faceit-app

CMD ["./cmd/faceit-app"]