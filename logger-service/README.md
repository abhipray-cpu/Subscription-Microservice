# Logger Service

This service write logs to a persistence storage in this case MonggoDB

## Table of Contents

- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)

## Project Structure

.

- ├── logger-service.dockerfile
- ├── api
    - └── main.go
    - └── handler.go
    - └── router.go
- ├── data
  - └── models.go

## Prerequisites

- go should be installed in your system

## Installation

- go mod tidy: Run this command to install al the required dependencies

## Usage

- This service is used to write logs from all the services to mongoDB
