FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV Symbol=MSFT \
    NDAYS=4

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .
COPY ./config/config.json /dist/config/config.json

# Build the application
RUN go build stocks.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/stocks .

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["/dist/stocks"]