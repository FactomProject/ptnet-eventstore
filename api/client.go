package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/FactomProject/ptnet-eventstore/event"
	pb "github.com/FactomProject/ptnet-eventstore/finite"
	"google.golang.org/grpc"
	"log"
)

func TestClient() *client {
	return NewClient("localhost:50051")
}

func NewClient(address string) *client {
	return &client{pbClient(address, "finite")}
}

func pbClient(address string, name string) pb.EventStoreClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewEventStoreClient(conn)
}

type client struct {
	pb.EventStoreClient
}

func (c *client) Dispatch(ctx context.Context, schema string, id string, action []string, multiple uint64, payload interface{}, state []uint64) (*pb.State, error) {
	j, _ := json.Marshal(payload)

	status, err := c.EventStoreClient.Dispatch(
		ctx,
		&pb.Command{
			Schema:   schema,
			Id:       id,
			Action:   action,
			Multiple: multiple,
			Payload:  j,
			State:    state,
		})
	if err == nil && status.Code != 0 {
		return status.State, errors.New(status.Message)
	} else {
		return status.State, err
	}
}

func (c *client) GetState(ctx context.Context, schema string, id string, uuid string) ([]*pb.State, error) {
	stateList, err := c.EventStoreClient.GetState(
		ctx,
		&pb.Query{
			Schema: schema,
			Id:     id,
			Uuid:   uuid,
		})

	return stateList.List, err

}

func (c *client) GetEvent(ctx context.Context, schema string, id string, uuid string) ([]*pb.Event, error) {
	eventList, err := c.EventStoreClient.GetEvent(
		ctx,
		&pb.Query{
			Schema: schema,
			Id:     id,
			Uuid:   uuid,
		})

	return eventList.List, err
}

func (c *client) GetMachine(ctx context.Context, schema string, id string, uuid string) (*pb.Machine, error) {
	_ = id   // REVIEW: unsure if oid/state_id is ever useful but wanted to keep same Query param
	_ = uuid // REVIEW: if machines are ever versioned use old event uuid to get the valid historical version

	return c.EventStoreClient.GetMachine(
		ctx,
		&pb.Query{
			Schema: schema,
			Id:     id,
			Uuid:   uuid,
		})
}

func (c *client) Ping(ctx context.Context) (ok bool) {
	nonce := event.NewUuid().String()
	pong, err := c.EventStoreClient.Status(ctx, &pb.Ping{Nonce: nonce})

	if err == nil && pong.Nonce == nonce {
		return true
	} else {
		return false
	}

}

func (c *client) ListMachines(ctx context.Context) ([]string, error) {

	ml, err := c.EventStoreClient.ListMachines(ctx, &pb.MachineQuery{})
	if err != nil {
		panic(err)
	}
	return ml.List, nil
}
