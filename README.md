# Project Name

A Go project for calling Python functions in runtime with cgo and protobuf

## Table of Contents

- [Introduction](#introduction)
- [Contributing](#contributing)
- [License](#license)

[//]: # (- [Installation]&#40;#installation&#41;)
[//]: # (- [Usage]&#40;#usage&#41;)
[//]: # (- [Example]&#40;#example&#41;)

## Introduction

Imagine yourself as a Python!

This project aims to ensure smooth integration of Go and Python by allowing Python
functions to be called at runtime in a Go application.
Integration is achieved through the use of cgo, which allows direct calls to Python
functions via a custom rpc client-server based on the classic protobuf.

[//]: # (To work correctly, you must have compiled Cpython and generated grpc services.)

[//]: # (## Installation)

[//]: # (1. **Clone the repository:**)

[//]: # ()
[//]: # (    ```bash)

[//]: # (    git clone https://github.com/yourusername/yourproject.git)

[//]: # (    cd yourproject)

[//]: # (    ```)

[//]: # ()
[//]: # (2. **Build the project:**)

[//]: # ()
[//]: # (    ```bash)

[//]: # (    go build)

[//]: # (    ```)

[//]: # (## Usage)

[//]: # ()
[//]: # (To use this project in your Go application, follow these steps:)

[//]: # ()
[//]: # (1. Import the package in your Go code:)

[//]: # ()
[//]: # (    ```go)

[//]: # (    import "github.com/yourusername/yourproject")

[//]: # (    ```)

[//]: # ()
[//]: # (2. Use the provided functions to call Python functions:)

[//]: # ()
[//]: # (    ```go)

[//]: # (    // Example Go code calling a Python function)

[//]: # (    result, err := yourproject.CallPythonFunction&#40;"your_python_function", args...&#41;)

[//]: # (    if err != nil {)

[//]: # (        log.Fatal&#40;err&#41;)

[//]: # (    })

[//]: # (    fmt.Println&#40;"Result from Python:", result&#41;)

[//]: # (    ```)

[//]: # ()
[//]: # (3. Start the Python runtime environment before making calls:)

[//]: # ()
[//]: # (    ```go)

[//]: # (    err := yourproject.StartPythonRuntime&#40;&#41;)

[//]: # (    if err != nil {)

[//]: # (        log.Fatal&#40;"Failed to start Python runtime:", err&#41;)

[//]: # (    })

[//]: # (    defer yourproject.StopPythonRuntime&#40;&#41;)

[//]: # (    ```)

[//]: # (## Example)

[//]: # ()

[//]: # (Check the `example` directory for a comprehensive example of using this project.)

## Contributing

Feel free to contribute to this project by opening issues or submitting pull requests. Your feedback and contributions
are highly appreciated.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.