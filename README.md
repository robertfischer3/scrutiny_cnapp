# Scrutiny-CNAPP
### A fischer3.net project
A open source structured Go application following best practices and modern architecture. This project provides a solid foundation for building scalable, maintainable CNAPP security applications.

## Features

- **Clean Architecture**: Separation of concerns with domain-driven design principles
- **API Ready**: RESTful API structure with middleware support
- **Configuration Management**: Environment-based configuration using Viper
- **Structured Logging**: JSON-formatted logging with different log levels
- **Robust Error Handling**: Consistent error types and handling patterns
- **Database Integration**: Ready-to-use database abstractions and migrations
- **Comprehensive Testing**: Unit, integration, and end-to-end testing setup
- **Docker Support**: Containerization for consistent development and deployment
- **CI/CD Integration**: GitHub Actions workflow for automated testing and deployment
- **Documentation**: Auto-generated API documentation
- **Monitoring**: Prometheus metrics and health check endpoints

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and docker-compose (optional, for containerized development)
- Make (for using the Makefile commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/project-name.git
cd project-name

# Install dependencies
go mod download

# Build the application
make build
```

### Running the Application

```bash
# Run directly with Go
make run

# Or build and run the binary
make build
./build/app

# Or using Docker
docker-compose up
```

## Development

This project follows standard Go project layout and best practices. The codebase is organized into logical modules with clear separation of concerns.

### Project Structure

- `cmd/`: Application entry points
- `internal/`: Private application code
- `pkg/`: Public libraries
- `api/`: API definitions and documentation
- `configs/`: Configuration files
- `scripts/`: Utility scripts

### Development Workflow

1. Create a feature branch
2. Implement changes with tests
3. Ensure all tests pass (`make test`)
4. Verify code quality (`make lint`)
5. Submit a pull request

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

