// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: pkgs/proto/submission.proto

package pkgs

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Submission_SubmitSnapshotSimulation_FullMethodName = "/submission.Submission/SubmitSnapshotSimulation"
	Submission_SubmitSnapshot_FullMethodName           = "/submission.Submission/SubmitSnapshot"
)

// SubmissionClient is the client API for Submission service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubmissionClient interface {
	SubmitSnapshotSimulation(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[SnapshotSubmission, SubmissionResponse], error)
	SubmitSnapshot(ctx context.Context, in *SnapshotSubmission, opts ...grpc.CallOption) (*SubmissionResponse, error)
}

type submissionClient struct {
	cc grpc.ClientConnInterface
}

func NewSubmissionClient(cc grpc.ClientConnInterface) SubmissionClient {
	return &submissionClient{cc}
}

func (c *submissionClient) SubmitSnapshotSimulation(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[SnapshotSubmission, SubmissionResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Submission_ServiceDesc.Streams[0], Submission_SubmitSnapshotSimulation_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SnapshotSubmission, SubmissionResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Submission_SubmitSnapshotSimulationClient = grpc.ClientStreamingClient[SnapshotSubmission, SubmissionResponse]

func (c *submissionClient) SubmitSnapshot(ctx context.Context, in *SnapshotSubmission, opts ...grpc.CallOption) (*SubmissionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SubmissionResponse)
	err := c.cc.Invoke(ctx, Submission_SubmitSnapshot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubmissionServer is the server API for Submission service.
// All implementations must embed UnimplementedSubmissionServer
// for forward compatibility.
type SubmissionServer interface {
	SubmitSnapshotSimulation(grpc.ClientStreamingServer[SnapshotSubmission, SubmissionResponse]) error
	SubmitSnapshot(context.Context, *SnapshotSubmission) (*SubmissionResponse, error)
	mustEmbedUnimplementedSubmissionServer()
}

// UnimplementedSubmissionServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSubmissionServer struct{}

func (UnimplementedSubmissionServer) SubmitSnapshotSimulation(grpc.ClientStreamingServer[SnapshotSubmission, SubmissionResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SubmitSnapshotSimulation not implemented")
}
func (UnimplementedSubmissionServer) SubmitSnapshot(context.Context, *SnapshotSubmission) (*SubmissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitSnapshot not implemented")
}
func (UnimplementedSubmissionServer) mustEmbedUnimplementedSubmissionServer() {}
func (UnimplementedSubmissionServer) testEmbeddedByValue()                    {}

// UnsafeSubmissionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubmissionServer will
// result in compilation errors.
type UnsafeSubmissionServer interface {
	mustEmbedUnimplementedSubmissionServer()
}

func RegisterSubmissionServer(s grpc.ServiceRegistrar, srv SubmissionServer) {
	// If the following call pancis, it indicates UnimplementedSubmissionServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Submission_ServiceDesc, srv)
}

func _Submission_SubmitSnapshotSimulation_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SubmissionServer).SubmitSnapshotSimulation(&grpc.GenericServerStream[SnapshotSubmission, SubmissionResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Submission_SubmitSnapshotSimulationServer = grpc.ClientStreamingServer[SnapshotSubmission, SubmissionResponse]

func _Submission_SubmitSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SnapshotSubmission)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubmissionServer).SubmitSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Submission_SubmitSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubmissionServer).SubmitSnapshot(ctx, req.(*SnapshotSubmission))
	}
	return interceptor(ctx, in, info, handler)
}

// Submission_ServiceDesc is the grpc.ServiceDesc for Submission service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Submission_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "submission.Submission",
	HandlerType: (*SubmissionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitSnapshot",
			Handler:    _Submission_SubmitSnapshot_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubmitSnapshotSimulation",
			Handler:       _Submission_SubmitSnapshotSimulation_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pkgs/proto/submission.proto",
}
