package web

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// authenticator implements the API's single-user auth: a password
// configured via `api.auth.password` is exchanged at the login endpoint
// for a stateless HMAC-signed bearer token.
type authenticator struct {
	password string
	secret   []byte
	ttl      time.Duration
}

func newAuthenticator(password string, ttl time.Duration) *authenticator {
	a := &authenticator{password: password, ttl: ttl}
	if password != "" {
		sum := sha256.Sum256([]byte("mineserver-api/auth:" + password))
		a.secret = sum[:]
	}
	return a
}

// enabled reports whether a password is configured. When not, every
// authenticated endpoint answers 503 instead of silently allowing access.
func (a *authenticator) enabled() bool { return len(a.secret) > 0 }

// issue creates a token valid until the returned expiry time.
func (a *authenticator) issue(now time.Time) (string, time.Time, error) {
	exp := now.Add(a.ttl)
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", time.Time{}, fmt.Errorf("generating token nonce: %w", err)
	}
	body := fmt.Sprintf("%d.%s", exp.Unix(), hex.EncodeToString(nonce))
	mac := hmac.New(sha256.New, a.secret)
	_, _ = mac.Write([]byte(body))
	token := base64.RawURLEncoding.EncodeToString(
		fmt.Appendf(nil, "%s.%s", body, hex.EncodeToString(mac.Sum(nil))),
	)
	return token, exp, nil
}

// verify reports whether token is well-formed, unexpired and authentic.
func (a *authenticator) verify(token string, now time.Time) bool {
	if !a.enabled() {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	parts := strings.Split(string(raw), ".")
	if len(parts) != 3 {
		return false
	}
	exp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || now.Unix() >= exp {
		return false
	}
	mac := hmac.New(sha256.New, a.secret)
	_, _ = mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := hex.EncodeToString(mac.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(expected), []byte(parts[2])) == 1
}

// middleware protects the API routes: it requires a valid bearer token,
// answering 401 for missing/invalid tokens and 503 when auth is not
// configured.
func (a *authenticator) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.enabled() {
			writeError(w, http.StatusServiceUnavailable, "api auth is not configured: set api.auth.password in the config file")
			return
		}
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || !a.verify(strings.TrimSpace(token), time.Now()) {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type loginRequest struct {
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// loginHandler exchanges the configured password for a bearer token.
func (a *authenticator) loginHandler(w http.ResponseWriter, r *http.Request) {
	if !a.enabled() {
		writeError(w, http.StatusServiceUnavailable, "api auth is not configured: set api.auth.password in the config file")
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(a.password)) != 1 {
		writeError(w, http.StatusUnauthorized, "invalid password")
		return
	}
	token, exp, err := a.issue(time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, loginResponse{Token: token, ExpiresAt: exp})
}
