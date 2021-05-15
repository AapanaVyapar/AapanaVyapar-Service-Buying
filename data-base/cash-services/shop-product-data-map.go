package data_base

import (
	"context"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (dataBase *CashDataBase) AddShopProductMapDataToCash(ctx context.Context, shopId string, data interface{}) error {

	err := dataBase.Cash.LPush(ctx, "shopProductMap_"+shopId, data).Err()
	if err != nil {
		return status.Errorf(codes.Internal, "unable to add data to hash of Cash  : %w", err)
	}

	return nil

}

func (dataBase *CashDataBase) CheckIsProductBelongsToShopFromCash(ctx context.Context, shopId string, productId string) (string, error) {

	index, err := dataBase.Cash.LPos(ctx, "shopProductMap_"+shopId, productId, redis.LPosArgs{}).Result()
	if err != nil {
		return "", status.Errorf(codes.NotFound, "Shop Not Exist %v", err)
	}

	val, err := dataBase.Cash.LIndex(ctx, "shopProductMap_"+shopId, index).Result()
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

func (dataBase *CashDataBase) GetShopProductMapDataFromCash(ctx context.Context, shopId string) ([]string, error) {

	val, err := dataBase.Cash.LRange(ctx, "shopProductMap_"+shopId, 0, -1).Result()
	switch {
	case err == redis.Nil:
		return []string{}, status.Errorf(codes.NotFound, "Value Not Exist %v", err)
	case err != nil:
		return []string{}, status.Errorf(codes.Internal, "Unable To Fetch Value %v", err)
	case len(val) == 0:
		return []string{}, status.Errorf(codes.Unknown, "Empty  %v", err)
	}

	return val, nil

}

func (dataBase *CashDataBase) DelProductFromShopProductMapDataFromCash(ctx context.Context, shopId string, productId string) error {
	err := dataBase.Cash.LRem(ctx, "shopProductMap_"+shopId, 0, productId).Err()
	if err != nil {
		return status.Errorf(codes.Unknown, "Unable To Delete Data From Hash Of Cash", err)
	}
	return nil

}

func (dataBase *CashDataBase) DelShopFromShopProductMapDataFromCash(ctx context.Context, shopId string) error {
	err := dataBase.Cash.Del(ctx, "shopProductMap_"+shopId).Err()
	if err != nil {
		return status.Errorf(codes.Unknown, "Unable To Delete Data From Hash Of Cash", err)
	}
	return nil

}
