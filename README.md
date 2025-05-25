# Stori Challenge - Transaction Processing System

A robust Go application that processes CSV transaction files and sends email summaries, built with SOLID principles and clean architecture.

## Features

✅ **Core Requirements**
- CSV transaction file processing with validation
- Email summary generation with HTML formatting
- Monthly transaction grouping and statistics
- Average credit/debit amount calculations
- Docker containerization

✅ **Bonus Features**
- PostgreSQL database integration for transaction persistence
- Styled HTML emails with responsive design
- File watching for automated processing
- Comprehensive logging and error handling
- Configuration management via environment variables

## Architecture

The system follows **SOLID principles** with clean architecture:

- **Domain Layer**: Core business logic and models
- **Services Layer**: Business services and interfaces
- **Infrastructure Layer**: External dependencies (database, email, file system)
- **Configuration**: Environment-based configuration management

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)

### 1. Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd stori-challenge

# Copy environment configuration
cp .env.example .env

# Edit .env with your SMTP credentials
# SMTP_USERNAME=your-email@gmail.com
# SMTP_PASSWORD=your-app-password

# Start the application
docker-compose up --build
```

### 2. Local Development

```bash
# Install dependencies
go mod download

# Set environment variables
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export RECIPIENT_EMAIL="recipient@example.com"
export WATCH_DIRECTORY="./data"

# Run the application
go run cmd/processor/main.go
```

## Usage

### Processing Transaction Files

1. **Automatic Processing**: Place CSV files in the watched directory (`/data` by default)
2. **Manual Processing**: Copy the sample file to trigger processing:
   ```bash
   cp data/transactions.csv /data/new_transactions.csv
   ```

### CSV File Format

```csv
Id,Date,Transaction
0,7/15,+60.5
1,7/28,-10.3
2,8/2,-20.46
```

**Requirements:**
- Header: `Id,Date,Transaction`
- Date format: `M/D` (assumes current year)
- Transaction: `+amount` (credit) or `-amount` (debit)

### Email Summary

The system sends HTML emails containing:
- **Total account balance**
- **Monthly transaction counts**
- **Average credit and debit amounts**
- **Responsive styling with Stori branding**

## Configuration

Configure the application using environment variables:

### Required Settings

```bash
# Email Configuration
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
RECIPIENT_EMAIL=recipient@example.com
```

### Optional Settings

```bash
# Application
ENVIRONMENT=development          # development|production
LOG_LEVEL=info                  # debug|info|warn|error

# File Processing
WATCH_DIRECTORY=/data
PROCESSED_DIRECTORY=/data/processed

# Email Settings
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=noreply@stori.com

# Database (Optional)
DB_HOST=localhost
DB_PORT=5432
DB_USER=stori_user
DB_PASSWORD=stori_password
DB_NAME=stori_challenge
```

## Development

### Project Structure

```
├── cmd/processor/             # Application entry point
├── internal/
│   ├── config/                # Configuration management
│   ├── domain/                # Domain models and business logic
│   ├── services/              # Business services
│   └── infrastructure/        # External dependencies
│       ├── database/          # Database implementation
│       ├── email/             # Email service implementation
│       └── file/              # File processing implementation
├── data/                      # Sample data and test files
├── docker-compose.yml         # Local development environment
└── Dockerfile                 # Production container
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/domain/... -v
```

### Building

```bash
# Build binary
go build -o processor cmd/processor/main.go

# Build Docker image
docker build -t stori-processor .
```

## Email Configuration

### Gmail Setup

1. Enable 2-factor authentication
2. Generate an App Password:
   - Google Account → Security → App passwords
   - Select "Mail" and generate password
3. Use the generated password as `SMTP_PASSWORD`

### Other SMTP Providers

Update the SMTP configuration for your provider:

```bash
# Example: Outlook
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587

# Example: SendGrid
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
```

## Database Features (Bonus)

When database configuration is provided:

- **Transaction Persistence**: All transactions are saved to PostgreSQL
- **Account Management**: Automatic account creation and linking
- **Data Integrity**: Database transactions ensure consistency
- **Migrations**: Automatic schema setup

### Database Access

Access the database using the included Adminer:
- URL: http://localhost:8080
- Server: postgres
- Username: stori_user
- Password: stori_password
- Database: stori_challenge

## Monitoring and Logging

### Log Levels

- **DEBUG**: Detailed diagnostic information
- **INFO**: General operational messages
- **WARN**: Warning conditions
- **ERROR**: Error conditions that don't stop the application
- **FATAL**: Critical errors that cause application shutdown

### Health Checks

The Docker container includes health checks:
```bash
# Check container health
docker-compose ps

# View logs
docker-compose logs processor
```

## Production Deployment

### Docker

```bash
# Build production image
docker build -t stori-processor:latest .

# Run with production configuration
docker run -d \
  --name stori-processor \
  -e ENVIRONMENT=production \
  -e SMTP_USERNAME=your-email@gmail.com \
  -e SMTP_PASSWORD=your-password \
  -e RECIPIENT_EMAIL=notifications@yourcompany.com \
  -v /host/data:/data \
  stori-processor:latest
```

### AWS Lambda (Bonus)

For serverless deployment:
1. Modify the main function to handle Lambda events
2. Use AWS SES for email delivery
3. Use S3 for file storage and triggers
4. Use RDS for database persistence

## Troubleshooting

### Common Issues

**Email not sending:**
```bash
# Check SMTP credentials
docker-compose logs processor | grep "SMTP"

# Test with debug logging
docker-compose up --build -e LOG_LEVEL=debug
```

**File not processed:**
```bash
# Check file format and permissions
# Ensure CSV has correct headers: Id,Date,Transaction
# Verify watch directory is mounted correctly
```

**Database connection failed:**
```bash
# Check database is running
docker-compose ps postgres

# Verify connection settings
docker-compose logs postgres
```

### Debug Mode

Enable debug logging for detailed information:

```bash
# Docker Compose
LOG_LEVEL=debug docker-compose up

# Local development
LOG_LEVEL=debug go run cmd/processor/main.go
```

## API Documentation

This application processes files automatically and doesn't expose HTTP endpoints. The main interface is:

1. **Input**: CSV files in the watched directory
2. **Output**: HTML email summaries to configured recipients
3. **Storage**: Optional PostgreSQL database for persistence

## Contributing

1. Follow the existing code style and architecture
2. Add tests for new functionality
3. Update documentation for significant changes
4. Ensure Docker builds successfully

## License

This project is part of the Stori technical challenge.

---

**Built with ❤️ using Go, following SOLID principles and clean architecture patterns.**