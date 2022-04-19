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

func (v VidMode) equals(other VidMode) bool {
	if v.RedBits+v.GreenBits+v.BlueBits != other.RedBits+other.GreenBits+other.BlueBits {
		return false
	}

	if v.Width != other.Width {
		return false
	}

	if v.Height != other.Height {
		return false
	}

	if v.RefreshRate != other.RefreshRate {
		return false
	}

	return true
}

func splitBPP(bpp int) (red, green, blue int) {
	// We assume that by 32 the user really meant 24
	if bpp == 32 {
		bpp = 24
	}

	// Convert "bits per pixel" to red, green & blue sizes
	red = bpp / 3
	green = bpp / 3
	blue = bpp / 3
	delta := bpp - (red * 3)
	if delta >= 1 {
		green++
	}
	if delta == 2 {
		red++
	}
	return
}
