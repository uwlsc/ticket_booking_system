package user

import (
	"context"

	"ticket-api/pkg"

	"ticket-api/internal/middlewares"

	"go.uber.org/fx"
)

// Routes struct
type Routes struct {
	logger               pkg.Logger
	handler              pkg.HTTPServer
	userHandler          *Handler
	paginationMiddleware middlewares.PaginationMiddleware
	uploadMiddleware     middlewares.UploadMiddleware
	authMiddleware       middlewares.FirebaseAuthMiddleware
}

func NewRoutes(
	lc fx.Lifecycle, logger pkg.Logger, handler pkg.HTTPServer, userHandler *Handler,
	p middlewares.PaginationMiddleware,
	uploadMiddleware middlewares.UploadMiddleware, a middlewares.FirebaseAuthMiddleware,
) *Routes {
	routes := &Routes{
		handler:              handler,
		logger:               logger,
		userHandler:          userHandler,
		paginationMiddleware: p,
		uploadMiddleware:     uploadMiddleware,
		authMiddleware:       a,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				routes.Setup()
				return nil
			},
		},
	)

	return routes
}

// Setup user routes
func (s *Routes) Setup() {
	s.logger.Info("Setting up User routes")

	api := s.handler.Group("/user")

	api.GET("", s.authMiddleware.HandleAuthWithRole(""), s.paginationMiddleware.Handle(), s.userHandler.GetAllUsers)

	api.GET("/less", s.authMiddleware.HandleAuthWithRole(""), s.userHandler.GetUserWithLessDetail)

	api.POST("", s.userHandler.SaveUser)

	api.GET("/:id", s.userHandler.GetUserById)

	api.PUT("/:id", s.authMiddleware.HandleAuthWithRole(""), s.userHandler.UpdateUser)

	api.GET(
		"/search", s.authMiddleware.HandleAuthWithRole("admin"), s.paginationMiddleware.Handle(),
		s.userHandler.SearchUsers,
	)

	api.PUT(
		"/profile-picture/:id",
		s.authMiddleware.HandleAuthWithRole(""),
		s.uploadMiddleware.Push(
			s.uploadMiddleware.Config().
				WebpEnable(true).
				Folder("user").
				ThumbEnable(true).
				Extension(".jpg", ".jpeg", ".png", ".webp", ".heic", ".heif", "avif").
				Field("profile_picture"),
		).Handle(),
		s.userHandler.UpdateProfilePicture,
	)

	api.DELETE("/:id", s.authMiddleware.HandleAuthWithRole(""), s.userHandler.DeleteUser)
}
