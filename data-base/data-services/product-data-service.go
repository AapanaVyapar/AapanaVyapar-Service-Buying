package data_base

import (
	"aapanavyapar-service-buying/configurations/mongodb"
	"aapanavyapar-service-buying/data-base/helpers"
	"aapanavyapar-service-buying/data-base/structs"
	"aapanavyapar-service-buying/pb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"time"
)

func (dataBase *MongoDataBase) CreateProduct(context context.Context, dataInsert structs.ProductData) (primitive.ObjectID, error) {

	if err := helpers.Validate(dataInsert); err != nil {
		return primitive.ObjectID{}, err
	}

	name, err := dataBase.GetNameFromShopData(context, dataInsert.ShopId)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("shop with specified id does not exist")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataInsert.Timestamp = time.Now().UTC()
	dataInsert.Likes = 0
	dataInsert.ShopName = name

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	id, err := productData.InsertOne(context, dataInsert)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return id.InsertedID.(primitive.ObjectID), nil
}

func (dataBase *MongoDataBase) GetAllProductsOfShopByFunctionFromProductData(context context.Context, shopId string, sendData func(data structs.ProductData) error) error {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{{"shop_id", shopId}}
	cursor, err := productData.Find(context, filter)

	if err != nil {
		return err
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		result := structs.ProductData{}
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

func (dataBase *MongoDataBase) GetAllProductsFromProductData(context context.Context, sendData func(data structs.ProductData) error) error {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{}
	cursor, err := productData.Find(context, filter)

	if err != nil {
		return err
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		result := structs.ProductData{}
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

func (dataBase *MongoDataBase) GetAllProductsByCategoryOfShopsByFunctionFromProductData(context context.Context, shopIds []string, category []pb.Category, sendData func(data structs.ProductData) error) error {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{
		{"shop_id", bson.M{"$in": shopIds}},
		{"category", bson.M{"$in": category}},
	}
	cursor, err := productData.Find(context, filter)

	if err != nil {
		return err
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		result := structs.ProductData{}
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

func (dataBase *MongoDataBase) GetAllProductsOfShopByArrayFromProductData(context context.Context, shopId string) ([]structs.ProductData, error) {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{{"shop_id", shopId}}
	cursor, err := productData.Find(context, filter)

	if err != nil {
		return []structs.ProductData{}, err
	}
	defer cursor.Close(context)

	if err := cursor.Err(); err != nil {
		return []structs.ProductData{}, err
	}

	var results []structs.ProductData
	err = cursor.All(context, &results)
	if err != nil {
		return []structs.ProductData{}, err
	}

	return results, nil

}

func (dataBase *MongoDataBase) GetSpecificProductsOfShopFromProductData(context context.Context, shopId string, productId primitive.ObjectID) (structs.ProductData, error) {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{{"shop_id", shopId}, {"_id", productId}}

	var data structs.ProductData
	err := productData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return structs.ProductData{}, err
	}

	return data, nil

}

func (dataBase *MongoDataBase) GetProductFromProductData(context context.Context, productId primitive.ObjectID) (structs.ProductData, error) {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{{"_id", productId}}

	var data structs.ProductData
	err := productData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return structs.ProductData{}, err
	}

	return data, nil

}

func (dataBase *MongoDataBase) GetShopDetailsFromProductData(context context.Context, productId primitive.ObjectID) (*structs.ProductData, error) {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{{"_id", productId}}

	var data structs.ProductData
	err := productData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}

func (dataBase *MongoDataBase) IsExistProductExist(context context.Context, key string, value interface{}) bool {
	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.D{{key, value}}
	singleCursor := productData.FindOne(context, filter)

	if singleCursor.Err() != nil {
		return false
	}

	return true

}

func (dataBase *MongoDataBase) DelProductFromProductData(context context.Context, shopId string, productId primitive.ObjectID) error {
	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.M{"shop_id": shopId, "_id": productId}

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	_, err := productData.DeleteOne(context, filter)
	if err != nil {
		return err
	}

	return nil
}

func (dataBase *MongoDataBase) DelProductsOfShopFromProductData(context context.Context, shopId string) error {
	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	filter := bson.M{"shop_id": shopId}

	_, err := productData.DeleteMany(context, filter)
	if err != nil {
		return err
	}

	return nil
}

func (dataBase *MongoDataBase) AddProductImageInProductData(context context.Context, shopId string, productId primitive.ObjectID, imageURL string) error {

	if _, err := url.ParseRequestURI(imageURL); err != nil {
		return fmt.Errorf("invalid image url")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$push": bson.M{
				"images": imageURL,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update stock")

}

func (dataBase *MongoDataBase) DelProductImageFromProductData(context context.Context, shopId string, productId primitive.ObjectID, imageURL string) error {

	if _, err := url.ParseRequestURI(imageURL); err != nil {
		return fmt.Errorf("invalid image url")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$pull": bson.M{
				"images": imageURL,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update stock")

}

func (dataBase *MongoDataBase) UpdateProductTitleInProductData(context context.Context, shopId string, productId primitive.ObjectID, title string) error {

	if title == "" {
		return fmt.Errorf("title can not be empty")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$set": bson.M{
				"title": title,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update product title")
}

func (dataBase *MongoDataBase) UpdateProductCategoryInProductData(context context.Context, shopId string, productId primitive.ObjectID, category []pb.Category) error {

	if len(category) == 0 {
		return fmt.Errorf("category can not be empty")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$set": bson.M{
				"category": category,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update product category")
}

func (dataBase *MongoDataBase) UpdateProductDescriptionInProductData(context context.Context, shopId string, productId primitive.ObjectID, description string) error {

	if description == "" {
		return fmt.Errorf("description can not be empty")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$set": bson.M{
				"description": description,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update product description")
}

func (dataBase *MongoDataBase) UpdateProductShortDescriptionInProductData(context context.Context, shopId string, productId primitive.ObjectID, shortDescription string) error {

	if shortDescription == "" {
		return fmt.Errorf("short description can not be empty")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$set": bson.M{
				"short_description": shortDescription,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update product short description")
}

func (dataBase *MongoDataBase) UpdateProductShippingInfoInProductData(context context.Context, shopId string, productId primitive.ObjectID, shippingInfo string) error {

	if shippingInfo == "" {
		return fmt.Errorf("shipping info can not be empty")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$set": bson.M{
				"shipping_info": shippingInfo,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update shipping info")

}

func (dataBase *MongoDataBase) UpdateProductStockInfoInProductData(context context.Context, shopId string, productId primitive.ObjectID, stock uint32) error {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
			//"$expr":   bson.M{"$lte": bson.A{"max_stock", stock}},
		},
		bson.M{
			"$set": bson.M{
				"stock": stock,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update stock")
}

func (dataBase *MongoDataBase) UpdateProductPriceInProductData(context context.Context, shopId string, productId primitive.ObjectID, price float64) error {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$set": bson.M{
				"price": price,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update product price")

}

func (dataBase *MongoDataBase) UpdateProductOfferInProductData(context context.Context, shopId string, productId primitive.ObjectID, offer uint8) error {

	if offer > 100 {
		return fmt.Errorf("offer should not exceed 100")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := productData.UpdateOne(context,
		bson.M{
			"shop_id": shopId,
			"_id":     productId,
		},
		bson.M{
			"$set": bson.M{
				"offer": offer,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update offer")

}

func (dataBase *MongoDataBase) DecreaseStockToMakeOrderFromProductData(context context.Context, productId primitive.ObjectID, quantity uint32) (float32, uint32, error) {

	if quantity == 0 {
		return 0, 0, fmt.Errorf("quantity can not be zero")
	}

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	var data structs.ProductData

	err := productData.FindOneAndUpdate(context,
		bson.M{
			"_id":   productId,
			"stock": bson.M{"$gte": quantity},
		},
		bson.A{
			bson.M{
				"$set": bson.M{
					"stock": bson.M{"$subtract": bson.A{"$stock", quantity}},
				},
			},
		},
	).Decode(&data)

	if err != nil {
		return 0, 0, err
	}

	return data.Price, data.Offer, nil

}

func (dataBase *MongoDataBase) IncreaseStockFromProductData(context context.Context, productId primitive.ObjectID) error {

	productData := mongodb.OpenProductDataCollection(dataBase.Data)

	result, err := productData.UpdateOne(context,
		bson.M{
			"_id": productId,
		},
		bson.M{
			"$inc": bson.M{
				"stock": 1,
			},
		},
	)

	if err != nil {
		return err
	}

	fmt.Println(result.ModifiedCount)

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("max product limit reach") // Check for inconsistency

}

/*
	dataProduct := structs.ProductData{
		_id:          primitive.NewObjectID(), //product id
		ShopId:       dataInsert.ShopId,
		Title:        "Yellow Shirt",
		Description:  "Best in Class Size XL",
		ShippingInfo: "200x70x10",
		Stock:        10,
		Price:        1000,
		Offer:        10,
		Category:     []constants.Categories{constants.MENS_CLOTHING},
		Images:       []string{"https://image.com"},
		Timestamp:    time.Now().UTC(),
	}

*/
