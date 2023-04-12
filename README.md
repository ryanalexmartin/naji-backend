# Go Omegle Clone

A simple and extensible Omegle-like anonymous chat service built with Go and WebSockets.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Introduction

This project aims to provide a clean, scalable, and extensible backend for an Omegle-like chat service using Go's concurrency 
features and the Gorilla WebSocket package. It demonstrates how to handle multiple user connections, manage chat sessions, 
and broadcast messages to users in a concurrent and efficient manner.

## Features

- Concurrent handling of user connections and chat sessions
- User matching and chat room management
- Message broadcasting to users within a chat room
- Graceful handling of WebSocket connections, messages, and disconnections
- Modular project structure for maintainability and scalability

## Requirements

- Go 1.16 or higher
- Gorilla WebSocket package
- (Optional) Gin or Echo web framework




## Initial outline plan of action

    Plan your architecture:
        Use a concurrent, event-driven model, which is well-suited for handling multiple users and connections.
        Implement a mechanism to match users and manage chat sessions.
        Separate concerns by dividing the code into packages/modules, such as handlers, models, and utilities.

    Set up your dependencies:
        Install the Gorilla WebSocket package (github.com/gorilla/websocket) to manage WebSocket connections.
        Consider using a framework like Gin (github.com/gin-gonic/gin) or Echo (github.com/labstack/echo) to simplify routing and middleware.

    Create data structures:
        Define a User struct to represent connected users, with fields like ID, Conn (WebSocket connection), and Room (chat session ID).
        Define a Room struct to represent chat sessions, with fields like ID, Users, and a message buffer.
        Use a concurrent-safe data structure like sync.Map to manage user and room collections.

    Implement WebSocket handlers:
        Create a function to upgrade incoming HTTP connections to WebSocket connections.
        Implement a handleConnection function that reads messages from users, validates them, and dispatches them to the appropriate room.
        Implement a handleDisconnection function to gracefully close WebSocket connections and remove users from rooms.

    Implement chat session management:
        Create a matchUsers function to pair users and assign them to a room.
        Implement a function to broadcast messages to all users in a room.
        Optionally, implement a mechanism for users to leave or end chat sessions.

    Implement the main application:
        Set up your server with appropriate routes, middleware, and handlers.
        Start a background goroutine to continuously match users.
        Run the server and listen for incoming connections.

    Testing and optimization:
        Test your service with multiple concurrent connections and ensure that the matching and messaging are working as expected.
        Optimize your code for performance, memory usage, and error handling.
        Consider using monitoring tools like Prometheus (prometheus.io) to track performance metrics.





## More detailed architecture

    Concurrent and event-driven model:
        Use goroutines to handle multiple WebSocket connections and chat sessions concurrently.
        Use channels to communicate between goroutines and handle events like user connections, disconnections, and messages.

    Modular design:
        Separate your code into packages/modules for better organization and maintainability:
            main: Entry point of your application, setting up routes, and starting the server.
            handlers: Functions to handle incoming WebSocket connections, messages, and disconnections.
            models: Structs for User, Room, and other necessary data structures.
            services: Functions for chat session management, user matching, and message broadcasting.
            utils: Helper functions, constants, and error handling.

    User and Room management:
        Use sync.Map or similar concurrent-safe data structures to store and manage users and rooms.
        Implement functions to add, remove, and update users and rooms in the maps.

    WebSocket handling:
        Use the Gorilla WebSocket package to manage WebSocket connections.
        Implement handlers for WebSocket connections, messages, and disconnections.
        Ensure that all WebSocket operations are handled gracefully to prevent data races and deadlocks.

    Chat session management:
        Implement a user-matching mechanism that pairs users and assigns them to a room.
        Implement message broadcasting to send messages to all users in a room.
        Optionally, allow users to leave or end chat sessions, and handle room cleanup accordingly.

    Server setup:
        Use a framework like Gin or Echo to set up routes, middleware, and handlers.
        Start a background goroutine to continuously match users.
        Run the server and listen for incoming connections.

By following this architecture plan, you will create a well-structured, extensible, and elegant backend for your 
Omegle-like service. This will allow you to easily modify and expand your application in the future while minimizing 
the chances of introducing race conditions and other concurrency-related issues.



### Tests to implement:

    User and Room management (models package):
        TestUserCreation: Test creating a new user with a unique ID and valid WebSocket connection.
        TestRoomCreation: Test creating a new room with a unique ID and no users initially.
        TestAddUserToRoom: Test adding a user to a room and updating the user's room ID.
        TestRemoveUserFromRoom: Test removing a user from a room and updating the user's room ID.

    WebSocket handling (handlers package):
        TestWebSocketUpgrade: Test upgrading an HTTP request to a WebSocket connection.
        TestHandleMessage: Test handling and dispatching incoming messages from users.
        TestHandleDisconnection: Test gracefully closing a WebSocket connection and removing the user from the room.

    Chat session management (services package):
        TestMatchUsers: Test pairing users and assigning them to a room.
        TestBroadcastMessage: Test sending a message to all users in a room.
        TestLeaveChatSession: Test allowing users to leave or end chat sessions, and handling room cleanup accordingly.
