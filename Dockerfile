# setup the base image
FROM golang 

# setup the directory for the app
WORKDIR /usr/src/app

# cpy all files to the container
# COPY ..
# ADD 

# install system wise and app-specific dependencies
# RUN commands

# expose our port number
EXPOSE 5000

# cmds once build start running
# CMD []
