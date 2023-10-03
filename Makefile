.PHONY: all clean

# Define the targets and their dependencies
all: client server

# Build the client binary
client: client.go
	go build -o client client.go

# Build the server binary
server: server.go
	go build -o server server.go

# Clean up the generated executables
clean:
	rm -f client server