// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/crypto/secp256k1/keys.proto

package secp256k1

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

// PubKey defines a secp256k1 public key
// Key is the compressed form of the pubkey. The first byte depends is a 0x02 byte
// if the y-coordinate is the lexicographically largest of the two associated with
// the x-coordinate. Otherwise the first byte is a 0x03.
// This prefix is followed with the x-coordinate.
type PubKey struct {
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *PubKey) Reset()      { *m = PubKey{} }
func (*PubKey) ProtoMessage() {}
func (*PubKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_e0835e68ebdcb224, []int{0}
}
func (m *PubKey) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PubKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PubKey.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PubKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PubKey.Merge(m, src)
}
func (m *PubKey) XXX_Size() int {
	return m.Size()
}
func (m *PubKey) XXX_DiscardUnknown() {
	xxx_messageInfo_PubKey.DiscardUnknown(m)
}

var xxx_messageInfo_PubKey proto.InternalMessageInfo

func (m *PubKey) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

// PrivKey defines a secp256k1 private key.
type PrivKey struct {
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *PrivKey) Reset()         { *m = PrivKey{} }
func (m *PrivKey) String() string { return proto.CompactTextString(m) }
func (*PrivKey) ProtoMessage()    {}
func (*PrivKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_e0835e68ebdcb224, []int{1}
}
func (m *PrivKey) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PrivKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PrivKey.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PrivKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivKey.Merge(m, src)
}
func (m *PrivKey) XXX_Size() int {
	return m.Size()
}
func (m *PrivKey) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivKey.DiscardUnknown(m)
}

var xxx_messageInfo_PrivKey proto.InternalMessageInfo

func (m *PrivKey) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func init() {
	proto.RegisterType((*PubKey)(nil), "cosmos.crypto.secp256k1.PubKey")
	proto.RegisterType((*PrivKey)(nil), "cosmos.crypto.secp256k1.PrivKey")
}

func init() {
	proto.RegisterFile("cosmos/crypto/secp256k1/keys.proto", fileDescriptor_e0835e68ebdcb224)
}

var fileDescriptor_e0835e68ebdcb224 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4a, 0xce, 0x2f, 0xce,
	0xcd, 0x2f, 0xd6, 0x4f, 0x2e, 0xaa, 0x2c, 0x28, 0xc9, 0xd7, 0x2f, 0x4e, 0x4d, 0x2e, 0x30, 0x32,
	0x35, 0xcb, 0x36, 0xd4, 0xcf, 0x4e, 0xad, 0x2c, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12,
	0x87, 0xa8, 0xd1, 0x83, 0xa8, 0xd1, 0x83, 0xab, 0x91, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0xab,
	0xd1, 0x07, 0xb1, 0x20, 0xca, 0x95, 0x14, 0xb8, 0xd8, 0x02, 0x4a, 0x93, 0xbc, 0x53, 0x2b, 0x85,
	0x04, 0xb8, 0x98, 0xb3, 0x53, 0x2b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x78, 0x82, 0x40, 0x4c, 0x2b,
	0x96, 0x19, 0x0b, 0xe4, 0x19, 0x94, 0xa4, 0xb9, 0xd8, 0x03, 0x8a, 0x32, 0xcb, 0xb0, 0x2a, 0x71,
	0x0a, 0x3d, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c,
	0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xeb, 0xf4, 0xcc, 0x92,
	0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xcc, 0xa2, 0xcc, 0xe2, 0xbc, 0xd4, 0x12, 0xfd,
	0xe4, 0xfc, 0xa2, 0x54, 0xdd, 0xe2, 0x94, 0x6c, 0xdd, 0xf4, 0x7c, 0xfd, 0xe4, 0xfc, 0xdc, 0xdc,
	0xfc, 0x3c, 0x98, 0x57, 0x40, 0x1e, 0x40, 0xf8, 0x27, 0x89, 0x0d, 0xec, 0x38, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xfd, 0x7f, 0xff, 0x02, 0xf1, 0x00, 0x00, 0x00,
}

func (m *PubKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PubKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PubKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintKeys(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PrivKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PrivKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PrivKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintKeys(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintKeys(dAtA []byte, offset int, v uint64) int {
	offset -= sovKeys(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PubKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovKeys(uint64(l))
	}
	return n
}

func (m *PrivKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovKeys(uint64(l))
	}
	return n
}

func sovKeys(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozKeys(x uint64) (n int) {
	return sovKeys(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PubKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeys
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
			return fmt.Errorf("proto: PubKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PubKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeys
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthKeys
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeys
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipKeys(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthKeys
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthKeys
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
func (m *PrivKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeys
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
			return fmt.Errorf("proto: PrivKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PrivKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeys
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthKeys
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeys
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipKeys(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthKeys
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthKeys
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
func skipKeys(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowKeys
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
					return 0, ErrIntOverflowKeys
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowKeys
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
				return 0, ErrInvalidLengthKeys
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupKeys
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthKeys
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthKeys        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowKeys          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupKeys = fmt.Errorf("proto: unexpected end of group")
)
