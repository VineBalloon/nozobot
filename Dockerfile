FROM golang:1.12.1

# Don't go to go/src to enable go modules
WORKDIR ~/nozobot
COPY . .

# Install dependencies
RUN go get
RUN go install

# Make sure you pass in env vars when running
ENV CARUDO=""
ENV WAIFU=""

# Expose websocket port
EXPOSE 443

# Run nozobot
CMD ["nozobot"]