module aapanavyapar-service-buying

go 1.15

require (
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.8.0
	github.com/google/uuid v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/o1egl/paseto/v2 v2.1.1
	github.com/razorpay/razorpay-go v0.0.0-20201204135735-096d3be7d2df
	go.mongodb.org/mongo-driver v1.5.0
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.25.0

)

replace github.com/razorpay/razorpay-go => ../razorpay-go

