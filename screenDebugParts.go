package izapple2

import (
	"image"
)

// SnapshotParts the currently visible screen
func (a *Apple2) SnapshotParts() *image.RGBA {
	videoMode := getCurrentVideoMode(a)
	isSecondPage := (videoMode & videoSecondPage) != 0
	videoBase := videoMode & videoBaseMask
	mixMode := videoMode & videoMixTextMask
	modifiers := videoMode & videoModifiersMask

	snapScreen := snapshotByMode(a, videoMode, ScreenModePlain)
	snapPage1 := snapshotByMode(a, videoMode&^videoSecondPage, ScreenModePlain)
	snapPage2 := snapshotByMode(a, videoMode|videoSecondPage, ScreenModePlain)
	var snapAux *image.RGBA

	/*
		if videoBase == videoRGBMix {
		_, mask := snapshotDoubleHiResModeMono(a, isSecondPage, true /*isRGBMixMode*/ /*, color.White)
		snapAux = filterMask(mask)
	}*/

	if videoBase == videoText40RGB {
		snapAux = snapshotText40RGBModeColors(a, isSecondPage)
	} else {
		switch mixMode {
		case videoMixText80:
			snapAux = snapshotByMode(a, videoText80|modifiers, ScreenModePlain)
		case videoMixText40RGB:
			snapAux = snapshotByMode(a, videoText40RGB|modifiers, ScreenModePlain)
		default:
			snapAux = snapshotByMode(a, videoText40|modifiers, ScreenModePlain)
		}
	}

	return mixFourSnapshots([]*image.RGBA{snapScreen, snapAux, snapPage1, snapPage2})
}

// VideoModeName returns the name of the current video mode
func (a *Apple2) VideoModeName() string {
	videoMode := getCurrentVideoMode(a)
	videoBase := videoMode & videoBaseMask
	mixMode := videoMode & videoMixTextMask

	var name string

	switch videoBase {
	case videoText40:
		name = "TEXT40COL"
	case videoText80:
		name = "TEXT80COL"
	case videoText40RGB:
		name = "TEXT40COLRGB"
	case videoGR:
		name = "GR"
	case videoDGR:
		name = "DGR"
	case videoHGR:
		name = "HGR"
	case videoDHGR:
		name = "DHGR"
	case videoMono560:
		name = "Mono560"
	case videoRGBMix:
		name = "RGMMIX"
	case videoRGB160:
		name = "RGB160"
	case videoSHR:
		name = "SHR"
	default:
		name = "Unknown video mode"
	}

	if (videoMode & videoSecondPage) != 0 {
		name += "-PAGE2"
	}

	switch mixMode {
	case videoMixText40:
		name += "-MIX40"
	case videoMixText80:
		name += "-MIX80"
	case videoMixText40RGB:
		name += "-MIX40RGB"
	}

	return name
}

func mixFourSnapshots(snaps []*image.RGBA) *image.RGBA {
	width := snaps[0].Rect.Dx()
	height := snaps[0].Rect.Dy()
	size := image.Rect(0, 0, width*2, height*2)
	out := image.NewRGBA(size)

	for i := 1; i < 4; i++ {
		if snaps[i].Bounds().Dx() < width {
			snaps[i] = doubleWidthFilter(snaps[i])
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			out.Set(x, y, snaps[0].At(x, y))
			out.Set(x+width, y, snaps[1].At(x, y))
			out.Set(x, y+height, snaps[2].At(x, y))
			out.Set(x+width, y+height, snaps[3].At(x, y))
		}
	}

	return out
}

func doubleWidthFilter(in *image.RGBA) *image.RGBA {
	b := in.Bounds()
	size := image.Rect(0, 0, 2*b.Dx(), b.Dy())
	out := image.NewRGBA(size)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := in.At(x, y)
			out.Set(2*x, y, c)
			out.Set(2*x+1, y, c)
		}
	}
	return out
}
