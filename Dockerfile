FROM golang:alpine

# Set necessary environmet variables needed for our image for testing removed for prod build
# ENV Symbol=MSFT \
#     NDAYS=4

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .

# Create config for testing, removed for prod
# COPY ./config/config.txt /dist/config/config.txt

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