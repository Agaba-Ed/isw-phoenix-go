package isw

// Configuration holds all environment variables
type Configuration struct {
	APIURL           string
	BillersAPIURL    string
	ClientID         string
	OwnerPhoneNumber string
	Passphrase       string
	ClientSecretKey  string
	TerminalID       string
	AppVersion       string
	SerialID         string
}

// KeyPair represents RSA key pair
type KeyPair struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// ECDHKeyPair represents ECDH key pair
type ECDHKeyPair struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// ClientRegistrationRequest represents client registration payload
type ClientRegistrationRequest struct {
	RequestReference string `json:"requestReference"`
	GPRSCoordinate   string `json:"gprsCoordinate"`
	Name             string `json:"name"`
	PhoneNumber      string `json:"phoneNumber"`
	NIN              string `json:"nin"`
	Gender           string `json:"gender"`
	EmailAddress     string `json:"emailAddress"`
	OwnerPhoneNumber string `json:"ownerPhoneNumber"`
}

// ClientRegistrationPayload represents the full registration payload
type ClientRegistrationPayload struct {
	TerminalID             string `json:"terminalId"`
	AppVersion             string `json:"appVersion"`
	SerialID               string `json:"serialId"`
	RequestReference       string `json:"requestReference"`
	GPRSCoordinate         string `json:"gprsCoordinate"`
	Name                   string `json:"name"`
	PhoneNumber            string `json:"phoneNumber"`
	NIN                    string `json:"nin"`
	Gender                 string `json:"gender"`
	EmailAddress           string `json:"emailAddress"`
	OwnerPhoneNumber       string `json:"ownerPhoneNumber"`
	PublicKey              string `json:"publicKey"`
	ClientSessionPublicKey string `json:"clientSessionPublicKey"`
}

// CompleteRegistrationPayload represents completion payload
type CompleteRegistrationPayload struct {
	TerminalID           string `json:"terminalId"`
	Password             string `json:"password"`
	TransactionReference string `json:"transactionReference"`
	RequestReference     string `json:"requestReference"`
	SerialID             string `json:"serialId"`
}

// KeyExchangePayload represents key exchange payload
type KeyExchangePayload struct {
	TerminalID             string `json:"terminalId"`
	AppVersion             string `json:"appVersion"`
	SerialID               string `json:"serialId"`
	RequestReference       string `json:"requestReference"`
	ClientSessionPublicKey string `json:"clientSessionPublicKey"`
	Password               string `json:"password"`
}

// CustomerValidationRequest represents customer validation payload
type CustomerValidationRequest struct {
	RequestReference    string `json:"requestReference"`
	PaymentCode         string `json:"paymentCode"`
	CustomerID          string `json:"customerId"`
	CurrencyCode        string `json:"currencyCode"`
	Amount              string `json:"amount"`
	AlternateCustomerID string `json:"alternateCustomerId,omitempty"`
	CustomerToken       string `json:"customerToken,omitempty"`
	TransactionCode     string `json:"transactionCode,omitempty"`
	AdditionalData      string `json:"additionalData,omitempty"`
}

// CustomerValidationPayload represents the full validation payload
type CustomerValidationPayload struct {
	TerminalID          string `json:"terminalId"`
	RequestReference    string `json:"requestReference"`
	PaymentCode         string `json:"paymentCode"`
	CustomerID          string `json:"customerId"`
	CurrencyCode        string `json:"currencyCode"`
	Amount              string `json:"amount"`
	AlternateCustomerID string `json:"alternateCustomerId,omitempty"`
	CustomerToken       string `json:"customerToken,omitempty"`
	TransactionCode     string `json:"transactionCode,omitempty"`
	AdditionalData      string `json:"additionalData,omitempty"`
}

// PaymentRequest represents payment request
type PaymentRequest struct {
	RequestReference         string `json:"requestReference"`
	Amount                   string `json:"amount"`
	CustomerID               string `json:"customerId"`
	PhoneNumber              string `json:"phoneNumber"`
	PaymentCode              string `json:"paymentCode"`
	CustomerName             string `json:"customerName"`
	SourceOfFunds            string `json:"sourceOfFunds"`
	Narration                string `json:"narration"`
	DepositorName            string `json:"depositorName"`
	Location                 string `json:"location"`
	AlternateCustomerID      string `json:"alternateCustomerId,omitempty"`
	TransactionCode          string `json:"transactionCode,omitempty"`
	CustomerToken            string `json:"customerToken,omitempty"`
	AdditionalData           string `json:"additionalData,omitempty"`
	CollectionsAccountNumber string `json:"collectionsAccountNumber,omitempty"`
	PIN                      string `json:"pin,omitempty"`
	OTP                      string `json:"otp,omitempty"`
	CurrencyCode             string `json:"currencyCode"`
}

// PaymentPayload represents the full payment payload
type PaymentPayload struct {
	TerminalID               string `json:"terminalId"`
	RequestReference         string `json:"requestReference"`
	Amount                   string `json:"amount"`
	CustomerID               string `json:"customerId"`
	PhoneNumber              string `json:"phoneNumber"`
	PaymentCode              string `json:"paymentCode"`
	CustomerName             string `json:"customerName"`
	SourceOfFunds            string `json:"sourceOfFunds"`
	Narration                string `json:"narration"`
	DepositorName            string `json:"depositorName"`
	Location                 string `json:"location"`
	AlternateCustomerID      string `json:"alternateCustomerId,omitempty"`
	TransactionCode          string `json:"transactionCode,omitempty"`
	CustomerToken            string `json:"customerToken,omitempty"`
	AdditionalData           string `json:"additionalData,omitempty"`
	CollectionsAccountNumber string `json:"collectionsAccountNumber,omitempty"`
	PIN                      string `json:"pin,omitempty"`
	OTP                      string `json:"otp,omitempty"`
	CurrencyCode             string `json:"currencyCode"`
}

// APIResponse represents the standard API response structure
type APIResponse struct {
	ResponseCode     string      `json:"responseCode"`
	ResponseMessage  string      `json:"responseMessage,omitempty"`
	RequestReference string      `json:"requestReference,omitempty"`
	Response         interface{} `json:"response"`
	Success          bool        `json:"success,omitempty"`
	Error            string      `json:"error,omitempty"`
}

// HTTPHeaders represents HTTP headers for requests
type HTTPHeaders struct {
	Authorization string
	Signature     string
	ContentType   string
	Nonce         string
	Timestamp     string
	AuthToken     string
}

// AdditionalParameters represents additional parameters for signature generation
type AdditionalParameters struct {
	Amount           string
	TerminalID       string
	RequestReference string
	CustomerID       string
	PaymentCode      string
}
