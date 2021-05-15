package helpers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/o1egl/paseto/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"
)

const External = 7

func IsHasAccessTo(accessGroup []int, access int) bool {
	for _, a := range accessGroup {
		if a == access {
			return true
		}
	}
	return false
}

func CheckForAPIKey(toTestKey string) bool {
	if toTestKey == os.Getenv("API_KEY_FOR_MOBILE") || toTestKey == os.Getenv("API_KEY_FOR_WEB") {
		return true
	}
	return false
}

func Validate(s interface{}) error {

	validate := validator.New()

	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		myErr := ""
		for _, err := range err.(validator.ValidationErrors) {
			myErr += " " + err.Field()
		}

		fmt.Println(err)
		return fmt.Errorf("Error in validating :" + myErr)
	}

	return nil
}

func ValidateToken(ctx context.Context, tokenString, key string, access int) (*paseto.JSONToken, error) {
	var receivedToken paseto.JSONToken
	var newFooter string

	fmt.Println("Decrypting Token")
	err := paseto.Decrypt(tokenString, []byte(key), &receivedToken, &newFooter)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Token %v", err)
	}

	fmt.Println("Validating Token")
	err = receivedToken.Validate(
		paseto.ValidAt(time.Now()),
		paseto.IssuedBy(os.Getenv("TOKEN_ISSUER")),
	)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid Token %v", err)
	}

	fmt.Println("Checking Access Group For Token")
	var accessGroup []int
	err = receivedToken.Get("accessGroup", &accessGroup)
	fmt.Println("Access Group  : ", accessGroup)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Token Is Not Valid %v", err)
	}

	if !IsHasAccessTo(accessGroup, access) {
		return nil, status.Errorf(codes.PermissionDenied, "Token Does Not Have Valid Permission To Access Resources %v", err)
	}

	fmt.Println("Checking Refresh Token ID Token")
	_, err = uuid.Parse(receivedToken.Jti)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument ", err)
	}

	fmt.Println("Checking User UUID Of Token")
	_, err = uuid.Parse(receivedToken.Audience)
	if err != nil {
		fmt.Println("Checking User UUID Of Token : ", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument ", err)
	}

	fmt.Println("Done With Token")
	return &receivedToken, nil
}

func ContextError(ctx context.Context) error {

	switch ctx.Err() {
	case context.Canceled:
		log.Println("Request Canceled")
		return status.Error(codes.DeadlineExceeded, "Request Canceled")
	case context.DeadlineExceeded:
		log.Println("DeadLine Exceeded")
		return status.Error(codes.DeadlineExceeded, "DeadLine Exceeded")
	default:
		return nil
	}
}

func VerifySignature(toCheckMessage, originalMessage, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(toCheckMessage))
	expectedMAC := mac.Sum(nil)

	if hex.EncodeToString(expectedMAC) == originalMessage {
		return true
	}

	return false
}
