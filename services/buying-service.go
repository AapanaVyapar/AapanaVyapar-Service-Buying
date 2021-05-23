package services

import (
	redisDataBase "aapanavyapar-service-buying/data-base/cash-services"
	mongoDataBase "aapanavyapar-service-buying/data-base/data-services"
	"aapanavyapar-service-buying/data-base/helpers"
	"aapanavyapar-service-buying/data-base/mapper"
	"aapanavyapar-service-buying/data-base/structs"
	"aapanavyapar-service-buying/pb"
	"context"
	"github.com/razorpay/razorpay-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type BuyingService struct {
	Data *mongoDataBase.MongoDataBase
	Cash *redisDataBase.CashDataBase
}

func NewBuyingService() *BuyingService {
	mongoData := mongoDataBase.NewDataBase()
	redisData := redisDataBase.NewDataBase()

	order := BuyingService{
		Data: mongoData,
		Cash: redisData,
	}

	return &order
}

func (buyingService *BuyingService) PlaceOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	if !helpers.CheckForAPIKey(request.GetApiKey()) {
		return nil, status.Errorf(codes.Unauthenticated, "No API Key Is Specified")
	}

	token, err := helpers.ValidateToken(ctx, request.GetToken(), os.Getenv("AUTH_TOKEN_SECRETE"), helpers.External)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request With Invalid Token : %v", err)
	}

	if request.GetQuantity() <= 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Quantity Can Not Be Zero : %v", err)
	}

	address := structs.Address{
		FullName:      request.GetAddress().GetFullName(),
		HouseDetails:  request.GetAddress().GetHouseDetails(),
		StreetDetails: request.GetAddress().GetStreetDetails(),
		LandMark:      request.GetAddress().GetLandMark(),
		PinCode:       request.GetAddress().GetPinCode(),
		City:          request.GetAddress().GetCity(),
		State:         request.GetAddress().GetState(),
		Country:       request.GetAddress().GetCountry(),
		PhoneNo:       request.GetAddress().GetPhoneNo(),
	}

	err = helpers.Validate(address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Address : ", err)
	}

	productId, err := primitive.ObjectIDFromHex(request.GetProductId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Product Id : %v", err)
	}

	product, err := buyingService.Data.GetSpecificProductsOfShopFromProductData(ctx, request.GetShopId(), productId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Product Not Found : %v", err)
	}

	if product.Stock < 1 {
		return nil, status.Errorf(codes.ResourceExhausted, "Product No Longer Available")
	}

	client := razorpay.NewClient(os.Getenv("APT_KEY_RAZORPAY"), os.Getenv("SECRET_KEY_RAZORPAY"))

	distance := 100
	deliveryCost := mapper.CalculateDeliveryCost(distance, &address)

	price := product.Price
	price = price - ((price / 100) * float32(product.Offer))
	price += deliveryCost

	deliveryTimeStamp := mapper.CalculateDeliveryTime(distance)

	data := map[string]interface{}{
		"amount":          price,
		"currency":        "INR",
		"receipt":         "some_receipt_id",
		"payment_capture": 0, // Capture Payment Manually = 0 So That hole cycle of payment is successful
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error While Creating Order : %v", err)
	}

	cacheData := structs.CacheOrder{
		Address:           address,
		Id:                token.Audience,
		Quantity:          request.GetQuantity(),
		ProductId:         productId,
		DeliveryCost:      deliveryCost,
		Offer:             product.Offer,
		DeliveryTimeStamp: deliveryTimeStamp,
	}

	err = buyingService.Cash.SetDataToCash(ctx, body["id"].(string), cacheData.Marshal(), time.Hour*12) //Set Auto Refund Dead Line in Manual Payment Capture to less than 12 hour.
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable To Add Order To Cache : %v", err)
	}

	return &pb.CreateOrderResponse{
		OrderId:      body["id"].(string),
		Currency:     body["currency"].(string),
		Amount:       float32(body["amount"].(float64)),
		ProductName:  product.Title,
		ProductImage: product.Images[0],
		ProductId:    request.GetProductId(),
		ShopId:       request.GetShopId(),
	}, nil
}

func (buyingService *BuyingService) CapturePayment(ctx context.Context, request *pb.CapturePaymentRequest) (*pb.CapturePaymentResponse, error) {
	if !helpers.CheckForAPIKey(request.GetApiKey()) {
		return nil, status.Errorf(codes.Unauthenticated, "No API Key Is Specified")
	}

	token, err := helpers.ValidateToken(ctx, request.GetToken(), os.Getenv("AUTH_TOKEN_SECRETE"), helpers.External)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Request With Invalid Token : %v", err)
	}

	paymentId := request.GetRazorpayPaymentId()
	orderId := request.GetRazorpayOrderId()

	bytes, err := buyingService.Cash.GetDataFromCash(ctx, orderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable To Get Order From Cache : %v", err)
	}

	var cache = structs.CacheOrder{}
	structs.UnmarshalCacheOrder([]byte(bytes), &cache)

	if token.Audience != cache.Id {
		return nil, status.Errorf(codes.Unauthenticated, "Order Does Not Belongs To You")
	}

	client := razorpay.NewClient(os.Getenv("APT_KEY_RAZORPAY"), os.Getenv("SECRET_KEY_RAZORPAY"))
	payment, err := client.Payment.Fetch(paymentId, nil, nil)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Error of Payment :  %v ", err)
	}

	if payment["status"].(string) != "authorized" {
		return nil, status.Errorf(codes.Unknown, "Invalid Payment Status : "+payment["status"].(string))
	}

	oId, err := buyingService.Data.CreateOrder(ctx, token.Audience, paymentId, cache, payment, orderId, client)
	if err != nil {
		return nil, err
	}

	return &pb.CapturePaymentResponse{
		Status:  true,
		OrderId: oId.Hex(),
	}, nil

}
