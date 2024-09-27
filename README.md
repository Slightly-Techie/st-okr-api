# OKR Project

Welcome to the OKR (Objectives and Key Results) project! We're excited to have you here.

## What are OKRs?

OKRs are a goal-setting framework used to define objectives and the key results that indicate success. They help teams align their efforts and focus on achieving specific, measurable outcomes.

## API Setup Guide
Follow the steps below to set up the API locally:

### Prerequisites
- Ensure you have the following installed on your machine:

- Docker
- Go
- Git

### Getting Started
#### Clone the repository

Use `git clone` to clone the repository and navigate to the project directory.

#### Copy environment variables

Create a .env file by copying .env.public and rename it to .env.
Provide your API keys and other required environment variables in the .env file.

#### Spin up the database

Use Docker to start the database container by running `docker-compose up`.

#### Install dependencies

Run `go mod tidy` to install the necessary dependencies.

#### Start the development server

Use `go run ./cmd/main.go` to start the API server.
The server should now be running on http://localhost:{PORT}.

#### Testing the API

You can test the API by making a `GET` request to:
```
http://localhost:{PORT}
```
- Replace {PORT} with the port number specified in your .env file.

#### Common Issues
- Ensure Docker is running before starting the database.
- Check that the environment variables in .env are correctly set.
