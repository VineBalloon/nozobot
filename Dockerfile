FROM golang:1.12.1

WORKDIR /go/src/nozobot
COPY . .

# Install dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# Get API keys
ENV CARUDO=""
ENV WAIFU=""

# Expose websocket
EXPOSE 443

CMD ["nozobot"]
