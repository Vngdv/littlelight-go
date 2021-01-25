FROM golang:alpine as build
# Create app directory and copy everything over there
RUN apk add --update-cache git

RUN mkdir /app
ADD . /app

WORKDIR /app

RUN go get github.com/bwmarrin/discordgo
RUN go build -o main .

# Production stage
FROM alpine:latest

COPY --from=build /app/main .

# Add user for more security
RUN addgroup -S littlelight && adduser -s /bin/false -S littlelight -G littlelight
RUN chown -R littlelight:littlelight /main

USER littlelight

ENTRYPOINT ["/main"]
