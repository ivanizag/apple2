package screen

import (
	"image"
	"image/color"
)

const (
	doubleHiResWidth = 2 * hiResWidth
	rgb160Width      = 4 * 160
)

func snapshotDoubleHiRes(vs VideoSource, isSecondPage bool, getNTSCMask bool, light color.Color) (*image.RGBA, *image.Alpha) {
	dataMain := vs.GetVideoMemory(isSecondPage, false)
	dataAux := vs.GetVideoMemory(isSecondPage, true)
	return renderDoubleHiRes(dataMain, dataAux, getNTSCMask, light)
}

func renderDoubleHiRes(dataMain []uint8, dataAux []uint8, getNTSCMask bool, light color.Color) (*image.RGBA, *image.Alpha) {

	// As described in "Inside the Apple IIe"
	size := image.Rect(0, 0, doubleHiResWidth, hiResHeight)
	img := image.NewRGBA(size)

	// To support RGB-mode 14 we will have a mask to mark where we should not have the NTSC filter applied
	// See: https://apple2online.com/web_documents/Video-7%20Manual%20KB.pdf
	var ntscMask *image.Alpha
	if getNTSCMask {
		ntscMask = image.NewAlpha(size)
	}

	for y := 0; y < hiResHeight; y++ {
		offset := getHiResLineOffset(y)
		lineParts := [][]uint8{
			dataAux[offset : offset+hiResLineBytes],
			dataMain[offset : offset+hiResLineBytes],
		}
		x := 0
		// For the NTSC filter to work we have to insert an initial black pixel and skip the last one ¿?
		img.Set(x, y, color.Black)
		if getNTSCMask {
			ntscMask.Set(x, y, color.Opaque)
		}
		x++
		for iByte := 0; iByte < hiResLineBytes; iByte++ {
			for iPart := 0; iPart < 2; iPart++ {
				b := lineParts[iPart][iByte]

				mask := color.Transparent // Apply the NTSC filter
				if getNTSCMask && b&0x80 == 0 {
					mask = color.Opaque // Do not apply the NTSC filter
				}

				for j := uint(0); j < 7; j++ {
					// Set color
					bit := (b >> j) & 1
					colour := light
					if bit == 0 {
						colour = color.Black
					}
					img.Set(x, y, colour)

					// Set mask if requested
					if getNTSCMask {
						ntscMask.Set(x, y, mask)
					}
					x++
				}
			}
		}
	}
	return img, ntscMask
}

func snapshotDoubleHiRes160(vs VideoSource, isSecondPage bool, light color.Color) *image.RGBA {
	dataMain := vs.GetVideoMemory(isSecondPage, false)
	dataAux := vs.GetVideoMemory(isSecondPage, true)
	return renderDoubleHiRes160(dataMain, dataAux, light)
}

func renderDoubleHiRes160(dataMain []uint8, dataAux []uint8, light color.Color) *image.RGBA {
	size := image.Rect(0, 0, rgb160Width, hiResHeight)
	img := image.NewRGBA(size)

	for y := 0; y < hiResHeight; y++ {
		offset := getHiResLineOffset(y)
		lineParts := [][]uint8{
			dataAux[offset : offset+hiResLineBytes],
			dataMain[offset : offset+hiResLineBytes],
		}
		x := 0
		for iByte := 0; iByte < hiResLineBytes; iByte++ {
			for iPart := 0; iPart < 2; iPart++ {
				b := lineParts[iPart][iByte]
				for j := uint(0); j < 8; j++ {
					// Set color
					bit := (b >> j) & 1
					colour := light
					if bit == 0 {
						colour = color.Black
					}
					img.Set(x, y, colour)
					x++
				}
			}
		}
	}
	return img
}
