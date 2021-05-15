package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func InitMongo() *mongo.Client {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	credential := options.Credential{
		Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetAuth(credential).SetMinPoolSize(100).SetMaxPoolSize(200))
	if err != nil {
		panic(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, nil)

	if err != nil {
		panic(err)
	}

	return client

}

func OpenDefaultDataCollection(client *mongo.Client) *mongo.Collection {
	database := client.Database("db_aapanavypar")
	defaultData := database.Collection("defaultData")
	return defaultData
}

func OpenUserDataCollection(client *mongo.Client) *mongo.Collection {
	database := client.Database("db_aapanavypar")
	userData := database.Collection("userData")
	return userData
}

func OpenOrderDataCollection(client *mongo.Client) *mongo.Collection {
	database := client.Database("db_aapanavypar")
	orderData := database.Collection("orderData")
	return orderData
}

func OpenShopDataCollection(client *mongo.Client) *mongo.Collection {
	database := client.Database("db_aapanavypar")
	shopData := database.Collection("shopData")
	return shopData
}

func OpenProductDataCollection(client *mongo.Client) *mongo.Collection {
	database := client.Database("db_aapanavypar")
	productData := database.Collection("productData")
	return productData
}

func OpenAnalyticalDataCollection(client *mongo.Client) *mongo.Collection {
	database := client.Database("db_aapanavypar")
	analyticalData := database.Collection("analyticalData")
	return analyticalData
}
