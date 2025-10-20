package isw

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"isw-phoenix-go/internal/logger"
	"path/filepath"
)

// Client represents the ISW Phoenix client
type Client struct {
	config     *Configuration
	httpClient *HTTPClient
	privateKey string
	publicKey  string
	authToken  string
	sessionKey string
}

// NewClient creates a new ISW client
func NewClient(config *Configuration) *Client {
	client := &Client{
		config: config,
	}
	client.initializeKeys()
	// Create HTTP client with private key after keys are loaded
	client.httpClient = NewHTTPClient(config, client.privateKey)
	return client
}

// initializeKeys loads saved key values
func (c *Client) initializeKeys() error {
	var err error

	c.privateKey, err = c.readFromFile("private.txt")
	if err != nil {
		return err
	}

	c.publicKey, err = c.readFromFile("public.txt")
	if err != nil {
		return err
	}

	c.authToken, err = c.readFromFile("authToken.txt")
	if err != nil {
		return err
	}

	c.sessionKey, err = c.readFromFile("sessionKey.txt")
	if err != nil {
		return err
	}

	return nil
}

// writeToFile writes data to a file in the isw directory
func (c *Client) writeToFile(filename string, data string) error {
	filePath := filepath.Join("isw", filename)
	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

// readFromFile reads data from a file in the isw directory
func (c *Client) readFromFile(filename string) (string, error) {
	filePath := filepath.Join("isw", filename)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GenerateRSAKeyPair generates RSA key pair and saves to files
func (c *Client) GenerateRSAKeyPair() (*KeyPair, error) {
	keyPair, err := GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	// Save keys to files
	c.writeToFile("public.txt", keyPair.PublicKey)
	c.writeToFile("private.txt", keyPair.PrivateKey)

	// Update client keys
	c.publicKey = keyPair.PublicKey
	c.privateKey = keyPair.PrivateKey

	return keyPair, nil
}

// GenerateECDHKeyPair generates ECDH key pair
func (c *Client) GenerateECDHKeyPair() (*ECDHKeyPair, error) {
	return GenerateECDHKeyPair()
}

// ClientRegistration performs client registration
func (c *Client) ClientRegistration(data *ClientRegistrationRequest) (*APIResponse, error) {
	ecdhKeys, err := c.GenerateECDHKeyPair()
	if err != nil {
		return nil, err
	}

	payload := &ClientRegistrationPayload{
		TerminalID:             c.config.TerminalID,
		AppVersion:             c.config.AppVersion,
		SerialID:               c.config.SerialID,
		RequestReference:       data.RequestReference,
		GPRSCoordinate:         data.GPRSCoordinate,
		Name:                   data.Name,
		PhoneNumber:            data.PhoneNumber,
		NIN:                    data.NIN,
		Gender:                 data.Gender,
		EmailAddress:           data.EmailAddress,
		OwnerPhoneNumber:       data.OwnerPhoneNumber,
		PublicKey:              c.publicKey,
		ClientSessionPublicKey: ecdhKeys.PublicKey,
	}

	resp, err := c.httpClient.Post("/api/v1/phoenix/client/clientRegistration", payload, c.authToken, c.sessionKey, nil)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode == "90000" {
		// Handle successful registration
		responseData := resp.Response.(map[string]interface{})

		serverSessionPublicKey := responseData["serverSessionPublicKey"].(string)
		transactionReference := responseData["transactionReference"].(string)

		decryptedServerSessionPublicKey, err := DecryptWithPrivateKey(serverSessionPublicKey, c.privateKey)
		if err != nil {
			return nil, err
		}

		sessionKey, err := DoECDH(decryptedServerSessionPublicKey, ecdhKeys.PrivateKey)
		if err != nil {
			return nil, err
		}

		// Complete registration
		completeResp, err := c.completeClientRegistration(sessionKey, transactionReference, data.RequestReference)
		if err != nil {
			return nil, err
		}

		if completeResp.ResponseCode == "90000" {
			completeData := completeResp.Response.(map[string]interface{})
			clientSecret := completeData["clientSecret"].(string)
			newAuthToken := completeData["authToken"].(string)

			newDecryptedAuthToken, err := DecryptWithPrivateKey(newAuthToken, c.privateKey)
			if err != nil {
				return nil, err
			}

			decryptedClientSecret, err := DecryptWithPrivateKey(clientSecret, c.privateKey)
			if err != nil {
				return nil, err
			}

			// Save to files
			c.writeToFile("authToken.txt", newDecryptedAuthToken)
			c.writeToFile("secrete.txt", decryptedClientSecret)
			c.writeToFile("sessionKey.txt", sessionKey)

			// Update client state
			c.authToken = newDecryptedAuthToken
			c.sessionKey = sessionKey

			completeResp.Response = map[string]interface{}{
				"response":         completeResp.Response,
				"DECRYPTED_SECRET": decryptedClientSecret,
			}
		}

		return completeResp, nil
	}

	return resp, nil
}

// completeClientRegistration completes client registration
func (c *Client) completeClientRegistration(sessionKey, transactionReference, requestReference string) (*APIResponse, error) {
	encryptedPassword, err := EncryptPassword(c.config.Passphrase, sessionKey)
	if err != nil {
		return nil, err
	}

	payload := &CompleteRegistrationPayload{
		TerminalID:           c.config.TerminalID,
		Password:             encryptedPassword,
		TransactionReference: transactionReference,
		RequestReference:     requestReference,
		SerialID:             c.config.SerialID,
	}

	return c.httpClient.Post("/api/v1/phoenix/client/completeClientRegistration", payload, c.authToken, c.sessionKey, nil)
}

// DoKeyExchange performs key exchange
func (c *Client) DoKeyExchange() (*APIResponse, error) {
	// Log the passphrase sourced from environment/config for visibility during key exchange
	logger.Log.WithField("passphrase", c.config.Passphrase).Info("Starting key exchange with env passphrase")

	ecdhKeys, err := c.GenerateECDHKeyPair()
	if err != nil {
		return nil, err
	}

	requestReference := generateRequestReference()

	// Hash password (Node parity): base64(SHA512(passphrase))
	passwordHash := hashPassword(c.config.Passphrase)

	passwordCipher := passwordHash + requestReference + c.config.SerialID

	password, err := SignMessage(passwordCipher, c.privateKey)
	if err != nil {
		return nil, err
	}

	payload := &KeyExchangePayload{
		TerminalID:             c.config.TerminalID,
		AppVersion:             c.config.AppVersion,
		SerialID:               c.config.SerialID,
		RequestReference:       requestReference,
		ClientSessionPublicKey: ecdhKeys.PublicKey,
		Password:               password,
	}

	resp, err := c.httpClient.Post("/api/v1/phoenix/client/doKeyExchange", payload, c.authToken, c.sessionKey, nil)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode == "90000" {
		responseData := resp.Response.(map[string]interface{})

		serverSessionPublicKey := responseData["serverSessionPublicKey"].(string)
		authToken := responseData["authToken"].(string)

		decryptedAuthToken, err := DecryptWithPrivateKey(authToken, c.privateKey)
		if err != nil {
			return nil, err
		}

		decryptedServerSessionPublicKey, err := DecryptWithPrivateKey(serverSessionPublicKey, c.privateKey)
		if err != nil {
			return nil, err
		}

		sessionKey, err := DoECDH(decryptedServerSessionPublicKey, ecdhKeys.PrivateKey)
		if err != nil {
			return nil, err
		}

		// Save to files
		c.writeToFile("authToken.txt", decryptedAuthToken)
		c.writeToFile("sessionKey.txt", sessionKey)

		// Update client state
		c.authToken = decryptedAuthToken
		c.sessionKey = sessionKey
	}

	return resp, nil
}

// GetCategories retrieves biller categories
func (c *Client) GetCategories() (*APIResponse, error) {
	url := fmt.Sprintf("qt-api/Biller/categories-by-client/%s/%s", c.config.TerminalID, c.config.TerminalID)
	return c.httpClient.GetV2(url, nil, c.authToken, c.sessionKey)
}

// GetCategoryBillers retrieves billers by category
func (c *Client) GetCategoryBillers(categoryID string) (*APIResponse, error) {
	url := fmt.Sprintf("qt-api/Biller/biller-by-category/%s", categoryID)
	return c.httpClient.GetV2(url, nil, c.authToken, c.sessionKey)
}

// GetPaymentItems retrieves payment items for a biller
func (c *Client) GetPaymentItems(billerID string) (*APIResponse, error) {
	url := fmt.Sprintf("qt-api/Biller/items/biller-id/%s", billerID)
	return c.httpClient.GetV2(url, nil, c.authToken, c.sessionKey)
}

// AccountBalance retrieves account balance
func (c *Client) AccountBalance() (*APIResponse, error) {
	requestReference := generateRequestReference()
	url := fmt.Sprintf("/api/v1/phoenix/sente/accountBalance?terminalId=%s&requestReference=%s", c.config.TerminalID, requestReference)
	return c.httpClient.Get(url, nil, c.authToken, c.sessionKey)
}

// TransactionInquiry retrieves transaction status
func (c *Client) TransactionInquiry(requestReference string) (*APIResponse, error) {
	url := fmt.Sprintf("/api/v1/phoenix/sente/status?terminalId=%s&requestReference=%s", c.config.TerminalID, requestReference)
	return c.httpClient.Get(url, nil, c.authToken, c.sessionKey)
}

// ValidateCustomer validates customer
func (c *Client) ValidateCustomer(data *CustomerValidationRequest) (*APIResponse, error) {
	payload := &CustomerValidationPayload{
		TerminalID:          c.config.TerminalID,
		RequestReference:    data.RequestReference,
		PaymentCode:         data.PaymentCode,
		CustomerID:          data.CustomerID,
		CurrencyCode:        data.CurrencyCode,
		Amount:              data.Amount,
		AlternateCustomerID: data.AlternateCustomerID,
		CustomerToken:       data.CustomerToken,
		TransactionCode:     data.TransactionCode,
		AdditionalData:      data.AdditionalData,
	}

	return c.httpClient.Post("/api/v1/phoenix/sente/customerValidation", payload, c.authToken, c.sessionKey, nil)
}

// MakePayment makes a payment
func (c *Client) MakePayment(data *PaymentRequest) (*APIResponse, error) {
	payload := &PaymentPayload{
		TerminalID:               c.config.TerminalID,
		RequestReference:         data.RequestReference,
		Amount:                   data.Amount,
		CustomerID:               data.CustomerID,
		PhoneNumber:              data.PhoneNumber,
		PaymentCode:              data.PaymentCode,
		CustomerName:             data.CustomerName,
		SourceOfFunds:            data.SourceOfFunds,
		Narration:                data.Narration,
		DepositorName:            data.DepositorName,
		Location:                 data.Location,
		AlternateCustomerID:      data.AlternateCustomerID,
		TransactionCode:          data.TransactionCode,
		CustomerToken:            data.CustomerToken,
		AdditionalData:           data.AdditionalData,
		CollectionsAccountNumber: data.CollectionsAccountNumber,
		PIN:                      data.PIN,
		OTP:                      data.OTP,
		CurrencyCode:             data.CurrencyCode,
	}

	additionalParams := &AdditionalParameters{
		Amount:           data.Amount,
		TerminalID:       c.config.TerminalID,
		RequestReference: data.RequestReference,
		CustomerID:       data.CustomerID,
		PaymentCode:      data.PaymentCode,
	}

	return c.httpClient.Post("/api/v1/phoenix/sente/xpayment", payload, c.authToken, c.sessionKey, additionalParams)
}

// DecryptWithPrivateKey decrypts data using private key
func (c *Client) DecryptWithPrivateKey(encryptedData string) (string, error) {
	return DecryptWithPrivateKey(encryptedData, c.privateKey)
}

// Helper functions
func generateRequestReference() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
func hashPassword(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashed := hasher.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashed)
}
