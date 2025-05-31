# Bus Shuttle Checker

Automated checker for bus shuttle availability. The application runs on a schedule using GitHub Actions and displays its status on GitHub Pages.

## Features

- Checks shuttle availability every 30 minutes
- Sends email notifications when slots become available
- Displays last check status on GitHub Pages
- Browser-agnostic implementation (works with Chrome/Chromium)
- Retry logic for improved reliability

## Setup

### 1. Repository Setup

1. Fork this repository
2. Go to repository Settings > Pages
3. Set the source to "GitHub Actions"

### 2. Configure Secrets

Add the following secrets in your repository's Settings > Secrets and variables > Actions:

- `MAILGUN_API_KEY`: Your Mailgun API key
- `MAILGUN_DOMAIN`: Your Mailgun domain
- `RECIPIENT_EMAIL`: Email address to receive notifications
- `SENDER_EMAIL`: Email address to send notifications from

### 3. Enable GitHub Actions

The workflow is already configured in `.github/workflows/shuttle-checker.yml`. It will:

- Run every 30 minutes
- Check shuttle availability
- Update the status page
- Send email notifications if slots are available

## Status Page

The status page is available at: `https://<your-github-username>.github.io/bus-shuttle-checker/`

## Local Development

### Prerequisites

- Go 1.21 or later
- Chrome or Chromium browser

### Running Locally

1. Clone the repository:
```bash
git clone https://github.com/<your-username>/bus-shuttle-checker.git
cd bus-shuttle-checker
```

2. Set environment variables:
```bash
export MAILGUN_API_KEY=your_api_key
export MAILGUN_DOMAIN=your_domain
export RECIPIENT_EMAIL=your_email
```

3. Run the checker:
```bash
go run main.go check-all
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request