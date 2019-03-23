FROM golang:1.12.1

WORKDIR /go/src/nozobot
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV CARUDO $CARUDO
ENV WAIFU $WAIFU

CMD ["nozobot"]
