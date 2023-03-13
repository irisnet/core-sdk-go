// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/store/v1beta1/commit_info.proto

package store

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

// CommitInfo defines commit information used by the multi-store when committing
// a version/height.
type CommitInfo struct {
	Version    int64       `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	StoreInfos []StoreInfo `protobuf:"bytes,2,rep,name=store_infos,json=storeInfos,proto3" json:"store_infos"`
}

func (m *CommitInfo) Reset()         { *m = CommitInfo{} }
func (m *CommitInfo) String() string { return proto.CompactTextString(m) }
func (*CommitInfo) ProtoMessage()    {}
func (*CommitInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f8c656cdef8c524, []int{0}
}
func (m *CommitInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CommitInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CommitInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CommitInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitInfo.Merge(m, src)
}
func (m *CommitInfo) XXX_Size() int {
	return m.Size()
}
func (m *CommitInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitInfo.DiscardUnknown(m)
}

var xxx_messageInfo_CommitInfo proto.InternalMessageInfo

func (m *CommitInfo) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *CommitInfo) GetStoreInfos() []StoreInfo {
	if m != nil {
		return m.StoreInfos
	}
	return nil
}

// StoreInfo defines store-specific commit information. It contains a reference
// between a store name and the commit ID.
type StoreInfo struct {
	Name     string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CommitId CommitID `protobuf:"bytes,2,opt,name=commit_id,json=commitId,proto3" json:"commit_id"`
}

func (m *StoreInfo) Reset()         { *m = StoreInfo{} }
func (m *StoreInfo) String() string { return proto.CompactTextString(m) }
func (*StoreInfo) ProtoMessage()    {}
func (*StoreInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f8c656cdef8c524, []int{1}
}
func (m *StoreInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StoreInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StoreInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StoreInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StoreInfo.Merge(m, src)
}
func (m *StoreInfo) XXX_Size() int {
	return m.Size()
}
func (m *StoreInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_StoreInfo.DiscardUnknown(m)
}

var xxx_messageInfo_StoreInfo proto.InternalMessageInfo

func (m *StoreInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StoreInfo) GetCommitId() CommitID {
	if m != nil {
		return m.CommitId
	}
	return CommitID{}
}

// CommitID defines the committment information when a specific store is
// committed.
type CommitID struct {
	Version int64  `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Hash    []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *CommitID) Reset()      { *m = CommitID{} }
func (*CommitID) ProtoMessage() {}
func (*CommitID) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f8c656cdef8c524, []int{2}
}
func (m *CommitID) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CommitID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CommitID.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CommitID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitID.Merge(m, src)
}
func (m *CommitID) XXX_Size() int {
	return m.Size()
}
func (m *CommitID) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitID.DiscardUnknown(m)
}

var xxx_messageInfo_CommitID proto.InternalMessageInfo

func (m *CommitID) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *CommitID) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func init() {
	proto.RegisterType((*CommitInfo)(nil), "cosmos.base.store.v1beta1.CommitInfo")
	proto.RegisterType((*StoreInfo)(nil), "cosmos.base.store.v1beta1.StoreInfo")
	proto.RegisterType((*CommitID)(nil), "cosmos.base.store.v1beta1.CommitID")
}

func init() {
	proto.RegisterFile("cosmos/store/v1beta1/commit_info.proto", fileDescriptor_5f8c656cdef8c524)
}

var fileDescriptor_5f8c656cdef8c524 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xb1, 0x4e, 0xf3, 0x30,
	0x14, 0x85, 0xe3, 0x36, 0xfa, 0xff, 0xd6, 0x65, 0xb2, 0x18, 0x02, 0x83, 0x5b, 0x95, 0x0a, 0x75,
	0xa9, 0xad, 0x96, 0x8d, 0xa1, 0x43, 0x41, 0x48, 0x15, 0x5b, 0xd8, 0x58, 0x50, 0x92, 0xba, 0x89,
	0x85, 0x92, 0x8b, 0x72, 0x4d, 0x9f, 0x83, 0x91, 0x91, 0xc7, 0xe9, 0xd8, 0x91, 0x09, 0xa1, 0xe4,
	0x45, 0x50, 0x9c, 0x84, 0x8d, 0x6e, 0xe7, 0xe6, 0x9e, 0x7b, 0x3e, 0x9d, 0x98, 0x5e, 0x46, 0x80,
	0x29, 0xa0, 0x44, 0x03, 0xb9, 0x92, 0xbb, 0x79, 0xa8, 0x4c, 0x30, 0x97, 0x11, 0xa4, 0xa9, 0x36,
	0x4f, 0x3a, 0xdb, 0x82, 0x78, 0xc9, 0xc1, 0x00, 0x3b, 0xab, 0x7d, 0x22, 0x0c, 0x50, 0x09, 0x6b,
	0x16, 0x8d, 0xf9, 0xfc, 0x34, 0x86, 0x18, 0xac, 0x4b, 0x56, 0xaa, 0x3e, 0x18, 0x23, 0xa5, 0x37,
	0x36, 0x65, 0x9d, 0x6d, 0x81, 0x79, 0xf4, 0xff, 0x4e, 0xe5, 0xa8, 0x21, 0xf3, 0xc8, 0x88, 0x4c,
	0xbb, 0x7e, 0x3b, 0xb2, 0x7b, 0x3a, 0xb0, 0x71, 0x16, 0x86, 0x5e, 0x67, 0xd4, 0x9d, 0x0e, 0x16,
	0x13, 0xf1, 0x27, 0x4e, 0x3c, 0x54, 0x53, 0x15, 0xba, 0x72, 0xf7, 0x5f, 0x43, 0xc7, 0xa7, 0xd8,
	0x7e, 0xc0, 0x71, 0x4c, 0xfb, 0xbf, 0x6b, 0xc6, 0xa8, 0x9b, 0x05, 0xa9, 0xb2, 0xc0, 0xbe, 0x6f,
	0x35, 0xbb, 0xa3, 0xfd, 0xb6, 0xdb, 0xc6, 0xeb, 0x8c, 0xc8, 0x74, 0xb0, 0xb8, 0x38, 0xc2, 0x6a,
	0x1a, 0xdc, 0x36, 0xa8, 0x5e, 0x7d, 0xbb, 0xde, 0x8c, 0x97, 0xb4, 0xd7, 0xee, 0x8e, 0x74, 0x63,
	0xd4, 0x4d, 0x02, 0x4c, 0x2c, 0xe8, 0xc4, 0xb7, 0xfa, 0xda, 0x7d, 0xff, 0x18, 0x3a, 0xab, 0xe5,
	0xbe, 0xe0, 0xe4, 0x50, 0x70, 0xf2, 0x5d, 0x70, 0xf2, 0x56, 0x72, 0xe7, 0x50, 0x72, 0xe7, 0xb3,
	0xe4, 0xce, 0xe3, 0x24, 0xd6, 0x26, 0x79, 0x0d, 0x45, 0x04, 0xa9, 0xd4, 0xb9, 0xc6, 0x4c, 0x19,
	0x19, 0x41, 0xae, 0x66, 0xb8, 0x79, 0x9e, 0xc5, 0x50, 0x3f, 0x54, 0xf8, 0xcf, 0xfe, 0xe4, 0xab,
	0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcb, 0xf1, 0xa3, 0x51, 0xbf, 0x01, 0x00, 0x00,
}

func (m *CommitInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CommitInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CommitInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.StoreInfos) > 0 {
		for iNdEx := len(m.StoreInfos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.StoreInfos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCommitInfo(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Version != 0 {
		i = encodeVarintCommitInfo(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *StoreInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StoreInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StoreInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.CommitId.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintCommitInfo(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintCommitInfo(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CommitID) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CommitID) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CommitID) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintCommitInfo(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x12
	}
	if m.Version != 0 {
		i = encodeVarintCommitInfo(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommitInfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommitInfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *CommitInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Version != 0 {
		n += 1 + sovCommitInfo(uint64(m.Version))
	}
	if len(m.StoreInfos) > 0 {
		for _, e := range m.StoreInfos {
			l = e.Size()
			n += 1 + l + sovCommitInfo(uint64(l))
		}
	}
	return n
}

func (m *StoreInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovCommitInfo(uint64(l))
	}
	l = m.CommitId.Size()
	n += 1 + l + sovCommitInfo(uint64(l))
	return n
}

func (m *CommitID) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Version != 0 {
		n += 1 + sovCommitInfo(uint64(m.Version))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovCommitInfo(uint64(l))
	}
	return n
}

func sovCommitInfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommitInfo(x uint64) (n int) {
	return sovCommitInfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CommitInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommitInfo
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
			return fmt.Errorf("proto: CommitInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CommitInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StoreInfos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitInfo
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
				return ErrInvalidLengthCommitInfo
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommitInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StoreInfos = append(m.StoreInfos, StoreInfo{})
			if err := m.StoreInfos[len(m.StoreInfos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommitInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommitInfo
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
func (m *StoreInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommitInfo
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
			return fmt.Errorf("proto: StoreInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StoreInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommitInfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommitInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommitId", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitInfo
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
				return ErrInvalidLengthCommitInfo
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommitInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CommitId.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommitInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommitInfo
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
func (m *CommitID) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommitInfo
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
			return fmt.Errorf("proto: CommitID: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CommitID: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitInfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitInfo
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
				return ErrInvalidLengthCommitInfo
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCommitInfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommitInfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommitInfo
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
func skipCommitInfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommitInfo
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
					return 0, ErrIntOverflowCommitInfo
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
					return 0, ErrIntOverflowCommitInfo
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
				return 0, ErrInvalidLengthCommitInfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommitInfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommitInfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommitInfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommitInfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommitInfo = fmt.Errorf("proto: unexpected end of group")
)
