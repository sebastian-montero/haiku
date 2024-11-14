# Flow ꩜

Flow is a real-time writing platform where users can create notebooks, start writing sessions, and view live updates from other contributors. The application leverages WebSockets to enable real-time interaction in writing sessions, with all data stored in a PostgreSQL database.

## Features

- **User Management**: Register, update, and delete users.
- **Notebook Creation**: Users can create and manage notebooks.
- **Real-Time Writing**: A single users can write in a session with updates visible in real-time.
- **Session Management**: Start, view, and end writing sessions. Session states are updated and saved to maintain content integrity.
- **Data Persistence**: All data is stored in PostgreSQL with clear structures for users, notebooks, sessions, and content.

## Setup and Installation

### Prerequisites

- **Docker**: Used to set up and run PostgreSQL and other dependencies.
- **Go**: To build and run the server.

### Instructions

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/flow.git
   cd flow
   ```
2. Ensure docker-compose.yml is configured with PostgreSQL, server, and client services.
  ```bash
   docker-compose up -d
   ```

## Server functionality

## User Endpoints

### Create User
- **URL**: `/users`
- **Method**: `POST`
- **Description**: Register a new user.

### Update User
- **URL**: `/users`
- **Method**: `PUT`
- **Description**: Update an existing user’s details.

### Get User by ID
- **URL**: `/users/{id}`
- **Method**: `GET`
- **Description**: Retrieve a user by their unique ID.

### Delete User by ID
- **URL**: `/users/{id}`
- **Method**: `DELETE`
- **Description**: Delete a user by their unique ID.

---

## Notebook Endpoints

### Create Notebook
- **URL**: `/notebooks`
- **Method**: `POST`
- **Description**: Create a new notebook.

### Update Notebook
- **URL**: `/notebooks`
- **Method**: `PUT`
- **Description**: Update an existing notebook.

### Get Notebook by ID
- **URL**: `/notebooks/{id}`
- **Method**: `GET`
- **Description**: Retrieve a notebook by its unique ID.

### Get Notebooks by Owner ID
- **URL**: `/notebooks/by_owner/{owner_id}`
- **Method**: `GET`
- **Description**: Retrieve all notebooks created by a specific owner.

### Delete Notebook by ID
- **URL**: `/notebooks/{id}`
- **Method**: `DELETE`
- **Description**: Delete a notebook by its unique ID.

---

## Session Endpoints

### Create Session
- **URL**: `/sessions`
- **Method**: `POST`
- **Description**: Start a new writing session for a notebook.

### Get Active Sessions
- **URL**: `/sessions`
- **Method**: `GET`
- **Description**: Retrieve all currently active sessions.

### Get Session by ID
- **URL**: `/sessions/{id}`
- **Method**: `GET`
- **Description**: Retrieve details of a specific session.

### End Session
- **URL**: `/sessions/{id}/end`
- **Method**: `PUT`
- **Description**: Mark a session as inactive and save the final content.

### Delete Session by ID
- **URL**: `/sessions/{id}`
- **Method**: `DELETE`
- **Description**: Delete a session by its unique ID.

---

## Content Endpoints

### Create Content
- **URL**: `/content`
- **Method**: `POST`
- **Description**: Create content within an active session.

### Get Latest Content by Session ID
- **URL**: `/content/by_session/{session_id}`
- **Method**: `GET`
- **Description**: Retrieve the most recent content for a specific session.

---

## WebSocket Endpoint

### Session WebSocket Connection
- **URL**: `/ws/{conn_type}/{notebook_id}`
- **Method**: `GET`
- **Description**: Establish a WebSocket connection for real-time collaboration.

