# Subscription Service

This is the entrypoint of the entire system this listens to incoming request and communicate with other service to provide required functionalities

## Table of Contents

- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)

## Project Structure

.

- ├── subscription-service.dockerfile
- ├── cmd/api
  - └── main.go
  - └── producer.go
  - └── handler.go
  - └── router.go
  - └── middleware.go
- ├── auth
  - └── authenticator.go
  - └── github_authenticator.go
- ├── clients
  - └── sns_client.go
  - └── twilio_client.go
- ├── data
  - └── models.go
  - └── reddis_store.go
  - └── reddis_client.go
- ├── util
  - └── util.go
- ├── api
  - └── worker.go
  - └── workflow
    - └── otp_workflow.go
    - └── welcome_workflow.go
  - └── activities
    - └── activity.go
    - └── mail_activity.go
    - └── otp_activity.go
    - └── sns_activity.go


## Prerequisites

- go should be installed in your system

## Installation

- go mod tidy: Run this command to install al the required dependencies

## Usage

- This is the entrypoint of the entire system this listens to incoming request and communicate with other service to provide required functionalities

- auth: this package is responsible for providing github authentication
- clients: this package provides and initializes all the clients like ses and twilio
- cmd/api: this is the main application that intilizes the main fiel and the application configuration
- data: this package initializes all the storage interfaces
- util: this provides all the utilities functionalities
- worker: this package is for handling temporal workflows and activities
