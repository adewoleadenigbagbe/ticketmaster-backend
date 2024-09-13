# Use an official Golang runtime as a parent image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY ./apis .

# Installs Go dependencies
RUN go mod tidy && go mod download

# Tells Docker which network port your container listens on
EXPOSE 8071

# Define the command to run your application
CMD ["go run main.go"]