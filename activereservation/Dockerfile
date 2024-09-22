# Use an official Golang runtime as a parent image
FROM golang:1.20

RUN mkdir /shared
COPY ./shared /shared


# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY ./activereservation .

# Installs Go dependencies
RUN go mod tidy && go mod download

RUN go build

# Tells Docker which network port your container listens on
EXPOSE 8072

# Define the command to run your application
CMD ["go","run","main.go"]