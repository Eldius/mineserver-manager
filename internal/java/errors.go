package java

import "errors"

var (
	// ErrVersionNotSupported is returned when a java version is not supported
	ErrVersionNotSupported = errors.New("java version not supported")
	// ErrPlatformNotSupported is returned when a platform (OS/Arch) is not supported
	ErrPlatformNotSupported = errors.New("platform not supported")
)
