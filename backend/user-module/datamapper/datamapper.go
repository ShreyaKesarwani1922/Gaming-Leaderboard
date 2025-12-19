package datamapper

import "context"

type IDataMapper interface {
	MapUserResponse(ctx context.Context) error
}

type DataMapper struct {
}

func (d DataMapper) MapUserResponse(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
