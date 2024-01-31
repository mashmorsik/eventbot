package grpc

import (
	"errors"
	"eventbot/api/grpc/grpc_proto/github.com/mashmorsik/eventbot"
	"go.uber.org/multierr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateEventValidation(req *eventbotv1.CreateEventRequest) error {
	var validationErr error

	if req.GetName() == "" {
		validationErr = multierr.Append(validationErr, errors.New("EventName is required"))
	}

	if req.GetChatId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("ChatId is required"))
	}

	if req.GetUserId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("UserId is required"))
	}

	if req.GetTimeDate() == nil {
		validationErr = multierr.Append(validationErr, errors.New("TimeDate is required"))
	}

	if req.GetCron() == "" {
		validationErr = multierr.Append(validationErr, errors.New("cron is required"))
	}

	if validationErr != nil {
		return status.Error(codes.InvalidArgument, validationErr.Error())
	}

	return nil
}

func EditEventValidation(req *eventbotv1.EditEventRequest) error {
	var validationErr error

	if req.GetEventId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("EventId is required"))
	}

	if req.GetName() == "" {
		validationErr = multierr.Append(validationErr, errors.New("EventName is required"))
	}

	if req.GetTimeDate() == nil {
		validationErr = multierr.Append(validationErr, errors.New("TimeDate is required"))
	}

	if req.GetCron() == "" {
		validationErr = multierr.Append(validationErr, errors.New("cron is required"))
	}

	if validationErr != nil {
		return status.Errorf(codes.InvalidArgument, validationErr.Error())
	}

	return nil
}

func DeleteEventValidation(req *eventbotv1.DeleteEventRequest) error {
	var validationErr error

	if req.GetEventId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("EventId is required"))
	}

	if validationErr != nil {
		return status.Errorf(codes.InvalidArgument, validationErr.Error())
	}

	return nil
}

func GetEventsValidation(req *eventbotv1.GetEventsRequest) error {
	var validationErr error

	if req.GetUserId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("UserId is required"))
	}

	if validationErr != nil {
		return status.Errorf(codes.InvalidArgument, validationErr.Error())
	}

	return nil
}

func DisableEventValidation(req *eventbotv1.DisableEventRequest) error {
	var validationErr error

	if req.GetEventId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("EventId is required"))
	}

	if validationErr != nil {
		return status.Errorf(codes.InvalidArgument, validationErr.Error())
	}

	return nil
}

func EnableEventValidation(req *eventbotv1.EnableEventRequest) error {
	var validationErr error

	if req.GetEventId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("EventId is required"))
	}

	if validationErr != nil {
		return status.Errorf(codes.InvalidArgument, validationErr.Error())
	}

	return nil
}

func DeleteAllEventsValidation(req *eventbotv1.DeleteAllRequest) error {
	var validationErr error

	if req.GetUserId() == 0 {
		validationErr = multierr.Append(validationErr, errors.New("UserId is required"))
	}

	if validationErr != nil {
		return status.Errorf(codes.InvalidArgument, validationErr.Error())
	}

	return nil
}
