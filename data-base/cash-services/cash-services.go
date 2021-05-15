package data_base

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (dataBase *CashDataBase) GetDataFromCash(ctx context.Context, key string) (string, error) {

	fmt.Println("Searching Data For : " + key + " in cash")

	val, err := dataBase.Cash.Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return "", status.Errorf(codes.NotFound, "Token Not Exist %v", err)
	case err != nil:
		return "", status.Errorf(codes.Internal, "Unable To Fetch Value %v", err)
	case val == "":
		return "", status.Errorf(codes.Unknown, "Empty Value %v", err)
	}

	return val, nil
}

func (dataBase *CashDataBase) SetDataToCash(ctx context.Context, key string, value interface{}, expiration time.Duration) error {

	err := dataBase.Cash.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return status.Errorf(codes.Internal, "unable to add data to cash  : %w", err)
	}

	fmt.Println("Adding Data For : " + key + " in cash")

	return nil
}

func (dataBase *CashDataBase) DelDataFromCash(ctx context.Context, key string) error {

	err := dataBase.Cash.Del(ctx, key).Err()

	if err != nil {
		return status.Errorf(codes.NotFound, "No Data Is Found", err)
	}

	return nil
}
