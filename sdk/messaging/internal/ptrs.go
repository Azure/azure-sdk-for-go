// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package internal

// PtrBool takes a boolean and returns a pointer to that bool. For use in literal pointers, ptrBool(true) -> *bool
func PtrBool(toPtr bool) *bool {
	return &toPtr
}

// PtrString takes a string and returns a pointer to that string. For use in literal pointers,
// PtrString(fmt.Sprintf("..", foo)) -> *string
func PtrString(toPtr string) *string {
	return &toPtr
}

// PtrInt32 takes a int32 and returns a pointer to that int32. For use in literal pointers, ptrInt32(1) -> *int32
func PtrInt32(number int32) *int32 {
	return &number
}

// PtrInt64 takes a int64 and returns a pointer to that int64. For use in literal pointers, ptrInt64(1) -> *int64
func PtrInt64(number int64) *int64 {
	return &number
}
