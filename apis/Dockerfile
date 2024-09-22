# Use an official Golang runtime as a parent image
FROM golang:1.20

RUN apt-get update \
    && apt-get install -y \
        curl \
        libxrender1 \
        libfontconfig \
        libxtst6 \
        xz-utils

RUN curl "https://github.com/wkhtmltopdf/packaging/archive/refs/tags/0.12.6.1-3.tar.gz" -L -o "wkhtmltopdf.tar.gz"
RUN tar xzf wkhtmltopdf.tar.gz
RUN mv wkhtmltopdf.tar.gz /usr/local/bin/wkhtmltopdf
# RUN chmod +x /usr/local/bin/wkhtmltopdf

#Copy Shared folder
RUN mkdir /shared
COPY ./shared /shared

#Copy Infastructure folder
RUN mkdir /infastructure
COPY ./infastructure /infastructure


# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY ./apis .

# Installs Go dependencies
RUN go mod tidy && go mod download

RUN go build

# Tells Docker which network port your container listens on
EXPOSE 8071

# Define the command to run your application
CMD ["go","run","main.go"]