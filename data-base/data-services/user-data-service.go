package data_base

import (
	"aapanavyapar-service-buying/configurations/mongodb"
	"aapanavyapar-service-buying/data-base/helpers"
	"aapanavyapar-service-buying/data-base/structs"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (dataBase *MongoDataBase) CreateUser(context context.Context, userId string, userName string) (string, error) {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	if dataBase.IsExistInUserData(context, "_id", userId) {
		return "", fmt.Errorf("already exist")
	}

	dataInsert := structs.UserData{
		UserId:    userId,
		UserName:  userName,
		Address:   nil,
		Cart:      nil,
		Favorites: nil,
		Orders:    nil,
	}

	id, err := userData.InsertOne(context, dataInsert)
	if err != nil {
		return "", err
	}

	return id.InsertedID.(string), nil
}

func (dataBase *MongoDataBase) IsExistInUserData(context context.Context, key string, value interface{}) bool {
	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	filter := bson.D{{key, value}}
	singleCursor := userData.FindOne(context, filter)

	if singleCursor.Err() != nil {
		return false
	}

	return true
}

func (dataBase *MongoDataBase) SetAddressInUserData(context context.Context, userId string, userName string, address structs.Address) error {

	if err := helpers.Validate(address); err != nil {
		return err
	}

	_, _ = dataBase.CreateUser(context, userId, userName) // If already exist not going to create.

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := userData.UpdateOne(context,
		bson.M{"_id": userId},
		bson.D{
			{"$set", bson.D{{"address", address}}},
		},
	)
	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to set the address")
}

func (dataBase *MongoDataBase) DelAddressInUserData(context context.Context, userId string) error {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := userData.UpdateOne(context,
		bson.M{
			"_id": userId,
		},
		bson.D{
			{"$set",
				bson.D{
					{"address", nil},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to delete address")
}

func (dataBase *MongoDataBase) GetAddressUserData(context context.Context, userId string) (*structs.Address, error) {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	filter := bson.D{{"_id", userId}}

	data := structs.UserData{}
	err := userData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data.Address, nil
}

func (dataBase *MongoDataBase) GetCartUserData(context context.Context, userId string) (*structs.ProductIdsForCart, error) {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	filter := bson.D{{"_id", userId}}

	data := structs.UserData{}
	err := userData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data.Cart, nil
}

func (dataBase *MongoDataBase) GetFavoritesUserData(context context.Context, userId string) (*structs.ProductIdsForFavAndOrd, error) {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	filter := bson.D{{"_id", userId}}

	data := structs.UserData{}
	err := userData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data.Favorites, nil
}

func (dataBase *MongoDataBase) GetOrdersUserData(context context.Context, userId string) (*structs.ProductIdsForFavAndOrd, error) {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	filter := bson.D{{"_id", userId}}

	data := structs.UserData{}
	err := userData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data.Orders, nil
}

func (dataBase *MongoDataBase) AddToCartUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	//Checking In Cash For Existence Of Product So No Need To Check In DataBase And If Some Inconsistency Occurs
	//Then What When User Add Other Items In Cart It Get Automatically Deleted And In Front End We Make That Product Id Invisible If Product Not Exist
	//
	//if !dataBase.IsExistProductExist(context, "_id", productId) {
	//	return fmt.Errorf("product does not exist")
	//}
	//

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result := userData.FindOne(context, bson.M{"_id": userId, "cart.products": productId})

	// Error will be thrown if favorites is null or product is not in favorites in both cases we have to just add product
	if result.Err() != nil {
		res, err := userData.UpdateOne(context,
			bson.M{
				"_id": userId,
			},
			bson.D{
				{"$push",
					bson.M{
						"cart.products": bson.M{
							"$each":  bson.A{productId},
							"$slice": -15,
						},
					},
				},
			},
		)
		if err != nil {
			return err
		}

		if res.ModifiedCount > 0 || res.MatchedCount > 0 {
			return nil
		}

		return fmt.Errorf("unable to add to cart")
	}

	return fmt.Errorf("alredy exist in cart")
}

func (dataBase *MongoDataBase) AddToFavoritesUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	if !dataBase.IsExistProductExist(context, "_id", productId) {
		return fmt.Errorf("product does not exist")
	}

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result := userData.FindOne(context, bson.M{"_id": userId, "favorites.products": productId})

	// Error will be thrown if favorites is null or product is not in favorites in both cases we have to just add product
	if result.Err() != nil {
		res, err := userData.UpdateOne(context,
			bson.M{
				"_id": userId,
			},
			bson.D{
				{"$push",
					bson.M{
						"favorites.products": bson.M{
							"$each":  bson.A{productId},
							"$slice": -20,
						},
					},
				},
			},
		)
		if err != nil {
			return err
		}

		if res.ModifiedCount > 0 || res.MatchedCount > 0 {
			return nil
		}

		return fmt.Errorf("unable to add to faviroute")
	}

	return fmt.Errorf("alredy exist in faviroute")
}

func (dataBase *MongoDataBase) AddToOrdersUserData(context context.Context, userId string, orderId primitive.ObjectID) error {

	if !dataBase.IsExistOrderExist(context, "_id", orderId) {
		return fmt.Errorf("order is not created")
	}

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	//dataBase.mutex.Lock()
	//defer dataBase.mutex.Unlock()

	result := userData.FindOne(context, bson.M{"_id": userId, "orders.products": orderId})

	// Error will be thrown if favorites in null or product is not in favorites in both cases we have to just add product
	if result.Err() != nil {
		res, err := userData.UpdateOne(context,
			bson.M{
				"_id": userId,
			},
			bson.D{
				{"$push",
					bson.M{
						"orders.products": orderId,
					},
				},
			},
		)
		if err != nil {
			return err
		}

		if res.ModifiedCount > 0 || res.MatchedCount > 0 {
			return nil
		}

		return fmt.Errorf("unable to add order")
	}

	return fmt.Errorf("order alredy exist")
}

func (dataBase *MongoDataBase) DelFromCartUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := userData.UpdateOne(context,
		bson.M{
			"_id": userId,
		},
		bson.M{
			"$pull": bson.M{
				"cart.products": productId,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to delete from cart")
}

func (dataBase *MongoDataBase) DelFromFavoritesUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := userData.UpdateOne(context,
		bson.M{
			"_id": userId,
		},
		bson.M{
			"$pull": bson.M{
				"favorites.products": productId,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to delete from faviroute")
}

func (dataBase *MongoDataBase) DelFromOrdersUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := userData.UpdateOne(context,
		bson.M{
			"_id": userId,
		},
		bson.M{
			"$pull": bson.M{
				"orders.products": productId,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to delete from order")
}

/*
	dataInsert := structs.UserData{
		_Id:        "1", //userid
		UserName:  "Shitij",
		Address:   structs.Address{
			FullName:      "Shitij Shailendra Agrawal",
			HouseDetails:  "B.K Road Chopda",
			StreetDetails: "B.K Road Chopda",
			LandMark:      "HDFC Bank",
			PinCode:       "425107",
			City:          "Chopda",
			State:         "Maharastra",
			Country:       "India",
			PhoneNo:       "9172879779",
		},
		Cart:      structs.ShopAndProductIds{
			Products: []string{"1", "2"},
		},
		Favorites: structs.ShopAndProductIds{
			Products: []string{"1", "2"},
		},
		Orders:    structs.ShopAndProductIds{
			Products: []string{"1", "2"},
		},
	}

*/

/*

func (dataBase *MongoDataBase) AddToCartUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	if !dataBase.IsExistProductExist(context, "_id", productId) {
		return fmt.Errorf("product does not exist")
	}

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result := userData.FindOne(context, bson.M{"_id": userId, "cart.products.product_id": productId})

	// Error will be thrown if cart is null or product is not in cart in both cases we have to just add product
	if result.Err() != nil {
		result, err := userData.UpdateOne(context,
			bson.M{
				"_id": userId,
			},
			bson.M{
				"$push": bson.M{
					"cart.products": bson.D{
						{"product_id", productId},
						{"no_product", 1},
					},
				},
			},
		)
		if err != nil {
			return err
		}

		if result.ModifiedCount > 0 || result.MatchedCount > 0 {
			return nil
		}

		return fmt.Errorf("unable to add to cart")

	}

	res, err := userData.UpdateOne(context,
		bson.M{
			"_id":                      userId,
			"cart.products.product_id": productId,
		},
		bson.M{
			"$inc": bson.M{
				"cart.products.$.no_product": 1,
			},
		},
	)
	if err != nil {
		return err
	}

	if res.ModifiedCount > 0 || res.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to add to cart")
}

func (dataBase *MongoDataBase) DelFromCartUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	result, err := userData.UpdateOne(context,
		bson.M{
			"_id": userId,
		},
		bson.M{
			"$pull": bson.M{
				"cart.products": bson.D{
					{"product_id", productId},
				},
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to delete from cart")
}

func (dataBase *MongoDataBase) RemoveFromCartUserData(context context.Context, userId string, productId primitive.ObjectID) error {

	userData := mongodb.OpenUserDataCollection(dataBase.Data)

	result, err := userData.UpdateOne(context,
		bson.M{
			"_id":                      userId,
			"cart.products.product_id": productId,
			"cart.products.no_product": bson.M{"$gt": 1},
		},
		bson.M{
			"$inc": bson.M{
				"cart.products.$.no_product": -1,
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

	return fmt.Errorf("minimum quantity for cart is one")

}

*/
