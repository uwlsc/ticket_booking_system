package seeds

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"ticket-api/config"
	"ticket-api/internal/constants"
	"ticket-api/internal/domain"
	"ticket-api/pkg"
	"ticket-api/pkg/firebase"
)

type AdminSeed struct {
	logger      pkg.Logger
	config      config.Config
	service     firebase.Service
	userService domain.UserService
}

func NewAdminSeed(l pkg.Logger, c config.Config, s firebase.Service, lc fx.Lifecycle, u domain.UserService) AdminSeed {
	admin := AdminSeed{l, c, s, u}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return admin.Setup()
			},
		},
	)

	return admin
}

func (a AdminSeed) Setup() error {
	// Create email manually in firebase
	email := a.config.AdminEmail
	password := a.config.AdminPassword

	a.logger.Info("🌱 seeding  admin data...")

	if email == "" || password == "" {
		a.logger.Error(
			"Got empty admin email and password from environment variables. " +
				"Admin seed not executed.",
		)
		return errors.New("no email and password specified")
	}

	_, err := a.service.RetrieveUserByEmail(email)

	if err != nil {
		adminClaim := gin.H{
			constants.RoleIsAdmin: true,
		}

		user, err := a.service.CreateUserWithClaims(email, password, adminClaim)
		if err != nil {
			a.logger.Error("Firebase Admin user can't be created: ", err.Error())
			return nil
		}

		username := "admin"

		err = a.userService.CreateUser(
			&domain.User{
				FirebaseId: &user, Email: &email, Username: &username,
				StripeID: &username,
			},
		)
		if err != nil {
			_ = a.service.DeleteUser(context.Background(), user)
			a.logger.Error("Firebase Admin user can't be created: ", err.Error())
			return nil
		}

		a.logger.Info(
			"Firebase Admin User Created, email: ", email, " password: ", password,
		)

		return nil
	}

	a.logger.Info("Admin already exist")
	return nil
}
