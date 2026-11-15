package middleware

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"hemingway-kayak-web/pkg/auth"
)

func newAuthedApp(secret []byte, revoked *auth.Revocation) *fiber.App {
	app := fiber.New()
	app.Use(JWTAuth(secret, revoked))
	app.Get("/whoami", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"userId": UserID(c)})
	})
	return app
}

func do(app *fiber.App, token string) (int, string) {
	req := httptest.NewRequest("GET", "/whoami", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(resp.Body)
	return resp.StatusCode, buf.String()
}

func TestJWTAuthOK(t *testing.T) {
	secret := []byte("secret")
	token, err := auth.Generate(secret, "user-9", "a@b.com", time.Hour)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	status, body := do(newAuthedApp(secret, nil), token)
	if status != 200 {
		t.Fatalf("status = %d, want 200 (body %s)", status, body)
	}
	if !strings.Contains(body, "user-9") {
		t.Errorf("body should carry the user id, got %s", body)
	}
}

func TestJWTAuthMissingHeader(t *testing.T) {
	if status, _ := do(newAuthedApp([]byte("s"), nil), ""); status != 401 {
		t.Fatalf("status = %d, want 401", status)
	}
}

func TestJWTAuthInvalidToken(t *testing.T) {
	if status, _ := do(newAuthedApp([]byte("s"), nil), "not-a-token"); status != 401 {
		t.Fatalf("status = %d, want 401", status)
	}
}

func TestJWTAuthWrongSecret(t *testing.T) {
	token, _ := auth.Generate([]byte("other"), "u", "e", time.Hour)
	if status, _ := do(newAuthedApp([]byte("s"), nil), token); status != 401 {
		t.Fatalf("status = %d, want 401", status)
	}
}

func TestJWTAuthRevoked(t *testing.T) {
	secret := []byte("secret")
	token, _ := auth.Generate(secret, "u", "e", time.Hour)

	revoked := auth.NewRevocation()
	claims, err := auth.Parse(secret, token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	revoked.Revoke(token, claims.ExpiresAt.Time)

	if status, _ := do(newAuthedApp(secret, revoked), token); status != 401 {
		t.Fatalf("status = %d, want 401 for a revoked token", status)
	}
}

func TestRequireOwnership(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userId", "owner-1")
		return c.Next()
	})
	app.Get("/users/:id", RequireOwnership("id"), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	if status, _ := do2(app, "/users/owner-1"); status != 200 {
		t.Fatalf("owner status = %d, want 200", status)
	}
	if status, _ := do2(app, "/users/someone-else"); status != 403 {
		t.Fatalf("other status = %d, want 403", status)
	}
}

func do2(app *fiber.App, path string) (int, string) {
	req := httptest.NewRequest("GET", path, nil)
	resp, err := app.Test(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(resp.Body)
	return resp.StatusCode, buf.String()
}
