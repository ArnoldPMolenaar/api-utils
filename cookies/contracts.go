package cookies

import (
	"fmt"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"
)

const defaultMaxSerializedCookieSize = 4096

// WriteConfig defines behavior for serialized cookie size enforcement.
type WriteConfig struct {
	MaxSerializedSize int
}

// DefaultWriteConfig returns default write constraints.
func DefaultWriteConfig() WriteConfig {
	return WriteConfig{MaxSerializedSize: defaultMaxSerializedCookieSize}
}

func (cfg WriteConfig) maxSerializedSize() int {
	if cfg.MaxSerializedSize <= 0 {
		return defaultMaxSerializedCookieSize
	}

	return cfg.MaxSerializedSize
}

// CookiePolicy defines reusable cookie attributes for API-wide defaults.
type CookiePolicy struct {
	Path        string
	Domain      string
	Expires     time.Time
	MaxAge      int
	Secure      bool
	HTTPOnly    bool
	SameSite    string
	SessionOnly bool
	Partitioned bool
}

// DefaultCookiePolicy returns a secure default cookie policy.
func DefaultCookiePolicy() CookiePolicy {
	return CookiePolicy{
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	}
}

// NewCookie builds a Fiber cookie using a shared policy.
func NewCookie(name, value string, policy CookiePolicy) *fiber.Cookie {
	if policy.Path == "" {
		policy.Path = "/"
	}

	return &fiber.Cookie{
		Name:        name,
		Value:       value,
		Path:        policy.Path,
		Domain:      policy.Domain,
		Expires:     policy.Expires,
		MaxAge:      policy.MaxAge,
		Secure:      policy.Secure,
		HTTPOnly:    policy.HTTPOnly,
		SameSite:    policy.SameSite,
		SessionOnly: policy.SessionOnly,
		Partitioned: policy.Partitioned,
	}
}

// SignedConfig defines HMAC configuration for signed cookies.
type SignedConfig struct {
	SecretKey []byte
	Write     WriteConfig
}

func (cfg SignedConfig) writeConfig() WriteConfig {
	if cfg.Write.MaxSerializedSize <= 0 {
		return DefaultWriteConfig()
	}

	return cfg.Write
}

// EncryptedConfig defines encryption configuration for encrypted cookies.
type EncryptedConfig struct {
	SecretKey string
	Write     WriteConfig
}

func (cfg EncryptedConfig) writeConfig() WriteConfig {
	if cfg.Write.MaxSerializedSize <= 0 {
		return DefaultWriteConfig()
	}

	return cfg.Write
}

// CryptoDeps defines encryption/decryption dependencies for cookie helpers.
type CryptoDeps struct {
	Encrypt func(secretKey, plaintext string) (string, error)
	Decrypt func(secretKey, ciphertext string) (string, error)
}

// DefaultCryptoDeps returns production crypto implementations.
func DefaultCryptoDeps() CryptoDeps {
	return CryptoDeps{
		Encrypt: utils.Encrypt,
		Decrypt: utils.Decrypt,
	}
}

func (d CryptoDeps) validateEncrypt() error {
	if d.Encrypt == nil {
		return fmt.Errorf("%w: Encrypt dependency cannot be nil", ErrInvalidConfiguration)
	}

	return nil
}

func (d CryptoDeps) validateDecrypt() error {
	if d.Decrypt == nil {
		return fmt.Errorf("%w: Decrypt dependency cannot be nil", ErrInvalidConfiguration)
	}

	return nil
}


