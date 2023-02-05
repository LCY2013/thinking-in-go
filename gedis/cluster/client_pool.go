package cluster

import (
	"context"
	"errors"
	"github.com/LCY2013/thinking-in-go/gedis/tcp/client"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type connectionFactory struct {
	peer string
}

func (c connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	client, err := client.MakeClient(c.peer)
	if err != nil {
		return nil, err
	}
	client.Start()
	return pool.NewPooledObject(client), err
}

func (c connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	client, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	client.Close()
	return nil
}

func (c connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// do validate
	return true
}

func (c connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do activate
	return nil
}

func (c connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate
	return nil
}
