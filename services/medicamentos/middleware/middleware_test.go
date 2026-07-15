package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_DentroDelLimite(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	if !rl.Allow("192.168.1.1") {
		t.Error("Primer request deberia ser permitido")
	}
	if !rl.Allow("192.168.1.1") {
		t.Error("Segundo request deberia ser permitido")
	}
	if !rl.Allow("192.168.1.1") {
		t.Error("Tercer request deberia ser permitido")
	}
}

func TestRateLimiter_FueraDelLimite(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	rl.Allow("10.0.0.1")
	rl.Allow("10.0.0.1")
	if rl.Allow("10.0.0.1") {
		t.Error("Cuarto request deberia ser bloqueado")
	}
}

func TestRateLimiter_IPsDiferentes(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	rl.Allow("1.1.1.1")
	if !rl.Allow("2.2.2.2") {
		t.Error("IPs diferentes no deberian compartir contador")
	}
}
