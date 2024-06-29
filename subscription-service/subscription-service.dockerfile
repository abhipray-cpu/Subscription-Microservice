# This Dockerfile is used to build a Docker image for the subscription service microservice.
# It uses the golang:bookworm base image and copies the subscriptionApp binary into the /app directory.
# The CMD instruction specifies the command to run when the container starts.

FROM golang:bookworm

RUN mkdir /app

COPY subscriptionApp /app

CMD [ "/app/subscriptionApp"]