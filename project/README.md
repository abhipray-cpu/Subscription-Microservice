# Project

This project is a full-stack, server-rendered application using Django, featuring a sleek user interface, a powerful backend, and efficient database handling. The project is containerized using Docker and includes Docker Compose and a Makefile for automation.

## Table of Contents

- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Usage](#usage)
- [Docker Compose](#docker-compose)
  - [Starting the Services](#starting-the-services)
  - [Stopping the Services](#stopping-the-services)
- [Makefile](#makefile)
  - [Available Commands](#available-commands)

## Project Structure

.
- └──  docker-compose.yml
- └──  Makefile
- └── README.md

## Prerequisites

- Docker and makefile extensions should be installed in you system and run the commands with root permissions 

## Usage

- This dir is focused solely on building the application bianries and deploy docker containers to run the services

## Docker 

- docker-compose up: this command spins up all the containers

- docker-compose down: this command stops all the running containers

## Makefile

- make up: this command spins up all the containers

- make down: this command stops all the running containers

- make up_build: this command stops all the running containers, builds the required app binaries, and then spins up all the containers