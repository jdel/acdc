# ACDC: A Continuous Docker Compose

## Goal

ACDC aims at creating a way to use docker-compose as a lightweight continuous deployment endpoint.

## How

ACDC provides a REST API to trigger Docker Compose actions as well as dedicated webhook receivers for, but not limited to, slack and docker hub.

For ease of development and deployment alike, ACDC is written in Go, and is shipped as a single statically linked binary.

## Warning

ACDC has not been develpped with maximum security in mind. In its current state, I do not recomment to run in a production environment.
