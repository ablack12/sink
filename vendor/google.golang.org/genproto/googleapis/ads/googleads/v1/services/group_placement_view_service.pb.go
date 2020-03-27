// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/services/group_placement_view_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v1/resources"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Request message for [GroupPlacementViewService.GetGroupPlacementView][google.ads.googleads.v1.services.GroupPlacementViewService.GetGroupPlacementView].
type GetGroupPlacementViewRequest struct {
	// The resource name of the Group Placement view to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetGroupPlacementViewRequest) Reset()         { *m = GetGroupPlacementViewRequest{} }
func (m *GetGroupPlacementViewRequest) String() string { return proto.CompactTextString(m) }
func (*GetGroupPlacementViewRequest) ProtoMessage()    {}
func (*GetGroupPlacementViewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_722a806b0135cef4, []int{0}
}

func (m *GetGroupPlacementViewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetGroupPlacementViewRequest.Unmarshal(m, b)
}
func (m *GetGroupPlacementViewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetGroupPlacementViewRequest.Marshal(b, m, deterministic)
}
func (m *GetGroupPlacementViewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetGroupPlacementViewRequest.Merge(m, src)
}
func (m *GetGroupPlacementViewRequest) XXX_Size() int {
	return xxx_messageInfo_GetGroupPlacementViewRequest.Size(m)
}
func (m *GetGroupPlacementViewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetGroupPlacementViewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetGroupPlacementViewRequest proto.InternalMessageInfo

func (m *GetGroupPlacementViewRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetGroupPlacementViewRequest)(nil), "google.ads.googleads.v1.services.GetGroupPlacementViewRequest")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/services/group_placement_view_service.proto", fileDescriptor_722a806b0135cef4)
}

var fileDescriptor_722a806b0135cef4 = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xcd, 0xaa, 0xd3, 0x40,
	0x14, 0x26, 0x11, 0x04, 0x83, 0x6e, 0x02, 0x62, 0x8d, 0x45, 0x4a, 0xed, 0x42, 0xba, 0x98, 0x21,
	0x4a, 0x11, 0x46, 0x2b, 0xa4, 0x5d, 0xc4, 0x95, 0x94, 0x0a, 0x59, 0x48, 0x20, 0x8c, 0xc9, 0x21,
	0x04, 0x92, 0x99, 0x98, 0x99, 0xa4, 0x0b, 0x71, 0x23, 0x08, 0xee, 0x7d, 0x03, 0x97, 0x3e, 0x4a,
	0xc1, 0x95, 0xaf, 0xe0, 0xca, 0xa7, 0xb8, 0xa4, 0x93, 0x49, 0x6f, 0xe9, 0xcd, 0xed, 0xee, 0x63,
	0xce, 0xf7, 0x33, 0xe7, 0x9b, 0xb1, 0xd6, 0x29, 0xe7, 0x69, 0x0e, 0x98, 0x26, 0x02, 0x2b, 0xd8,
	0xa2, 0xc6, 0xc5, 0x02, 0xaa, 0x26, 0x8b, 0x41, 0xe0, 0xb4, 0xe2, 0x75, 0x19, 0x95, 0x39, 0x8d,
	0xa1, 0x00, 0x26, 0xa3, 0x26, 0x83, 0x5d, 0xd4, 0x4d, 0x51, 0x59, 0x71, 0xc9, 0xed, 0x89, 0x52,
	0x22, 0x9a, 0x08, 0xd4, 0x9b, 0xa0, 0xc6, 0x45, 0xda, 0xc4, 0x79, 0x33, 0x14, 0x53, 0x81, 0xe0,
	0x75, 0x35, 0x94, 0xa3, 0xfc, 0x9d, 0xb1, 0x56, 0x97, 0x19, 0xa6, 0x8c, 0x71, 0x49, 0x65, 0xc6,
	0x99, 0xe8, 0xa6, 0x8f, 0xae, 0x4d, 0xe3, 0x3c, 0x03, 0x26, 0xd5, 0x60, 0xba, 0xb6, 0xc6, 0x3e,
	0x48, 0xbf, 0xf5, 0xdd, 0x68, 0xdb, 0x20, 0x83, 0xdd, 0x16, 0x3e, 0xd7, 0x20, 0xa4, 0xfd, 0xcc,
	0x7a, 0xa0, 0xe3, 0x23, 0x46, 0x0b, 0x18, 0x19, 0x13, 0xe3, 0xf9, 0xbd, 0xed, 0x7d, 0x7d, 0xf8,
	0x9e, 0x16, 0xf0, 0xe2, 0x87, 0x69, 0x3d, 0x3e, 0xb7, 0xf8, 0xa0, 0x16, 0xb3, 0xff, 0x18, 0xd6,
	0xc3, 0x1b, 0x33, 0xec, 0xb7, 0xe8, 0x52, 0x29, 0xe8, 0xb6, 0xcb, 0x39, 0x8b, 0x41, 0x7d, 0x5f,
	0x19, 0x3a, 0x57, 0x4f, 0x97, 0xdf, 0xfe, 0xfe, 0xfb, 0x69, 0xbe, 0xb2, 0x17, 0x6d, 0xb9, 0x5f,
	0x4e, 0xd6, 0x5b, 0xc6, 0xb5, 0x90, 0xbc, 0x80, 0x4a, 0xe0, 0xb9, 0x6a, 0xfb, 0x44, 0x2a, 0xf0,
	0xfc, 0xab, 0xf3, 0x64, 0xef, 0x8d, 0x8e, 0x61, 0x1d, 0x2a, 0x33, 0x81, 0x62, 0x5e, 0xac, 0xbe,
	0x9b, 0xd6, 0x2c, 0xe6, 0xc5, 0xc5, 0xc5, 0x56, 0x4f, 0x07, 0x0b, 0xdb, 0xb4, 0x0f, 0xb3, 0x31,
	0x3e, 0xbe, 0xeb, 0x3c, 0x52, 0x9e, 0x53, 0x96, 0x22, 0x5e, 0xa5, 0x38, 0x05, 0x76, 0x78, 0x36,
	0x7c, 0x4c, 0x1d, 0xfe, 0x95, 0xaf, 0x35, 0xf8, 0x65, 0xde, 0xf1, 0x3d, 0xef, 0xb7, 0x39, 0xf1,
	0x95, 0xa1, 0x97, 0x08, 0xa4, 0x60, 0x8b, 0x02, 0x17, 0x75, 0xc1, 0x62, 0xaf, 0x29, 0xa1, 0x97,
	0x88, 0xb0, 0xa7, 0x84, 0x81, 0x1b, 0x6a, 0xca, 0x7f, 0x73, 0xa6, 0xce, 0x09, 0xf1, 0x12, 0x41,
	0x48, 0x4f, 0x22, 0x24, 0x70, 0x09, 0xd1, 0xb4, 0x4f, 0x77, 0x0f, 0xf7, 0x7c, 0x79, 0x15, 0x00,
	0x00, 0xff, 0xff, 0xb0, 0xfe, 0x1f, 0x2e, 0x3c, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GroupPlacementViewServiceClient is the client API for GroupPlacementViewService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GroupPlacementViewServiceClient interface {
	// Returns the requested Group Placement view in full detail.
	GetGroupPlacementView(ctx context.Context, in *GetGroupPlacementViewRequest, opts ...grpc.CallOption) (*resources.GroupPlacementView, error)
}

type groupPlacementViewServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupPlacementViewServiceClient(cc grpc.ClientConnInterface) GroupPlacementViewServiceClient {
	return &groupPlacementViewServiceClient{cc}
}

func (c *groupPlacementViewServiceClient) GetGroupPlacementView(ctx context.Context, in *GetGroupPlacementViewRequest, opts ...grpc.CallOption) (*resources.GroupPlacementView, error) {
	out := new(resources.GroupPlacementView)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v1.services.GroupPlacementViewService/GetGroupPlacementView", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupPlacementViewServiceServer is the server API for GroupPlacementViewService service.
type GroupPlacementViewServiceServer interface {
	// Returns the requested Group Placement view in full detail.
	GetGroupPlacementView(context.Context, *GetGroupPlacementViewRequest) (*resources.GroupPlacementView, error)
}

// UnimplementedGroupPlacementViewServiceServer can be embedded to have forward compatible implementations.
type UnimplementedGroupPlacementViewServiceServer struct {
}

func (*UnimplementedGroupPlacementViewServiceServer) GetGroupPlacementView(ctx context.Context, req *GetGroupPlacementViewRequest) (*resources.GroupPlacementView, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupPlacementView not implemented")
}

func RegisterGroupPlacementViewServiceServer(s *grpc.Server, srv GroupPlacementViewServiceServer) {
	s.RegisterService(&_GroupPlacementViewService_serviceDesc, srv)
}

func _GroupPlacementViewService_GetGroupPlacementView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupPlacementViewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupPlacementViewServiceServer).GetGroupPlacementView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v1.services.GroupPlacementViewService/GetGroupPlacementView",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupPlacementViewServiceServer).GetGroupPlacementView(ctx, req.(*GetGroupPlacementViewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GroupPlacementViewService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v1.services.GroupPlacementViewService",
	HandlerType: (*GroupPlacementViewServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGroupPlacementView",
			Handler:    _GroupPlacementViewService_GetGroupPlacementView_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v1/services/group_placement_view_service.proto",
}