# Dog App Backend (BFF)

[![Go](https://github.com/agugliotta/dog-app-bff/actions/workflows/go.yml/badge.svg)](https://github.com/agugliotta/dog-app-bff/actions/workflows/go.yml)

This is the backend for a Dog App, acting as a Backend For Frontend (BFF). It provides APIs to manage dog breeds and pet information.

## Overview

The backend is built using Go and leverages a PostgreSQL database for data persistence. It exposes RESTful APIs to be consumed by frontend applications (e.g., a mobile app).

### Key Features

* **Breed Management:**
    * Retrieve a list of all dog breeds.
    * Retrieve details for a specific dog breed.
* **Pet Management:**
    * Retrieve a list of all registered pets.
    * Retrieve details for a specific pet.
    * Create new pet records.

## Getting Started

These instructions will guide you on how to set up and run the backend application locally for development and testing.

### Prerequisites

* Go (version 1.24 or higher)
* Docker (for running the PostgreSQL database)
* Make (for running common development tasks)

### Running Locally

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd dog-app-bff
    ```

2.  **Start the PostgreSQL database using Docker:**
    ```bash
    make db-start
    ```
    This command will start a PostgreSQL container named `dog-app-bff-postgres-test`.

3.  **Run the backend application:**
    ```bash
    make run
    ```
    The backend server will start on port `8080`.

### API Endpoints

* `GET /api/v1/breeds`: Get all dog breeds.
* `GET /api/v1/breeds/{id}`: Get a specific dog breed by ID.
* `GET /api/v1/pets`: Get all registered pets.
* `GET /api/v1/pets/{id}`: Get a specific pet by ID.
* `POST /api/v1/pets`: Create a new pet.

### Running Tests

The project includes unit and integration tests to ensure the reliability of the codebase.

1.  **Ensure the test database is running:**
    If you haven't already, run the test database:
    ```bash
    make db-start
    ```

2.  **Run all tests:**
    ```bash
    make test
    ```
    This command will execute both unit and integration tests.

### GitHub Actions

The project is configured with GitHub Actions for Continuous Integration. The workflow defined in `.github/workflows/go.yml` automatically builds and tests the application on push and pull requests to the `main` branch.

The badge at the top of this `README` indicates the current build status of the `main` branch. A green badge signifies a successful build.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for any bugs or feature requests.

## License

This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.