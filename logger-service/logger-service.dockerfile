# This Dockerfile is used to build a Docker image for the logger service.
# It starts with the base image "golang:bookworm" and copies the loggerApp binary into the container.
# The CMD instruction specifies the command to run when the container starts.

FROM golang:bookworm

RUN mkdir /app

COPY loggerApp /app

CMD [ "/app/loggerApp"]
