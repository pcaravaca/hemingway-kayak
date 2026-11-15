package auth

import (
	"sync"
	"time"
)

// Revocation mantiene en memoria la lista de tokens invalidados por logout.
// Es intencionalmente simple: si el proceso se reinicia, los tokens siguen
// siendo válidos hasta su expiración natural.
type Revocation struct {
	mu     sync.Mutex
	tokens map[string]time.Time // token -> expiración
}

// NewRevocation crea una lista de revocación vacía.
func NewRevocation() *Revocation {
	return &Revocation{tokens: make(map[string]time.Time)}
}

// Revoke anota un token como inválido hasta que expire.
func (r *Revocation) Revoke(token string, expiresAt time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token] = expiresAt
	r.cleanupLocked()
}

// IsRevoked indica si el token fue invalidado y aún no expiró.
func (r *Revocation) IsRevoked(token string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	exp, ok := r.tokens[token]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(r.tokens, token)
		return false
	}
	return true
}

func (r *Revocation) cleanupLocked() {
	now := time.Now()
	for token, exp := range r.tokens {
		if now.After(exp) {
			delete(r.tokens, token)
		}
	}
}
