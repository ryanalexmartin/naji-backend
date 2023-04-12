    TestServerIsRunning (main_test.go)
        Start server in a goroutine
        Send an HTTP request to the server
        Check if the server is running and responding with status code 200

    startServer (main.go)
        Set up Gin or Echo router
        Define an HTTP GET route for the root path
        Start the server

    TestUserCreation (models/user_test.go)
        Create a new user with a WebSocket connection
        Check if the user has a unique ID, valid connection, and an empty room ID

    NewUser (models/user.go)
        Define User struct with ID, connection, and room ID fields
        Create a function to create a new user with a unique ID and a WebSocket connection

    TestRoomCreation (models/room_test.go)
        Create a new room
        Check if the room has a unique ID and no users initially

    NewRoom (models/room.go)
        Define Room struct with ID and Users fields (sync.Map)
        Create a function to create a new room with a unique ID and an empty Users map

You can use this compact representation to communicate the progress to future ChatGPT prompts concisely, staying within the token limit.

