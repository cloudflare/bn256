// +build amd64,!generic arm64,!generic

package bn256

import "golang.org/x/sys/cpu"

// This file contains forward declarations for the architecture-specific
// assembly implementations of these functions, provided that they exist.

var hasBMI2 = cpu.X86.HasBMI2

// go:noescape
func gfpNeg(c, a *gfP)

//go:noescape
func gfpAdd(c, a, b *gfP)

//go:noescape
func gfpSub(c, a, b *gfP)

//go:noescape
func gfpMul(c, a, b *gfP)
