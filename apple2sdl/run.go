package apple2sdl

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	"go6502/apple2"
)

// SDLRun starts the Apple2 emulator on SDL
func SDLRun(a *apple2.Apple2) {
	window, renderer, err := sdl.CreateWindowAndRenderer(800, 600,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic("Failed to create window")
	}
	window.SetResizable(true)

	defer window.Destroy()
	defer renderer.Destroy()
	window.SetTitle("Apple2")

	kp := newSDLKeyBoard()
	a.SetKeyboardProvider(&kp)
	go a.Run(false, false)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				kp.putKey(t)
			case *sdl.TextInputEvent:
				//fmt.Printf("[%d ms] TextInput\ttype:%d\texts:%s\n",
				//	t.Timestamp, t.Type, t.GetText())
				kp.putText(t)
			}
		}

		img := *apple2.Snapshot(a)
		surface, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&img.Pix[0]), 40*7, 24*8, 32, 40*7*4,
			0xff000000, 0x00ff0000, 0x0000ff00, 0x000000ff)
		if err != nil {
			panic(err)
		}

		texture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err)
		}

		renderer.Clear()
		w, h := window.GetSize()
		renderer.Copy(texture, nil, &sdl.Rect{X: 0, Y: 0, W: w, H: h})
		renderer.Present()

		surface.Free()
		texture.Destroy()

		sdl.Delay(1000 / 60)
	}

}