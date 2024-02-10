# Project Overview

This project is structured around five core components, designed to handle email notifications based on currency exchange rates from the Coinbase service. Below is a brief overview of each component and their functionalities:

### Core Components:
- **Handler** (`/internal/handlers/handler.go`): Manages incoming requests and interfaces with the EmailService to execute the main logic which includes fetching the current exchange rate, subscribing users to notifications, and sending emails.
- **EmailService** (`/internal/service/email.go`): Acts as the central business logic layer, coordinating with Fstorage, CoinbaseClient, and SmtpClient to manage rates, subscriptions, and email notifications.
- **Fstorage** (`/internal/storage/fstorage.go`): Handles data persistence by saving email entries to a file, leveraging a hash table for efficient lookup.
- **CoinbaseClient** (`/internal/client/coinbase.go`): Fetches the current exchange rate from Coinbase's API.
- **SmtpClient** (`/internal/client/smtp.go`): Sends out email notifications using an SMTP server.

### Handler Functions:
- **getRate()**: Retrieves the current rate from the EmailService. In case of an issue, it returns a 400 status code.
- **Subscribe()**: Registers an email for notifications, checking for duplicates and validating the email format.
- **SendEmails()**: Triggers the sending of email notifications to all subscribed users, tracking any errors.

### EmailService Methods:
- **Rate()**: Contacts CoinbaseClient for the latest exchange rate.
- **SendEmail()**: Gathers emails from Fstorage and sends out notifications, collecting any unsent emails due to errors.
- **Subscribe()**: Adds a new email to the subscription list, ensuring no duplicates.

### Fstorage:
- Manages a file and a hash table for storing email addresses.
- **Append()**: Adds a new email, ensuring uniqueness.
- **GetEmails()**: Retrieves a list of subscribed emails.
- **buildIndex()**: Initializes the hash table during startup.

### External Communication:
- **CoinbaseClient** makes a GET request to Coinbase for the latest exchange rate in UAH.
- **SmtpClient** interfaces with an SMTP server to send emails.

### Application Setup and Launch:
1. Clone the repository.
2. Edit `config.yml` with your SMTP server details (e.g., Gmail).
3. Build a Docker image: `docker build --tag go-email .`
4. Run the container: `docker run -p 1313:1313 go-email` (ensure the first port matches `bind_ip` in the config).
