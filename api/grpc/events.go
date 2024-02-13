package grpc

import (
	"context"
	"errors"
	"eventbot/Logger"
	eventbotv12 "eventbot/api/grpc/grpc_proto/github.com/mashmorsik/eventbot"
	"eventbot/internal/command"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventsClient struct {
	eventbotv12.EventsServer
	c command.EventInterface
}

//func HandleEvents(gRPC *grpc_proto.Server) {
//	eventbotv1.EventsServer(gRPC)
//}

func (ec *EventsClient) Create(ctx context.Context, req *eventbotv12.CreateEventRequest,
) (*eventbotv12.CreateEventResponse, error) {
	if req == nil {
		err := errors.New("CreateEventRequest is nil")
		Logger.Sugar.Errorln(err)
		return nil, err
	}

	if validationErr := CreateEventValidation(req); validationErr != nil {
		return nil, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	eventId, err := ec.c.CreateNewEvent(req.GetUserId(), req.GetChatId(), req.GetName(),
		req.GetTimeDate().AsTime(), req.GetCron())
	if err != nil {
		Logger.Sugar.Errorf("CreateEvent failed, error:%v\n", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("CreateNewEvent failed, error:%s\n", err.Error()))
	}

	return &eventbotv12.CreateEventResponse{
		EventId: int32(eventId),
	}, nil
}

func (ec *EventsClient) Edit(ctx context.Context, req *eventbotv12.EditEventRequest) (*emptypb.Empty, error) {
	if req == nil {
		err := errors.New("EditEventRequest is nil")
		Logger.Sugar.Errorln(err)
		return nil, err
	}

	if validationErr := EditEventValidation(req); validationErr != nil {
		return nil, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	err := ec.c.UpdateEvent(int(req.GetEventId()), req.GetName(), req.GetTimeDate().AsTime(), req.GetCron())
	if err != nil {
		Logger.Sugar.Errorln("UpdateEvent failed, error:%v\n", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("UpdateEvent failed, error:%v\n", err.Error()))
	}

	return &emptypb.Empty{}, nil
}

func (ec *EventsClient) Delete(ctx context.Context, req *eventbotv12.DeleteEventRequest) (*emptypb.Empty, error) {
	if req == nil {
		err := errors.New("DeleteEventRequest is nil")
		Logger.Sugar.Errorln(err)
		return nil, err
	}

	if validationErr := DeleteEventValidation(req); validationErr != nil {
		return nil, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	err := ec.c.DeleteEvent(int(req.GetEventId()))
	if err != nil {
		Logger.Sugar.Errorln("DeleteEvent failed, error:%v\n", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("DeleteEvent failed, error:%v\n", err.Error()))
	}

	return &emptypb.Empty{}, nil
}

func (ec *EventsClient) Get(ctx context.Context, req *eventbotv12.GetEventsRequest) (*eventbotv12.GetEventsResponse, error) {
	if req == nil {
		err := errors.New("GetEventsRequest is nil")
		Logger.Sugar.Errorln(err)
		return nil, err
	}

	if validationErr := GetEventsValidation(req); validationErr != nil {
		return nil, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	events, err := ec.c.GetEvents(req.GetUserId())
	if err != nil {
		Logger.Sugar.Errorln("GetEvents failed, error:%v\n", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("GetEvents failed, error:%v\n", err.Error()))
	}

	eventsRes := make([]*eventbotv12.GetEventsResponse_EventsMap, 0, len(events))
	for key, value := range events {
		eventsRes = append(eventsRes, &eventbotv12.GetEventsResponse_EventsMap{
			Key: int32(key),
			Value: &eventbotv12.Event{
				EventId:   int32(value.EventId),
				UserId:    value.UserId,
				Name:      value.Name,
				TimeDate:  timestamppb.New(value.TimeDate),
				Cron:      value.Cron,
				LastFired: timestamppb.New(value.LastFired),
				Disabled:  value.Disabled,
			},
		})
	}

	return &eventbotv12.GetEventsResponse{Events: eventsRes}, nil
}

func (ec *EventsClient) Disable(ctx context.Context, req *eventbotv12.DisableEventRequest) (*emptypb.Empty, error) {
	if req == nil {
		err := errors.New("DisableEventRequest is nil")
		Logger.Sugar.Errorln(err)
		return nil, err
	}

	if validationErr := DisableEventValidation(req); validationErr != nil {
		return nil, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	err := ec.c.DisableEvent(int(req.GetEventId()))
	if err != nil {
		Logger.Sugar.Errorln("DisableEvent failed, error:%v\n", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("DisableEvent failed, error:%v\n", err.Error()))
	}

	return &emptypb.Empty{}, nil
}

func (ec *EventsClient) Enable(ctx context.Context, req *eventbotv12.EnableEventRequest) (*emptypb.Empty, error) {
	if req == nil {
		err := errors.New("EnableEventRequest is nil")
		Logger.Sugar.Errorln(err)
		return nil, err
	}

	if validationErr := EnableEventValidation(req); validationErr != nil {
		return nil, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	err := ec.c.EnableEvent(int(req.GetEventId()))
	if err != nil {
		Logger.Sugar.Errorln("EnableEvent failed, error:%v\n", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("EnableEvent failed, error:%v\n", err.Error()))
	}

	return &emptypb.Empty{}, nil
}

func (ec *EventsClient) DeleteAll(ctx context.Context, req *eventbotv12.DeleteAllRequest) (*emptypb.Empty, error) {
	if req == nil {
		err := errors.New("EnableEventRequest is nil")
		Logger.Sugar.Errorln(err)
		return nil, err
	}

	if validationErr := DeleteAllEventsValidation(req); validationErr != nil {
		return nil, status.Error(codes.InvalidArgument, validationErr.Error())
	}

	err := ec.c.DeleteAllEvents(req.GetUserId())
	if err != nil {
		Logger.Sugar.Errorln("DeleteAllEvents failed, error:%v\n", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("DeleteAllEvents failed, error:%v\n", err.Error()))
	}

	return &emptypb.Empty{}, nil
}
