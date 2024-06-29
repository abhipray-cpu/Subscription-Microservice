# Go Microservice Template

This project is a template for building microservices in Go. It provides a basic setup for developing, testing, and deploying Go-based microservices with an emphasis on clean architecture and best practices.

## Features

- **Clean Architecture**: Organized according to the Clean Architecture principles for easy maintenance and scalability.
- **RESTful API**: A simple RESTful API setup using Gorilla Mux for routing.
- **Docker Support**: Includes Dockerfile for building and running the microservice in a Docker container.
- **Logging**: Integrated logging for easy debugging and monitoring.
- **Configuration Management**: Utilizes Viper for managing configurations from files and environment variables.
- **Health Check Endpoint**: A health check endpoint for checking the service status.
- **Unit and Integration Tests**: Basic setup for writing unit and integration tests using the testing package.

## Getting Started

### Prerequisites

- Go 1.15 or higher
- Docker (optional for containerization)

### Running the Service

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/yourprojectname.git
   ```

### Managing Docker Containers with Makefile Commands

The project includes a Makefile for easy Docker container management. Here are the commands you can use:

- **Start Containers**: To start all containers in the background without forcing a build, run:
  ```bash
  make up
  ```
  This command starts the Docker images.

- **Build and Start Containers**: To stop any running containers, build all projects (when required), and start the containers, use:
  ```bash
  make up_build
  ```
  This command stops any running Docker images, builds (when necessary), and starts the Docker images.

- **Stop Containers**: To stop all running containers, execute:
  ```bash
  make down
  ```
  This command stops the Docker compose and cleans up.

These commands provide a convenient way to manage your Docker containers directly from the command line.