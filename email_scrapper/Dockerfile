# setup the base image
FROM golang:1.22

# setup the directory for the app
# so that all following commands are using the same locations
# it syntax-sugars your mkdir/cd commands
WORKDIR /usr/src/app

# precopy to setup go dependency
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# cpy all files to the container
COPY . .
# ADD 

# install system wise and app-specific dependencies
# RUN commands

# expose our port number
EXPOSE 5000

# cmds once build start running
# CMD []
