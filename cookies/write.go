package cookies

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"

	"time"
)

// Write encodes the cookie value using base64 and writes it to the response.
func Write(c fiber.Ctx, cookie *fiber.Cookie) error {
	// Encode the cookie value using base64.
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	// Measure the exact serialized Set-Cookie value size, matching Fiber's
	// cookie normalization/serialization logic.
	if serializedCookieSize(cookie) > 4096 {
		return ErrValueTooLong
	}

	// Write the cookie as normal.
	c.Cookie(cookie)

	return nil
}

// serializedCookieSize returns the exact serialized Set-Cookie byte length
// using the same normalization and serialization rules Fiber v3 applies.
func serializedCookieSize(cookie *fiber.Cookie) int {
	cloned := *cookie

	if cloned.Path == "" {
		cloned.Path = "/"
	}

	if cloned.SessionOnly {
		cloned.MaxAge = 0
		cloned.Expires = time.Time{}
	}

	var sameSite http.SameSite
	switch {
	case strings.EqualFold(cloned.SameSite, fiber.CookieSameSiteStrictMode):
		sameSite = http.SameSiteStrictMode
	case strings.EqualFold(cloned.SameSite, fiber.CookieSameSiteNoneMode):
		sameSite = http.SameSiteNoneMode
		cloned.Secure = true
	case strings.EqualFold(cloned.SameSite, fiber.CookieSameSiteDisabled):
		sameSite = 0
	case strings.EqualFold(cloned.SameSite, fiber.CookieSameSiteLaxMode):
		sameSite = http.SameSiteLaxMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	if cloned.Partitioned {
		cloned.Secure = true
	}

	hc := &http.Cookie{
		Name:        cloned.Name,
		Value:       cloned.Value,
		Path:        cloned.Path,
		Domain:      cloned.Domain,
		Expires:     cloned.Expires,
		MaxAge:      cloned.MaxAge,
		Secure:      cloned.Secure,
		HttpOnly:    cloned.HTTPOnly,
		SameSite:    sameSite,
		Partitioned: cloned.Partitioned,
	}

	if err := hc.Valid(); err != nil {
		// Fiber ignores invalid cookies; size check should not block that behavior.
		return 0
	}

	fcookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(fcookie)

	fcookie.SetKey(hc.Name)
	fcookie.SetValue(hc.Value)
	fcookie.SetPath(hc.Path)
	fcookie.SetDomain(hc.Domain)

	if !cloned.SessionOnly {
		fcookie.SetMaxAge(hc.MaxAge)
		fcookie.SetExpire(hc.Expires)
	}

	fcookie.SetSecure(hc.Secure)
	fcookie.SetHTTPOnly(hc.HttpOnly)

	switch sameSite {
	case http.SameSiteLaxMode:
		fcookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	case http.SameSiteStrictMode:
		fcookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	case http.SameSiteNoneMode:
		fcookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
	case http.SameSiteDefaultMode:
		fcookie.SetSameSite(fasthttp.CookieSameSiteDefaultMode)
	default:
		fcookie.SetSameSite(fasthttp.CookieSameSiteDisabled)
	}

	fcookie.SetPartitioned(hc.Partitioned)

	return len(fcookie.Cookie())
}

