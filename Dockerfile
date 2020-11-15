FROM golang:alpine
# Create app directory and copy everything over there
RUN apk add --update-cache git

RUN mkdir /app
ADD . /app

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]