package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	TARGET_FPS    = 60
	SECOND        = 1000 // 1000 milliseconds = 1 second
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600
	BOX_SPEED     = 10 // 10 pixels

	TIME_FORMAT = "2006-01-02 15:04:05"
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	PanicOnError(err)

	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Deltatime",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		SCREEN_WIDTH, SCREEN_HEIGHT,
		sdl.WINDOW_OPENGL,
	)
	PanicOnError(err)

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	PanicOnError(err)

	defer renderer.Destroy()

	box := &sdl.Rect{X: 100, Y: 100, W: 100, H: 100}

	var dt float64
	var frame int
	var prevTime float64
	var sec int

	S := &SecondController{t: time.Now()}

	var restarts int

GameLoop:
	for {

		window.SetTitle(fmt.Sprintf("%f", math.Abs(dt)))

		startFrame := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.QUIT:
				break GameLoop
			}
		}

		err = renderer.SetDrawColor(0, 255, 0, 0)
		PanicOnError(err)
		err = renderer.Clear()
		PanicOnError(err)

		err = renderer.SetDrawColor(255, 0, 0, sdl.ALPHA_TRANSPARENT)
		PanicOnError(err)
		err = renderer.FillRect(box)
		PanicOnError(err)

		box.X += int32(dt * BOX_SPEED)

		if box.X > SCREEN_WIDTH {
			box.X = -(0 + box.W/2)
			restarts++
		}

		renderer.Present()

		elapsed := time.Since(startFrame)
		dt = elapsed.Seconds() * TARGET_FPS

		fmt.Printf("%ds - Frame: %02d, StartFrame=%s, PrevTime=%f, FrameTime=%f, DeltaTime=%f, Expected=%f, X=%d, Restarts=%d\n",
			sec, frame, startFrame.Format(TIME_FORMAT), prevTime, elapsed.Seconds(), dt, (1.0 / TARGET_FPS), box.X, restarts)

		prevTime = elapsed.Seconds()

		if dt > 1.5 {
			dt = 1.5
		}

		frame++

		if S.HasSecondElapsed() {
			sec++
		}
	}
}

type SecondController struct {
	t time.Time
}

func (sc *SecondController) HasSecondElapsed() bool {
	elapsed := time.Since(sc.t)
	if elapsed.Seconds() >= 1 && elapsed.Seconds() < 2 {
		sc.t = time.Now()
		return true
	}
	return false
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
