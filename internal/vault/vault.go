// Package vault provides encrypted storage for OAuth tokens and API credentials.
// Uses AES-256-GCM for encryption at rest.
package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Vault manages encrypted credential storage.
type Vault struct {
	mu      sync.RWMutex
	dir     string
	key     []byte
	entries map[string]*Entry
}

// Entry represents a stored credential.
type Entry struct {
	Provider  string    `json:"provider"`
	Account   string    `json:"account"`
	Data      []byte    `json:"data"`       // encrypted payload
	Nonce     []byte    `json:"nonce"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// New creates a vault with the given directory and passphrase.
// The passphrase is stretched via SHA-256 to produce a 32-byte AES key.
func New(dir string, passphrase string) (*Vault, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("vault: cannot create directory: %w", err)
	}

	hash := sha256.Sum256([]byte(passphrase))
	v := &Vault{
		dir:     dir,
		key:     hash[:],
		entries: make(map[string]*Entry),
	}

	if err := v.load(); err != nil {
		return nil, err
	}

	return v, nil
}

// Store encrypts and persists a credential.
func (v *Vault) Store(provider, account string, data []byte, expiresAt time.Time) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	block, err := aes.NewCipher(v.key)
	if err != nil {
		return fmt.Errorf("vault: cipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("vault: gcm error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("vault: nonce error: %w", err)
	}

	encrypted := gcm.Seal(nil, nonce, data, nil)

	key := provider + ":" + account
	v.entries[key] = &Entry{
		Provider:  provider,
		Account:   account,
		Data:      encrypted,
		Nonce:     nonce,
		UpdatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	return v.persist()
}

// Retrieve decrypts and returns a stored credential.
func (v *Vault) Retrieve(provider, account string) ([]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	key := provider + ":" + account
	entry, ok := v.entries[key]
	if !ok {
		return nil, fmt.Errorf("vault: credential not found: %s", key)
	}

	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		return nil, fmt.Errorf("vault: credential expired: %s", key)
	}

	block, err := aes.NewCipher(v.key)
	if err != nil {
		return nil, fmt.Errorf("vault: cipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("vault: gcm error: %w", err)
	}

	plaintext, err := gcm.Open(nil, entry.Nonce, entry.Data, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: decrypt error: %w", err)
	}

	return plaintext, nil
}

// List returns all stored credential keys (provider:account).
func (v *Vault) List() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	keys := make([]string, 0, len(v.entries))
	for k := range v.entries {
		keys = append(keys, k)
	}
	return keys
}

// Delete removes a credential.
func (v *Vault) Delete(provider, account string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	key := provider + ":" + account
	delete(v.entries, key)
	return v.persist()
}

// PurgeExpired removes all expired entries.
func (v *Vault) PurgeExpired() int {
	v.mu.Lock()
	defer v.mu.Unlock()

	count := 0
	now := time.Now()
	for key, entry := range v.entries {
		if !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
			delete(v.entries, key)
			count++
		}
	}

	if count > 0 {
		v.persist()
	}
	return count
}

func (v *Vault) vaultFile() string {
	return filepath.Join(v.dir, "vault.enc")
}

func (v *Vault) persist() error {
	data, err := json.Marshal(v.entries)
	if err != nil {
		return fmt.Errorf("vault: marshal error: %w", err)
	}
	return os.WriteFile(v.vaultFile(), data, 0600)
}

func (v *Vault) load() error {
	path := v.vaultFile()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("vault: read error: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &v.entries)
}

// KeyFingerprint returns a short hex fingerprint of the vault key (for display).
func (v *Vault) KeyFingerprint() string {
	h := sha256.Sum256(v.key)
	return hex.EncodeToString(h[:4])
}
