# Ethereum Transaction Parser

## Overview

This project is an Ethereum transaction parser designed to monitor and process transactions for subscribed Ethereum addresses. The parser is built with flexibility and extensibility in mind, allowing for customizable components and concurrent processing to efficiently handle blockchain data.

## Features

- **Customizable HTTP Client**: The parser uses an HTTP client for interacting with the Ethereum blockchain via JSON-RPC. This client is customizable, allowing you to set timeouts, delays, retries, and other configurations to suit your needs.

- **Pluggable Storage Layer**: Storage is abstracted behind an interface, enabling you to replace the default in-memory storage with other implementations (e.g., database storage) without changing the parser logic.

- **Concurrent Parsing and Processing**: The parser runs in a separate goroutine to continuously fetch and parse new blocks. It updates the storage concurrently, allowing for real-time transaction data retrieval for subscribed addresses.

- **Thread-Safe Operations**: All operations, including subscription management and transaction retrieval, are thread-safe, ensuring data integrity in a concurrent environment.

## Testing
- **Unit Tests**: The codebase includes unit tests for core components, ensuring reliability and correctness.

- **Mock Implementations**: Mocks are used for testing external dependencies, such as the Ethereum client and storage layer, facilitating comprehensive testing without relying on actual external services.