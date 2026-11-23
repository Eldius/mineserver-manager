package web

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPassword = "s3cr3t-p4ss"

func TestTokenRoundtrip(t *testing.T) {
	a := newAuthenticator(testPassword, time.Hour)

	token, exp, err := a.issue(time.Now())
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(time.Hour), exp, 5*time.Second)
	assert.True(t, a.verify(token, time.Now()))
}

func TestVerifyRejectsBadTokens(t *testing.T) {
	a := newAuthenticator(testPassword, time.Hour)
	token, _, err := a.issue(time.Now())
	require.NoError(t, err)

	other := newAuthenticator("other-password", time.Hour)

	assert.False(t, a.verify("not-a-token", time.Now()), "garbage")
	assert.False(t, a.verify(token+"x", time.Now()), "tampered")
	assert.False(t, a.verify(token, time.Now().Add(2*time.Hour)), "expired")
	assert.False(t, other.verify(token, time.Now()), "signed by another password")

	raw, _ := base64.RawURLEncoding.DecodeString(token)
	parts := strings.Split(string(raw), ".")
	swapped := base64.RawURLEncoding.EncodeToString(fmt.Appendf(nil, "%d.%s.%s", time.Now().Add(48*time.Hour).Unix(), parts[1], parts[2]))
	assert.False(t, a.verify(swapped, time.Now()), "re-signed expiry without the secret")
}

func TestDisabledAuthenticatorRejectsEverything(t *testing.T) {
	a := newAuthenticator("", time.Hour)
	assert.False(t, a.enabled())
	assert.False(t, a.verify("whatever", time.Now()))
}

func TestLoginHandler(t *testing.T) {
	a := newAuthenticator(testPassword, time.Hour)

	t.Run("valid password returns usable token", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"password":"s3cr3t-p4ss"}`))
		a.loginHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"token"`)
	})

	t.Run("wrong password is unauthorized", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"password":"nope"}`))
		a.loginHandler(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("bad body is a 400", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{`))
		a.loginHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestLoginHandlerWithoutConfiguredPassword(t *testing.T) {
	a := newAuthenticator("", time.Hour)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"password":"anything"}`))
	a.loginHandler(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestMiddleware(t *testing.T) {
	a := newAuthenticator(testPassword, time.Hour)
	protected := a.middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	token, _, err := a.issue(time.Now())
	require.NoError(t, err)

	t.Run("no token is unauthorized", func(t *testing.T) {
		rec := httptest.NewRecorder()
		protected.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/instances", nil))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("valid token passes", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/instances", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		protected.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("disabled auth answers 503", func(t *testing.T) {
		open := newAuthenticator("", time.Hour).middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		rec := httptest.NewRecorder()
		open.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/instances", nil))
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})
}
