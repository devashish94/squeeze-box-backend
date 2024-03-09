FROM golang:1.22.0

# Install NGINX
RUN apt-get update && \
    apt-get install -y nginx zip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create NGINX configuration
COPY ./config/nginx.conf /etc/nginx/nginx.conf

# Set up a simple Go server
WORKDIR /app
COPY . .

# Expose ports
EXPOSE 80 4000

RUN go mod tidy

# Start NGINX and the Go server
# CMD ["service", "nginx", "start"] && ["go run ."]
# CMD ["service", "nginx", "start"] && ["echo Hello World"]
# CMD echo Hello World && tail -f /dev/null
CMD service nginx start && go build -o server.out . && ./server.out

