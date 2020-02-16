FROM golang:1.13-alpine

# Install Git
RUN apk --update add git && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

# Copy app directory
COPY . /srv/app

# Install app dependencies
WORKDIR /srv/app
RUN go get github.com/Masterminds/sprig

# Ports
EXPOSE 80

# Launch
CMD [ "go", "run", "server.go" ]
