// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gateway.proto

package proto

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type GatewayRegisterRequest struct {
	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *GatewayRegisterRequest) Reset()      { *m = GatewayRegisterRequest{} }
func (*GatewayRegisterRequest) ProtoMessage() {}
func (*GatewayRegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{0}
}
func (m *GatewayRegisterRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GatewayRegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GatewayRegisterRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GatewayRegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GatewayRegisterRequest.Merge(m, src)
}
func (m *GatewayRegisterRequest) XXX_Size() int {
	return m.Size()
}
func (m *GatewayRegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GatewayRegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GatewayRegisterRequest proto.InternalMessageInfo

func (m *GatewayRegisterRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto.RegisterType((*GatewayRegisterRequest)(nil), "proto.GatewayRegisterRequest")
}

func init() { proto.RegisterFile("gateway.proto", fileDescriptor_f1a937782ebbded5) }

var fileDescriptor_f1a937782ebbded5 = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x4f, 0x2c, 0x49,
	0x2d, 0x4f, 0xac, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xba, 0xe9,
	0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0xe9, 0xf9, 0xe9, 0xf9, 0xfa, 0x60,
	0xe1, 0xa4, 0xd2, 0x34, 0x30, 0x0f, 0xcc, 0x01, 0xb3, 0x20, 0xba, 0xa4, 0xf8, 0x8a, 0x52, 0x8b,
	0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x21, 0x7c, 0x25, 0x0d, 0x2e, 0x31, 0x77, 0x88, 0xb1, 0x41, 0xa9,
	0xe9, 0x99, 0xc5, 0x25, 0xa9, 0x45, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x7c, 0x5c,
	0x4c, 0x99, 0x29, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xac, 0x41, 0x4c, 0x99, 0x29, 0x4e, 0x26, 0x17,
	0x1e, 0xca, 0x31, 0xdc, 0x78, 0x28, 0xc7, 0xf0, 0xe1, 0xa1, 0x1c, 0x63, 0xc3, 0x23, 0x39, 0xc6,
	0x15, 0x8f, 0xe4, 0x18, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39,
	0xc6, 0x17, 0x8f, 0xe4, 0x18, 0x3e, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63,
	0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0x92, 0xd8, 0xc0, 0xd6, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0xdf, 0x97, 0x92, 0x29, 0xbd, 0x00, 0x00, 0x00,
}

func (this *GatewayRegisterRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GatewayRegisterRequest)
	if !ok {
		that2, ok := that.(GatewayRegisterRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Id != that1.Id {
		return false
	}
	return true
}
func (this *GatewayRegisterRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&proto.GatewayRegisterRequest{")
	s = append(s, "Id: "+fmt.Sprintf("%#v", this.Id)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringGateway(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *GatewayRegisterRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GatewayRegisterRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GatewayRegisterRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintGateway(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintGateway(dAtA []byte, offset int, v uint64) int {
	offset -= sovGateway(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GatewayRegisterRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovGateway(uint64(m.Id))
	}
	return n
}

func sovGateway(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGateway(x uint64) (n int) {
	return sovGateway(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *GatewayRegisterRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&GatewayRegisterRequest{`,
		`Id:` + fmt.Sprintf("%v", this.Id) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGateway(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *GatewayRegisterRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGateway
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GatewayRegisterRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GatewayRegisterRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGateway
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGateway(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGateway
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthGateway
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGateway(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGateway
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGateway
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGateway
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGateway
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthGateway
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGateway
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGateway(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthGateway
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGateway = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGateway   = fmt.Errorf("proto: integer overflow")
)
