package middlewares

import (
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"

	"ticket-api/internal/constants"
	"ticket-api/internal/domain"
	"ticket-api/internal/responses"
	"ticket-api/pkg/firebase"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// FirebaseAuthMiddleware structure
type FirebaseAuthMiddleware struct {
	service     firebase.Service
	userService domain.UserService
}

// NewFirebaseAuthMiddleware creates new firebase authentication
func NewFirebaseAuthMiddleware(service firebase.Service, u domain.UserService) FirebaseAuthMiddleware {
	return FirebaseAuthMiddleware{
		service:     service,
		userService: u,
	}
}

// HandleAuthWithRole handles multiple roles
func (m FirebaseAuthMiddleware) HandleAuthWithRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.getTokenFromHeader(c)
		if err != nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if role != "" && token.Claims[role] == nil {
			responses.ErrorJSON(c, http.StatusUnauthorized, "auth-not-authorized-user")
			c.Abort()
			return
		}

		user, err := m.userService.GetUserByFirebaseId(token.UID)
		if err != nil {
			responses.ErrorJSON(
				c, http.StatusInternalServerError, "user not registered in the system",
			)
			c.Abort()
			return
		}

		c.Set(constants.Claims, token.Claims)
		c.Set(constants.FirebaseId, token.UID)
		// c.Set(constants.StripeId, *user.StripeID)
		c.Set(constants.UID, user.Id)

		// set email to the sentry message
		id := user.Id.String()
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.ConfigureScope(
				func(scope *sentry.Scope) {
					scope.SetUser(sentry.User{ID: id})
				},
			)
		}

		c.Next()
	}
}

// HandleAuthOptionally handles multiple roles
func (m FirebaseAuthMiddleware) HandleAuthOptionally() gin.HandlerFunc {
	return func(c *gin.Context) {
		noAuth := false
		token, err := m.getTokenFromHeader(c)

		if err != nil {
			if err.Error() == "ID token must be a non-empty string" {
				noAuth = true
			} else {
				responses.ErrorJSON(c, http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
		}

		if !noAuth {
			user, err := m.userService.GetUserByFirebaseId(token.UID)
			if err != nil {
				responses.ErrorJSON(
					c, http.StatusInternalServerError, "user not registered in the system",
				)
				c.Abort()
				return
			}

			c.Set(constants.Claims, token.Claims)
			c.Set(constants.FirebaseId, token.UID)
			c.Set(constants.StripeId, user.StripeID)
			c.Set(constants.UID, user.Id)

			id := user.Id.String()
			if hub := sentrygin.GetHubFromContext(c); hub != nil {
				hub.ConfigureScope(
					func(scope *sentry.Scope) {
						scope.SetUser(sentry.User{ID: id})
					},
				)
			}
		}

		c.Next()
	}
}

// getTokenFromHeader gets token from header
func (m FirebaseAuthMiddleware) getTokenFromHeader(c *gin.Context) (*auth.Token, error) {
	header := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	token, err := m.service.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
