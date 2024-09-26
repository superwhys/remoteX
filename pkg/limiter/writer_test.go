// File:		writer_test.go
// Created by:	Hoven
// Created on:	2024-08-02
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package limiter

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"golang.org/x/time/rate"
)

func TestLimitWriterWrite(t *testing.T) {
	src := make([]byte, int(12.5*8192))
	if _, err := rand.Reader.Read(src); err != nil {
		t.Fatal(err)
	}

	dst := new(bytes.Buffer)
	cw := &countingWriter{w: dst}
	lw := &LimitWriter{
		writer:     cw,
		baseWaiter: (*baseWaiter)(rate.NewLimiter(rate.Limit(42), limiterBurstSize)),
	}
	if _, err := io.Copy(lw, bytes.NewReader(src)); err != nil {
		t.Fatal(err)
	}

	t.Logf("write count: %d", len(src))

	// 限流器每秒允许 42 K
	// 根据 LimitWriter 单次写入计算, 一次写入大小为1024K
	// 因此要写入12.5 * 8192 字节，即 102400 字节的数据
	// 最小预期写入和最大预期写入为 10*8 和 15*8
	if cw.writeCount < 10*8 {
		t.Error("expected lots of smaller writes")
	}
	if cw.writeCount > 15*8 {
		t.Error("expected fewer larger writes")
	}
	if !bytes.Equal(src, dst.Bytes()) {
		t.Error("results should be equal")
	}
}

type countingWriter struct {
	w          io.Writer
	writeCount int
}

func (w *countingWriter) Write(data []byte) (int, error) {
	w.writeCount++
	return w.w.Write(data)
}
