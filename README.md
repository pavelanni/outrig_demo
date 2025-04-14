# Outrig Demo Program

This is a demo program designed to showcase the capabilities of [Outrig](https://outrig.run/), a debugging utility for Go programs. The program implements a simple HTTP server that demonstrates various features that can be monitored and debugged using Outrig.

## Features Demonstrated

* **Logging**: The program uses `log/slog` for structured logging
* **Goroutines**: Background goroutine for periodic logging
* **Variable Watching**: Global variables that can be monitored
* **Memory Management**: Dynamic memory allocation and deallocation
* **Runtime Statistics**: Request counting and memory usage tracking

## API Endpoints

### GET `/stats`

Returns current statistics:

* Request count
* Memory allocated (in MB)
* Debug mode status

### POST `/config`

Updates server configuration:

```json
{
  "max_memory_mb": 100,
  "debug_mode": true
}
```

### POST `/memory`

Manages memory allocation:

```json
{
  "action": "allocate",
  "size_mb": 50
}
```

or

```json
{
  "action": "release"
}
```

## Running the Demo

1. Run the demo program:

   ```bash
   go run main.go
   ```

2. Use Outrig to monitor:

   * Logs in stdout
   * Goroutine states
   * Variable changes
   * Memory usage
   * CPU statistics

   Start the Outrig UI:

   ```bash
   outrig server
   ```

   Open the Outrig UI in your browser:

   ```bash
   open http://localhost:5005
   ```

## Example Usage

1. Start the server:

   ```bash
   go run main.go
   ```

2. Allocate memory:

   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"action":"allocate","size_mb":50}' http://localhost:8080/memory
   ```

3. Check stats:

   ```bash
   curl http://localhost:8080/stats
   ```

4. Update config:

   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"max_memory_mb":200,"debug_mode":true}' http://localhost:8080/config
   ```

5. Release memory:

   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"action":"release"}' http://localhost:8080/memory
   ```

While running these commands, use Outrig to observe:

* Log messages in real-time
* Goroutine creation and state changes
* Memory allocation and deallocation
* Variable value changes
* Runtime statistics

## Running as a Container

You can also run the demo as a container:

1. Build the container image:

   ```bash
   podman build -t outrig-demo .
   ```

2. Run the container:

   ```bash
   podman run -d -p 8080:8080 --name outrig-demo outrig-demo
   ```

3. Access the application:

   ```bash
   curl http://localhost:8080/stats
   ```

4. To stop and remove the container:

   ```bash
   podman stop outrig-demo
   podman rm outrig-demo
   ```

Note: When running in a container, you'll need to ensure the Outrig UI is running on your host machine and properly configured to connect to the containerized application.