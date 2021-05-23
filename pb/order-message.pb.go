// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.13.0
// source: order-message.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiKey    string   `protobuf:"bytes,47,opt,name=apiKey,proto3" json:"apiKey,omitempty"`
	Token     string   `protobuf:"bytes,48,opt,name=token,proto3" json:"token,omitempty"`
	ProductId string   `protobuf:"bytes,49,opt,name=productId,proto3" json:"productId,omitempty"`
	ShopId    string   `protobuf:"bytes,50,opt,name=shopId,proto3" json:"shopId,omitempty"`
	Address   *Address `protobuf:"bytes,51,opt,name=address,proto3" json:"address,omitempty"`
	Quantity  uint32   `protobuf:"varint,52,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *CreateOrderRequest) Reset() {
	*x = CreateOrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderRequest) ProtoMessage() {}

func (x *CreateOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_order_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderRequest.ProtoReflect.Descriptor instead.
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return file_order_message_proto_rawDescGZIP(), []int{0}
}

func (x *CreateOrderRequest) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *CreateOrderRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *CreateOrderRequest) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *CreateOrderRequest) GetShopId() string {
	if x != nil {
		return x.ShopId
	}
	return ""
}

func (x *CreateOrderRequest) GetAddress() *Address {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *CreateOrderRequest) GetQuantity() uint32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type CreateOrderResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId      string  `protobuf:"bytes,53,opt,name=orderId,proto3" json:"orderId,omitempty"`
	Currency     string  `protobuf:"bytes,54,opt,name=currency,proto3" json:"currency,omitempty"`
	Amount       float32 `protobuf:"fixed32,55,opt,name=amount,proto3" json:"amount,omitempty"`
	ProductName  string  `protobuf:"bytes,56,opt,name=productName,proto3" json:"productName,omitempty"`
	ProductId    string  `protobuf:"bytes,57,opt,name=productId,proto3" json:"productId,omitempty"`
	ShopId       string  `protobuf:"bytes,58,opt,name=shopId,proto3" json:"shopId,omitempty"`
	ProductImage string  `protobuf:"bytes,59,opt,name=productImage,proto3" json:"productImage,omitempty"`
}

func (x *CreateOrderResponse) Reset() {
	*x = CreateOrderResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderResponse) ProtoMessage() {}

func (x *CreateOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_order_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderResponse.ProtoReflect.Descriptor instead.
func (*CreateOrderResponse) Descriptor() ([]byte, []int) {
	return file_order_message_proto_rawDescGZIP(), []int{1}
}

func (x *CreateOrderResponse) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *CreateOrderResponse) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *CreateOrderResponse) GetAmount() float32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *CreateOrderResponse) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *CreateOrderResponse) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *CreateOrderResponse) GetShopId() string {
	if x != nil {
		return x.ShopId
	}
	return ""
}

func (x *CreateOrderResponse) GetProductImage() string {
	if x != nil {
		return x.ProductImage
	}
	return ""
}

var File_order_message_proto protoreflect.FileDescriptor

var file_order_message_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2d, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb8, 0x01, 0x0a, 0x12,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x18, 0x2f, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x30, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x31, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x64, 0x18, 0x32, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x68, 0x6f, 0x70, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x33, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75,
	0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x34, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x71, 0x75,
	0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0xdf, 0x01, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x35, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x18, 0x36, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x37,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x38, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x39, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x68, 0x6f, 0x70, 0x49, 0x64, 0x18, 0x3a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x68,
	0x6f, 0x70, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x18, 0x3b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x42, 0x30, 0x0a, 0x26, 0x63, 0x6f, 0x6d, 0x2e,
	0x61, 0x61, 0x70, 0x61, 0x6e, 0x61, 0x76, 0x79, 0x61, 0x70, 0x61, 0x72, 0x2e, 0x61, 0x61, 0x70,
	0x61, 0x6e, 0x61, 0x76, 0x79, 0x61, 0x70, 0x61, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x50, 0x01, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_order_message_proto_rawDescOnce sync.Once
	file_order_message_proto_rawDescData = file_order_message_proto_rawDesc
)

func file_order_message_proto_rawDescGZIP() []byte {
	file_order_message_proto_rawDescOnce.Do(func() {
		file_order_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_order_message_proto_rawDescData)
	})
	return file_order_message_proto_rawDescData
}

var file_order_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_order_message_proto_goTypes = []interface{}{
	(*CreateOrderRequest)(nil),  // 0: CreateOrderRequest
	(*CreateOrderResponse)(nil), // 1: CreateOrderResponse
	(*Address)(nil),             // 2: Address
}
var file_order_message_proto_depIdxs = []int32{
	2, // 0: CreateOrderRequest.address:type_name -> Address
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_order_message_proto_init() }
func file_order_message_proto_init() {
	if File_order_message_proto != nil {
		return
	}
	file_common_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_order_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateOrderRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_order_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateOrderResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_order_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_order_message_proto_goTypes,
		DependencyIndexes: file_order_message_proto_depIdxs,
		MessageInfos:      file_order_message_proto_msgTypes,
	}.Build()
	File_order_message_proto = out.File
	file_order_message_proto_rawDesc = nil
	file_order_message_proto_goTypes = nil
	file_order_message_proto_depIdxs = nil
}
