package data_base

import (
	"aapanavyapar-service-buying/configurations/mongodb"
	"aapanavyapar-service-buying/data-base/helpers"
	"aapanavyapar-service-buying/data-base/mapper"
	"aapanavyapar-service-buying/data-base/structs"
	"aapanavyapar-service-buying/pb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"net/url"
	"time"
)

func (dataBase *MongoDataBase) CreateShop(context context.Context, dataInsert *structs.ShopData) (string, error) {

	if err := helpers.Validate(dataInsert); err != nil {
		return "", err
	}

	for _, i := range dataInsert.Images {
		if _, err := url.ParseRequestURI(i); err != nil {
			return "", fmt.Errorf("invalid url")
		}
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataInsert.SectorNo = mapper.MapLocationToSector(dataInsert.Location)
	dataInsert.Timestamp = time.Now().UTC()

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	id, err := shopData.InsertOne(context, dataInsert)
	if err != nil {
		return "", err
	}

	return id.InsertedID.(string), nil
}

func (dataBase *MongoDataBase) GetAllShopsFromShopData(context context.Context, sendData func(data structs.ShopData) error) error {

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	filter := bson.D{}
	cursor, err := shopData.Find(context, filter)

	if err != nil {
		return err
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		result := structs.ShopData{}
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

func (dataBase *MongoDataBase) GetShopFromShopData(context context.Context, shopId string) (structs.ShopData, error) {

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	filter := bson.D{{"_id", shopId}}

	data := structs.ShopData{}
	err := shopData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return structs.ShopData{}, err
	}

	return data, nil
}

func (dataBase *MongoDataBase) AddRatingInShopData(context context.Context, shopId string, rating structs.Rating) error {

	if err := helpers.Validate(rating); err != nil {
		return err
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result := shopData.FindOne(context, bson.M{"_id": shopId, "ratings.user_id": rating.UserId})

	// Error will be thrown if rating is null or rating is already present in both cases we have to just add product
	if result.Err() != nil {
		res, err := shopData.UpdateOne(context,
			bson.M{
				"_id": shopId,
			},
			bson.D{
				{"$push",
					bson.M{
						"ratings": rating,
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

		return fmt.Errorf("unable to add rating to shop")
	}

	return fmt.Errorf("you already rated to the shop")

}

func (dataBase *MongoDataBase) GetRatingsFromShopData(context context.Context, shopId string) (*[]structs.Rating, error) {

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	filter := bson.D{{"_id", shopId}}

	data := structs.ShopData{}
	err := shopData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data.Ratings, nil
}

func (dataBase *MongoDataBase) GetNameFromShopData(context context.Context, shopId string) (string, error) {

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	filter := bson.D{{"_id", shopId}}

	data := structs.ShopData{}
	err := shopData.FindOne(context, filter).Decode(&data)

	if err != nil {
		return "", err
	}

	return data.ShopName, nil
}

func (dataBase *MongoDataBase) IsExistShopExist(context context.Context, key string, value interface{}) bool {
	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	filter := bson.D{{key, value}}
	singleCursor := shopData.FindOne(context, filter)

	if singleCursor.Err() != nil {
		return false
	}

	return true

}

func (dataBase *MongoDataBase) DelShopFromShopData(context context.Context, shopId string) (int64, error) {
	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	filter := bson.M{"_id": shopId}

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := dataBase.Data.StartSession()
	if err != nil {
		return 0, err
	}
	defer session.EndSession(context)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := dataBase.DelProductsOfShopFromProductData(sessCtx, shopId)
		if err != nil {
			return 0, err
		}

		count, err := shopData.DeleteOne(sessCtx, filter)
		if err != nil {
			return 0, err
		}

		return count.DeletedCount, nil
	}

	result, err := session.WithTransaction(context, callback, txnOpts)
	if err != nil {
		return 0, err
	}

	return result.(int64), nil
}

func (dataBase *MongoDataBase) DelShopImageFromShopData(context context.Context, shopId string, imageURL string) error {

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
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

	return fmt.Errorf("unable to remove image from the shop")
}

func (dataBase *MongoDataBase) AddShopImageInShopData(context context.Context, shopId string, imageURL string) error {

	if _, err := url.ParseRequestURI(imageURL); err != nil {
		return fmt.Errorf("invalid image url")
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
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

	return fmt.Errorf("unable to add image to the shop")

}

func (dataBase *MongoDataBase) UpdateShopPrimaryImageInShopData(context context.Context, shopId string, imageURL string) error {

	if _, err := url.ParseRequestURI(imageURL); err != nil {
		return fmt.Errorf("invalid image url")
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
		},
		bson.M{
			"$set": bson.M{
				"primary_image": imageURL,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update primary image")
}

func (dataBase *MongoDataBase) UpdateShopKeeperNameInShopData(context context.Context, shopId string, name string) error {

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
		},
		bson.M{
			"$set": bson.M{
				"shop_keeper_name": name,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update shopkeepers name")

}

func (dataBase *MongoDataBase) UpdateShopAddressAndLocationInShopData(context context.Context, shopId string, address structs.Address, location structs.Location) error {

	if err := helpers.Validate(address); err != nil {
		return err
	}
	if err := helpers.Validate(location); err != nil {
		return err
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)
	sectorNo := mapper.MapLocationToSector(&location)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
		},
		bson.M{
			"$set": bson.M{
				"address":   address,
				"location":  location,
				"sector_no": sectorNo,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update address")
}

func (dataBase *MongoDataBase) UpdateShopLocationInShopData(context context.Context, shopId string, location structs.Location) error {

	if err := helpers.Validate(location); err != nil {
		return err
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
		},
		bson.M{
			"$set": bson.M{
				"location": location,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update location")
}

func (dataBase *MongoDataBase) UpdateCategoryInShopData(context context.Context, shopId string, category []pb.Category) error {

	if len(category) == 0 {
		return fmt.Errorf("category can not be empty")
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
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

	return fmt.Errorf("unable to update category")
}

func (dataBase *MongoDataBase) GetCategoryFromShopData(context context.Context, shopId string) ([]pb.Category, error) {

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	shop := structs.ShopData{}
	err := shopData.FindOne(context,
		bson.M{
			"_id": shopId,
		},
	).Decode(&shop)

	if err != nil {
		return []pb.Category{}, err
	}

	return shop.Category, nil
}

func (dataBase *MongoDataBase) UpdateBusinessInfoInShopData(context context.Context, shopId string, info string) error {

	if info == "" {
		return fmt.Errorf("can not accept empty information")
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
		},
		bson.M{
			"$set": bson.M{
				"business_information": info,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update business information")

}

func (dataBase *MongoDataBase) UpdateOperationalHoursInShopData(context context.Context, shopId string, operationalHours structs.OperationalHours) error {

	if err := helpers.Validate(operationalHours); err != nil {
		return err
	}

	shopData := mongodb.OpenShopDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result, err := shopData.UpdateOne(context,
		bson.M{
			"_id": shopId,
		},
		bson.M{
			"$set": bson.M{
				"operational_hours": operationalHours,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 || result.MatchedCount > 0 {
		return nil
	}

	return fmt.Errorf("unable to update operational hours")

}

/*
{
    _id: '123',
    shop_name: 'Testing Stores',
    shop_keeper_name: 'ABC Person',
    images: [
        'https://store.com'
    ],
    primary_image: 'https://www.primary.store.com',
    address: {
        full_name: 'ABC Person',
        house_details: 'Testing Store',
        street_details: 'Mustufa Chishti Colony Main Rd, Panchshil Nagar',
        land_mark: 'Milap Store',
        code: '425107',
        city: 'chopda',
        state: 'maharastra',
        country: 'india',
        phone_no: '9890713171'
    },
    location: {
        longitude: '21.246435522726177',
        latitude: '75.29615236552934'
    },
    sector_no: 10,
    category: [
        7,
        4
    ],
    business_information: 'Famous Seller Of Cloths In Chopda',
    operational_hours: {
        sunday: [
            '0AM',
            '0PM'
        ],
        monday: [
            '9AM',
            '9PM'
        ],
        tuesday: [
            '9AM',
            '9PM'
        ],
        wednesday: [
            '9AM',
            '9PM'
        ],
        thursday: [
            '9AM',
            '9PM'
        ],
        friday: [
            '9AM',
            '9PM'
        ],
        saturday: [
            '9AM',
            '9PM'
        ]
    },
    timestamp: ISODate('2021-04-12T12:19:12.058Z')
}
*/
/*
	Declarations : Once the shop is created you can not change its name but you can delete your shop.

*/
