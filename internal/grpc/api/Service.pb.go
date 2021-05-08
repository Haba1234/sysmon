// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: Service.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubscriptionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Period *durationpb.Duration `protobuf:"bytes,1,opt,name=period,proto3" json:"period,omitempty"`
	Depth  int64                `protobuf:"varint,2,opt,name=depth,proto3" json:"depth,omitempty"`
}

func (x *SubscriptionRequest) Reset() {
	*x = SubscriptionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscriptionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscriptionRequest) ProtoMessage() {}

func (x *SubscriptionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscriptionRequest.ProtoReflect.Descriptor instead.
func (*SubscriptionRequest) Descriptor() ([]byte, []int) {
	return file_Service_proto_rawDescGZIP(), []int{0}
}

func (x *SubscriptionRequest) GetPeriod() *durationpb.Duration {
	if x != nil {
		return x.Period
	}
	return nil
}

func (x *SubscriptionRequest) GetDepth() int64 {
	if x != nil {
		return x.Depth
	}
	return 0
}

type LoadAverage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OneMin     string `protobuf:"bytes,1,opt,name=oneMin,proto3" json:"oneMin,omitempty"`
	FiveMin    string `protobuf:"bytes,2,opt,name=fiveMin,proto3" json:"fiveMin,omitempty"`
	FifteenMin string `protobuf:"bytes,3,opt,name=fifteenMin,proto3" json:"fifteenMin,omitempty"`
}

func (x *LoadAverage) Reset() {
	*x = LoadAverage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadAverage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadAverage) ProtoMessage() {}

func (x *LoadAverage) ProtoReflect() protoreflect.Message {
	mi := &file_Service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadAverage.ProtoReflect.Descriptor instead.
func (*LoadAverage) Descriptor() ([]byte, []int) {
	return file_Service_proto_rawDescGZIP(), []int{1}
}

func (x *LoadAverage) GetOneMin() string {
	if x != nil {
		return x.OneMin
	}
	return ""
}

func (x *LoadAverage) GetFiveMin() string {
	if x != nil {
		return x.FiveMin
	}
	return ""
}

func (x *LoadAverage) GetFifteenMin() string {
	if x != nil {
		return x.FifteenMin
	}
	return ""
}

type CPUAverage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserMode string `protobuf:"bytes,1,opt,name=userMode,proto3" json:"userMode,omitempty"`
	SysMode  string `protobuf:"bytes,2,opt,name=sysMode,proto3" json:"sysMode,omitempty"`
	Idle     string `protobuf:"bytes,3,opt,name=idle,proto3" json:"idle,omitempty"`
}

func (x *CPUAverage) Reset() {
	*x = CPUAverage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CPUAverage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CPUAverage) ProtoMessage() {}

func (x *CPUAverage) ProtoReflect() protoreflect.Message {
	mi := &file_Service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CPUAverage.ProtoReflect.Descriptor instead.
func (*CPUAverage) Descriptor() ([]byte, []int) {
	return file_Service_proto_rawDescGZIP(), []int{2}
}

func (x *CPUAverage) GetUserMode() string {
	if x != nil {
		return x.UserMode
	}
	return ""
}

func (x *CPUAverage) GetSysMode() string {
	if x != nil {
		return x.SysMode
	}
	return ""
}

func (x *CPUAverage) GetIdle() string {
	if x != nil {
		return x.Idle
	}
	return ""
}

type StatisticsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string       `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	La     *LoadAverage `protobuf:"bytes,2,opt,name=la,proto3" json:"la,omitempty"`
	Cp     *CPUAverage  `protobuf:"bytes,3,opt,name=cp,proto3" json:"cp,omitempty"`
}

func (x *StatisticsResponse) Reset() {
	*x = StatisticsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatisticsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatisticsResponse) ProtoMessage() {}

func (x *StatisticsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatisticsResponse.ProtoReflect.Descriptor instead.
func (*StatisticsResponse) Descriptor() ([]byte, []int) {
	return file_Service_proto_rawDescGZIP(), []int{3}
}

func (x *StatisticsResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *StatisticsResponse) GetLa() *LoadAverage {
	if x != nil {
		return x.La
	}
	return nil
}

func (x *StatisticsResponse) GetCp() *CPUAverage {
	if x != nil {
		return x.Cp
	}
	return nil
}

var File_Service_proto protoreflect.FileDescriptor

var file_Service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5e, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x31, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x70, 0x65, 0x72, 0x69,
	0x6f, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x22, 0x5f, 0x0a, 0x0b, 0x4c, 0x6f, 0x61, 0x64,
	0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x6e, 0x65, 0x4d, 0x69,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x6e, 0x65, 0x4d, 0x69, 0x6e, 0x12,
	0x18, 0x0a, 0x07, 0x66, 0x69, 0x76, 0x65, 0x4d, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x66, 0x69, 0x76, 0x65, 0x4d, 0x69, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x69, 0x66,
	0x74, 0x65, 0x65, 0x6e, 0x4d, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x66,
	0x69, 0x66, 0x74, 0x65, 0x65, 0x6e, 0x4d, 0x69, 0x6e, 0x22, 0x56, 0x0a, 0x0a, 0x43, 0x50, 0x55,
	0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4d,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4d,
	0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79, 0x73, 0x4d, 0x6f, 0x64, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x79, 0x73, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x69, 0x64, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x64, 0x6c,
	0x65, 0x22, 0x77, 0x0a, 0x12, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x24, 0x0a, 0x02, 0x6c, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67,
	0x65, 0x52, 0x02, 0x6c, 0x61, 0x12, 0x23, 0x0a, 0x02, 0x63, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x50, 0x55, 0x41,
	0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x02, 0x63, 0x70, 0x32, 0x5b, 0x0a, 0x0a, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x4d, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x1c, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x12, 0x5a, 0x10, 0x2e, 0x2e, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_Service_proto_rawDescOnce sync.Once
	file_Service_proto_rawDescData = file_Service_proto_rawDesc
)

func file_Service_proto_rawDescGZIP() []byte {
	file_Service_proto_rawDescOnce.Do(func() {
		file_Service_proto_rawDescData = protoimpl.X.CompressGZIP(file_Service_proto_rawDescData)
	})
	return file_Service_proto_rawDescData
}

var file_Service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_Service_proto_goTypes = []interface{}{
	(*SubscriptionRequest)(nil), // 0: service.SubscriptionRequest
	(*LoadAverage)(nil),         // 1: service.LoadAverage
	(*CPUAverage)(nil),          // 2: service.CPUAverage
	(*StatisticsResponse)(nil),  // 3: service.StatisticsResponse
	(*durationpb.Duration)(nil), // 4: google.protobuf.Duration
}
var file_Service_proto_depIdxs = []int32{
	4, // 0: service.SubscriptionRequest.period:type_name -> google.protobuf.Duration
	1, // 1: service.StatisticsResponse.la:type_name -> service.LoadAverage
	2, // 2: service.StatisticsResponse.cp:type_name -> service.CPUAverage
	0, // 3: service.Statistics.ListStatistics:input_type -> service.SubscriptionRequest
	3, // 4: service.Statistics.ListStatistics:output_type -> service.StatisticsResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_Service_proto_init() }
func file_Service_proto_init() {
	if File_Service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscriptionRequest); i {
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
		file_Service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadAverage); i {
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
		file_Service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CPUAverage); i {
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
		file_Service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatisticsResponse); i {
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
			RawDescriptor: file_Service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_Service_proto_goTypes,
		DependencyIndexes: file_Service_proto_depIdxs,
		MessageInfos:      file_Service_proto_msgTypes,
	}.Build()
	File_Service_proto = out.File
	file_Service_proto_rawDesc = nil
	file_Service_proto_goTypes = nil
	file_Service_proto_depIdxs = nil
}
