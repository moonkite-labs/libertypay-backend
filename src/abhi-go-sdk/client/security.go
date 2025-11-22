package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// CredentialEncryption handles encryption and decryption of sensitive credentials
type CredentialEncryption struct {
	key []byte
}

// NewCredentialEncryption creates a new credential encryption instance with a password
func NewCredentialEncryption(password string) *CredentialEncryption {
	// Generate key from password using SHA-256
	hash := sha256.Sum256([]byte(password))
	return &CredentialEncryption{
		key: hash[:],
	}
}

// Encrypt encrypts plaintext using AES-GCM
func (ce *CredentialEncryption) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Create AES cipher
	block, err := aes.NewCipher(ce.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Use GCM mode for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return base64 encoded result
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func (ce *CredentialEncryption) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Decode from base64
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(ce.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Use GCM mode for authenticated decryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Check minimum length
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// SecureCredentials holds encrypted credentials
type SecureCredentials struct {
	EncryptedUsername string `json:"encrypted_username"`
	EncryptedPassword string `json:"encrypted_password"`
	Salt              string `json:"salt"`
}

// CredentialStore interface for different storage backends
type CredentialStore interface {
	Store(key string, credentials *SecureCredentials) error
	Retrieve(key string) (*SecureCredentials, error)
	Delete(key string) error
	Exists(key string) bool
}

// MemoryCredentialStore implements in-memory credential storage
type MemoryCredentialStore struct {
	store map[string]*SecureCredentials
}

// NewMemoryCredentialStore creates a new in-memory credential store
func NewMemoryCredentialStore() *MemoryCredentialStore {
	return &MemoryCredentialStore{
		store: make(map[string]*SecureCredentials),
	}
}

// Store stores credentials in memory
func (ms *MemoryCredentialStore) Store(key string, credentials *SecureCredentials) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if credentials == nil {
		return errors.New("credentials cannot be nil")
	}
	
	ms.store[key] = credentials
	return nil
}

// Retrieve retrieves credentials from memory
func (ms *MemoryCredentialStore) Retrieve(key string) (*SecureCredentials, error) {
	credentials, exists := ms.store[key]
	if !exists {
		return nil, fmt.Errorf("credentials not found for key: %s", key)
	}
	return credentials, nil
}

// Delete removes credentials from memory
func (ms *MemoryCredentialStore) Delete(key string) error {
	delete(ms.store, key)
	return nil
}

// Exists checks if credentials exist for the given key
func (ms *MemoryCredentialStore) Exists(key string) bool {
	_, exists := ms.store[key]
	return exists
}

// CredentialManager manages secure credential operations
type CredentialManager struct {
	encryption *CredentialEncryption
	store      CredentialStore
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager(encryptionPassword string, store CredentialStore) *CredentialManager {
	if store == nil {
		store = NewMemoryCredentialStore()
	}

	return &CredentialManager{
		encryption: NewCredentialEncryption(encryptionPassword),
		store:      store,
	}
}

// StoreCredentials encrypts and stores credentials
func (cm *CredentialManager) StoreCredentials(key, username, password string) error {
	encryptedUsername, err := cm.encryption.Encrypt(username)
	if err != nil {
		return fmt.Errorf("failed to encrypt username: %w", err)
	}

	encryptedPassword, err := cm.encryption.Encrypt(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Generate salt for additional security
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	credentials := &SecureCredentials{
		EncryptedUsername: encryptedUsername,
		EncryptedPassword: encryptedPassword,
		Salt:              base64.URLEncoding.EncodeToString(salt),
	}

	return cm.store.Store(key, credentials)
}

// RetrieveCredentials retrieves and decrypts credentials
func (cm *CredentialManager) RetrieveCredentials(key string) (username, password string, err error) {
	credentials, err := cm.store.Retrieve(key)
	if err != nil {
		return "", "", err
	}

	username, err = cm.encryption.Decrypt(credentials.EncryptedUsername)
	if err != nil {
		return "", "", fmt.Errorf("failed to decrypt username: %w", err)
	}

	password, err = cm.encryption.Decrypt(credentials.EncryptedPassword)
	if err != nil {
		return "", "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	return username, password, nil
}

// DeleteCredentials removes stored credentials
func (cm *CredentialManager) DeleteCredentials(key string) error {
	return cm.store.Delete(key)
}

// CredentialsExist checks if credentials exist for the given key
func (cm *CredentialManager) CredentialsExist(key string) bool {
	return cm.store.Exists(key)
}

// ClearCredentials securely clears credentials from memory
func (cm *CredentialManager) ClearCredentials() {
	// This would clear any in-memory credentials
	// Implementation depends on the store type
	if memStore, ok := cm.store.(*MemoryCredentialStore); ok {
		for key := range memStore.store {
			delete(memStore.store, key)
		}
	}
}