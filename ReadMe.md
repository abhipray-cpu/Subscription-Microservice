# Go Microservice Template

## Overview

This Go Microservice Template is meticulously engineered to streamline the management of subscriptions, serving as the backbone for subscription-based services. It is crafted with the intent to provide a robust, scalable solution that caters to the dynamic needs of subscription management, including user subscriptions, billing cycles, and notification services.

Leveraging the power of Go, this template not only promises high performance and reliability but also emphasizes clean architecture and best practices. It is equipped to handle high volumes of transactions with ease, ensuring data integrity and security are never compromised.



## Features

- **Scalability at Its Core**: Designed from the ground up to scale seamlessly with your needs. Whether you're handling hundreds or millions of users, this microservice adapts to your growth, ensuring that your subscription services remain fast, efficient, and reliable.

- **CockroachDB Integration**: Utilizes CockroachDB, the resilient, distributed SQL database, for managing user data and payment information. This ensures your data is always available, consistent, and secured, even in the face of hardware failures or data center outages.

- **Global Log Collection Service**: Features a sophisticated log collection service that aggregates logs from all related services, including user management and payment processing. This centralized approach simplifies monitoring and debugging, providing clear insights into the health and performance of your subscription services.

- **Lemon Squeezy for Subscription Management**: Employs Lemon Squeezy for an effortless subscription management experience. This integration allows for easy setup of subscription plans, handling of recurring payments, and management of customer accounts, all within a user-friendly interface.

- **Temporal for Workflow Orchestration**: Incorporates Temporal to architect scalable and resilient workflows. Whether it's automating billing cycles, managing user subscriptions, or orchestrating communication between services, Temporal ensures that your workflows are robust, fault-tolerant, and easy to manage.


## Architecture

![System Design](./project/Subscription%20Microservice.png)

## Prerequisites

Before you begin, ensure you have met the following requirements:

- **Go**: The microservice is developed in Go. Ensure you have Go installed on your machine. The minimum version required is Go 1.15. You can download it from [the Go website](https://golang.org/dl/).

- **Docker**: For containerization and running dependencies such as databases or other services in isolated environments. Download Docker from [the Docker website](https://www.docker.com/get-started).

- **CockroachDB**: This microservice uses CockroachDB for managing user data and payments. You can run CockroachDB locally as a standalone binary, in a Docker container, or use CockroachDB Cloud. Instructions can be found on [the CockroachDB documentation page](https://www.cockroachlabs.com/docs/stable/start-a-local-cluster.html).

- **Temporal**: For orchestrating workflows, Temporal must be set up. You can run Temporal locally using Docker. Find the setup instructions on [the Temporal documentation site](https://docs.temporal.io/docs/server/quick-install).

- **Git**: To clone the repository and manage version control. Download Git from [the Git website](https://git-scm.com/downloads).

- **An IDE or Text Editor**: While not strictly necessary, having an Integrated Development Environment (IDE) like Visual Studio Code, GoLand, or a text editor such as Sublime Text or Atom can be very helpful for development.

Ensure all the above prerequisites are installed and properly configured before proceeding with the installation of the microservice.


## Getting Started

- **Clone the directory**: [git clone](https://github.com/abhipray-cpu/Subscription-Microservice.git) 

- **Install the dependencies**: go mod tidy

- **Set the env variables for both project and other services**

- **Run ngrok for webhook testing**

### Installation

### Installation

Follow this step-by-step guide to install all necessary dependencies for the Subscription Microservice. This guide assumes you have already cloned the repository to your local machine.

1. **Open a Terminal**: Navigate to the cloned repository's directory on your machine.

2. **Install Go Dependencies**: Run the following command to ensure all Go dependencies are correctly installed:

    ```bash
    go mod tidy
    ```

    This command will download and install the necessary Go modules and dependencies defined in `go.mod` and `go.sum` files.

3. **Install Ngrok**: To test webhooks locally, you will need Ngrok. If you haven't installed Ngrok yet, download it from [Ngrok's official website](https://ngrok.com/download) and follow the installation instructions.

4. **Database Setup (Optional)**: If your microservice requires a database, ensure you have the database running either locally or in a container. Follow the database's official guide to install and set it up as needed.

5. **Environment Variables**: Before running the service, you will need to set up the required environment variables. Create a `.env` file in the root directory of your project and populate it with the necessary variables. Refer to the **Configuration** section for details on the required environment variables.

6. **Additional Tools**: Depending on your microservice's needs, you might need to install additional tools or services. Ensure these are installed and properly configured before proceeding.

By following these steps, you will have installed all the necessary dependencies and are ready to configure and run the Subscription Microservice.


### Running the Service

Instructions on how to start the service, including any build steps or commands to run.

### Running the Service

To start the Subscription Microservice, follow these steps:

1. **Navigate to the Project Directory**: Open a terminal and change to the project directory:

    ```bash
    cd ./project
    ```

2. **Build and Start the Service**: Use the provided Makefile to build and start the service:

    ```bash
    make build_up
    ```

    This command compiles the service and starts it, along with any necessary dependencies defined in the Makefile.

3. **Expose the Service with Ngrok**: To make your local service accessible externally (useful for testing webhooks or sharing with others), use Ngrok to expose it:

    ```bash
    ngrok http localhost:8085
    ```

    Replace `8085` with the actual port your service is running on if different. Ngrok will provide you with a public URL that forwards to your local service.

s
