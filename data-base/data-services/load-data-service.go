package data_base

import (
	"aapanavyapar-service-buying/configurations/mongodb"
	"aapanavyapar-service-buying/data-base/structs"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func (dataBase *MongoDataBase) GetAllBasicCategories(context context.Context, sendData func(data structs.BasicCategoriesData) error) error {

	defaultData := mongodb.OpenDefaultDataCollection(dataBase.Data)

	cursor, err := defaultData.Find(context, bson.D{})

	if err != nil {
		return err
	}
	defer cursor.Close(context)

	for cursor.Next(context) {
		result := structs.BasicCategoriesData{}
		err = cursor.Decode(&result)

		fmt.Println(result.Category)
		fmt.Println(result.SubCategories)

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

/*
{
	_id: ObjectId(),
	category: "food",
	sub_categories:[
		"",
		"",
		""
	]
}
*/
