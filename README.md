# ISW Phoenix Go Middleware

A Go implementation of the ISW Phoenix middleware for handling client registration, key exchange, bill payments, and various financial operations.

## Prerequisites

- **Go 1.21+**: Go programming language
- **Git**: For version control

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd isw-phoenix-go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create the `isw` directory for storing keys:
```bash
mkdir isw
```

## Configuration

Copy the sample environment file and update it with your credentials:

```bash
cp config/env.sample config/.env
```

Update the `config/.env` file with your credentials:

```env
API_URL=https://dev.interswitch.io
BILLERS_API_URL=https://iswapigateway-develop.azurewebsites.net/
OWNER_PHONE_NUMBER=0777245670
CLIENT_ID=YOUR_CLIENT_ID
CLIENT_SECRET_KEY=YOUR_CLIENT_SECRET_KEY
TERMINAL_ID=YOUR_TERMINAL_ID
APP_VERSION=1
serialId=12345
PASSPHRASE=YOUR_PASSPHRASE
PORT=3000
```

## Running the Application

### Start the Middleware Server
```bash
go run main.go
```

### Run the Example Client
```bash
go run example_client.go
```

## API Endpoints

- `POST /validateCustomer` - Validate customer
- `GET /getbillerbycategory/:id` - Get billers by category
- `GET /getcategories` - Get categories
- `GET /billeritems/:id` - Get payment items
- `POST /payment` - Make payment
- `GET /transStatus/:id` - Get transaction status
- `GET /accountBalance` - Get account balance
- `GET /keyPair` - Generate RSA key pair
- `POST /clientRegistration` - Register client
- `POST /doKeyExchange` - Perform key exchange
- `GET /generateECDHKeyPair` - Generate ECDH key pair
- `GET /util` - Utility operations

## Features

- **Client Registration**: Register new clients with the system
- **Key Exchange**: Secure key exchange between clients and server
- **Retrieve Categories**: Fetch available service categories
- **Get Billers by Category**: List billers for a specific category
- **Retrieve Payment Items**: Fetch payment items for a biller
- **Make Payments**: Process payments securely
- **Check Account Balance**: Retrieve account balance
- **Transaction Status**: Check transaction status

## Project Structure

```
isw-phoenix-go/
├── main.go                 # Main application entry point
├── example_client.go       # Example client implementation
├── go.mod                  # Go module file
├── internal/
│   ├── isw/               # ISW client implementation
│   │   ├── client.go      # Main client logic
│   │   ├── crypto.go      # Cryptographic operations
│   │   ├── http.go        # HTTP client
│   │   └── types.go       # Type definitions
│   └── server/            # HTTP server handlers
│       └── handlers.go    # Request handlers
├── config/
│   └── env.sample         # Sample environment configuration
└── isw/                  # Key storage directory
    ├── private.txt
    ├── public.txt
    ├── authToken.txt
    └── sessionKey.txt
```

## Usage Example

### Client Registration
```go
registrationData := map[string]interface{}{
    "requestReference":  "unique-reference",
    "gprsCoordinate":    "1.234,5.678",
    "name":              "John Doe",
    "phoneNumber":       "1234567890",
    "nin":               "NIN123456789",
    "gender":            "M",
    "emailAddress":      "john@example.com",
    "ownerPhoneNumber":  "1234567890",
}

resp, err := client.RegisterClient(registrationData)
```

### Making a Payment
```go
paymentData := map[string]interface{}{
    "requestReference": "PAY123456",
    "amount":           "1000",
    "customerId":       "CUST123",
    "phoneNumber":      "1234567890",
    "paymentCode":      "PAY123",
    "customerName":     "John Doe",
    "sourceOfFunds":    "CASH",
    "narration":        "Bill payment",
    "depositorName":    "John Doe",
    "location":         "Lagos",
    "currencyCode":     "NGN",
}

// First validate customer
validationResp, err := client.ValidateCustomer(paymentData)

// Then make payment if validation succeeds
if validationResp.ResponseCode == "90000" {
    paymentResp, err := client.MakePayment(paymentData)
}
```

## Cryptographic Operations

The middleware supports various cryptographic operations:

- **RSA Key Generation**: Generate RSA key pairs for secure communication
- **ECDH Key Exchange**: Perform Elliptic Curve Diffie-Hellman key exchange
- **AES Encryption**: Encrypt sensitive data using AES-256-CBC
- **Digital Signatures**: Sign messages using RSA-SHA256
- **Password Hashing**: Hash passwords using SHA-512

## Error Handling

The middleware includes comprehensive error handling:

- HTTP request errors
- Cryptographic operation errors
- API response validation
- File I/O errors for key storage

## Security Features

- Secure key storage in local files
- Encrypted communication with ISW APIs
- Digital signatures for request authentication
- Session key management
- Password encryption

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.

## Support

For support and questions, please refer to the ISW Phoenix API documentation or contact the development team.
