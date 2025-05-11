# 🚀 My Golang Learning Journey 🚀

<div align="center">
  <img src="https://golang.org/doc/gopher/doc.png" width="300" />
  <h3>Learning Go One Line at a Time</h3>
</div>

## 📚 Project Overview

This project serves as my personal playground for learning the Go programming language (Golang). I've built a RESTful API with user authentication, post management, and proper API documentation using Swagger.

> "The best way to learn is by doing" - and that's exactly what this project represents!

## 🛠️ Tech Stack

- **Go** (v1.24.0) - Main programming language
- **Gin** - Web framework
- **Swagger** - API documentation
- **JWT** - Authentication
- **Rate Limiting** - API protection
- **Middleware** - For authentication and caching

## 🌟 Features

- 🔐 **User Authentication**

  - Registration
  - Login
  - JWT-based authentication

- 📝 **Post Management**

  - Create, read, update, delete posts
  - Post ownership and permissions

- 📋 **API Documentation**

  - Comprehensive Swagger documentation
  - Interactive API testing

- 🛡️ **Security Features**
  - Password hashing
  - JWT authentication
  - Rate limiting for DDoS protection

## 🏗️ Project Structure

```
.
├── config/           # Database and authentication configuration
├── controllers/      # Request handlers and business logic
├── docs/             # Swagger API documentation
├── middleware/       # Authentication and caching middleware
├── models/           # Data models
├── routes/           # API routing configuration
├── utils/            # Helper functions (JWT, hashing)
├── go.mod            # Go module dependencies
├── go.sum            # Checksums for dependencies
└── main.go           # Application entry point
```

## 🚀 Getting Started

### Prerequisites

- Go 1.23+ installed on your machine
- Git for version control
- Basic understanding of REST APIs

### Running the Project

1. Clone this repository:

   ```bash
   git clone <repository-url>
   cd project-go
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run the application:

   ```bash
   go run main.go
   ```

4. Access the API documentation:
   ```
   http://localhost:8080/swagger/index.html
   ```

## 📝 What I've Learned

- ✅ Go syntax and language fundamentals
- ✅ Web service development with Gin
- ✅ API documentation with Swagger
- ✅ Authentication with JWT
- ✅ Middleware implementation
- ✅ Testing in Go
- ✅ Project structuring and organization

## 🔮 Future Improvements

- [ ] Add comprehensive unit and integration tests
- [ ] Implement database migrations
- [ ] Add user profiles and avatar uploads
- [ ] Implement social features like comments and likes
- [ ] Deploy to a cloud provider

## 💭 Reflections

This project has been an incredible journey into the world of Go. I've particularly enjoyed the language's simplicity and performance. The concurrency model with goroutines is powerful yet intuitive.

> "Go is not meant to innovate programming theory. It's meant to innovate programming practice." - Rob Pike

<div align="center">
  <p>Happy Coding! 💻</p>
  <p>Made with ❤️ and lots of ☕</p>
</div>
