// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: grpc/proto.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Auction_Bidding_FullMethodName      = "/AuctionSystem.Auction/Bidding"
	Auction_AskForResult_FullMethodName = "/AuctionSystem.Auction/AskForResult"
)

// AuctionClient is the client API for Auction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuctionClient interface {
	Bidding(ctx context.Context, in *BidAmount, opts ...grpc.CallOption) (*BidSuccess, error)
	AskForResult(ctx context.Context, in *ResultRequest, opts ...grpc.CallOption) (*Result, error)
}

type auctionClient struct {
	cc grpc.ClientConnInterface
}

func NewAuctionClient(cc grpc.ClientConnInterface) AuctionClient {
	return &auctionClient{cc}
}

func (c *auctionClient) Bidding(ctx context.Context, in *BidAmount, opts ...grpc.CallOption) (*BidSuccess, error) {
	out := new(BidSuccess)
	err := c.cc.Invoke(ctx, Auction_Bidding_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) AskForResult(ctx context.Context, in *ResultRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, Auction_AskForResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuctionServer is the server API for Auction service.
// All implementations must embed UnimplementedAuctionServer
// for forward compatibility
type AuctionServer interface {
	Bidding(context.Context, *BidAmount) (*BidSuccess, error)
	AskForResult(context.Context, *ResultRequest) (*Result, error)
	mustEmbedUnimplementedAuctionServer()
}

// UnimplementedAuctionServer must be embedded to have forward compatible implementations.
type UnimplementedAuctionServer struct {
}

func (UnimplementedAuctionServer) Bidding(context.Context, *BidAmount) (*BidSuccess, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bidding not implemented")
}
func (UnimplementedAuctionServer) AskForResult(context.Context, *ResultRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AskForResult not implemented")
}
func (UnimplementedAuctionServer) mustEmbedUnimplementedAuctionServer() {}

// UnsafeAuctionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuctionServer will
// result in compilation errors.
type UnsafeAuctionServer interface {
	mustEmbedUnimplementedAuctionServer()
}

func RegisterAuctionServer(s grpc.ServiceRegistrar, srv AuctionServer) {
	s.RegisterService(&Auction_ServiceDesc, srv)
}

func _Auction_Bidding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidAmount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).Bidding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auction_Bidding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).Bidding(ctx, req.(*BidAmount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_AskForResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).AskForResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auction_AskForResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).AskForResult(ctx, req.(*ResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Auction_ServiceDesc is the grpc.ServiceDesc for Auction service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auction_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AuctionSystem.Auction",
	HandlerType: (*AuctionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bidding",
			Handler:    _Auction_Bidding_Handler,
		},
		{
			MethodName: "AskForResult",
			Handler:    _Auction_AskForResult_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/proto.proto",
}
