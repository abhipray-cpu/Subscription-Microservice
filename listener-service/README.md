# Listener Service

This service subscribes to a kafka queue and listens for all the new events and then post the new logs to logger service

## Table of Contents

- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)

## Project Structure

.

- ├── listener-service.dockerfile
- ├── api
    - └── main.go
    - └── listener.go
    - └── handler.go
  - └── router.go


## Prerequisites

- go should be installed in your system

## Installation

- go mod tidy: Run this command to install al the required dependencies

## Usage

- This service listens to a kafka queue to listen for all the new logs generated by different services and then forward these logs to logger service over http