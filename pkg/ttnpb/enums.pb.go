// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lorawan-stack/api/enums.proto

package ttnpb // import "go.thethings.network/lorawan-stack/pkg/ttnpb"

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strconv "strconv"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type DownlinkPathConstraint int32

const (
	// Indicates that the gateway can be selected for downlink without constraints by the Network Server.
	DOWNLINK_PATH_CONSTRAINT_NONE DownlinkPathConstraint = 0
	// Indicates that the gateway can be selected for downlink only if no other or better gateway can be selected.
	DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER DownlinkPathConstraint = 1
	// Indicates that this gateway will never be selected for downlink, even if that results in no available downlink path.
	DOWNLINK_PATH_CONSTRAINT_NEVER DownlinkPathConstraint = 2
)

var DownlinkPathConstraint_name = map[int32]string{
	0: "DOWNLINK_PATH_CONSTRAINT_NONE",
	1: "DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER",
	2: "DOWNLINK_PATH_CONSTRAINT_NEVER",
}
var DownlinkPathConstraint_value = map[string]int32{
	"DOWNLINK_PATH_CONSTRAINT_NONE":         0,
	"DOWNLINK_PATH_CONSTRAINT_PREFER_OTHER": 1,
	"DOWNLINK_PATH_CONSTRAINT_NEVER":        2,
}

func (DownlinkPathConstraint) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_enums_ac75038b2469ccb1, []int{0}
}

// State enum defines states that an entity can be in.
type State int32

const (
	// Denotes that the entity has been requested and is pending review by an admin.
	STATE_REQUESTED State = 0
	// Denotes that the entity has been reviewed and approved by an admin.
	STATE_APPROVED State = 1
	// Denotes that the entity has been reviewed and rejected by an admin.
	STATE_REJECTED State = 2
	// Denotes that the entity has been flagged and is pending review by an admin.
	STATE_FLAGGED State = 3
	// Denotes that the entity has been reviewed and suspended by an admin.
	STATE_SUSPENDED State = 4
)

var State_name = map[int32]string{
	0: "STATE_REQUESTED",
	1: "STATE_APPROVED",
	2: "STATE_REJECTED",
	3: "STATE_FLAGGED",
	4: "STATE_SUSPENDED",
}
var State_value = map[string]int32{
	"STATE_REQUESTED": 0,
	"STATE_APPROVED":  1,
	"STATE_REJECTED":  2,
	"STATE_FLAGGED":   3,
	"STATE_SUSPENDED": 4,
}

func (State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_enums_ac75038b2469ccb1, []int{1}
}

func init() {
	proto.RegisterEnum("ttn.lorawan.v3.DownlinkPathConstraint", DownlinkPathConstraint_name, DownlinkPathConstraint_value)
	golang_proto.RegisterEnum("ttn.lorawan.v3.DownlinkPathConstraint", DownlinkPathConstraint_name, DownlinkPathConstraint_value)
	proto.RegisterEnum("ttn.lorawan.v3.State", State_name, State_value)
	golang_proto.RegisterEnum("ttn.lorawan.v3.State", State_name, State_value)
}
func (x DownlinkPathConstraint) String() string {
	s, ok := DownlinkPathConstraint_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x State) String() string {
	s, ok := State_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

func init() {
	proto.RegisterFile("lorawan-stack/api/enums.proto", fileDescriptor_enums_ac75038b2469ccb1)
}
func init() {
	golang_proto.RegisterFile("lorawan-stack/api/enums.proto", fileDescriptor_enums_ac75038b2469ccb1)
}

var fileDescriptor_enums_ac75038b2469ccb1 = []byte{
	// 422 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xbf, 0x4f, 0xdc, 0x40,
	0x10, 0x85, 0x67, 0x08, 0x49, 0xb1, 0x52, 0x88, 0xe3, 0x48, 0x29, 0x90, 0x18, 0x29, 0x91, 0x52,
	0x04, 0x05, 0xbb, 0xe0, 0x2f, 0x70, 0xce, 0xcb, 0x8f, 0x04, 0xd9, 0xce, 0xda, 0x10, 0x29, 0x8d,
	0xe5, 0x43, 0x17, 0x9f, 0x75, 0xb0, 0x7b, 0xf2, 0xed, 0xe5, 0x5a, 0xca, 0x2b, 0x29, 0x53, 0x46,
	0x49, 0x43, 0x49, 0x49, 0x49, 0x49, 0x49, 0x49, 0x89, 0x77, 0x1b, 0x4a, 0x4a, 0xca, 0x28, 0x76,
	0x24, 0x44, 0x41, 0xb7, 0xef, 0xd3, 0xb7, 0x6f, 0x8a, 0xc7, 0x56, 0x0e, 0x54, 0x5d, 0xcc, 0x0a,
	0xb9, 0x36, 0xd1, 0xc5, 0xfe, 0xc8, 0x2f, 0xc6, 0x95, 0x3f, 0x90, 0xd3, 0xc3, 0x89, 0x37, 0xae,
	0x95, 0x56, 0xee, 0x92, 0xd6, 0xd2, 0xfb, 0xaf, 0x78, 0x3f, 0xd6, 0x97, 0xd7, 0xca, 0x4a, 0x0f,
	0xa7, 0x7d, 0x6f, 0x5f, 0x1d, 0xfa, 0xa5, 0x2a, 0x95, 0xdf, 0x6a, 0xfd, 0xe9, 0xf7, 0x36, 0xb5,
	0xa1, 0x7d, 0x75, 0xdf, 0x57, 0x8f, 0x91, 0xbd, 0x0e, 0xd5, 0x4c, 0x1e, 0x54, 0x72, 0x94, 0x14,
	0x7a, 0xd8, 0x53, 0x72, 0xa2, 0xeb, 0xa2, 0x92, 0xda, 0x7d, 0xc3, 0x56, 0xc2, 0xf8, 0x6b, 0xb4,
	0xb3, 0x1d, 0x7d, 0xce, 0x93, 0x20, 0xdb, 0xca, 0x7b, 0x71, 0x94, 0x66, 0x22, 0xd8, 0x8e, 0xb2,
	0x3c, 0x8a, 0x23, 0xee, 0x80, 0xfb, 0x9e, 0xbd, 0x7b, 0x54, 0x49, 0x04, 0xdf, 0xe0, 0x22, 0x8f,
	0xb3, 0x2d, 0x2e, 0x1c, 0x74, 0xdf, 0x32, 0x7a, 0xbc, 0x8d, 0xef, 0x71, 0xe1, 0x2c, 0x2c, 0x2f,
	0xce, 0xff, 0x10, 0xac, 0xd6, 0xec, 0x69, 0xaa, 0x0b, 0x3d, 0x70, 0x5f, 0xb1, 0x17, 0x69, 0x16,
	0x64, 0x3c, 0x17, 0xfc, 0xcb, 0x2e, 0x4f, 0x33, 0x1e, 0x3a, 0xe0, 0xba, 0x6c, 0xa9, 0x83, 0x41,
	0x92, 0x88, 0x78, 0x8f, 0x87, 0x0e, 0xde, 0x33, 0xc1, 0x3f, 0xf1, 0xde, 0x3f, 0x6f, 0xc1, 0x7d,
	0xc9, 0x9e, 0x77, 0x6c, 0x63, 0x27, 0xd8, 0xdc, 0xe4, 0xa1, 0xf3, 0xe4, 0xbe, 0x2f, 0xdd, 0x4d,
	0x13, 0x1e, 0x85, 0x3c, 0x74, 0x16, 0xbb, 0x9b, 0x1f, 0x7f, 0xe3, 0x45, 0x43, 0x78, 0xd9, 0x10,
	0x5e, 0x35, 0x04, 0xd7, 0x0d, 0xc1, 0x4d, 0x43, 0x70, 0xdb, 0x10, 0xdc, 0x35, 0x84, 0x47, 0x86,
	0x70, 0x6e, 0x08, 0x4e, 0x0c, 0xe1, 0xa9, 0x21, 0x38, 0x33, 0x04, 0xe7, 0x86, 0xe0, 0xc2, 0x10,
	0x5e, 0x1a, 0xc2, 0x2b, 0x43, 0x70, 0x6d, 0x08, 0x6f, 0x0c, 0xc1, 0xad, 0x21, 0xbc, 0x33, 0x04,
	0x47, 0x96, 0x60, 0x6e, 0x09, 0x8f, 0x2d, 0xc1, 0x4f, 0x4b, 0xf8, 0xcb, 0x12, 0x9c, 0x58, 0x82,
	0x53, 0x4b, 0x78, 0x66, 0x09, 0xcf, 0x2d, 0xe1, 0xb7, 0x0f, 0xa5, 0xf2, 0xf4, 0x70, 0xa0, 0x87,
	0x95, 0x2c, 0x27, 0x9e, 0x1c, 0xe8, 0x99, 0xaa, 0x47, 0xfe, 0xc3, 0xc1, 0xc7, 0xa3, 0xd2, 0xd7,
	0x5a, 0x8e, 0xfb, 0xfd, 0x67, 0xed, 0x64, 0xeb, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x3c, 0xb3,
	0x8a, 0x41, 0x12, 0x02, 0x00, 0x00,
}
