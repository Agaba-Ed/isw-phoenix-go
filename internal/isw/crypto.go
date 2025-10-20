package isw

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/url"
)

// GenerateRSAKeyPair generates RSA key pair
func GenerateRSAKeyPair() (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Encode private key
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return &KeyPair{
		PublicKey:  string(publicKeyPEM),
		PrivateKey: string(privateKeyPEM),
	}, nil
}

// GenerateECDHKeyPair generates ECDH key pair
func GenerateECDHKeyPair() (*ECDHKeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Get public key in compressed format
	publicKeyBytes := elliptic.MarshalCompressed(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)

	// Get private key as hex
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	return &ECDHKeyPair{
		PublicKey:  publicKeyBase64,
		PrivateKey: privateKeyHex,
	}, nil
}

// DoECDH performs ECDH key exchange
func DoECDH(serverPublicKeyBase64, ecdhPrivate string) (string, error) {
	// Decode server public key
	serverPublicKeyBytes, err := base64.StdEncoding.DecodeString(serverPublicKeyBase64)
	if err != nil {
		return "", err
	}

	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), serverPublicKeyBytes)
	if x == nil {
		return "", fmt.Errorf("invalid public key")
	}

	// Decode private key
	privateKeyBytes, err := hex.DecodeString(ecdhPrivate)
	if err != nil {
		return "", err
	}

	privateKey := new(big.Int).SetBytes(privateKeyBytes)

	// Perform ECDH and use only the X coordinate as the shared secret
	sharedX, _ := elliptic.P256().ScalarMult(x, y, privateKey.Bytes())
	sessionKeyHex := sharedX.Text(16) // hex string of X coordinate
	if len(sessionKeyHex) < 64 {
		sessionKeyHex = fmt.Sprintf("%064s", sessionKeyHex)
	}
	return sessionKeyHex, nil
}

// Helper function to parse private key from PEM string
func ParsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format if PKCS1 fails
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		privateKey = key.(*rsa.PrivateKey)
	}

	return privateKey, nil
}

func SignMessage(message, privateKeyPEM string) (string, error) {
	// Decode PEM block
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block: invalid PEM format")
	}

	// Try parsing as PKCS8 first (modern standard)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fallback to PKCS1 format
		rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("failed to parse private key: %w", err)
		}
		privateKey = rsaKey
	}

	// Type assert to RSA private key
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("key is not an RSA private key, got %T", privateKey)
	}

	// Hash the message with SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(message))
	hashed := hasher.Sum(nil)

	// Sign using PKCS1v15 with SHA-256
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %w", err)
	}

	// Encode to base64
	return base64.StdEncoding.EncodeToString(signature), nil
}

// DecryptWithPrivateKey decrypts data using RSA private key
func DecryptWithPrivateKey(encryptedData, privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("not an RSA private key")
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encryptedBytes, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// EncryptAuthToken encrypts auth token using session key
func EncryptAuthToken(authToken, sessionKey string) (string, error) {
	sessionKeyBytes, err := hex.DecodeString(sessionKey)
	if err != nil {
		return "", err
	}

	// Use first 32 bytes for AES-256
	if len(sessionKeyBytes) < 32 {
		return "", fmt.Errorf("session key too short")
	}
	key := sessionKeyBytes[:32]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, 16) // Zero IV
	mode := cipher.NewCBCEncrypter(block, iv)

	// Pad the data
	data := []byte(authToken)
	padding := aes.BlockSize - len(data)%aes.BlockSize
	paddedData := make([]byte, len(data)+padding)
	copy(paddedData, data)
	for i := len(data); i < len(paddedData); i++ {
		paddedData[i] = byte(padding)
	}

	encrypted := make([]byte, len(paddedData))
	mode.CryptBlocks(encrypted, paddedData)

	// Combine IV and encrypted data
	combined := make([]byte, 16+len(encrypted))
	copy(combined[:16], iv)
	copy(combined[16:], encrypted)

	return base64.StdEncoding.EncodeToString(combined), nil
}

// EncryptPassword encrypts password using session key
func EncryptPassword(password, sessionKey string) (string, error) {
	// Hash password with SHA-512 and get hex string
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashedPasswordHex := hex.EncodeToString(hasher.Sum(nil))

	hashedPasswordBytes, err := hex.DecodeString(hashedPasswordHex)
	if err != nil {
		return "", err
	}
	base64Hash := base64.StdEncoding.EncodeToString(hashedPasswordBytes)

	sessionKeyBytes, err := hex.DecodeString(sessionKey)
	if err != nil {
		return "", err
	}

	// Use first 32 bytes for AES-256
	if len(sessionKeyBytes) < 32 {
		return "", fmt.Errorf("session key too short")
	}
	key := sessionKeyBytes[:32]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, 16) // Zero IV
	mode := cipher.NewCBCEncrypter(block, iv)

	// Pad the data
	data := []byte(base64Hash)
	padding := aes.BlockSize - len(data)%aes.BlockSize
	paddedData := make([]byte, len(data)+padding)
	copy(paddedData, data)
	for i := len(data); i < len(paddedData); i++ {
		paddedData[i] = byte(padding)
	}

	encrypted := make([]byte, len(paddedData))
	mode.CryptBlocks(encrypted, paddedData)

	combined := make([]byte, 16+len(encrypted))
	copy(combined[:16], iv)
	copy(combined[16:], encrypted)

	return base64.StdEncoding.EncodeToString(combined), nil
}

// GenerateAuthorizationHeader generates authorization header
func GenerateAuthorizationHeader(clientID string) string {
	encodedClientID := base64.StdEncoding.EncodeToString([]byte(clientID))
	return fmt.Sprintf("InterswitchAuth %s", encodedClientID)
}

// GenerateSignatureHeader generates signature header using RSA private key
func GenerateSignatureHeader(httpMethod, resourceURL, timestamp, nonce, clientID, clientSecretKey, privateKey string, additionalParams *AdditionalParameters) (string, error) {

	encodedURI := url.QueryEscape(resourceURL)

	baseString := fmt.Sprintf("%s&%s&%s&%s&%s&%s",
		httpMethod,
		encodedURI,
		timestamp,
		nonce,
		clientID,
		clientSecretKey)

	if additionalParams != nil {
		baseString += fmt.Sprintf("&%s&%s&%s&%s&%s",
			additionalParams.Amount,
			additionalParams.TerminalID,
			additionalParams.RequestReference,
			additionalParams.CustomerID,
			additionalParams.PaymentCode)
	}

	return SignMessage(baseString, privateKey)
}
