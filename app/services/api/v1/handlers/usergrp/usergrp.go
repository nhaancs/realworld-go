// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhaancs/bhms/foundation/sms"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nhaancs/bhms/app/services/api/v1/request"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user  *user.Core
	auth  *auth.Auth
	sms   *sms.SMS
	keyID string
}

// New constructs a handlers for route access.
func New(
	user *user.Core,
	auth *auth.Auth,
	keyID string,
	sms *sms.SMS,
) *Handlers {
	return &Handlers{
		user:  user,
		auth:  auth,
		keyID: keyID,
		sms:   sms,
	}
}

// CheckOTP verify user OTP.
func (h *Handlers) CheckOTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// verify user's OTP
	// updated user status to Active
	return nil
}

// Register adds a new user to the system.
// todo: do rate limit for this api to prevent sending to many sms
func (h *Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var dto RegisterDTO
	if err := web.Decode(r, &dto); err != nil {
		return request.NewError(err, http.StatusBadRequest)
	}

	e, err := toRegisterEntity(dto)
	if err != nil {
		return request.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Register(ctx, e)
	if err != nil {
		if errors.Is(err, user.ErrUniquePhone) {
			return request.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("register: usr[%+v]: %+v", usr, err)
	}

	if _, err = h.sms.SendOTP(ctx, sms.OTPInfo{Phone: usr.Phone}); err != nil {
		return fmt.Errorf("senotp: usr[%+v]: %+v", usr, err)
	}

	return web.Respond(ctx, w, toUserDTO(usr), http.StatusCreated)
}

// Token provides an API token for the authenticated user.
func (h *Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	phone, pass, ok := r.BasicAuth()
	if !ok {
		return auth.NewAuthError("must provide email and password in Basic auth")
	}

	usr, err := h.user.Authenticate(ctx, phone, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return request.NewError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure):
			return auth.NewAuthError(err.Error())
		default:
			return fmt.Errorf("authenticate: %w", err)
		}
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "bhms",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	token, err := h.auth.GenerateToken(ctx, h.keyID, claims)
	if err != nil {
		return fmt.Errorf("generatetoken: %w", err)
	}

	return web.Respond(ctx, w, toToken(token), http.StatusOK)
}
