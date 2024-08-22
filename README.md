# Go Microservice with MongoDB

This repository contains a Go microservice using the Gin framework and MongoDB. The microservice provides basic CRUD operations for managing users.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or later)
- [MongoDB](https://www.mongodb.com/try/download/community) (installed and running)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [curl](https://curl.se/) or [PowerShell's `Invoke-RestMethod`](https://learn.microsoft.com/en-us/powershell/scripting/overview)

## Getting Started

### Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/anassajja/Training_Session_Microservice.git
cd Training_Session_Microservice

### Setup Environment Variables

Create a .env file in the root directory with the following content:

MONGODB_URI=mongodb://localhost:27017


### Install Dependencies

Run the following command to install the required Go packages:

bash go mod tidy

## Run the Server
Start the server by running:

bash go run cmd/main.go

The server should now be running on http://localhost:8080.


### Testing the API
### Recommanded to Test API endpoints using Postman Tool 

You can test the API using curl or PowerShell’s Invoke-RestMethod

# Test GET Users

Retrieve a list of users:

curl -X GET http://localhost:8080/users

Or in PowerShell:

powershell

Invoke-RestMethod -Method Get -Uri http://localhost:8080/users

# Test POST User

Create a new user:

curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "email": "john.doe@example.com"}'

Or in PowerShell:

powershell

$body = @{
    name  = "John Doe"
    email = "john.doe@example.com"
}
Invoke-RestMethod -Method Post -Uri http://localhost:8080/users -Body ($body | ConvertTo-Json) -ContentType "application/json"

# Running Tests
To run the tests for the project, use the following command:

bash

go test ./...
This command will execute all tests in your project. Make sure to create test files following Go’s testing conventions (*_test.go) to ensure coverage.


### **Summary of Updates**

- **Environment Variables Section**: Added instructions for creating a `.env` file.
- **Test Commands Section**: Added instructions for running tests with `go test ./...`.

By including these updates, your `README.md` will provide clear instructions for setting up, running, and testing your Go microservice project.


