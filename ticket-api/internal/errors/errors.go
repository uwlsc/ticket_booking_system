package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNothingToUpdate              = errors.New("nothing to update")
	ErrNothingToUpload              = errors.New("nothing to upload")
	ErrUnauthorizedAccess           = errors.New("unauthorized access")
	ErrInvalidUUID                  = errors.New("invalid uuid")
	ErrInvalidEventID               = errors.New("invalid event id")
	ErrEventNotFound                = errors.New("event not found")
	ErrInvalidEventParticipantID    = errors.New("invalid event participant id")
	ErrExtensionMismatch            = errors.New("file extension not supported")
	ErrThumbExtensionMismatch       = errors.New("file extension not supported for thumbnail")
	ErrFileRead                     = errors.New("file read error")
	ErrNoUserFound                  = errors.New("no user found for provided user id")
	ErrNoEventParticipantFound      = errors.New("no event participant found")
	ErrCannotCreateEventParticipant = errors.New("cannot create event participant")
	ErrDuplicateEventParticipant    = errors.New("user has already registered in the event")
	ErrNoAttachmentFound            = errors.New("no attachment found for provided id")
	ErrInternal                     = errors.New("internal server error")
	ErrAlreadySubscribed            = errors.New("user has already subscribed")
	ErrNoSubscription               = errors.New("no subscription yet")
	ErrInvalidSubscriptionTime      = errors.New("invalid subscription time")
	ErrInvalidPaymentMethod         = errors.New("invalid card details")
	ErrFailedToCreateStripeId       = errors.New("failed to create stripe customer id")
	ErrFailedToUnsubscribe          = errors.New("failed to unsubscribe")
	ErrPaymentMethodNotFound        = errors.New("event not found")
	ErrUrlRead                      = errors.New("url read error")
	ErrImageFormat                  = errors.New("invalid image format")
	ErrNoFile                       = errors.New("no file was sent")
)

// ErrTokenVerification for dynamic error
type ErrTokenVerification struct {
	id string
}

func NewErrTokenVerification(id string) error {
	return ErrTokenVerification{
		id: id,
	}
}

func (e ErrTokenVerification) Error() string {
	return fmt.Sprintf("error verifying id token %v\n", e.id)
}
