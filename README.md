# Go Microservice with MongoDB

This repository contains a Go microservice using the Gin framework and MongoDB. The microservice provides basic CRUD operations for managing users and includes several controllers for various functionalities.

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
```

### Setup Environment Variables

Create a `.env` file in the root directory with the following content:

```
MONGODB_URI=mongodb://localhost:27017
```

### Install Dependencies

Run the following command to install the required Go packages:

```bash
go mod tidy
```

## Run the Server

Start the server by running:

```bash
go run cmd/main.go
```

The server should now be running on [http://localhost:8080](http://localhost:8080).

### Testing the API

Recommended to test API endpoints using Postman tool.

#### Test GET Users

Retrieve a list of users:

```bash
curl -X GET http://localhost:8080/users
```

Or in PowerShell:

```powershell
Invoke-RestMethod -Method Get -Uri http://localhost:8080/users
```

#### Test POST User

Create a new user:

```bash
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "email": "john.doe@example.com"}'
```

Or in PowerShell:

```powershell
$body = @{
    name  = "John Doe"
    email = "john.doe@example.com"
}
Invoke-RestMethod -Method Post -Uri http://localhost:8080/users -Body ($body | ConvertTo-Json) -ContentType "application/json"
```

## Controllers Overview

### `SessionController`

- **Scheduling Sessions**: 
  - **Create Session**: Add new training sessions with details such as time, location, and capacity.
  - **Update Session**: Modify existing sessions, including changes to time, location, or other details.
  - **Delete Session**: Remove sessions from the schedule.

- **Getting Sessions**: 
  - **Retrieve All Sessions**: List all scheduled sessions.
  - **Retrieve Session by ID**: Get details of a specific session using its unique identifier.

- **Managing Sessions**: 
  - **Update Session**: Apply updates to sessions as required.
  - **Cancel Session**: Mark sessions as canceled and handle related notifications.
  - **Modify Session**: Make adjustments to session details as needed.

### `UserController`

- **User Enrollment**: 
  - **Register User**: Add users to the system and manage their registration.
  - **Participate in Sessions**: Handle user participation in sessions, including enrollments and withdrawals.

- **User Invitations**: 
  - **Send Invitations**: Manage invitations for private or special sessions.
  - **Manage Invitations**: Track and manage the status of invitations sent to users.

- **User Cancellation**: 
  - **Cancel Enrollment**: Process user requests to cancel their session participation.
  - **Process Refunds**: Handle refunds for canceled sessions as required.

- **Feedback**: 
  - **Submit Feedback**: Collect and process user feedback for sessions and coaches.
  - **Review Management**: Handle reviews submitted by users.

### `CoachController`

- **Coach Management**: 
  - **Manage Coaches**: Handle actions related to coach profiles and their sessions.

- **Coach Role Assignment**: 
  - **Assign Roles**: Allocate roles and responsibilities to coaches and assistants.

- **Coach Session Creation**: 
  - **Create Sessions**: Allow coaches to schedule and manage their own sessions.

### `BusinessController`

- **Business Management**: 
  - **Manage Sessions**: Oversee training sessions from a business perspective.

- **Business Dashboard**: 
  - **Manage Tools**: Provide tools for businesses to oversee and manage their training sessions.

- **Business Session Creation**: 
  - **Create and Manage Sessions**: Facilitate the creation and management of sessions from a business account.

### `ReservationController`

- **Pitch Reservations**: 
  - **Manage Reservations**: Handle reservations for pitches, including availability checks and booking confirmations.

### `NotificationController`

- **Session Notifications**: 
  - **Send Notifications**: Notify users about session updates, cancellations, and reminders.

- **User Notifications**: 
  - **Send Notifications**: Notify users about invitations, changes, and updates related to their sessions.

### `FeedbackController`

- **Session Feedback**: 
  - **Collect Feedback**: Manage the collection and processing of feedback for sessions.

- **Coach Feedback**: 
  - **Handle Coach Feedback**: Process feedback related to coaches and their performance.

### `QRCodeController`

- **QR Code Generation**: 
  - **Generate QR Codes**: Create QR codes for verifying session participation.

- **QR Code Validation**: 
  - **Validate QR Codes**: Ensure the integrity of session participation through QR code validation.

## Running Tests

To run the tests for the project, use the following command:

```bash
go test ./...
```

This command will execute all tests in your project. Make sure to create test files following Go’s testing conventions (`*_test.go`) to ensure coverage.


### **Summary of Updates**

- **Environment Variables Section**: Added instructions for creating a `.env` file.
- **Test Commands Section**: Added instructions for running tests with `go test ./...`.

By including these updates, your `README.md` will provide clear instructions for setting up, running, and testing your Go microservice project.


