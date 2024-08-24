# JWT Auth Service

## Overview

The **JWT Auth Service** is responsible for authenticating users via username and password and issuing JSON Web Tokens (JWT) to allow access to protected resources. 

## Features

- **User Authentication**: Validates user credentials (username and password) to ensure authenticity.
- **JWT Generation**: Issues JWTs upon successful authentication (Used to access protected resources)
- **Token-Based Authorization**: Secures endpoints and resources by requiring a valid JWT for access.
- **Token Expiration and Refresh**: (Optional) Manages token lifecycles, including expiration and refreshing tokens (Not yet implemented).
- **Cross-Origin Resource Sharing (CORS)**: Configured to handle requests from specific origins.

## Technologies

- **Language**: Go
- **Framework**: Fiber
- **Authentication**: JWT (RS256)
- **Database**:
  - **ORM: GORM
  - **Development**: SQLite
  - **Production & QA**: MySQL Server

## Setup and Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Elimists/auth-proreact.git
   cd jwt-auth-service
