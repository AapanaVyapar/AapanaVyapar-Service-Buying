package data_base

import (
	"aapanavyapar-service-buying/configurations/mongodb"
	"aapanavyapar-service-buying/data-base/structs"
	"aapanavyapar-service-buying/pb"
	"context"
	"fmt"
	"github.com/razorpay/razorpay-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

func (dataBase *MongoDataBase) CreateOrder(context context.Context, userId string, paymentId string, cache structs.CacheOrder, payment map[string]interface{}, razorpaySignature string, razorpayClient *razorpay.Client) (primitive.ObjectID, error) {
	if !dataBase.IsExistInUserData(context, "_id", userId) {
		return primitive.ObjectID{}, fmt.Errorf("user does not exist")
	}

	productData, err := dataBase.GetShopDetailsFromProductData(context, cache.ProductId)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("product does not exist")
	}

	if productData.Stock < 1 {
		return primitive.ObjectID{}, fmt.Errorf("product out of stock")
	}

	var order structs.OrderData

	orderCollection := mongodb.OpenOrderDataCollection(dataBase.Data)
	order.ShopId = productData.ShopId
	order.ProductName = productData.Title
	order.ProductImage = productData.Images[0]
	order.OrderTimeStamp = time.Now().UTC()
	order.UserId = userId
	order.Quantity = cache.Quantity
	order.ProductId = cache.ProductId
	order.Address = &cache.Address
	order.RazorpaySignature = razorpaySignature
	order.DeliveryTimeStamp = cache.DeliveryTimeStamp
	order.DeliveryCost = cache.DeliveryCost
	order.Price = float32(payment["amount"].(float64))
	order.Offer = cache.Offer

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := dataBase.Data.StartSession()
	if err != nil {
		return primitive.ObjectID{}, err
	}
	defer session.EndSession(context)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		_, _, err := dataBase.DecreaseStockToMakeOrderFromProductData(sessCtx, cache.ProductId, cache.Quantity)
		if err != nil {
			return primitive.ObjectID{}, err
		}

		order.Status = pb.Status_PENDING

		id, err := orderCollection.InsertOne(sessCtx, order)
		if err != nil {
			return primitive.ObjectID{}, err
		}

		err = dataBase.AddToOrdersUserData(sessCtx, userId, id.InsertedID.(primitive.ObjectID))
		if err != nil {
			return primitive.ObjectID{}, err
		}

		resp, err := razorpayClient.Payment.Capture(paymentId, payment["amount"].(float64), map[string]interface{}{
			"currency": "INR",
		}, nil)
		if err != nil {
			return primitive.ObjectID{}, err
		}

		if resp["captured"].(bool) == false {
			return primitive.ObjectID{}, fmt.Errorf("unable to capture payment")
		}

		return id.InsertedID.(primitive.ObjectID), nil
	}

	result, err := session.WithTransaction(context, callback, txnOpts)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return result.(primitive.ObjectID), nil
}

func (dataBase *MongoDataBase) UpdateOrderStatusInOrderData(context context.Context, orderId primitive.ObjectID, status pb.Status) error {

	orderData := mongodb.OpenOrderDataCollection(dataBase.Data)

	result, err := orderData.UpdateOne(context,
		bson.M{
			"_id": orderId,
		},
		bson.M{
			"$set": bson.M{
				"status": status,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update order")

}

func (dataBase *MongoDataBase) GetOrderInfoFromOrderData(context context.Context, orderId primitive.ObjectID) (structs.OrderData, error) {

	orderData := mongodb.OpenOrderDataCollection(dataBase.Data)

	var data structs.OrderData

	err := orderData.FindOne(context,
		bson.M{
			"_id": orderId,
		},
	).Decode(&data)

	if err != nil {
		return structs.OrderData{}, err
	}

	return data, nil

}

func (dataBase *MongoDataBase) GetOrderInfoByShopIdFromOrderData(context context.Context, shopId string, sendData func(data structs.OrderData) error) error {

	orderData := mongodb.OpenOrderDataCollection(dataBase.Data)

	cursor, err := orderData.Find(context,
		bson.M{
			"shop_id": shopId,
		},
	)

	if err != nil {
		return err
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		result := structs.OrderData{}
		err = cursor.Decode(&result)

		if err != nil {
			return err
		}

		if err = sendData(result); err != nil {
			return err
		}

	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil

}

func (dataBase *MongoDataBase) IsExistOrderExist(context context.Context, key string, value interface{}) bool {
	productData := mongodb.OpenOrderDataCollection(dataBase.Data)

	filter := bson.D{{key, value}}
	singleCursor := productData.FindOne(context, filter)

	if singleCursor.Err() != nil {
		return false
	}

	return true

}

/*
	order := structs.OrderData{
		OrderId:   primitive.ObjectID{},
		UserId:    "",
		Status:    0,
		ProductId: primitive.ObjectID{},
		TimeStamp: time.Time{},
		Price:     0,
		Quantity:  0,
	}

*/
