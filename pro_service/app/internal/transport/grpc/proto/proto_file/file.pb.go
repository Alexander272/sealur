// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/file.proto

/*
Package proto_file is a generated protocol buffer package.

It is generated from these files:
	proto/file.proto

It has these top-level messages:
	MetaData
	File
	FileUploadRequest
	FileUploadResponse
	FileDownloadRequest
	GroupDownloadRequest
	FileDownloadResponse
	CopyFileRequest
	CopyGroupRequest
	FileDeleteRequest
	GroupDeleteRequest
	MessageResponse
	PingRequest
	PingResponse
*/
package proto_file

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MetaData struct {
	Name   string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Type   string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Size   int64  `protobuf:"varint,3,opt,name=size" json:"size,omitempty"`
	Group  string `protobuf:"bytes,4,opt,name=group" json:"group,omitempty"`
	Bucket string `protobuf:"bytes,5,opt,name=bucket" json:"bucket,omitempty"`
}

func (m *MetaData) Reset()                    { *m = MetaData{} }
func (m *MetaData) String() string            { return proto.CompactTextString(m) }
func (*MetaData) ProtoMessage()               {}
func (*MetaData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MetaData) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaData) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *MetaData) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *MetaData) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *MetaData) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

type File struct {
	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *File) Reset()                    { *m = File{} }
func (m *File) String() string            { return proto.CompactTextString(m) }
func (*File) ProtoMessage()               {}
func (*File) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *File) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type FileUploadRequest struct {
	// Types that are valid to be assigned to Request:
	//	*FileUploadRequest_Metadata
	//	*FileUploadRequest_File
	Request isFileUploadRequest_Request `protobuf_oneof:"request"`
}

func (m *FileUploadRequest) Reset()                    { *m = FileUploadRequest{} }
func (m *FileUploadRequest) String() string            { return proto.CompactTextString(m) }
func (*FileUploadRequest) ProtoMessage()               {}
func (*FileUploadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isFileUploadRequest_Request interface{ isFileUploadRequest_Request() }

type FileUploadRequest_Metadata struct {
	Metadata *MetaData `protobuf:"bytes,2,opt,name=metadata,oneof"`
}
type FileUploadRequest_File struct {
	File *File `protobuf:"bytes,1,opt,name=file,oneof"`
}

func (*FileUploadRequest_Metadata) isFileUploadRequest_Request() {}
func (*FileUploadRequest_File) isFileUploadRequest_Request()     {}

func (m *FileUploadRequest) GetRequest() isFileUploadRequest_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *FileUploadRequest) GetMetadata() *MetaData {
	if x, ok := m.GetRequest().(*FileUploadRequest_Metadata); ok {
		return x.Metadata
	}
	return nil
}

func (m *FileUploadRequest) GetFile() *File {
	if x, ok := m.GetRequest().(*FileUploadRequest_File); ok {
		return x.File
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FileUploadRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FileUploadRequest_OneofMarshaler, _FileUploadRequest_OneofUnmarshaler, _FileUploadRequest_OneofSizer, []interface{}{
		(*FileUploadRequest_Metadata)(nil),
		(*FileUploadRequest_File)(nil),
	}
}

func _FileUploadRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FileUploadRequest)
	// request
	switch x := m.Request.(type) {
	case *FileUploadRequest_Metadata:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Metadata); err != nil {
			return err
		}
	case *FileUploadRequest_File:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.File); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("FileUploadRequest.Request has unexpected type %T", x)
	}
	return nil
}

func _FileUploadRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FileUploadRequest)
	switch tag {
	case 2: // request.metadata
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MetaData)
		err := b.DecodeMessage(msg)
		m.Request = &FileUploadRequest_Metadata{msg}
		return true, err
	case 1: // request.file
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(File)
		err := b.DecodeMessage(msg)
		m.Request = &FileUploadRequest_File{msg}
		return true, err
	default:
		return false, nil
	}
}

func _FileUploadRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FileUploadRequest)
	// request
	switch x := m.Request.(type) {
	case *FileUploadRequest_Metadata:
		s := proto.Size(x.Metadata)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *FileUploadRequest_File:
		s := proto.Size(x.File)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type FileUploadResponse struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	OrigName string `protobuf:"bytes,3,opt,name=origName" json:"origName,omitempty"`
	Url      string `protobuf:"bytes,4,opt,name=url" json:"url,omitempty"`
}

func (m *FileUploadResponse) Reset()                    { *m = FileUploadResponse{} }
func (m *FileUploadResponse) String() string            { return proto.CompactTextString(m) }
func (*FileUploadResponse) ProtoMessage()               {}
func (*FileUploadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *FileUploadResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *FileUploadResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FileUploadResponse) GetOrigName() string {
	if m != nil {
		return m.OrigName
	}
	return ""
}

func (m *FileUploadResponse) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type FileDownloadRequest struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Bucket string `protobuf:"bytes,2,opt,name=bucket" json:"bucket,omitempty"`
	Group  string `protobuf:"bytes,3,opt,name=group" json:"group,omitempty"`
	Name   string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
}

func (m *FileDownloadRequest) Reset()                    { *m = FileDownloadRequest{} }
func (m *FileDownloadRequest) String() string            { return proto.CompactTextString(m) }
func (*FileDownloadRequest) ProtoMessage()               {}
func (*FileDownloadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *FileDownloadRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *FileDownloadRequest) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *FileDownloadRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *FileDownloadRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GroupDownloadRequest struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket" json:"bucket,omitempty"`
	Group  string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
}

func (m *GroupDownloadRequest) Reset()                    { *m = GroupDownloadRequest{} }
func (m *GroupDownloadRequest) String() string            { return proto.CompactTextString(m) }
func (*GroupDownloadRequest) ProtoMessage()               {}
func (*GroupDownloadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GroupDownloadRequest) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *GroupDownloadRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

type FileDownloadResponse struct {
	// Types that are valid to be assigned to Response:
	//	*FileDownloadResponse_Metadata
	//	*FileDownloadResponse_File
	Response isFileDownloadResponse_Response `protobuf_oneof:"response"`
}

func (m *FileDownloadResponse) Reset()                    { *m = FileDownloadResponse{} }
func (m *FileDownloadResponse) String() string            { return proto.CompactTextString(m) }
func (*FileDownloadResponse) ProtoMessage()               {}
func (*FileDownloadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type isFileDownloadResponse_Response interface{ isFileDownloadResponse_Response() }

type FileDownloadResponse_Metadata struct {
	Metadata *MetaData `protobuf:"bytes,1,opt,name=metadata,oneof"`
}
type FileDownloadResponse_File struct {
	File *File `protobuf:"bytes,2,opt,name=file,oneof"`
}

func (*FileDownloadResponse_Metadata) isFileDownloadResponse_Response() {}
func (*FileDownloadResponse_File) isFileDownloadResponse_Response()     {}

func (m *FileDownloadResponse) GetResponse() isFileDownloadResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *FileDownloadResponse) GetMetadata() *MetaData {
	if x, ok := m.GetResponse().(*FileDownloadResponse_Metadata); ok {
		return x.Metadata
	}
	return nil
}

func (m *FileDownloadResponse) GetFile() *File {
	if x, ok := m.GetResponse().(*FileDownloadResponse_File); ok {
		return x.File
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FileDownloadResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FileDownloadResponse_OneofMarshaler, _FileDownloadResponse_OneofUnmarshaler, _FileDownloadResponse_OneofSizer, []interface{}{
		(*FileDownloadResponse_Metadata)(nil),
		(*FileDownloadResponse_File)(nil),
	}
}

func _FileDownloadResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FileDownloadResponse)
	// response
	switch x := m.Response.(type) {
	case *FileDownloadResponse_Metadata:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Metadata); err != nil {
			return err
		}
	case *FileDownloadResponse_File:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.File); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("FileDownloadResponse.Response has unexpected type %T", x)
	}
	return nil
}

func _FileDownloadResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FileDownloadResponse)
	switch tag {
	case 1: // response.metadata
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MetaData)
		err := b.DecodeMessage(msg)
		m.Response = &FileDownloadResponse_Metadata{msg}
		return true, err
	case 2: // response.file
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(File)
		err := b.DecodeMessage(msg)
		m.Response = &FileDownloadResponse_File{msg}
		return true, err
	default:
		return false, nil
	}
}

func _FileDownloadResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FileDownloadResponse)
	// response
	switch x := m.Response.(type) {
	case *FileDownloadResponse_Metadata:
		s := proto.Size(x.Metadata)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *FileDownloadResponse_File:
		s := proto.Size(x.File)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type CopyFileRequest struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Bucket   string `protobuf:"bytes,2,opt,name=bucket" json:"bucket,omitempty"`
	Group    string `protobuf:"bytes,3,opt,name=group" json:"group,omitempty"`
	NewGroup string `protobuf:"bytes,4,opt,name=newGroup" json:"newGroup,omitempty"`
}

func (m *CopyFileRequest) Reset()                    { *m = CopyFileRequest{} }
func (m *CopyFileRequest) String() string            { return proto.CompactTextString(m) }
func (*CopyFileRequest) ProtoMessage()               {}
func (*CopyFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CopyFileRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CopyFileRequest) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *CopyFileRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *CopyFileRequest) GetNewGroup() string {
	if m != nil {
		return m.NewGroup
	}
	return ""
}

type CopyGroupRequest struct {
	Bucket   string `protobuf:"bytes,1,opt,name=bucket" json:"bucket,omitempty"`
	Group    string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
	NewGroup string `protobuf:"bytes,3,opt,name=newGroup" json:"newGroup,omitempty"`
}

func (m *CopyGroupRequest) Reset()                    { *m = CopyGroupRequest{} }
func (m *CopyGroupRequest) String() string            { return proto.CompactTextString(m) }
func (*CopyGroupRequest) ProtoMessage()               {}
func (*CopyGroupRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *CopyGroupRequest) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *CopyGroupRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *CopyGroupRequest) GetNewGroup() string {
	if m != nil {
		return m.NewGroup
	}
	return ""
}

type FileDeleteRequest struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Bucket string `protobuf:"bytes,2,opt,name=bucket" json:"bucket,omitempty"`
	Group  string `protobuf:"bytes,3,opt,name=group" json:"group,omitempty"`
	Name   string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
}

func (m *FileDeleteRequest) Reset()                    { *m = FileDeleteRequest{} }
func (m *FileDeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*FileDeleteRequest) ProtoMessage()               {}
func (*FileDeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *FileDeleteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *FileDeleteRequest) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *FileDeleteRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *FileDeleteRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GroupDeleteRequest struct {
	Bucket string `protobuf:"bytes,1,opt,name=bucket" json:"bucket,omitempty"`
	Group  string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
}

func (m *GroupDeleteRequest) Reset()                    { *m = GroupDeleteRequest{} }
func (m *GroupDeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*GroupDeleteRequest) ProtoMessage()               {}
func (*GroupDeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *GroupDeleteRequest) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *GroupDeleteRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

type MessageResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *MessageResponse) Reset()                    { *m = MessageResponse{} }
func (m *MessageResponse) String() string            { return proto.CompactTextString(m) }
func (*MessageResponse) ProtoMessage()               {}
func (*MessageResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *MessageResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type PingRequest struct {
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type PingResponse struct {
	Ping string `protobuf:"bytes,1,opt,name=ping" json:"ping,omitempty"`
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *PingResponse) GetPing() string {
	if m != nil {
		return m.Ping
	}
	return ""
}

func init() {
	proto.RegisterType((*MetaData)(nil), "proto_file.MetaData")
	proto.RegisterType((*File)(nil), "proto_file.File")
	proto.RegisterType((*FileUploadRequest)(nil), "proto_file.FileUploadRequest")
	proto.RegisterType((*FileUploadResponse)(nil), "proto_file.FileUploadResponse")
	proto.RegisterType((*FileDownloadRequest)(nil), "proto_file.FileDownloadRequest")
	proto.RegisterType((*GroupDownloadRequest)(nil), "proto_file.GroupDownloadRequest")
	proto.RegisterType((*FileDownloadResponse)(nil), "proto_file.FileDownloadResponse")
	proto.RegisterType((*CopyFileRequest)(nil), "proto_file.CopyFileRequest")
	proto.RegisterType((*CopyGroupRequest)(nil), "proto_file.CopyGroupRequest")
	proto.RegisterType((*FileDeleteRequest)(nil), "proto_file.FileDeleteRequest")
	proto.RegisterType((*GroupDeleteRequest)(nil), "proto_file.GroupDeleteRequest")
	proto.RegisterType((*MessageResponse)(nil), "proto_file.MessageResponse")
	proto.RegisterType((*PingRequest)(nil), "proto_file.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "proto_file.PingResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for FileService service

type FileServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (FileService_DownloadClient, error)
	GroupDownload(ctx context.Context, in *GroupDownloadRequest, opts ...grpc.CallOption) (FileService_GroupDownloadClient, error)
	Upload(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadClient, error)
	Copy(ctx context.Context, in *CopyFileRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	CopyGroup(ctx context.Context, in *CopyGroupRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	Delete(ctx context.Context, in *FileDeleteRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	GroupDelete(ctx context.Context, in *GroupDeleteRequest, opts ...grpc.CallOption) (*MessageResponse, error)
}

type fileServiceClient struct {
	cc *grpc.ClientConn
}

func NewFileServiceClient(cc *grpc.ClientConn) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/proto_file.FileService/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (FileService_DownloadClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_FileService_serviceDesc.Streams[0], c.cc, "/proto_file.FileService/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileService_DownloadClient interface {
	Recv() (*FileDownloadResponse, error)
	grpc.ClientStream
}

type fileServiceDownloadClient struct {
	grpc.ClientStream
}

func (x *fileServiceDownloadClient) Recv() (*FileDownloadResponse, error) {
	m := new(FileDownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileServiceClient) GroupDownload(ctx context.Context, in *GroupDownloadRequest, opts ...grpc.CallOption) (FileService_GroupDownloadClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_FileService_serviceDesc.Streams[1], c.cc, "/proto_file.FileService/GroupDownload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceGroupDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileService_GroupDownloadClient interface {
	Recv() (*FileDownloadResponse, error)
	grpc.ClientStream
}

type fileServiceGroupDownloadClient struct {
	grpc.ClientStream
}

func (x *fileServiceGroupDownloadClient) Recv() (*FileDownloadResponse, error) {
	m := new(FileDownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileServiceClient) Upload(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_FileService_serviceDesc.Streams[2], c.cc, "/proto_file.FileService/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceUploadClient{stream}
	return x, nil
}

type FileService_UploadClient interface {
	Send(*FileUploadRequest) error
	CloseAndRecv() (*FileUploadResponse, error)
	grpc.ClientStream
}

type fileServiceUploadClient struct {
	grpc.ClientStream
}

func (x *fileServiceUploadClient) Send(m *FileUploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileServiceUploadClient) CloseAndRecv() (*FileUploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileUploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileServiceClient) Copy(ctx context.Context, in *CopyFileRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := grpc.Invoke(ctx, "/proto_file.FileService/Copy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) CopyGroup(ctx context.Context, in *CopyGroupRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := grpc.Invoke(ctx, "/proto_file.FileService/CopyGroup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) Delete(ctx context.Context, in *FileDeleteRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := grpc.Invoke(ctx, "/proto_file.FileService/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) GroupDelete(ctx context.Context, in *GroupDeleteRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := grpc.Invoke(ctx, "/proto_file.FileService/GroupDelete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for FileService service

type FileServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Download(*FileDownloadRequest, FileService_DownloadServer) error
	GroupDownload(*GroupDownloadRequest, FileService_GroupDownloadServer) error
	Upload(FileService_UploadServer) error
	Copy(context.Context, *CopyFileRequest) (*MessageResponse, error)
	CopyGroup(context.Context, *CopyGroupRequest) (*MessageResponse, error)
	Delete(context.Context, *FileDeleteRequest) (*MessageResponse, error)
	GroupDelete(context.Context, *GroupDeleteRequest) (*MessageResponse, error)
}

func RegisterFileServiceServer(s *grpc.Server, srv FileServiceServer) {
	s.RegisterService(&_FileService_serviceDesc, srv)
}

func _FileService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_file.FileService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileDownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServiceServer).Download(m, &fileServiceDownloadServer{stream})
}

type FileService_DownloadServer interface {
	Send(*FileDownloadResponse) error
	grpc.ServerStream
}

type fileServiceDownloadServer struct {
	grpc.ServerStream
}

func (x *fileServiceDownloadServer) Send(m *FileDownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _FileService_GroupDownload_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GroupDownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServiceServer).GroupDownload(m, &fileServiceGroupDownloadServer{stream})
}

type FileService_GroupDownloadServer interface {
	Send(*FileDownloadResponse) error
	grpc.ServerStream
}

type fileServiceGroupDownloadServer struct {
	grpc.ServerStream
}

func (x *fileServiceGroupDownloadServer) Send(m *FileDownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _FileService_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).Upload(&fileServiceUploadServer{stream})
}

type FileService_UploadServer interface {
	SendAndClose(*FileUploadResponse) error
	Recv() (*FileUploadRequest, error)
	grpc.ServerStream
}

type fileServiceUploadServer struct {
	grpc.ServerStream
}

func (x *fileServiceUploadServer) SendAndClose(m *FileUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileServiceUploadServer) Recv() (*FileUploadRequest, error) {
	m := new(FileUploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileService_Copy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CopyFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).Copy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_file.FileService/Copy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).Copy(ctx, req.(*CopyFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_CopyGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CopyGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).CopyGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_file.FileService/CopyGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).CopyGroup(ctx, req.(*CopyGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_file.FileService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).Delete(ctx, req.(*FileDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_GroupDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GroupDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_file.FileService/GroupDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GroupDelete(ctx, req.(*GroupDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _FileService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto_file.FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _FileService_Ping_Handler,
		},
		{
			MethodName: "Copy",
			Handler:    _FileService_Copy_Handler,
		},
		{
			MethodName: "CopyGroup",
			Handler:    _FileService_CopyGroup_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _FileService_Delete_Handler,
		},
		{
			MethodName: "GroupDelete",
			Handler:    _FileService_GroupDelete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Download",
			Handler:       _FileService_Download_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GroupDownload",
			Handler:       _FileService_GroupDownload_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Upload",
			Handler:       _FileService_Upload_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/file.proto",
}

func init() { proto.RegisterFile("proto/file.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 593 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xad, 0x3f, 0x9a, 0x3a, 0x93, 0x86, 0x84, 0x21, 0x02, 0xcb, 0x85, 0x12, 0xed, 0x01, 0x45,
	0x42, 0x4a, 0x51, 0x38, 0x21, 0x2e, 0x28, 0x44, 0x6d, 0x41, 0x2a, 0x20, 0xa3, 0x5e, 0x10, 0x12,
	0x72, 0x93, 0xad, 0x65, 0xd5, 0xb1, 0x8d, 0xbd, 0xa1, 0x6a, 0x0f, 0xfc, 0x4d, 0xfe, 0x0e, 0xda,
	0x5d, 0x7f, 0xac, 0x9d, 0xb6, 0xb4, 0x08, 0x4e, 0x9e, 0x99, 0x1d, 0xbf, 0xb7, 0xfb, 0xe6, 0xed,
	0x42, 0x3f, 0x49, 0x63, 0x16, 0xef, 0x9d, 0x06, 0x21, 0x1d, 0x8b, 0x10, 0x41, 0x7c, 0xbe, 0xf1,
	0x0a, 0x61, 0x60, 0x1d, 0x51, 0xe6, 0xcd, 0x3c, 0xe6, 0x21, 0x82, 0x19, 0x79, 0x4b, 0x6a, 0x6b,
	0x43, 0x6d, 0xd4, 0x76, 0x45, 0xcc, 0x6b, 0xec, 0x22, 0xa1, 0xb6, 0x2e, 0x6b, 0x3c, 0xe6, 0xb5,
	0x2c, 0xb8, 0xa4, 0xb6, 0x31, 0xd4, 0x46, 0x86, 0x2b, 0x62, 0x1c, 0xc0, 0xa6, 0x9f, 0xc6, 0xab,
	0xc4, 0x36, 0x45, 0xa3, 0x4c, 0xf0, 0x21, 0xb4, 0x4e, 0x56, 0xf3, 0x33, 0xca, 0xec, 0x4d, 0x51,
	0xce, 0x33, 0x32, 0x04, 0x73, 0x3f, 0x08, 0x29, 0xda, 0xb0, 0x35, 0x8f, 0x23, 0x46, 0x23, 0x26,
	0x48, 0xb7, 0xdd, 0x22, 0x25, 0x97, 0x70, 0x9f, 0x77, 0x1c, 0x27, 0x61, 0xec, 0x2d, 0x5c, 0xfa,
	0x7d, 0x45, 0x33, 0x86, 0x13, 0xb0, 0x96, 0x94, 0x79, 0x0b, 0x8f, 0x79, 0x62, 0x43, 0x9d, 0xc9,
	0x60, 0x5c, 0x9d, 0x65, 0x5c, 0x1c, 0xe4, 0x70, 0xc3, 0x2d, 0xfb, 0xf0, 0x19, 0x98, 0x7c, 0x51,
	0xe0, 0x77, 0x26, 0x7d, 0xb5, 0x9f, 0x13, 0x1c, 0x6e, 0xb8, 0x62, 0x7d, 0xda, 0x86, 0xad, 0x54,
	0xd2, 0x90, 0x53, 0x40, 0x95, 0x3b, 0x4b, 0xe2, 0x28, 0xa3, 0x78, 0x0f, 0xf4, 0x60, 0x91, 0x6b,
	0xa3, 0x07, 0x8b, 0x52, 0x2d, 0x5d, 0x51, 0xcb, 0x01, 0x2b, 0x4e, 0x03, 0xff, 0x03, 0xaf, 0x1b,
	0xa2, 0x5e, 0xe6, 0xd8, 0x07, 0x63, 0x95, 0x86, 0xb9, 0x3e, 0x3c, 0x24, 0x3e, 0x3c, 0xe0, 0x3c,
	0xb3, 0xf8, 0x3c, 0x52, 0x4f, 0xd9, 0x24, 0xaa, 0x44, 0xd4, 0x55, 0x11, 0x2b, 0xc9, 0x0d, 0x55,
	0xf2, 0x62, 0x5b, 0x66, 0xb5, 0x2d, 0x32, 0x83, 0xc1, 0x01, 0x5f, 0x6c, 0x32, 0x55, 0xc8, 0xda,
	0xd5, 0xc8, 0xba, 0x82, 0x4c, 0x7e, 0xc2, 0xa0, 0xbe, 0xdd, 0x5c, 0x18, 0x75, 0x2a, 0xda, 0x1d,
	0xa7, 0xa2, 0xff, 0x61, 0x2a, 0x00, 0x56, 0x9a, 0xf3, 0x90, 0x33, 0xe8, 0xbd, 0x8d, 0x93, 0x0b,
	0xbe, 0xfe, 0x6f, 0xa4, 0x72, 0xc0, 0x8a, 0xe8, 0xf9, 0x81, 0x62, 0xdb, 0x32, 0x27, 0x5f, 0xa1,
	0xcf, 0xc9, 0x44, 0xf2, 0x57, 0x72, 0xd5, 0xd0, 0x8d, 0x06, 0x3a, 0x95, 0xee, 0x9e, 0xd1, 0x90,
	0x32, 0xfa, 0xff, 0xe6, 0x3e, 0x05, 0x94, 0x73, 0xaf, 0xf1, 0xdc, 0x6d, 0xea, 0xcf, 0xa1, 0x77,
	0x44, 0xb3, 0xcc, 0xf3, 0x69, 0x39, 0x70, 0x1b, 0xb6, 0x96, 0xb2, 0x94, 0x23, 0x14, 0x29, 0xe9,
	0x42, 0xe7, 0x53, 0x10, 0xf9, 0x39, 0x13, 0x21, 0xb0, 0x2d, 0xd3, 0xfc, 0x47, 0x04, 0x33, 0x09,
	0x22, 0xbf, 0x78, 0x60, 0x78, 0x3c, 0xf9, 0x65, 0x42, 0x87, 0x6b, 0xf1, 0x99, 0xa6, 0x3f, 0x82,
	0x39, 0xc5, 0x57, 0x60, 0xf2, 0x7f, 0xf0, 0x91, 0xea, 0x09, 0x05, 0xd4, 0xb1, 0xd7, 0x17, 0x72,
	0xf8, 0x8f, 0x60, 0x15, 0xe6, 0xc4, 0xa7, 0x4d, 0x4b, 0x35, 0xbc, 0xef, 0x0c, 0xaf, 0x6f, 0x90,
	0x70, 0x2f, 0x34, 0x3c, 0x86, 0x6e, 0xed, 0xde, 0x60, 0xed, 0xa7, 0xab, 0xae, 0xd4, 0xad, 0x60,
	0xdf, 0x41, 0x4b, 0xbe, 0x2d, 0xf8, 0xa4, 0xd9, 0x5d, 0x7b, 0xef, 0x9c, 0xdd, 0xeb, 0x96, 0x25,
	0xd4, 0x48, 0xc3, 0x37, 0x60, 0x72, 0x9b, 0xe2, 0x8e, 0xda, 0xd9, 0xb8, 0x25, 0xce, 0x4e, 0xfd,
	0x3a, 0xd6, 0x87, 0xb9, 0x0f, 0xed, 0xd2, 0xe8, 0xf8, 0xb8, 0x09, 0xa3, 0xfa, 0xff, 0x66, 0x9c,
	0x19, 0xb4, 0xa4, 0xcd, 0xd6, 0x0f, 0x55, 0xb3, 0xdf, 0xcd, 0x28, 0xef, 0xa1, 0xa3, 0x38, 0x16,
	0x77, 0xd7, 0xf5, 0xbe, 0x35, 0xd6, 0xb4, 0xf7, 0xa5, 0x3b, 0xde, 0x7b, 0x5d, 0x35, 0x9c, 0xb4,
	0x44, 0xfc, 0xf2, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xdb, 0x27, 0x93, 0x12, 0x07, 0x00,
	0x00,
}