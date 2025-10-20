package server

import (
	"net/http"
	"time"

	"isw-phoenix-go/internal/isw"
	"isw-phoenix-go/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	iswClient *isw.Client
}

// NewHandlers creates new handlers
func NewHandlers(iswClient *isw.Client) *Handlers {
	return &Handlers{
		iswClient: iswClient,
	}
}

// ValidateCustomer handles customer validation
func (h *Handlers) ValidateCustomer(c *gin.Context) {
	var req isw.CustomerValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.iswClient.ValidateCustomer(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetBillerCategories handles getting billers by category
func (h *Handlers) GetBillerCategories(c *gin.Context) {
	categoryID := c.Param("id")

	resp, err := h.iswClient.GetCategoryBillers(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCategories handles getting categories
func (h *Handlers) GetCategories(c *gin.Context) {
	resp, err := h.iswClient.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPaymentItems handles getting payment items
func (h *Handlers) GetPaymentItems(c *gin.Context) {
	billerID := c.Param("id")

	resp, err := h.iswClient.GetPaymentItems(billerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Payment handles payment processing
func (h *Handlers) Payment(c *gin.Context) {
	var req isw.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.iswClient.MakePayment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// TransStatus handles transaction status inquiry
func (h *Handlers) TransStatus(c *gin.Context) {
	requestReference := c.Param("id")

	resp, err := h.iswClient.TransactionInquiry(requestReference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AccountBalance handles account balance inquiry
func (h *Handlers) AccountBalance(c *gin.Context) {
	resp, err := h.iswClient.AccountBalance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GenerateRSAKeyPair handles RSA key pair generation
func (h *Handlers) GenerateRSAKeyPair(c *gin.Context) {
	keyPair, err := h.iswClient.GenerateRSAKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, keyPair)
}

// ClientRegistration handles client registration
func (h *Handlers) ClientRegistration(c *gin.Context) {
	startTime := time.Now()

	var req isw.ClientRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"endpoint": "/clientRegistration",
			"error":    err.Error(),
			"duration": time.Since(startTime),
		}).Error("Failed to bind JSON request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"endpoint": "/clientRegistration",
		"request":  req,
	}).Info("Processing client registration request")

	resp, err := h.iswClient.ClientRegistration(&req)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"endpoint": "/clientRegistration",
			"error":    err.Error(),
			"duration": time.Since(startTime),
		}).Error("Client registration failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"endpoint": "/clientRegistration",
		"response": resp,
		"duration": time.Since(startTime),
	}).Info("Client registration completed")

	c.JSON(http.StatusOK, resp)
}

// DoKeyExchange handles key exchange
func (h *Handlers) DoKeyExchange(c *gin.Context) {
	startTime := time.Now()

	logger.Log.WithFields(logrus.Fields{
		"endpoint": "/doKeyExchange",
	}).Info("Processing key exchange request")

	resp, err := h.iswClient.DoKeyExchange()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"endpoint": "/doKeyExchange",
			"error":    err.Error(),
			"duration": time.Since(startTime),
		}).Error("Key exchange failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"endpoint": "/doKeyExchange",
		"response": resp,
		"duration": time.Since(startTime),
	}).Info("Key exchange completed")

	c.JSON(http.StatusOK, resp)
}

// GenerateECDHKeyPair handles ECDH key pair generation
func (h *Handlers) GenerateECDHKeyPair(c *gin.Context) {
	keyPair, err := h.iswClient.GenerateECDHKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, keyPair)
}

// Util handles utility operations (decryption)
func (h *Handlers) Util(c *gin.Context) {
	clientSecret := "zpRc8p/dlBKkizlvenLUh/evRn+2lW1FfCyRGx/SOAUMcxrkyzsNibQgAqv2seuGmObIE1e8IBzpiN+urx9b4mHv7Dm4GSjYNmLfMUfJom2q4ehDg7K6u/2Pm2hAu6yW3iH/UmQwbIElx4W06moXWugWQX0gvxXmv4UwCAwxYZBtTIIv8WQaP9j9PEQPHFmSqFXy2qJwij+ytpp5eHJ69WQs2Q0Vj0F3CkbYd1fvrNRlTCgXqJy1v9dw5lOpyhYL7m76P1v8Xnz9C1CkCPmlQ0XXNa0wt85kbpmg66kTtcOpLNP/qZHjHTJy1I+BpdCEULV203KqjiirSmNlhJiCxQ=="

	decrypted, err := h.iswClient.DecryptWithPrivateKey(clientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"decrypted": decrypted})
}
