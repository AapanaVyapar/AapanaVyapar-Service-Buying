package data_base

import (
	"context"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (dataBase *CashDataBase) AddShopDataToCash(ctx context.Context, shopId string, data interface{}) error {

	err := dataBase.Cash.HSet(ctx, "shops", shopId, data).Err()
	if err != nil {
		return status.Errorf(codes.Internal, "unable to add data to hash of Cash  : %w", err)
	}
	return nil

}

func (dataBase *CashDataBase) GetShopDataFromCash(ctx context.Context, shopId string) (string, error) {

	val, err := dataBase.Cash.HGet(ctx, "shops", shopId).Result()
	switch {
	case err == redis.Nil:
		return "", status.Errorf(codes.NotFound, "Value Not Exist %v", err)
	case err != nil:
		return "", status.Errorf(codes.Internal, "Unable To Fetch Value %v", err)
	case val == "":
		return "", status.Errorf(codes.Unknown, "Empty Value %v", err)
	}
	return val, nil

}

func (dataBase *CashDataBase) DelShopDataFromCash(ctx context.Context, shopId string) error {
	err := dataBase.Cash.HDel(ctx, "shops", shopId).Err()
	if err != nil {
		return status.Errorf(codes.Unknown, "Unable To Delete Data From Hash Of Cash", err)
	}
	return nil

}
