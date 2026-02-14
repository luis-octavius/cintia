<h1 align="center">Cintia</h1>

<p align="center">
  <!-- Core Technologies -->
  <img src="https://img.shields.io/badge/Go-1.21-blue?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Gin-1.9.1-00ADD8?style=for-the-badge&logo=gin" alt="Gin">
  <img src="https://img.shields.io/badge/Cobra-1.8.0-purple?style=for-the-badge&logo=go" alt="Cobra">
  
  <!-- Data & Messaging -->
  <img src="https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/RabbitMQ-3.12-FF6600?style=for-the-badge&logo=rabbitmq" alt="RabbitMQ">
  
  <!-- Status -->
  <img src="https://img.shields.io/badge/status-active-success?style=for-the-badge" alt="Status">
  <img src="https://img.shields.io/badge/CLI-Cobra-purple?style=for-the-badge" alt="CLI">
</p>

**Cintia** is a comprehensive job search management platform built with Go, combining automated job scraping with intelligent application tracking.

## Core Features

- **Automated Job Scraping**: Continuously collects job listings from LinkedIn, Indeed, and other sources
- **Application Management**: Track all your job applications in one place with status updates
- **Smart Notifications**: RabbitMQ-powered alerts for upcoming interviews and application deadlines
- **Personal Dashboard**: Visualize your job search progress with statistics and insights
- **Multi-Interface Access**: REST API + Command-line interface (Cobra) for administration and automation

## Architecture

- **Go Backend**: RESTful API with Gin framework
- **PostgreSQL**: Primary database with sqlc for type-safe queries
- **RabbitMQ**: Async notifications and reminders
- **JWT Authentication**: Secure user management
- **Scraper Service**: Concurrent job collection from multiple sources
- **Cobra CLI**: Full-featured command-line tool for administration and automation

## Roadmap

- [x] User authentication & profiles
- [ ] Job scraping pipeline
- [ ] Application tracking
- [ ] RabbitMQ integration for reminders
- [ ] Email notifications
- [ ] Analytics dashboard
- [ ] CLI tool with Cobra

## Project Structure

```bash
cmd/
├── api/          # HTTP API server (Gin)
├── cli/          # Command-line interface (Cobra)
└── scraper/      # Scraping service workers

internal/
├── user/         # User domain
├── job/          # Job listings domain
├── application/  # Applications tracking
├── scraper/      # Scraping logic
└── notification/ # RabbitMQ consumers
```

## Contributing

Feel free to contribute!
