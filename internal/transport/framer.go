package transport

import "golang.org/x/net/http2"

type GRPCFramer interface {
	ErrorDetail() error
	ReadFrame() (http2.Frame, error)
	SetMaxReadFrameSize(uint32)
	SetReuseFrames()
	WriteContinuation(uint32, bool, []byte) error
	WriteData(uint32, bool, []byte) error
	WriteGoAway(uint32, http2.ErrCode, []byte) error
	WriteHeaders(http2.HeadersFrameParam) error
	WritePing(bool, [8]byte) error
	WriteRSTStream(uint32, http2.ErrCode) error
	WriteRawFrame(http2.FrameType, http2.Flags, uint32, []byte) error
	WriteSettings(settings ...http2.Setting) error
	WriteSettingsAck() error
	WriteWindowUpdate(uint32, uint32) error
}
