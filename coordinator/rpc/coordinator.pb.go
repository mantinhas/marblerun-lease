// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: coordinator.proto

package rpc

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

type ActivationReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// TODO: sending the quote via metadata/context would be cleaner.
	Quote      []byte `protobuf:"bytes,1,opt,name=Quote,proto3" json:"Quote,omitempty"`
	CSR        []byte `protobuf:"bytes,2,opt,name=CSR,proto3" json:"CSR,omitempty"`
	MarbleType string `protobuf:"bytes,3,opt,name=MarbleType,proto3" json:"MarbleType,omitempty"`
	UUID       string `protobuf:"bytes,4,opt,name=UUID,proto3" json:"UUID,omitempty"`
}

func (x *ActivationReq) Reset() {
	*x = ActivationReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivationReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivationReq) ProtoMessage() {}

func (x *ActivationReq) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivationReq.ProtoReflect.Descriptor instead.
func (*ActivationReq) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{0}
}

func (x *ActivationReq) GetQuote() []byte {
	if x != nil {
		return x.Quote
	}
	return nil
}

func (x *ActivationReq) GetCSR() []byte {
	if x != nil {
		return x.CSR
	}
	return nil
}

func (x *ActivationReq) GetMarbleType() string {
	if x != nil {
		return x.MarbleType
	}
	return ""
}

func (x *ActivationReq) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

type DeactivationSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TrustProtocol string `protobuf:"bytes,1,opt,name=TrustProtocol,proto3" json:"TrustProtocol,omitempty"`
	// Types that are assignable to Settings:
	//	*DeactivationSettings_PingSettings
	//	*DeactivationSettings_LeaseSettings
	Settings isDeactivationSettings_Settings `protobuf_oneof:"Settings"`
}

func (x *DeactivationSettings) Reset() {
	*x = DeactivationSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeactivationSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivationSettings) ProtoMessage() {}

func (x *DeactivationSettings) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivationSettings.ProtoReflect.Descriptor instead.
func (*DeactivationSettings) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{1}
}

func (x *DeactivationSettings) GetTrustProtocol() string {
	if x != nil {
		return x.TrustProtocol
	}
	return ""
}

func (m *DeactivationSettings) GetSettings() isDeactivationSettings_Settings {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (x *DeactivationSettings) GetPingSettings() *PingSettings {
	if x, ok := x.GetSettings().(*DeactivationSettings_PingSettings); ok {
		return x.PingSettings
	}
	return nil
}

func (x *DeactivationSettings) GetLeaseSettings() *LeaseSettings {
	if x, ok := x.GetSettings().(*DeactivationSettings_LeaseSettings); ok {
		return x.LeaseSettings
	}
	return nil
}

type isDeactivationSettings_Settings interface {
	isDeactivationSettings_Settings()
}

type DeactivationSettings_PingSettings struct {
	PingSettings *PingSettings `protobuf:"bytes,2,opt,name=PingSettings,proto3,oneof"`
}

type DeactivationSettings_LeaseSettings struct {
	LeaseSettings *LeaseSettings `protobuf:"bytes,3,opt,name=LeaseSettings,proto3,oneof"`
}

func (*DeactivationSettings_PingSettings) isDeactivationSettings_Settings() {}

func (*DeactivationSettings_LeaseSettings) isDeactivationSettings_Settings() {}

type PingSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestInterval string `protobuf:"bytes,1,opt,name=RequestInterval,proto3" json:"RequestInterval,omitempty"`
	RetryInterval   string `protobuf:"bytes,2,opt,name=RetryInterval,proto3" json:"RetryInterval,omitempty"`
	Retries         int32  `protobuf:"varint,3,opt,name=Retries,proto3" json:"Retries,omitempty"`
}

func (x *PingSettings) Reset() {
	*x = PingSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingSettings) ProtoMessage() {}

func (x *PingSettings) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingSettings.ProtoReflect.Descriptor instead.
func (*PingSettings) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{2}
}

func (x *PingSettings) GetRequestInterval() string {
	if x != nil {
		return x.RequestInterval
	}
	return ""
}

func (x *PingSettings) GetRetryInterval() string {
	if x != nil {
		return x.RetryInterval
	}
	return ""
}

func (x *PingSettings) GetRetries() int32 {
	if x != nil {
		return x.Retries
	}
	return 0
}

type LeaseSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestInterval string `protobuf:"bytes,1,opt,name=RequestInterval,proto3" json:"RequestInterval,omitempty"`
	RetryInterval   string `protobuf:"bytes,2,opt,name=RetryInterval,proto3" json:"RetryInterval,omitempty"`
	Retries         int32  `protobuf:"varint,3,opt,name=Retries,proto3" json:"Retries,omitempty"`
}

func (x *LeaseSettings) Reset() {
	*x = LeaseSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaseSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaseSettings) ProtoMessage() {}

func (x *LeaseSettings) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaseSettings.ProtoReflect.Descriptor instead.
func (*LeaseSettings) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{3}
}

func (x *LeaseSettings) GetRequestInterval() string {
	if x != nil {
		return x.RequestInterval
	}
	return ""
}

func (x *LeaseSettings) GetRetryInterval() string {
	if x != nil {
		return x.RetryInterval
	}
	return ""
}

func (x *LeaseSettings) GetRetries() int32 {
	if x != nil {
		return x.Retries
	}
	return 0
}

type ActivationResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parameters           *Parameters           `protobuf:"bytes,1,opt,name=Parameters,proto3" json:"Parameters,omitempty"`
	DeactivationSettings *DeactivationSettings `protobuf:"bytes,2,opt,name=DeactivationSettings,proto3" json:"DeactivationSettings,omitempty"`
}

func (x *ActivationResp) Reset() {
	*x = ActivationResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivationResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivationResp) ProtoMessage() {}

func (x *ActivationResp) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivationResp.ProtoReflect.Descriptor instead.
func (*ActivationResp) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{4}
}

func (x *ActivationResp) GetParameters() *Parameters {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *ActivationResp) GetDeactivationSettings() *DeactivationSettings {
	if x != nil {
		return x.DeactivationSettings
	}
	return nil
}

type Parameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files map[string][]byte `protobuf:"bytes,1,rep,name=Files,proto3" json:"Files,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Env   map[string][]byte `protobuf:"bytes,2,rep,name=Env,proto3" json:"Env,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Argv  []string          `protobuf:"bytes,3,rep,name=Argv,proto3" json:"Argv,omitempty"`
}

func (x *Parameters) Reset() {
	*x = Parameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Parameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Parameters) ProtoMessage() {}

func (x *Parameters) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Parameters.ProtoReflect.Descriptor instead.
func (*Parameters) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{5}
}

func (x *Parameters) GetFiles() map[string][]byte {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *Parameters) GetEnv() map[string][]byte {
	if x != nil {
		return x.Env
	}
	return nil
}

func (x *Parameters) GetArgv() []string {
	if x != nil {
		return x.Argv
	}
	return nil
}

type PingReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MarbleType string `protobuf:"bytes,1,opt,name=MarbleType,proto3" json:"MarbleType,omitempty"`
	UUID       string `protobuf:"bytes,2,opt,name=UUID,proto3" json:"UUID,omitempty"`
}

func (x *PingReq) Reset() {
	*x = PingReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingReq) ProtoMessage() {}

func (x *PingReq) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingReq.ProtoReflect.Descriptor instead.
func (*PingReq) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{6}
}

func (x *PingReq) GetMarbleType() string {
	if x != nil {
		return x.MarbleType
	}
	return ""
}

func (x *PingReq) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

type PingResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *PingResp) Reset() {
	*x = PingResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResp) ProtoMessage() {}

func (x *PingResp) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResp.ProtoReflect.Descriptor instead.
func (*PingResp) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{7}
}

func (x *PingResp) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type LeaseOffer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeaseDuration string `protobuf:"bytes,1,opt,name=leaseDuration,proto3" json:"leaseDuration,omitempty"`
}

func (x *LeaseOffer) Reset() {
	*x = LeaseOffer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaseOffer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaseOffer) ProtoMessage() {}

func (x *LeaseOffer) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaseOffer.ProtoReflect.Descriptor instead.
func (*LeaseOffer) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{8}
}

func (x *LeaseOffer) GetLeaseDuration() string {
	if x != nil {
		return x.LeaseDuration
	}
	return ""
}

type LeaseResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *LeaseResp) Reset() {
	*x = LeaseResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaseResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaseResp) ProtoMessage() {}

func (x *LeaseResp) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaseResp.ProtoReflect.Descriptor instead.
func (*LeaseResp) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{9}
}

func (x *LeaseResp) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type RemainingLeaseReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MarbleType string `protobuf:"bytes,1,opt,name=MarbleType,proto3" json:"MarbleType,omitempty"`
	UUID       string `protobuf:"bytes,2,opt,name=UUID,proto3" json:"UUID,omitempty"`
}

func (x *RemainingLeaseReq) Reset() {
	*x = RemainingLeaseReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemainingLeaseReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemainingLeaseReq) ProtoMessage() {}

func (x *RemainingLeaseReq) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemainingLeaseReq.ProtoReflect.Descriptor instead.
func (*RemainingLeaseReq) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{10}
}

func (x *RemainingLeaseReq) GetMarbleType() string {
	if x != nil {
		return x.MarbleType
	}
	return ""
}

func (x *RemainingLeaseReq) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

type RemainingLeaseOffer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok            bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	LeaseDuration string `protobuf:"bytes,2,opt,name=leaseDuration,proto3" json:"leaseDuration,omitempty"`
}

func (x *RemainingLeaseOffer) Reset() {
	*x = RemainingLeaseOffer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemainingLeaseOffer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemainingLeaseOffer) ProtoMessage() {}

func (x *RemainingLeaseOffer) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemainingLeaseOffer.ProtoReflect.Descriptor instead.
func (*RemainingLeaseOffer) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{11}
}

func (x *RemainingLeaseOffer) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *RemainingLeaseOffer) GetLeaseDuration() string {
	if x != nil {
		return x.LeaseDuration
	}
	return ""
}

type DeactivateReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeactivateReq) Reset() {
	*x = DeactivateReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeactivateReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateReq) ProtoMessage() {}

func (x *DeactivateReq) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateReq.ProtoReflect.Descriptor instead.
func (*DeactivateReq) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{12}
}

type DeactivateResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeactivateResp) Reset() {
	*x = DeactivateResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeactivateResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateResp) ProtoMessage() {}

func (x *DeactivateResp) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateResp.ProtoReflect.Descriptor instead.
func (*DeactivateResp) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{13}
}

type AppUsage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CpuTime float32 `protobuf:"fixed32,1,opt,name=cpuTime,proto3" json:"cpuTime,omitempty"`
}

func (x *AppUsage) Reset() {
	*x = AppUsage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppUsage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppUsage) ProtoMessage() {}

func (x *AppUsage) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppUsage.ProtoReflect.Descriptor instead.
func (*AppUsage) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{14}
}

func (x *AppUsage) GetCpuTime() float32 {
	if x != nil {
		return x.CpuTime
	}
	return 0
}

var File_coordinator_proto protoreflect.FileDescriptor

var file_coordinator_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x72, 0x70, 0x63, 0x22, 0x6b, 0x0a, 0x0d, 0x41, 0x63, 0x74, 0x69,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x51, 0x75, 0x6f,
	0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x43, 0x53, 0x52, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x43, 0x53,
	0x52, 0x12, 0x1e, 0x0a, 0x0a, 0x4d, 0x61, 0x72, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x4d, 0x61, 0x72, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x55, 0x55, 0x49, 0x44, 0x22, 0xbd, 0x01, 0x0a, 0x14, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x24,
	0x0a, 0x0d, 0x54, 0x72, 0x75, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x54, 0x72, 0x75, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x37, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x48, 0x00, 0x52,
	0x0c, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x3a, 0x0a,
	0x0d, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x4c, 0x65, 0x61, 0x73, 0x65,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x48, 0x00, 0x52, 0x0d, 0x4c, 0x65, 0x61, 0x73,
	0x65, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x42, 0x0a, 0x0a, 0x08, 0x53, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x78, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12,
	0x24, 0x0a, 0x0d, 0x52, 0x65, 0x74, 0x72, 0x79, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x52, 0x65, 0x74, 0x72, 0x79, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22,
	0x79, 0x0a, 0x0d, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x12, 0x28, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x24, 0x0a, 0x0d, 0x52, 0x65,
	0x74, 0x72, 0x79, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x52, 0x65, 0x74, 0x72, 0x79, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x12, 0x18, 0x0a, 0x07, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x90, 0x01, 0x0a, 0x0e, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2f, 0x0a,
	0x0a, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x52, 0x0a, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x4d,
	0x0a, 0x14, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x14, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0xf0, 0x01,
	0x0a, 0x0a, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x30, 0x0a, 0x05,
	0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x2a,
	0x0a, 0x03, 0x45, 0x6e, 0x76, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x45, 0x6e, 0x76,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x45, 0x6e, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x41, 0x72,
	0x67, 0x76, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x41, 0x72, 0x67, 0x76, 0x1a, 0x38,
	0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x36, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x3d, 0x0a, 0x07, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x0a, 0x4d,
	0x61, 0x72, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x4d, 0x61, 0x72, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x55,
	0x55, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44, 0x22,
	0x1a, 0x0a, 0x08, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x6f,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x32, 0x0a, 0x0a, 0x4c,
	0x65, 0x61, 0x73, 0x65, 0x4f, 0x66, 0x66, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x1b, 0x0a, 0x09, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x0e, 0x0a, 0x02,
	0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x47, 0x0a, 0x11,
	0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65,
	0x71, 0x12, 0x1e, 0x0a, 0x0a, 0x4d, 0x61, 0x72, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x4d, 0x61, 0x72, 0x62, 0x6c, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x55, 0x55, 0x49, 0x44, 0x22, 0x4b, 0x0a, 0x13, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69,
	0x6e, 0x67, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x4f, 0x66, 0x66, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02,
	0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x24, 0x0a, 0x0d,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x0f, 0x0a, 0x0d, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x22, 0x10, 0x0a, 0x0e, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x24, 0x0a, 0x08, 0x41, 0x70, 0x70, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x70, 0x75, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x07, 0x63, 0x70, 0x75, 0x54, 0x69, 0x6d, 0x65, 0x32, 0x90, 0x02, 0x0a, 0x06,
	0x4d, 0x61, 0x72, 0x62, 0x6c, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x12, 0x12, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x12, 0x23, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x0c, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x1a, 0x0d, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x35, 0x0a, 0x0a, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x12, 0x12,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x1a, 0x13, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x31, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x70, 0x61,
	0x67, 0x61, 0x74, 0x65, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x0f, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x4c, 0x65, 0x61, 0x73, 0x65, 0x4f, 0x66, 0x66, 0x65, 0x72, 0x1a, 0x0e, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x42, 0x0a, 0x0e, 0x52, 0x65,
	0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x16, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x4c, 0x65, 0x61, 0x73,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x6d, 0x61, 0x69,
	0x6e, 0x69, 0x6e, 0x67, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x4f, 0x66, 0x66, 0x65, 0x72, 0x42, 0x32,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x64, 0x67,
	0x65, 0x6c, 0x65, 0x73, 0x73, 0x73, 0x79, 0x73, 0x2f, 0x6d, 0x61, 0x72, 0x62, 0x6c, 0x65, 0x72,
	0x75, 0x6e, 0x2f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coordinator_proto_rawDescOnce sync.Once
	file_coordinator_proto_rawDescData = file_coordinator_proto_rawDesc
)

func file_coordinator_proto_rawDescGZIP() []byte {
	file_coordinator_proto_rawDescOnce.Do(func() {
		file_coordinator_proto_rawDescData = protoimpl.X.CompressGZIP(file_coordinator_proto_rawDescData)
	})
	return file_coordinator_proto_rawDescData
}

var file_coordinator_proto_msgTypes = make([]protoimpl.MessageInfo, 17)
var file_coordinator_proto_goTypes = []interface{}{
	(*ActivationReq)(nil),        // 0: rpc.ActivationReq
	(*DeactivationSettings)(nil), // 1: rpc.DeactivationSettings
	(*PingSettings)(nil),         // 2: rpc.PingSettings
	(*LeaseSettings)(nil),        // 3: rpc.LeaseSettings
	(*ActivationResp)(nil),       // 4: rpc.ActivationResp
	(*Parameters)(nil),           // 5: rpc.Parameters
	(*PingReq)(nil),              // 6: rpc.PingReq
	(*PingResp)(nil),             // 7: rpc.PingResp
	(*LeaseOffer)(nil),           // 8: rpc.LeaseOffer
	(*LeaseResp)(nil),            // 9: rpc.LeaseResp
	(*RemainingLeaseReq)(nil),    // 10: rpc.RemainingLeaseReq
	(*RemainingLeaseOffer)(nil),  // 11: rpc.RemainingLeaseOffer
	(*DeactivateReq)(nil),        // 12: rpc.DeactivateReq
	(*DeactivateResp)(nil),       // 13: rpc.DeactivateResp
	(*AppUsage)(nil),             // 14: rpc.AppUsage
	nil,                          // 15: rpc.Parameters.FilesEntry
	nil,                          // 16: rpc.Parameters.EnvEntry
}
var file_coordinator_proto_depIdxs = []int32{
	2,  // 0: rpc.DeactivationSettings.PingSettings:type_name -> rpc.PingSettings
	3,  // 1: rpc.DeactivationSettings.LeaseSettings:type_name -> rpc.LeaseSettings
	5,  // 2: rpc.ActivationResp.Parameters:type_name -> rpc.Parameters
	1,  // 3: rpc.ActivationResp.DeactivationSettings:type_name -> rpc.DeactivationSettings
	15, // 4: rpc.Parameters.Files:type_name -> rpc.Parameters.FilesEntry
	16, // 5: rpc.Parameters.Env:type_name -> rpc.Parameters.EnvEntry
	0,  // 6: rpc.Marble.Activate:input_type -> rpc.ActivationReq
	6,  // 7: rpc.Marble.Ping:input_type -> rpc.PingReq
	12, // 8: rpc.Marble.Deactivate:input_type -> rpc.DeactivateReq
	8,  // 9: rpc.Marble.PropagateLease:input_type -> rpc.LeaseOffer
	10, // 10: rpc.Marble.RemainingLease:input_type -> rpc.RemainingLeaseReq
	4,  // 11: rpc.Marble.Activate:output_type -> rpc.ActivationResp
	7,  // 12: rpc.Marble.Ping:output_type -> rpc.PingResp
	13, // 13: rpc.Marble.Deactivate:output_type -> rpc.DeactivateResp
	9,  // 14: rpc.Marble.PropagateLease:output_type -> rpc.LeaseResp
	11, // 15: rpc.Marble.RemainingLease:output_type -> rpc.RemainingLeaseOffer
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_coordinator_proto_init() }
func file_coordinator_proto_init() {
	if File_coordinator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_coordinator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivationReq); i {
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
		file_coordinator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeactivationSettings); i {
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
		file_coordinator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingSettings); i {
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
		file_coordinator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaseSettings); i {
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
		file_coordinator_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivationResp); i {
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
		file_coordinator_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Parameters); i {
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
		file_coordinator_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingReq); i {
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
		file_coordinator_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResp); i {
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
		file_coordinator_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaseOffer); i {
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
		file_coordinator_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaseResp); i {
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
		file_coordinator_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemainingLeaseReq); i {
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
		file_coordinator_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemainingLeaseOffer); i {
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
		file_coordinator_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeactivateReq); i {
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
		file_coordinator_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeactivateResp); i {
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
		file_coordinator_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppUsage); i {
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
	file_coordinator_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*DeactivationSettings_PingSettings)(nil),
		(*DeactivationSettings_LeaseSettings)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_coordinator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   17,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_coordinator_proto_goTypes,
		DependencyIndexes: file_coordinator_proto_depIdxs,
		MessageInfos:      file_coordinator_proto_msgTypes,
	}.Build()
	File_coordinator_proto = out.File
	file_coordinator_proto_rawDesc = nil
	file_coordinator_proto_goTypes = nil
	file_coordinator_proto_depIdxs = nil
}
