package user

import (
	"strings"
	"time"

	// "github.com/google/uuid"

	"ticket-api/config"
	"ticket-api/internal/constants"
	"ticket-api/internal/domain"
	"ticket-api/internal/errors"
	"ticket-api/internal/responses"
	"ticket-api/internal/types"
	// "ticket-api/internal/upload"
	"ticket-api/pkg"
	"ticket-api/pkg/gcp"
	"ticket-api/pkg/stripe"

	"github.com/danhper/structomap"

	"ticket-api/pkg/firebase"

	"github.com/gin-gonic/gin"
)

// Handler data type
type Handler struct {
	service             domain.UserService
	logger              pkg.Logger
	bucketService       gcp.BucketService
	config              config.Config
	firebase            firebase.Service
	stripeService       stripe.Service
}

// NewHandler creates new user handler
func NewHandler(
	s domain.UserService, l pkg.Logger, b gcp.BucketService, config config.Config,
	fs firebase.Service, st stripe.Service,
) *Handler {
	return &Handler{s, l, b, config, fs, st}
}

// SaveUser saves the user
func (u *Handler) SaveUser(c *gin.Context) {
	
	createUser := CreateUser{}

	if err := c.Bind(&createUser); err != nil {
		errors.HandleError(u.logger, c, err)
		return
	}

	user := createUser.User

	email := ""
	if user.TempEmail != nil {
		email = *user.Email
	}
	
	stripeCustomer, err := u.stripeService.CreateStripeCustomer(*user.Username, email)
	if err != nil {
		errors.HandleError(u.logger, c, err)
		return
	}

	user.StripeID = &stripeCustomer.ID

	if err := u.service.CreateUser(&user); err != nil {
		errors.HandleError(u.logger, c, err)
		return
	}

	c.JSON(200, gin.H{"data": DetailSerializer{user}})
}

func (u *Handler) UpdateProfilePicture(c *gin.Context) {
	paramID, err := types.ShouldParseUUID(c.Param("id"))
	if err != nil {
		errors.HandleValidationError(u.logger, c, errors.ErrInvalidUUID)
		return
	}

	metadata, isFile := c.Get(constants.File)

	user, err := u.service.GetUserById(paramID)
	if err != nil {
		errors.HandleValidationError(u.logger, c, errors.ErrNoUserFound)
		return
	}

	// Delete old image while uploading new
	if isFile {
		if user.ProfilePicture != nil {
			err := u.bucketService.RemoveObject(
				string(*user.ProfilePicture), u.config.StorageBucketName, true,
			)
			if err != nil {
				u.logger.Error(err)
			}
		}

		fileMetadata := metadata.(pkg.UploadedFiles)

		image := gcp.SignedURL(fileMetadata[0].URL)
		user.ProfilePicture = &image

		err := u.service.SavePicture(user, paramID)
		if err != nil {
			errors.HandleErrorWithStatus(u.logger, c, 500, errors.ErrInternal)
			return
		}

		c.JSON(200, gin.H{"data": "image updated"})
		return
	}

	errors.HandleError(u.logger, c, errors.ErrNothingToUpdate)
}

// UpdateUser updates user
func (u *Handler) UpdateUser(c *gin.Context) {
	paramID := c.Param("id")
	userID, err := types.ShouldParseUUID(paramID)
	if err != nil {
		errors.HandleValidationError(u.logger, c, errors.ErrInvalidUUID)
		return
	}

	var newUser domain.User
	if err := c.ShouldBind(&newUser); err != nil {
		errors.HandleValidationError(u.logger, c, err)
		return
	}

	userMap := structomap.New().UseSnakeCase().PickAll().
		Omit("Password").
		Omit("Id").
		Omit("UUID").
		Omit("BaseDate").
		Omit("StripeID").
		Omit("FirebaseId").
		OmitIf(
			func(u interface{}) bool {
				return newUser.ProfilePicture == nil
			}, "ProfilePicture",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.NickName == nil
			}, "NickName",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.Email == nil
			}, "Email",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.Username == nil
			}, "Username",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.Location == nil
			}, "Location",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.Bio == nil
			}, "Bio",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.Genre == nil
			}, "Genre",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.Gender == nil
			}, "Gender",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.Age == nil
			}, "Age",
		).
		OmitIf(
			func(u interface{}) bool {
				return newUser.TempEmail == nil
			}, "TempEmail",
		).
		Transform(newUser)

	if len(userMap) > 0 {
		userMap["updated_at"] = time.Now()
		userMap["id"] = userID
		err := u.service.UpdateUser(userID, &userMap)
		if err != nil {
			errors.HandleError(u.logger, c, err)
			return
		}

		c.JSON(200, gin.H{"data": userMap})
		return
	}

	errors.HandleError(u.logger, c, errors.ErrNothingToUpdate)
}

// GetAllUsers gets the user
func (u *Handler) GetAllUsers(c *gin.Context) {
	users, err := u.service.SetPaginationScope(pkg.Paginate(c)).GetAllUsers()
	if err != nil {
		u.logger.Error(err)
	}

	responses.JSONWithPagination(c, 200, users)
}

// GetUserById gets the user
func (u *Handler) GetUserById(c *gin.Context) {
	idType := c.Query("type")

	var paramID interface{}
	var response DetailSerializer
	var user domain.User
	var err error

	if idType == "" {
		paramID, err = types.ShouldParseUUID(c.Param("id"))
		if err != nil {
			errors.HandleValidationError(u.logger, c, errors.ErrInvalidUUID)
			return
		}
		user, err = u.service.GetUserById(paramID.(types.BinaryUUID))
	} else {
		paramID = c.Param("id")
		user, err = u.service.GetUserByFirebaseId(paramID.(string))
	}

	if err != nil {
		c.JSON(400, gin.H{"data": "no user with this id exists"})
		return
	}

	response = DetailSerializer{user}

	if err != nil {
		u.logger.Error(err)
		errors.HandleErrorWithStatus(u.logger, c, 400, errors.ErrNoUserFound)
		return
	}

	c.JSON(200, gin.H{"data": response})
}

// GetUserWithLessDetail gets the user with less detail
func (u *Handler) GetUserWithLessDetail(c *gin.Context) {
	ids := c.Query("ids")
	userIds := strings.Split(ids, ",")
	users, err := u.service.GetUserWithLessDetail(userIds)
	if err != nil {
		u.logger.Error(err)
		errors.HandleErrorWithStatus(u.logger, c, 400, errors.ErrNoUserFound)
		return
	}
	lessUserDetails := make(map[string]domain.LessUserDetail)

	for _, user := range users {
		lessUserDetails[user.FirebaseId] = domain.LessUserDetail{
			Id:             user.Id,
			ProfilePicture: user.ProfilePicture,
			NickName:       user.NickName,
			Username:       user.Username,
		}
	}
	c.JSON(200, gin.H{"data": lessUserDetails})
}

// DeleteUser deletes user
func (u *Handler) DeleteUser(c *gin.Context) {
	paramID, err := types.ShouldParseUUID(c.Param("id"))
	if err != nil {
		errors.HandleValidationError(u.logger, c, errors.ErrInvalidUUID)
		return
	}

	user, err := u.service.GetUserById(paramID)
	if err != nil {
		errors.HandleValidationError(u.logger, c, errors.ErrNoUserFound)
	}

	// make email feild nil before updating user
	if err := u.service.UpdateUser(paramID, &map[string]interface{}{"email": nil}); err != nil {
		errors.HandleError(u.logger, c, err)
		return
	}

	// delete user from database
	if err := u.service.DeleteUser(paramID); err != nil {
		errors.HandleError(u.logger, c, err)
		return
	}

	// delete firebase user account
	err = u.firebase.DeleteUserFromFirebase(*user.FirebaseId)
	if err != nil {
		u.service.CreateUser(&user)
		errors.HandleError(u.logger, c, err)
		return
	}

	c.JSON(200, gin.H{"data": "user deleted"})
}

// SearchUsers Search Event
func (u *Handler) SearchUsers(c *gin.Context) {
	searchQuery := strings.ToLower(c.Query("query"))
	users, err := u.service.SetPaginationScope(pkg.Paginate(c)).SearchUsers(searchQuery)
	if err != nil {
		errors.HandleValidationError(u.logger, c, errors.ErrNoUserFound)
		return
	}

	responses.JSONWithPagination(c, 200, users)
}
