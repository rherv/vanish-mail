# Vanish Mail

A temporary email inbox generator written in Go.
This project is inspired by: https://tmail.link

## Installation

Make sure you have Go installed. If not, you can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Clone the repository:
```bash
git clone https://github.com/rherv/vanish-mail.git
```
### Change into the project directory:
```bash
cd vanish-mail
```
### Build the application:
```bash
go build
```

## Usage

### Run the application:
```bash
./vanish-mail -domain example.com -http 8080 -smtp 25 -delay 10
```

### Command-line Options
```
-domain: the domain to accept emails for (default is "localhost").
-http: HTTP service address (default is 8080).
-smtp: SMTP service address (default is 25).
-delay: The time in minutes to keep an email for (default is 10).
```

## Contributing
Feel free to contribute by opening issues or pull requests.
