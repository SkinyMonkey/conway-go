package main

import "time"
import "github.com/veandco/go-sdl2/sdl"

const (
	RowSize = 8
	CellSize = 10
	UniverseSize = RowSize * RowSize
)

const (
	WindowWidth = RowSize * CellSize
	WindowHeight = RowSize * CellSize
)

const (
	Dead byte = 0
	Alive byte = 1
)

// These helpers should be inlined by the compiler
func nextRow(index int) int {
  return index + RowSize
}

func previousRow(index int) int {
  return index - RowSize
}

func nextCell(index int) int {
  return index + 1
}

func previousCell(index int) int {
  return index - 1
}

func aliveNbrAround(index int, universe *[UniverseSize]byte) int {
  // Pre computed indices
  var indices = [...]int {
    previousCell(index),
    nextCell(index),

    nextRow(index),
    previousRow(index),

    nextCell(previousRow(index)),
    previousCell(previousRow(index)),

    nextCell(nextRow(index)),
    previousCell(nextRow(index)),
  }

  var result int
  for _, cellAroundindex := range indices {
    if (cellAroundindex > 0 &&
        cellAroundindex < UniverseSize &&
        universe[cellAroundindex] == Alive) {
      result = result + 1
    }
  }

  return result
}

// A dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
// A live cell with two or three live neighbors lives on to the next generation.
// A live cell with more than three live neighbors dies, as if by overpopulation.
func deadOrAlive(cell byte, index int, universe *[UniverseSize]byte) byte {
  var nbr = aliveNbrAround(index, universe)

  if (cell == Alive && nbr > 1 && nbr < 4) {
    return Alive
  } else if (nbr == 3) {
    return Alive
  }

  return Dead
}

func play(universe *[UniverseSize]byte) {
  var tmp[UniverseSize] byte

  for index ,_ := range universe {
    tmp[index] = deadOrAlive(universe[index], index, universe)
  }

	copy(universe[:], tmp[:]);
}

func setup(universe *[UniverseSize]byte) {
	// test
//	entity := [...]int{10, 17, 18, 19, 26}
//	entity := [...]int{9, 10, 11, 17, 19, 25, 26, 27}

	// spaceship
	entity := [...]int{ 27, 28, 29, 37, 44}

// NOTE : arbitrary
	// entity := [...]int{16, 17, 18, 31, 32, 33, 34, 35, 36, 58, 59, 60, 61, 62, 63, 64, 65, 66}

	for _, cell := range entity {
		universe[cell] = Alive
	}
}

// UI -----------------------------------------------------------------

func eventHandling(running bool) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			running = false
			break
		}
	}

	return running
}

func initSdl() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
}

func getWindow() *sdl.Window {
	window, err := sdl.CreateWindow(
		"Go conway game of life",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WindowHeight,
		WindowWidth,
		sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	return window
}

func getSurface(window *sdl.Window) *sdl.Surface {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	return surface
}

func drawUniverse(universe *[UniverseSize]byte, window *sdl.Window, surface *sdl.Surface)  {
	surface.FillRect(nil, 0) // Clear the screen

	for index, cell := range universe {
		if (cell == Alive) {
			x := int32(index % RowSize) * CellSize
			y := int32(index / RowSize) * CellSize

			rect := sdl.Rect{x, y, CellSize, CellSize}
			surface.FillRect(&rect, 0xffffffff)
		}
	}

	window.UpdateSurface()
}

// main ---------------------------------------------------------------

func main() {
	initSdl()
	defer sdl.Quit()

	window := getWindow()
	defer window.Destroy()

	surface := getSurface(window)
	surface.FillRect(nil, 0)

  var universe[UniverseSize] byte

  setup(&universe)

	running := true
	for running {
		time.Sleep(100 * time.Millisecond)
		drawUniverse(&universe, window, surface);
		play(&universe)

		running = eventHandling(running)
	}
}
