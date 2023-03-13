// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/snapshots/v1beta1/snapshot.proto

package snapshots

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/gogo/protobuf/gogoproto"
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

// Snapshot contains Tendermint state sync snapshot info.
type Snapshot struct {
	Height   uint64   `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Format   uint32   `protobuf:"varint,2,opt,name=format,proto3" json:"format,omitempty"`
	Chunks   uint32   `protobuf:"varint,3,opt,name=chunks,proto3" json:"chunks,omitempty"`
	Hash     []byte   `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Metadata Metadata `protobuf:"bytes,5,opt,name=metadata,proto3" json:"metadata"`
}

func (m *Snapshot) Reset()         { *m = Snapshot{} }
func (m *Snapshot) String() string { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()    {}
func (*Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_d68466eaebf2e253, []int{0}
}
func (m *Snapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Snapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Snapshot.Merge(m, src)
}
func (m *Snapshot) XXX_Size() int {
	return m.Size()
}
func (m *Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_Snapshot proto.InternalMessageInfo

func (m *Snapshot) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Snapshot) GetFormat() uint32 {
	if m != nil {
		return m.Format
	}
	return 0
}

func (m *Snapshot) GetChunks() uint32 {
	if m != nil {
		return m.Chunks
	}
	return 0
}

func (m *Snapshot) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Snapshot) GetMetadata() Metadata {
	if m != nil {
		return m.Metadata
	}
	return Metadata{}
}

// Metadata contains SDK-specific snapshot metadata.
type Metadata struct {
	ChunkHashes [][]byte `protobuf:"bytes,1,rep,name=chunk_hashes,json=chunkHashes,proto3" json:"chunk_hashes,omitempty"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_d68466eaebf2e253, []int{1}
}
func (m *Metadata) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Metadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Metadata.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Metadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metadata.Merge(m, src)
}
func (m *Metadata) XXX_Size() int {
	return m.Size()
}
func (m *Metadata) XXX_DiscardUnknown() {
	xxx_messageInfo_Metadata.DiscardUnknown(m)
}

var xxx_messageInfo_Metadata proto.InternalMessageInfo

func (m *Metadata) GetChunkHashes() [][]byte {
	if m != nil {
		return m.ChunkHashes
	}
	return nil
}

func init() {
	proto.RegisterType((*Snapshot)(nil), "cosmos.base.snapshots.v1beta1.Snapshot")
	proto.RegisterType((*Metadata)(nil), "cosmos.base.snapshots.v1beta1.Metadata")
}

func init() {
	proto.RegisterFile("cosmos/snapshots/v1beta1/snapshot.proto", fileDescriptor_d68466eaebf2e253)
}

var fileDescriptor_d68466eaebf2e253 = []byte{
	// 297 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xb1, 0x4e, 0xeb, 0x30,
	0x14, 0x86, 0xe3, 0xdb, 0xdc, 0xaa, 0x72, 0xcb, 0x62, 0x21, 0x64, 0x21, 0x61, 0x42, 0x97, 0x66,
	0xa9, 0xad, 0x02, 0x4f, 0xd0, 0x09, 0x06, 0x18, 0xc2, 0xc6, 0x82, 0x9c, 0xd4, 0xc4, 0x51, 0x95,
	0x9e, 0x2a, 0xc7, 0x45, 0xe2, 0x2d, 0x78, 0x15, 0xde, 0xa2, 0x63, 0x47, 0x26, 0x84, 0xda, 0x17,
	0x41, 0x71, 0x42, 0xc4, 0xc4, 0x76, 0xfe, 0xcf, 0xdf, 0xd1, 0xb1, 0x7e, 0x3a, 0xc9, 0x00, 0x4b,
	0x40, 0x85, 0x2b, 0xbd, 0x46, 0x0b, 0x0e, 0xd5, 0xcb, 0x2c, 0x35, 0x4e, 0xcf, 0x3a, 0x22, 0xd7,
	0x15, 0x38, 0x60, 0x67, 0x8d, 0x28, 0x53, 0x8d, 0x46, 0x76, 0xb6, 0x6c, 0xed, 0xd3, 0xe3, 0x1c,
	0x72, 0xf0, 0xa6, 0xaa, 0xa7, 0x66, 0x69, 0xfc, 0x4e, 0xe8, 0xe0, 0xa1, 0x75, 0xd9, 0x09, 0xed,
	0x5b, 0x53, 0xe4, 0xd6, 0x71, 0x12, 0x91, 0x38, 0x4c, 0xda, 0x54, 0xf3, 0x67, 0xa8, 0x4a, 0xed,
	0xf8, 0xbf, 0x88, 0xc4, 0x47, 0x49, 0x9b, 0x6a, 0x9e, 0xd9, 0xcd, 0x6a, 0x89, 0xbc, 0xd7, 0xf0,
	0x26, 0x31, 0x46, 0x43, 0xab, 0xd1, 0xf2, 0x30, 0x22, 0xf1, 0x28, 0xf1, 0x33, 0xbb, 0xa5, 0x83,
	0xd2, 0x38, 0xbd, 0xd0, 0x4e, 0xf3, 0xff, 0x11, 0x89, 0x87, 0x97, 0x13, 0xf9, 0xe7, 0x87, 0xe5,
	0x5d, 0xab, 0xcf, 0xc3, 0xed, 0xe7, 0x79, 0x90, 0x74, 0xeb, 0xe3, 0x29, 0x1d, 0xfc, 0xbc, 0xb1,
	0x0b, 0x3a, 0xf2, 0x47, 0x9f, 0xea, 0x23, 0x06, 0x39, 0x89, 0x7a, 0xf1, 0x28, 0x19, 0x7a, 0x76,
	0xe3, 0xd1, 0xfc, 0x7e, 0xbb, 0x17, 0x64, 0xb7, 0x17, 0xe4, 0x6b, 0x2f, 0xc8, 0xdb, 0x41, 0x04,
	0xbb, 0x83, 0x08, 0x3e, 0x0e, 0x22, 0x78, 0xbc, 0xce, 0x0b, 0x67, 0x37, 0xa9, 0xcc, 0xa0, 0x54,
	0x45, 0x55, 0xe0, 0xca, 0x38, 0x95, 0x41, 0x65, 0xa6, 0xb8, 0x58, 0x4e, 0x73, 0x50, 0xe8, 0xa0,
	0x32, 0xca, 0xbd, 0xae, 0xcd, 0xaf, 0xfa, 0xd3, 0xbe, 0x6f, 0xee, 0xea, 0x3b, 0x00, 0x00, 0xff,
	0xff, 0x4b, 0xce, 0xad, 0x18, 0x99, 0x01, 0x00, 0x00,
}

func (m *Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Snapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Snapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Metadata.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintSnapshot(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintSnapshot(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x22
	}
	if m.Chunks != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.Chunks))
		i--
		dAtA[i] = 0x18
	}
	if m.Format != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.Format))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Metadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Metadata) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Metadata) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChunkHashes) > 0 {
		for iNdEx := len(m.ChunkHashes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ChunkHashes[iNdEx])
			copy(dAtA[i:], m.ChunkHashes[iNdEx])
			i = encodeVarintSnapshot(dAtA, i, uint64(len(m.ChunkHashes[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintSnapshot(dAtA []byte, offset int, v uint64) int {
	offset -= sovSnapshot(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Snapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovSnapshot(uint64(m.Height))
	}
	if m.Format != 0 {
		n += 1 + sovSnapshot(uint64(m.Format))
	}
	if m.Chunks != 0 {
		n += 1 + sovSnapshot(uint64(m.Chunks))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovSnapshot(uint64(l))
	}
	l = m.Metadata.Size()
	n += 1 + l + sovSnapshot(uint64(l))
	return n
}

func (m *Metadata) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ChunkHashes) > 0 {
		for _, b := range m.ChunkHashes {
			l = len(b)
			n += 1 + l + sovSnapshot(uint64(l))
		}
	}
	return n
}

func sovSnapshot(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSnapshot(x uint64) (n int) {
	return sovSnapshot(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Snapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSnapshot
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
			return fmt.Errorf("proto: Snapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Snapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Format", wireType)
			}
			m.Format = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Format |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chunks", wireType)
			}
			m.Chunks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Chunks |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
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
				return ErrInvalidLengthSnapshot
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSnapshot
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSnapshot
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSnapshot
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Metadata.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSnapshot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSnapshot
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
func (m *Metadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSnapshot
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
			return fmt.Errorf("proto: Metadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Metadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChunkHashes", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
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
				return ErrInvalidLengthSnapshot
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSnapshot
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChunkHashes = append(m.ChunkHashes, make([]byte, postIndex-iNdEx))
			copy(m.ChunkHashes[len(m.ChunkHashes)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSnapshot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSnapshot
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
func skipSnapshot(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSnapshot
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
					return 0, ErrIntOverflowSnapshot
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
					return 0, ErrIntOverflowSnapshot
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
				return 0, ErrInvalidLengthSnapshot
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSnapshot
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSnapshot
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSnapshot        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSnapshot          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSnapshot = fmt.Errorf("proto: unexpected end of group")
)
