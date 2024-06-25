/*
 *
 * Copyright 2024 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package transport

import (
	"io"

	"golang.org/x/net/http2"
)

type GRPCFramer interface {
	ErrorDetail() error
	ReadFrame() (http2.Frame, error)
	WriteContinuation(streamID uint32, endHeaders bool, headerBlockFragment []byte) error
	WriteData(streamID uint32, endStream bool, data []byte) error
	WriteGoAway(streamID uint32, code http2.ErrCode, debugData []byte) error
	WriteHeaders(p http2.HeadersFrameParam) error
	WritePing(bool, [8]byte) error
	WriteRSTStream(streamID uint32, code http2.ErrCode) error
	WriteRawFrame(t http2.FrameType, flags http2.Flags, streamID uint32, payload []byte) error
	WriteSettings(settings ...http2.Setting) error
	WriteSettingsAck() error
	WriteWindowUpdate(streamID, incr uint32) error
}

type Framer struct {
	w io.Writer

	stubFramer *http2.Framer
	hdrBuf     [9]byte
}

func newGRPCFramer(w io.Writer, r io.Reader) *Framer {
	return &Framer{
		w:          w,
		stubFramer: http2.NewFramer(w, r),
	}
}

func (f *Framer) writeFrameHeader(size uint32, fType http2.FrameType, flags http2.Flags, streamID uint32) error {
	f.hdrBuf[0] = byte(size >> 16)
	f.hdrBuf[1] = byte(size >> 8)
	f.hdrBuf[2] = byte(size)
	f.hdrBuf[3] = byte(fType)
	f.hdrBuf[4] = byte(flags)
	f.hdrBuf[5] = byte(streamID >> 24)
	f.hdrBuf[6] = byte(streamID >> 16)
	f.hdrBuf[7] = byte(streamID >> 8)
	f.hdrBuf[8] = byte(streamID)

	_, err := f.w.Write(f.hdrBuf[:])
	return err
}

func (f *Framer) ErrorDetail() error {
	return f.stubFramer.ErrorDetail()
}

func (f *Framer) ReadFrame() (http2.Frame, error) {
	return f.stubFramer.ReadFrame()
}

func (f *Framer) WriteContinuation(streamID uint32, endHeaders bool, headerBlockFragment []byte) error {
	return f.stubFramer.WriteContinuation(streamID, endHeaders, headerBlockFragment)
}

func (f *Framer) WriteData(streamID uint32, endStream bool, data []byte) error {
	var flags http2.Flags
	if endStream {
		flags = http2.FlagDataEndStream
	}

	err := f.writeFrameHeader(uint32(len(data)), http2.FrameData, flags, streamID)
	if err != nil {
		return err
	}

	f.w.Write(data)
	return nil
}

func (f *Framer) WriteDataN(streamID uint32, endStream bool, data ...[]byte) error {
	tSize := 0
	for _, slice := range data {
		tSize += len(slice)
	}

	var flags http2.Flags
	if endStream {
		flags = http2.FlagDataEndStream
	}

	err := f.writeFrameHeader(uint32(tSize), http2.FrameData, flags, streamID)
	if err != nil {
		return err
	}

	for _, slice := range data {
		_, err = f.w.Write(slice)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Framer) WriteGoAway(maxStreamID uint32, code http2.ErrCode, debugData []byte) error {
	return f.stubFramer.WriteGoAway(maxStreamID, code, debugData)
}

func (f *Framer) WriteHeaders(p http2.HeadersFrameParam) error {

	var flags http2.Flags
	if p.EndStream {
		flags |= http2.FlagHeadersEndStream
	}
	if p.EndHeaders {
		flags |= http2.FlagHeadersEndHeaders
	}

	err := f.writeFrameHeader(uint32(len(p.BlockFragment)), http2.FrameHeaders, flags, p.StreamID)
	if err != nil {
		return err
	}

	_, err = f.w.Write(p.BlockFragment)
	return err
}

func (f *Framer) WritePing(ack bool, data [8]byte) error {
	return f.stubFramer.WritePing(ack, data)
}

func (f *Framer) WriteRSTStream(streamID uint32, code http2.ErrCode) error {
	return f.stubFramer.WriteRSTStream(streamID, code)
}

func (f *Framer) WriteRawFrame(t http2.FrameType, flags http2.Flags, streamID uint32, payload []byte) error {
	return f.stubFramer.WriteRawFrame(t, flags, streamID, payload)
}

func (f *Framer) WriteSettings(settings ...http2.Setting) error {
	return f.stubFramer.WriteSettings(settings...)
}

func (f *Framer) WriteSettingsAck() error {
	return f.stubFramer.WriteSettingsAck()
}

func (f *Framer) WriteWindowUpdate(streamID uint32, incr uint32) error {
	return f.stubFramer.WriteWindowUpdate(streamID, incr)
}
