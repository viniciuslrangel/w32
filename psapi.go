// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"reflect"
	"syscall"
	"unsafe"
)

var (
	modpsapi = syscall.NewLazyDLL("psapi.dll")

	procEnumProcesses        = modpsapi.NewProc("EnumProcesses")
	procEnumProcessModulesEx = modpsapi.NewProc("EnumProcessModulesEx")
)

func EnumProcesses(processIds []uint32, cb uint32, bytesReturned *uint32) bool {
	ret, _, _ := procEnumProcesses.Call(
		uintptr(unsafe.Pointer(&processIds[0])),
		uintptr(cb),
		uintptr(unsafe.Pointer(bytesReturned)))

	return ret != 0
}

func EnumProcessModules(hProcess HANDLE, filterFlags DWORD) ([]HMODULE, error) {
	modules := make([]HMODULE, 128)
	size := uintptr(0)
	for {
		modulesSize := uintptr(len(modules)) * reflect.TypeOf(modules).Elem().Size()
		ret, _, _ := procEnumProcessModulesEx.Call(
			uintptr(hProcess),
			uintptr(unsafe.Pointer(&modules[0])),
			modulesSize,
			uintptr(unsafe.Pointer(&size)),
			uintptr(filterFlags),
		)
		if ret == 0 {
			return nil, syscall.Errno(GetLastError())
		}
		retSize := size / reflect.TypeOf(modules).Elem().Size()
		if size > modulesSize {
			modules = make([]HMODULE, retSize)
			continue
		}
		return modules[:retSize], nil
	}
}
