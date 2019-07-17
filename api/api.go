package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/FactomProject/ptnet-eventstore/event"
	"github.com/FactomProject/ptnet-eventstore/eventstore"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stackdump/gopflow/statemachine"

	pb "github.com/FactomProject/ptnet-eventstore/finite"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var apiService *grpc.Server

// KLUDGE everyone is a superuser
var roles = map[statemachine.Role]bool{
	"SuperUser": true,
}

type server struct {
	es *eventstore.EventStore
}

func (s *server) publish(ctx context.Context, evt *event.Event) {
	// REVIEW: send to google pub/sub - dispatch for async operations and/or pub/sub
	fmt.Printf("PUBLISH: %s\n", evt.String())
	/*
		e := pbEvent(evt)
		j, _ := json.Marshal(e)
		fmt.Printf("%s", j)
	*/
}

func (s *server) Dispatch(ctx context.Context, in *pb.Command) (res *pb.EventStatus, err error) {

	res = &pb.EventStatus{
		State: &pb.State{
			Id:      in.Id,
			Schema:  in.Schema,
			State:   in.State,
			Head:    "",
			Created: nil,
			Updated: nil,
		},
		Code:    -1,
		Message: "bad input",
	}

	if len(in.Action) == 0 {
		return res, errors.New("action is required")
	}

	j := json.RawMessage{}
	err = j.UnmarshalJSON(in.Payload)

	var evt *event.Event
	var state *event.State

	if err == nil {
		evt, err = event.PrepareEvent(in.Id, in.Schema, in.Action, j)
	}

	if evt != nil && err == nil {
		state, err = s.es.Commit(context.WithValue(ctx, "roles", roles), evt)

		created, _ := ptypes.TimestampProto(state.Created)
		updated, _ := ptypes.TimestampProto(state.Updated)

		res.State = &pb.State{
			Id:      state.Id.String(),
			Schema:  state.Schema,
			State:   event.PqArrayToUint(state.State),
			Head:    state.Head.String(),
			Created: created,
			Updated: updated,
		}
	}

	if err != nil {
		res.Code = 1
		res.Message = fmt.Sprintf("%v", err)
	} else {
		s.publish(ctx, evt)
		res.Code = 0
		res.Message = ""
	}

	return res, nil
}

func pbEvent(evt *event.Event) *pb.Event {
	ts, _ := ptypes.TimestampProto(evt.TS)

	payload := &any.Any{
		// REVIEW compare w/ example TypeUrl: "example.com/yaddayaddayadda/" + proto.MessageName(t1),
		// can this be anything?
		TypeUrl: "https://project.factom.com/ptnet-eventstore/event.proto#JsonPayload",
		Value:   evt.Payload,
	}

	cmd := []*pb.Action{}
	cmd = append(cmd, &pb.Action{Action: evt.Action, Multiple: evt.Multiple})

	return &pb.Event{
		Id:      evt.Id.String(),
		Schema:  evt.Schema,
		Action:  cmd,
		Payload: payload,
		State:   event.PqArrayToUint(evt.State),
		Ts:      ts,
		Uuid:    evt.Uuid.String(),
		Parent:  evt.Parent.String(),
	}
}

func (s *server) GetEvent(ctx context.Context, in *pb.Query) (l *pb.EventList, err error) {
	l = &pb.EventList{}
	if in.Uuid == "" {
		for _, e := range s.es.GetEvents(in.Schema, in.Id) {
			l.List = append(l.List, pbEvent(e))
		}
	} else {
		e := s.es.GetEvent(in.Schema, in.Uuid)
		l.List = append(l.List, pbEvent(e))
	}
	return l, nil
}

func pbState(st *event.State) *pb.State {
	tsCreated, _ := ptypes.TimestampProto(st.Created)
	tsUpdated, _ := ptypes.TimestampProto(st.Updated)

	return &pb.State{
		Id:      st.Id.String(),
		Schema:  st.Schema,
		Head:    st.Head.String(),
		State:   event.PqArrayToUint(st.State),
		Updated: tsUpdated,
		Created: tsCreated,
	}
}

func (s *server) GetState(ctx context.Context, in *pb.Query) (*pb.StateList, error) {
	l := &pb.StateList{}
	if in.Uuid == "" {
		st := s.es.GetState(in.Schema, in.Id)
		l.List = append(l.List, pbState(st))
	} else {
		panic("Select state by head/uuid not yet supported")
	}
	return l, nil
}

func (s *server) Status(ctx context.Context, in *pb.Ping) (*pb.Pong, error) {
	return &pb.Pong{Code: 0, Nonce: in.Nonce}, nil
}

func (s *server) GetMachine(ctx context.Context, in *pb.Query) (*pb.Machine, error) {
	m, ok := s.es.GetMachine(in.Schema, in.Uuid)

	if !ok {
		return nil, errors.New(fmt.Sprintf("Unknown Schema %v", in.Schema))
	}

	t := make(map[string]*pb.Transition)

	for k, v := range m.Transitions {
		gl := map[string]*pb.Guard{}
		for l, d := range v.Guards {
			gl[string(l)] = &pb.Guard{Delta: d}
		}
		t[string(k)] = &pb.Transition{Delta: v.Delta, Role: string(v.Role), Guards: gl}
	}

	return &pb.Machine{
		Schema:      in.Schema,
		Initial:     m.Initial,
		Capacity:    m.Capacity,
		Transitions: t,
	}, nil
}

func (s *server) ListMachines(ctx context.Context, in *pb.MachineQuery) (*pb.MachineList, error) {
	_ = in // no query params used
	ml := &pb.MachineList{}
	for _, m := range s.es.ListMachines() {
		ml.List = append(ml.List, m)
	}

	return ml, nil
}

func Serve() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	apiService = grpc.NewServer()
	es := eventstore.NewEventStore()
	es.Migrate()
	pb.RegisterEventStoreServer(apiService, &server{es})

	if err := apiService.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Stop() {
	defer func() {
		if err := recover(); err != nil {
			time.Sleep(time.Millisecond * 5)
			Stop()
		}
	}()
	apiService.GracefulStop()
}
