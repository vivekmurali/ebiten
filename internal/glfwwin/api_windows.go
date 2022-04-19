// Copyright 2022 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package glfwwin

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

// For the definitions, see https://github.com/wine-mirror/wine
const (
	_CCHDEVICENAME                  = 32
	_CCHFORMNAME                    = 32
	_CDS_TEST                       = 0x00000002
	_DISP_CHANGE_SUCCESSFUL         = 0
	_ENUM_CURRENT_SETTINGS   uint32 = 0xffffffff
	_LOGPIXELSX                     = 88
	_LOGPIXELSY                     = 90
	_TLS_OUT_OF_INDEXES      uint32 = 0xffffffff
	_USER_DEFAULT_SCREEN_DPI        = 96
)

type (
	_HCURSOR    windows.Handle
	_HDC        windows.Handle
	_HDEVNOTIFY windows.Handle
	_HGLRC      windows.Handle
	_HICON      windows.Handle
	_HMONITOR   windows.Handle
)

type _MONITOR_DPI_TYPE int32

const (
	_MDT_EFFECTIVE_DPI _MONITOR_DPI_TYPE = 0
	_MDT_ANGULAR_DPI   _MONITOR_DPI_TYPE = 1
	_MDT_RAW_DPI       _MONITOR_DPI_TYPE = 2
	_MDT_DEFAULT       _MONITOR_DPI_TYPE = _MDT_EFFECTIVE_DPI
)

type _DEVMODEW struct {
	dmDeviceName       [_CCHDEVICENAME]uint16
	dmSpecVersion      uint16
	dmDriverVersion    uint16
	dmSize             uint16
	dmDriverExtra      uint16
	dmFields           uint32
	_                  [16]byte // union
	dmColor            int16
	dmDuplex           int16
	dmYResolution      int16
	dmTTOption         int16
	dmCollate          int16
	dmFormName         [_CCHFORMNAME]uint32
	dmLogPixels        uint16
	dmBitsPerPel       uint32
	dmPelsWidth        uint32
	dmPelsHeight       uint32
	dmDisplayFlags     uint32 // union with DWORD dmNup
	dmDisplayFrequency uint32
	dmICMMethod        uint32
	dmICMIntent        uint32
	dmMediaType        uint32
	dmDitherType       uint32
	dmReserved1        uint32
	dmReserved2        uint32
	dmPanningWidth     uint32
	dmPanningHeight    uint32
}

type _MONITORINFO struct {
	cbSize    uint32
	rcMonitor _RECT
	rcWork    _RECT
	dwFlags   uint32
}

type _RAWINPUT struct {
	header _RAWINPUTHEADER
	mouse  _RAWMOUSE

	// RAWMOUSE is the biggest among RAWHID, RAWKEYBOARD, and RAWMOUSE.
	// Then, padding is not needed here.
}

type _RAWINPUTHEADER struct {
	dwType  uint32
	dwSize  uint32
	hDevice windows.Handle
	wParam  uintptr
}

type _RAWMOUSE struct {
	usFlags            uint16
	ulButtons          uint32 // TODO: Check alignments
	ulRawButtons       uint32
	lLastX             int32
	lLastY             int32
	ulExtraInformation uint32
}

type _RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shcore32 = windows.NewLazySystemDLL("shcore32.dll")
	user32   = windows.NewLazySystemDLL("user32.dll")

	procIsWindows8Point1OrGreater = kernel32.NewProc("IsWindows8Point1OrGreater")
	procTlsAlloc                  = kernel32.NewProc("TlsAlloc")
	procTlsFree                   = kernel32.NewProc("TlsFree")
	procTlsGetValue               = kernel32.NewProc("TlsGetValue")
	procTlsSetValue               = kernel32.NewProc("TlsSetValue")

	procGetDpiForMonitor = shcore32.NewProc("GetDpiForMonitor")

	procChangeDisplaySettingsExW = user32.NewProc("ChangeDisplaySettingsExW")
	procEnumDisplaySettingsW     = user32.NewProc("EnumDisplaySettingsW")
	procGetMonitorInfoW          = user32.NewProc("GetMonitorInfoW")
)

func _ChangeDisplaySettingsExW(deviceName string, lpDevMode *_DEVMODEW, hwnd windows.HWND, dwflags uint32, lParam unsafe.Pointer) int32 {
	lpszDeviceName, err := windows.UTF16PtrFromString(deviceName)
	if err != nil {
		panic("glfwwin: device name must not include a NUL character")
	}

	r, _, _ := procChangeDisplaySettingsExW.Call(uintptr(unsafe.Pointer(lpszDeviceName)), uintptr(unsafe.Pointer(lpDevMode)), uintptr(hwnd), uintptr(dwflags), uintptr(lParam))
	runtime.KeepAlive(lpszDeviceName)
	runtime.KeepAlive(lpDevMode)

	return int32(r)
}

func _EnumDisplaySettingsW(deviceName string, iModeNum uint32) (_DEVMODEW, bool) {
	lpszDeviceName, err := windows.UTF16PtrFromString(deviceName)
	if err != nil {
		panic("glfwwin: device name must not include a NUL character")
	}

	var dm _DEVMODEW
	dm.dmSize = uint16(unsafe.Sizeof(dm))

	r, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(lpszDeviceName)), uintptr(iModeNum), uintptr(unsafe.Pointer(&dm)))
	runtime.KeepAlive(lpszDeviceName)

	if r == 0 {
		return _DEVMODEW{}, false
	}
	return dm, true
}

func _GetMonitorInfoW(hMonitor _HMONITOR) (_MONITORINFO, bool) {
	var mi _MONITORINFO
	mi.cbSize = uint32(unsafe.Sizeof(mi))
	r, _, _ := procGetMonitorInfoW.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(&mi)))
	if r == 0 {
		return _MONITORINFO{}, false
	}
	return mi, true
}

func _GetDpiForMonitor(hmonitor _HMONITOR, dpiType _MONITOR_DPI_TYPE) (dpiX, dpiY uint32, err error) {
	r, _, e := procGetDpiForMonitor.Call(uintptr(hmonitor), uintptr(dpiType), uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
	if r != 0 {
		return 0, 0, fmt.Errorf("glfwwin: GetDpiForMonitor failed: %w", e)
	}
	return dpiX, dpiY, nil
}

func _IsWindows8Point1OrGreater() bool {
	r, _, _ := procIsWindows8Point1OrGreater.Call()
	return r != 0
}

func _TlsAlloc() (uint32, error) {
	r, _, e := procTlsAlloc.Call()
	if uint32(r) == _TLS_OUT_OF_INDEXES {
		return 0, fmt.Errorf("glfwwin: TlsAlloc failed: %w", e)
	}
	return uint32(r), nil
}

func _TlsFree(dwTlsIndex uint32) error {
	r, _, e := procTlsFree.Call(uintptr(dwTlsIndex))
	if r == 0 {
		return fmt.Errorf("glfwwin: TlsFree failed: %w", e)
	}
	return nil
}

func _TlsGetValue(dwTlsIndex uint32) (uintptr, error) {
	r, _, e := procTlsGetValue.Call(uintptr(dwTlsIndex))
	if r == 0 && e != windows.ERROR_SUCCESS {
		return 0, fmt.Errorf("glfwwin: TlsGetValue failed: %w", e)
	}
	return r, nil
}

func _TlsSetValue(dwTlsIndex uint32, lpTlsValue uintptr) error {
	r, _, e := procTlsSetValue.Call(uintptr(dwTlsIndex), lpTlsValue)
	if r == 0 {
		return fmt.Errorf("glfwwin: TlsSetValue failed: %w", e)
	}
	return nil
}
