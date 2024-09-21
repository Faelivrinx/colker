# Colker

Colker is a small Go application that listens to Docker container events (`create`, `start`, `die`) and sends notifications to your favorite channel (for now, it's limited to `stdout`). The tool is useful for monitoring the health status of Docker containers and notifying you when containers are created, started, or stopped. It also manages health checks for specific containers after starting.

## Features

- Listens for Docker container lifecycle events (`create`, `start`, `die`)
- Sends notifications to a notifier (currently, `stdout`).
- Automatically registers health check after launching a container (if status url provided)
- Configurable via a `config.yaml` file.

## Table of Contents

- [Installation](#installation)
- [Example Configuration](#example-configuration)
- [Contributing](#contributing)

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (1.17 or higher)
- [Docker](https://www.docker.com/) installed and running

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/Faelivrinx/colker.git
   cd colker
   ```

2. Build the application:

   ```bash
   go build -o colker
   ```

3. Run the application:
   ```bash
   ./colker
   ```

Make sure Docker is running on your system for the tool to listen to Docker events.

## Configuration

The application uses a `config.yaml` file to specify containers, health check URLs, messages, and secrets.

### Example Configuration

Save this file as `config.yaml` in the root of your project directory:

```yaml
containers:
  - name: "my_container_1"
    status_url: "http://localhost:8080/status"
  - name: "my_container_2"
    status_url: "http://localhost:8081/status"

messages:
  start_message: "Container %s (version: %s) started and health check is in progress."
  stop_message: "Container %s stopped."
  final_message: "Health check monitoring is completed."

secret:
  secret_value: "my_secret_value"
```

### Contributing

Contributions are very welcome! If you'd like to contribute to Colker, feel free to fork the repository, make changes, and submit a pull request.

You can easily extend the functionality by adding new notifiers (for example, Slack, email, etc.). To do this, you would need to implement the Notifier interface found in the internal package.
