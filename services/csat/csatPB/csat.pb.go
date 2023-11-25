//
//protoc --go_out=. --go_opt=paths=source_relative \
//--go-grpc_out=. --go-grpc_opt=paths=source_relative \
//csatPB/csat.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.11.2
// source: csatPB/csat.proto

package csatPB

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{0}
}

type UserID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *UserID) Reset() {
	*x = UserID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserID) ProtoMessage() {}

func (x *UserID) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserID.ProtoReflect.Descriptor instead.
func (*UserID) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{1}
}

func (x *UserID) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

type Question struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Question string `protobuf:"bytes,1,opt,name=question,proto3" json:"question,omitempty"`
}

func (x *Question) Reset() {
	*x = Question{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Question) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Question) ProtoMessage() {}

func (x *Question) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Question.ProtoReflect.Descriptor instead.
func (*Question) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{2}
}

func (x *Question) GetQuestion() string {
	if x != nil {
		return x.Question
	}
	return ""
}

type QuestionList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Questions []*Question `protobuf:"bytes,1,rep,name=questions,proto3" json:"questions,omitempty"`
}

func (x *QuestionList) Reset() {
	*x = QuestionList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuestionList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuestionList) ProtoMessage() {}

func (x *QuestionList) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuestionList.ProtoReflect.Descriptor instead.
func (*QuestionList) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{3}
}

func (x *QuestionList) GetQuestions() []*Question {
	if x != nil {
		return x.Questions
	}
	return nil
}

type Answer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Starts  int32  `protobuf:"varint,1,opt,name=starts,proto3" json:"starts,omitempty"`
	Comment string `protobuf:"bytes,2,opt,name=comment,proto3" json:"comment,omitempty"`
}

func (x *Answer) Reset() {
	*x = Answer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{4}
}

func (x *Answer) GetStarts() int32 {
	if x != nil {
		return x.Starts
	}
	return 0
}

func (x *Answer) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

type StarsNum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StarsNum int32 `protobuf:"varint,1,opt,name=starsNum,proto3" json:"starsNum,omitempty"`
	Count    int64 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *StarsNum) Reset() {
	*x = StarsNum{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StarsNum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StarsNum) ProtoMessage() {}

func (x *StarsNum) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StarsNum.ProtoReflect.Descriptor instead.
func (*StarsNum) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{5}
}

func (x *StarsNum) GetStarsNum() int32 {
	if x != nil {
		return x.StarsNum
	}
	return 0
}

func (x *StarsNum) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type QuestionComment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Comment string `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
}

func (x *QuestionComment) Reset() {
	*x = QuestionComment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuestionComment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuestionComment) ProtoMessage() {}

func (x *QuestionComment) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuestionComment.ProtoReflect.Descriptor instead.
func (*QuestionComment) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{6}
}

func (x *QuestionComment) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

type QuestionStatistics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AvgStars            float32            `protobuf:"fixed32,1,opt,name=avgStars,proto3" json:"avgStars,omitempty"`
	StarsNumList        []*StarsNum        `protobuf:"bytes,2,rep,name=starsNumList,proto3" json:"starsNumList,omitempty"`
	QuestionCommentList []*QuestionComment `protobuf:"bytes,3,rep,name=questionCommentList,proto3" json:"questionCommentList,omitempty"`
}

func (x *QuestionStatistics) Reset() {
	*x = QuestionStatistics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuestionStatistics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuestionStatistics) ProtoMessage() {}

func (x *QuestionStatistics) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuestionStatistics.ProtoReflect.Descriptor instead.
func (*QuestionStatistics) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{7}
}

func (x *QuestionStatistics) GetAvgStars() float32 {
	if x != nil {
		return x.AvgStars
	}
	return 0
}

func (x *QuestionStatistics) GetStarsNumList() []*StarsNum {
	if x != nil {
		return x.StarsNumList
	}
	return nil
}

func (x *QuestionStatistics) GetQuestionCommentList() []*QuestionComment {
	if x != nil {
		return x.QuestionCommentList
	}
	return nil
}

type Statistics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatisticsList []*QuestionStatistics `protobuf:"bytes,1,rep,name=statisticsList,proto3" json:"statisticsList,omitempty"`
}

func (x *Statistics) Reset() {
	*x = Statistics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_csatPB_csat_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Statistics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Statistics) ProtoMessage() {}

func (x *Statistics) ProtoReflect() protoreflect.Message {
	mi := &file_csatPB_csat_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Statistics.ProtoReflect.Descriptor instead.
func (*Statistics) Descriptor() ([]byte, []int) {
	return file_csatPB_csat_proto_rawDescGZIP(), []int{8}
}

func (x *Statistics) GetStatisticsList() []*QuestionStatistics {
	if x != nil {
		return x.StatisticsList
	}
	return nil
}

var File_csatPB_csat_proto protoreflect.FileDescriptor

var file_csatPB_csat_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x73, 0x61, 0x74, 0x50, 0x42, 0x2f, 0x63, 0x73, 0x61, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x73, 0x61, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x20, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x22, 0x26, 0x0a, 0x08, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3c, 0x0a, 0x0c,
	0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x09,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x63, 0x73, 0x61, 0x74, 0x2e, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x09, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x3a, 0x0a, 0x06, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x72, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x72, 0x74, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x3c, 0x0a, 0x08, 0x53, 0x74, 0x61, 0x72, 0x73, 0x4e,
	0x75, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72, 0x73, 0x4e, 0x75, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x74, 0x61, 0x72, 0x73, 0x4e, 0x75, 0x6d, 0x12, 0x14,
	0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x2b, 0x0a, 0x0f, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x22, 0xad, 0x01, 0x0a, 0x12, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x76, 0x67, 0x53,
	0x74, 0x61, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x61, 0x76, 0x67, 0x53,
	0x74, 0x61, 0x72, 0x73, 0x12, 0x32, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x73, 0x4e, 0x75, 0x6d,
	0x4c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x73, 0x61,
	0x74, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x73, 0x4e, 0x75, 0x6d, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x72,
	0x73, 0x4e, 0x75, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x47, 0x0a, 0x13, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x73, 0x61, 0x74, 0x2e, 0x51, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x13, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x22, 0x4e, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12,
	0x40, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x4c, 0x69, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x73, 0x61, 0x74, 0x2e, 0x51,
	0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63,
	0x73, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x4c, 0x69, 0x73,
	0x74, 0x32, 0x9a, 0x01, 0x0a, 0x04, 0x63, 0x73, 0x61, 0x74, 0x12, 0x32, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x0c, 0x2e, 0x63, 0x73, 0x61,
	0x74, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x12, 0x2e, 0x63, 0x73, 0x61, 0x74, 0x2e,
	0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x2d,
	0x0a, 0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72,
	0x12, 0x0c, 0x2e, 0x63, 0x73, 0x61, 0x74, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x1a, 0x0b,
	0x2e, 0x63, 0x73, 0x61, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2f, 0x0a,
	0x0c, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x12, 0x0b, 0x2e,
	0x63, 0x73, 0x61, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x63, 0x73, 0x61,
	0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x22, 0x00, 0x42, 0x1e,
	0x5a, 0x1c, 0x48, 0x6e, 0x48, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2f, 0x63, 0x73, 0x61, 0x74, 0x2f, 0x63, 0x73, 0x61, 0x74, 0x50, 0x42, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_csatPB_csat_proto_rawDescOnce sync.Once
	file_csatPB_csat_proto_rawDescData = file_csatPB_csat_proto_rawDesc
)

func file_csatPB_csat_proto_rawDescGZIP() []byte {
	file_csatPB_csat_proto_rawDescOnce.Do(func() {
		file_csatPB_csat_proto_rawDescData = protoimpl.X.CompressGZIP(file_csatPB_csat_proto_rawDescData)
	})
	return file_csatPB_csat_proto_rawDescData
}

var file_csatPB_csat_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_csatPB_csat_proto_goTypes = []interface{}{
	(*Empty)(nil),              // 0: csat.Empty
	(*UserID)(nil),             // 1: csat.UserID
	(*Question)(nil),           // 2: csat.Question
	(*QuestionList)(nil),       // 3: csat.QuestionList
	(*Answer)(nil),             // 4: csat.Answer
	(*StarsNum)(nil),           // 5: csat.StarsNum
	(*QuestionComment)(nil),    // 6: csat.QuestionComment
	(*QuestionStatistics)(nil), // 7: csat.QuestionStatistics
	(*Statistics)(nil),         // 8: csat.Statistics
}
var file_csatPB_csat_proto_depIdxs = []int32{
	2, // 0: csat.QuestionList.questions:type_name -> csat.Question
	5, // 1: csat.QuestionStatistics.starsNumList:type_name -> csat.StarsNum
	6, // 2: csat.QuestionStatistics.questionCommentList:type_name -> csat.QuestionComment
	7, // 3: csat.Statistics.statisticsList:type_name -> csat.QuestionStatistics
	1, // 4: csat.csat.GetQuestions:input_type -> csat.UserID
	4, // 5: csat.csat.RegisterAnswer:input_type -> csat.Answer
	0, // 6: csat.csat.GetStatistic:input_type -> csat.Empty
	3, // 7: csat.csat.GetQuestions:output_type -> csat.QuestionList
	0, // 8: csat.csat.RegisterAnswer:output_type -> csat.Empty
	8, // 9: csat.csat.GetStatistic:output_type -> csat.Statistics
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_csatPB_csat_proto_init() }
func file_csatPB_csat_proto_init() {
	if File_csatPB_csat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_csatPB_csat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_csatPB_csat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserID); i {
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
		file_csatPB_csat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Question); i {
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
		file_csatPB_csat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuestionList); i {
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
		file_csatPB_csat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Answer); i {
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
		file_csatPB_csat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StarsNum); i {
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
		file_csatPB_csat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuestionComment); i {
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
		file_csatPB_csat_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuestionStatistics); i {
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
		file_csatPB_csat_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Statistics); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_csatPB_csat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_csatPB_csat_proto_goTypes,
		DependencyIndexes: file_csatPB_csat_proto_depIdxs,
		MessageInfos:      file_csatPB_csat_proto_msgTypes,
	}.Build()
	File_csatPB_csat_proto = out.File
	file_csatPB_csat_proto_rawDesc = nil
	file_csatPB_csat_proto_goTypes = nil
	file_csatPB_csat_proto_depIdxs = nil
}
