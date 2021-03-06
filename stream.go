package gmf

/*

#cgo pkg-config: libavformat libavcodec

#include "libavformat/avformat.h"
#include "libavcodec/avcodec.h"

*/
import "C"

import (
	"fmt"
)

type Stream struct {
	AvStream *C.struct_AVStream
	SwsCtx   *SwsCtx
	SwrCtx   *SwrCtx
	AvFifo   *AVAudioFifo
	cc       *CodecCtx
	Pts      int64
	CgoMemoryManage
}

func (s *Stream) Free() {
	if s.SwsCtx != nil {
		s.SwsCtx.Free()
	}
	if s.SwrCtx != nil {
		s.SwrCtx.Free()
	}
	if s.AvFifo != nil {
		s.AvFifo.Free()
	}
}

func (s *Stream) DumpContexCodec(codec *CodecCtx) {
	ret := C.avcodec_parameters_from_context(s.AvStream.codecpar, codec.avCodecCtx)
	if ret < 0 {
		panic("Failed to copy context from input to output stream codec context\n")
	}
}

func (s *Stream) SetCodecFlags() {
	s.AvStream.codec.flags |= C.AV_CODEC_FLAG_GLOBAL_HEADER
}

func (s *Stream) CodecCtx() *CodecCtx {
	// Supposed that output context is set and opened by user
	if s.IsCodecCtxSet() {
		return s.cc
	}

	// Open input codec context
	c, err := FindDecoder(int(s.AvStream.codec.codec_id))
	if err != nil {
		return nil
	}

	if s.cc = NewCodecCtx(c); s.cc == nil {
		panic("error allocating codec context")
	}

	ret := int(C.avcodec_parameters_to_context(s.cc.avCodecCtx, s.AvStream.codecpar))
	if ret < 0 {
		panic("error copying parameters to codec context")
	}

	if err := s.cc.Open(nil); err != nil {
		panic("error opening codec context")
	}

	s.cc.avCodecCtx.time_base = s.AvStream.codec.time_base

	return s.cc
}

func (s *Stream) SetCodecCtx(cc *CodecCtx) {
	s.cc = cc
}

func (s *Stream) SetCodecParameters(cp *CodecParameters) error {
	if cp == nil || cp.avCodecParameters == nil {
		return fmt.Errorf("codec parameters are not initialized")
	}

	s.AvStream.codecpar = cp.avCodecParameters
	return nil
}

func (s *Stream) IsCodecCtxSet() bool {
	return (s.cc != nil)
}

func (s *Stream) Index() int {
	return int(s.AvStream.index)
}

func (s *Stream) Id() int {
	return int(s.AvStream.id)
}

func (s *Stream) NbFrames() int {
	if int(s.AvStream.nb_frames) == 0 {
		return 1
	}

	return int(s.AvStream.nb_frames)
}

func (s *Stream) TimeBase() AVRational {
	return AVRational(s.AvStream.time_base)
}

func (s *Stream) Type() int32 {
	return s.CodecCtx().Type()
}

func (s *Stream) IsAudio() bool {
	return (s.Type() == AVMEDIA_TYPE_AUDIO)
}

func (s *Stream) IsVideo() bool {
	return (s.Type() == AVMEDIA_TYPE_VIDEO)
}

func (s *Stream) Duration() int64 {
	return int64(s.AvStream.duration)
}

func (s *Stream) SetTimeBase(val AVR) *Stream {
	s.AvStream.time_base.num = C.int(val.Num)
	s.AvStream.time_base.den = C.int(val.Den)
	return s
}

func (s *Stream) GetRFrameRate() AVRational {
	return AVRational(s.AvStream.r_frame_rate)
}

func (s *Stream) SetRFrameRate(val AVR) {
	s.AvStream.r_frame_rate.num = C.int(val.Num)
	s.AvStream.r_frame_rate.den = C.int(val.Den)
}

func (s *Stream) SetAvgFrameRate(val AVR) {
	s.AvStream.avg_frame_rate.num = C.int(val.Num)
	s.AvStream.avg_frame_rate.den = C.int(val.Den)
}

func (s *Stream) GetAvgFrameRate() AVRational {
	return AVRational(s.AvStream.avg_frame_rate)
}

func (s *Stream) GetStartTime() int64 {
	return int64(s.AvStream.start_time)
}

func (s *Stream) GetCodecPar() *CodecParameters {
	cp := NewCodecParameters()
	cp.avCodecParameters = s.AvStream.codecpar

	return cp
}

func (s *Stream) CopyCodecPar(cp *CodecParameters) error {
	ret := int(C.avcodec_parameters_copy(s.AvStream.codecpar, cp.avCodecParameters))
	if ret < 0 {
		return AvError(ret)
	}

	return nil
}
