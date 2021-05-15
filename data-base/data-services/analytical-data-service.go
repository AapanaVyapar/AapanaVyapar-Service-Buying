package data_base

import (
	"aapanavyapar-service-buying/configurations/mongodb"
	"aapanavyapar-service-buying/data-base/structs"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func (dataBase *MongoDataBase) AddAnalyticalDataToAnalyticalData(context context.Context, userId string, data structs.AnalyticalClickData) error {

	analyticalData := mongodb.OpenAnalyticalDataCollection(dataBase.Data)

	dataBase.mutex.Lock()
	defer dataBase.mutex.Unlock()

	result := analyticalData.FindOne(context, bson.M{"_id": userId})

	if result.Err() != nil {
		_, err := analyticalData.InsertOne(context,
			bson.M{
				"_id": userId,
				"most_visited": bson.M{
					"product": bson.A{data},
				},
			},
		)
		if err != nil {
			return err
		}

		return nil
	}

	res, err := analyticalData.UpdateOne(context,
		bson.M{
			"_id": userId,
		},
		bson.M{
			"$push": bson.M{
				"most_visited.product": bson.M{
					"$each":  bson.A{data},
					"$slice": -5,
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

	return fmt.Errorf("unable to add to analytical data")
}

/*
AnalyticalData =
{
	_Id:      "",
	MostVisited: {
		Products: [
			{
				productId: primitive.ObjectID,
				timestamp: time.Time,
				category: []constants.Category{}
			}
		]
	}
}

*/
