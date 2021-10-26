# syntax=docker/dockerfile:1
# mostly via https://docs.docker.com/language/golang/build-images/

# TODO: use a multi-stage build instead of including the full golang image

FROM golang:1.17-alpine

WORKDIR /app

COPY . ./

RUN go build -o /bin/app

EXPOSE 8080

CMD ["/bin/app"]