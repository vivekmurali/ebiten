// Copyright 2019 The Ebiten Authors
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

package atlas

import (
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

const (
	BaseCountToPutOnAtlas = baseCountToPutOnAtlas
	PaddingSize           = paddingSize
)

func PutImagesOnAtlasForTesting(graphicsDriver graphicsdriver.Graphics) error {
	return putImagesOnAtlas(graphicsDriver)
}

var (
	oldMinSize int
	oldMaxSize int
)

func SetImageSizeForTesting(min, max int) {
	oldMinSize = min
	oldMaxSize = max
	minSize = min
	maxSize = max
}

func ResetImageSizeForTesting() {
	minSize = oldMinSize
	maxSize = oldMaxSize
}

func (i *Image) IsOnAtlasForTesting() bool {
	backendsM.Lock()
	defer backendsM.Unlock()
	return i.isOnAtlas()
}

func (i *Image) EnsureIsolatedForTesting() {
	backendsM.Lock()
	defer backendsM.Unlock()
	i.ensureIsolated()
}

func ResolveDeferredForTesting() {
	resolveDeferred()
}
