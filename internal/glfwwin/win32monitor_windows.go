// Copyright 2002-2006 Marcus Geelnard
// Copyright 2006-2019 Camilla Löwy
// Copyright 2022 The Ebiten Authors
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.

package glfwwin

// TODO: Add more functions!

func getMonitorContentScaleWin32(handle _HMONITOR) (xscale, yscale float32, err error) {
	var xdpi, ydpi uint32

	if _IsWindows8Point1OrGreater() {
		var err error
		xdpi, ydpi, err = _GetDpiForMonitor(handle, _MDT_EFFECTIVE_DPI)
		if err != nil {
			return 0, 0, err
		}
	} else {
		dc := _GetDC(nil)
		xdpi = _GetDeviceCaps(dc, _LOGPIXELSX)
		ydpi = _GetDeviceCaps(dc, _LOGPIXELSY)
		_ReleaseDC(nil, dc)
	}

	xscale = float32(xdpi) / _USER_DEFAULT_SCREEN_DPI
	yscale = float32(ydpi) / _USER_DEFAULT_SCREEN_DPI
	return
}

// TODO: Add more functions!

func (m *Monitor) platformGetMonitorContentScale() (xscale, yscale float32, err error) {
	return getMonitorContentScaleWin32(m.win32.handle)
}

func (m *Monitor) platformGetMonitorWorkarea() (xpos, ypos, width, height int) {
	mi, ok := _GetMonitorInfoW(m.win32.handle)
	if !ok {
		return 0, 0, 0, 0
	}
	return int(mi.rcWork.left), int(mi.rcWork.top), int(mi.rcWork.right - mi.rcWork.left), int(mi.rcWork.bottom - mi.rcWork.top)
}

func (m *Monitor) platformAppendVideoModes(monitors []VidMode) ([]VidMode, error) {
	origLen := len(monitors)
	var modeIndex uint32
loop:
	for {
		dm, ok := _EnumDisplaySettingsW(m.win32.adapterName, modeIndex)
		if !ok {
			break
		}
		modeIndex++

		// Skip modes with less than 15 BPP
		if dm.dmBitsPerPel < 15 {
			continue
		}

		r, g, b := splitBPP(int(dm.dmBitsPerPel))
		mode := VidMode{
			Width:       int(dm.dmPelsWidth),
			Height:      int(dm.dmPelsHeight),
			RefreshRate: int(dm.dmDisplayFrequency),
			RedBits:     r,
			GreenBits:   g,
			BlueBits:    b,
		}

		// Skip duplicate modes
		for _, m := range monitors[origLen:] {
			if m.equals(mode) {
				continue loop
			}
		}

		if m.win32.modesPruned {
			// Skip modes not supported by the connected displays
			if _ChangeDisplaySettingsExW(m.win32.adapterName, &dm, 0, _CDS_TEST, nil) != _DISP_CHANGE_SUCCESSFUL {
				continue
			}
		}

		monitors = append(monitors, mode)
	}

	if len(monitors) == origLen {
		// HACK: Report the current mode if no valid modes were found
		if m, ok := m.platformGetVideoMode(); ok {
			monitors = append(monitors, m)
		}
	}

	return monitors, nil
}

func (m *Monitor) platformGetVideoMode() (VidMode, bool) {
	dm, ok := _EnumDisplaySettingsW(m.win32.adapterName, _ENUM_CURRENT_SETTINGS)
	if !ok {
		return VidMode{}, false
	}
	r, g, b := splitBPP(int(dm.dmBitsPerPel))
	return VidMode{
		Width:       int(dm.dmPelsWidth),
		Height:      int(dm.dmPelsHeight),
		RefreshRate: int(dm.dmDisplayFrequency),
		RedBits:     r,
		GreenBits:   g,
		BlueBits:    b,
	}, true
}

func (m *Monitor) platformGetGammaRamp(ramp *GammaRamp) {
	panic("glfwwin: platformGetGammaRamp is not implemented")
}

func (m *Monitor) platformSetGammaRamp(ramp *GammaRamp) {
	panic("glfwwin: platformSetGammaRamp is not implemented")
}

func (m *Monitor) in32Adapter() (string, error) {
	if !_glfw.initialized {
		return "", NotInitialized
	}
	return m.win32.adapterName, nil
}

func (m *Monitor) win32Monitor() (string, error) {
	if !_glfw.initialized {
		return "", NotInitialized
	}
	return m.win32.displayName, nil
}
