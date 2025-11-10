package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:3000"

// Client represents the example client
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new example client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RegisterClient performs client registration
func (c *Client) RegisterClient() error {
	registrationData := map[string]interface{}{
		"requestReference": generateUUID(),
		"gprsCoordinate":   "1.234,5.678",
		"name":             "Test user",
		"phoneNumber":      "877889975",
		"nin":              "CM6567JNJF",
		"gender":           "M",
		"emailAddress":     "testusr.gmail.com",
		"ownerPhoneNumber": "5877889975",
	}

	resp, err := c.makeRequest("POST", "/clientRegistration", registrationData)
	if err != nil {
		return err
	}

	fmt.Println("Registration Response:", resp)
	return nil
}

// DoKeyExchange performs key exchange
func (c *Client) DoKeyExchange() error {
	resp, err := c.makeRequest("POST", "/doKeyExchange", nil)
	if err != nil {
		return err
	}

	fmt.Println("Key Exchange Response:", resp)
	return nil
}

// GetCategories retrieves categories
func (c *Client) GetCategories() error {
	resp, err := c.makeRequest("GET", "/getcategories", nil)
	if err != nil {
		return err
	}

	fmt.Println("Categories:", resp)
	return nil
}

// GetBillersByCategory retrieves billers by category
func (c *Client) GetBillersByCategory(categoryID string) error {
	url := fmt.Sprintf("/getbillerbycategory/%s", categoryID)
	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return err
	}

	fmt.Println("Billers:", resp)
	return nil
}

// GetPaymentItems retrieves payment items
func (c *Client) GetPaymentItems(billerID string) error {
	url := fmt.Sprintf("/billeritems/%s", billerID)
	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return err
	}

	fmt.Println("Payment Items:", resp)
	return nil
}

// MakePayment makes a payment
func (c *Client) MakePayment() error {
	paymentData := map[string]interface{}{
		"requestReference": fmt.Sprintf("%d", time.Now().Unix()),
		"amount":           "1000",
		"customerId":       "0774528787",
		"phoneNumber":      "0774528787",
		"paymentCode":      "54048546968",
		"customerName":     "John Doe",
		"sourceOfFunds":    "CASH",
		"narration":        "Bill payment",
		"depositorName":    "John Doe",
		"location":         "Lagos",
		"currencyCode":     "NGN",
	}

	// First validate customer
	validationResp, err := c.makeRequest("POST", "/validateCustomer", paymentData)
	if err != nil {
		return err
	}

	fmt.Println("Validation Response:", validationResp)

	// Check if validation was successful
	if respMap, ok := validationResp.(map[string]interface{}); ok {
		if responseCode, exists := respMap["responseCode"]; exists && responseCode == "90000" {
			// Then make payment
			paymentResp, err := c.makeRequest("POST", "/payment", paymentData)
			if err != nil {
				return err
			}

			fmt.Println("Payment Response:", paymentResp)
		}
	}

	return nil
}

// AccountBalance retrieves account balance
func (c *Client) AccountBalance() error {
	resp, err := c.makeRequest("GET", "/accountBalance", nil)
	if err != nil {
		return err
	}

	fmt.Println("Account Balance:", resp)
	return nil
}

// makeRequest makes HTTP request
func (c *Client) makeRequest(method, endpoint string, data interface{}) (interface{}, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return string(responseBody), nil
	}

	return result, nil
}

// generateUUID generates a simple UUID
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func main() {
	client := NewClient()

	//fmt.Println("Step 1: Client Registration")
	//if err := client.RegisterClient(); err != nil {
	//	fmt.Printf("Registration failed: %v\n", err)
	//}

	// Uncomment other steps as needed
	//fmt.Println("\nStep 2: Key Exchange")
	//client.DoKeyExchange()

	//fmt.Println("\nStep 3: Get Categories")
	//client.GetCategories()

	// fmt.Println("\nStep 4: Get Billers")
	//client.GetBillersByCategory("60046")

	//fmt.Println("\nStep 5: Get Payment Items")
	//client.GetPaymentItems("110061")

	//fmt.Println("\nStep 6: Make Payment")
	//client.MakePayment()

	fmt.Println("\nStep 7: Get Account Balance")
	client.AccountBalance()
}
