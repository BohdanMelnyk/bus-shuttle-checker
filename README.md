# Bus Shuttle Checker

An automated tool that monitors Parks Canada shuttle bus availability for popular destinations like Moraine Lake and Lake O'Hara. The application helps you secure shuttle tickets by automatically checking availability and sending notifications when slots open up.

## What it does

- **Continuous Monitoring**: Automatically checks shuttle availability for multiple locations and specific dates
- **Date-Specific Tracking**: Monitors availability for particular dates you're interested in
- **Smart Notifications**: Sends email alerts when slots become available, including the specific dates and locations
- **Multiple Destinations**: Tracks availability for:
  - Moraine Lake Morning
  - Moraine Lake Midday
  - Lake O'Hara
- **Real-time Updates**: Runs checks every 30 minutes to ensure you don't miss any openings
- **API Integration**: Direct integration with Parks Canada reservation system for reliable results

## Features

- Uses official Parks Canada API for accurate availability checks
- Runs checks every 30 minutes
- Sends email notifications when slots become available
- HTTP endpoints for manual checks
- Efficient and lightweight implementation
- Retry logic for improved reliability

## Setup

### 1. Environment Variables

The application requires the following environment variables:

- `MAILGUN_API_KEY`: Your Mailgun API key
- `MAILGUN_DOMAIN`: Your Mailgun domain
- `RECIPIENT_EMAIL`: Email address to receive notifications
- `SENDER_EMAIL`: Email address to send notifications from
- `PORT`: (Optional) Port for the HTTP server (default: 8080)

### 2. Running with Docker

1. Build the Docker image:
```bash
docker build -t bus-shuttle-checker .
```

2. Run the container:
```bash
docker run -d \
  -e MAILGUN_API_KEY=your_api_key \
  -e MAILGUN_DOMAIN=your_domain \
  -e RECIPIENT_EMAIL=your_email \
  -e SENDER_EMAIL=your_sender_email \
  -p 8080:8080 \
  bus-shuttle-checker
```

### 3. Running Locally

1. Clone the repository:
```bash
git clone https://github.com/<your-username>/bus-shuttle-checker.git
cd bus-shuttle-checker
```

2. Create a `.env` file with your configuration:
```bash
MAILGUN_API_KEY=your_api_key
MAILGUN_DOMAIN=your_domain
RECIPIENT_EMAIL=your_email
SENDER_EMAIL=your_sender_email
PORT=8080
```

3. Run the application:
```bash
go run main.go
```

## API Endpoints

- `GET /health` - Check if the service is running
- `GET /check-all` - Manually trigger an availability check for all locations

## Supported Locations

1. **Moraine Lake Morning**
   - Resource IDs: [-2147476652, -2147476634, -2147476641, -2147476655]
   - Booking Category: 9

2. **Moraine Lake Midday**
   - Resource IDs: [-2147476651, -2147476653]
   - Booking Category: 9

3. **Lake O'Hara**
   - Resource IDs: [-2147479230, -2147479229]
   - Booking Category: 10

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request