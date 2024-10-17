// File:		adler.go
// Created by:	Hoven
// Created on:	2024-10-12
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package hash

const (
	CharOffset = uint32(0)
)

func CheckAdlerSum(buf1 []byte) uint32 {
	var s1, s2 uint32
	s1, s2 = 0, 0
	length := len(buf1)
	
	i := 0
	for i = 0; i < length-4; i += 4 {
		s2 += 4*(s1+uint32(buf1[i])) + 3*uint32(buf1[i+1]) + 2*uint32(buf1[i+2]) + uint32(buf1[i+3]) + 10*CharOffset
		s1 += uint32(buf1[i]) + uint32(buf1[i+1]) + uint32(buf1[i+2]) + uint32(buf1[i+3]) + 4*CharOffset
	}
	
	for ; i < length; i++ {
		s1 += uint32(buf1[i]) + CharOffset
		s2 += s1
	}
	
	return (s1 & 0xffff) + (s2 << 16)
}

func RollingUpdate(oldChecksum uint32, oldByte, newByte byte, blockLength uint32) uint32 {
	s1 := uint32(oldChecksum & 0xffff)
	s2 := uint32(oldChecksum >> 16)
	
	s1 = (s1 + uint32(newByte) - uint32(oldByte)) & 0xffff
	s2 = (s2 + s1 - blockLength*(uint32(oldByte)+CharOffset)) & 0xffff
	return (s1 & 0xffff) + (s2 << 16)
}
