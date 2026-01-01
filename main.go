package main 

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ALIVE 	= 1
	DEAD 	= 0
)

const (
	defaultAliveProbability = 0.2
	headerMargin            = 40
)


type Game struct {
	N			int
	cellSize	int
	margin		int
	grid		[][]int
	nextGrid	[][]int
	generation	int

	aliveImg	*ebiten.Image

}

func initGrid(N int, aliveProb float64) [][]int {

	grid := make([][]int, N)
	for y := 0; y < N; y++ {
		row := make([]int, N)
		for x := 0; x < N; x++ {
			if float64(rand.Float64()) < aliveProb {
				row[x] = ALIVE
			} else {
				row[x] = DEAD
			}
		}
		grid[y] = row
	}
	return grid
}

func newGame(N, cellSize int) *Game {
	g := &Game{
		N:			N,
		cellSize: 	cellSize,
		margin:		headerMargin,
		grid:		initGrid(N, defaultAliveProbability),
		nextGrid: 	make([][]int, N),
	}

	for i := range g.nextGrid {
		g.nextGrid[i] = make([]int, N)
	}

	aliveImg := ebiten.NewImage(cellSize, cellSize)
	aliveImg.Fill(color.White)
	g.aliveImg = aliveImg

	return g

}

func (g *Game) countAliveNeighbors(x, y int) int {
	
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx := x + dx
			ny := y + dy

			// Avoid index panic cuz by neighbors outside the grid N
			if nx < 0 || ny < 0 || nx >= g.N || ny >= g.N {
				continue
			}
			// Avoid their own pixel
			if nx == x && ny == y {
				continue
			}
			// Simple counter
			if g.grid[ny][nx] == ALIVE {
				count++
			}
		}
	}
	return count
}


// step applies one Game of Life generation
func (g *Game) step() {
	for y := 0; y < g.N; y++ {
		for x := 0; x < g.N; x++ {
			aliveNeighbors := g.countAliveNeighbors(x, y)
			cell := g.grid[y][x]

			switch cell {
				case ALIVE:
					if aliveNeighbors < 2 || aliveNeighbors > 3 {
						g.nextGrid[y][x] = DEAD
					} else {
						g.nextGrid[y][x] = ALIVE
					}
				case DEAD:
					if aliveNeighbors == 3 {
						g.nextGrid[y][x] = ALIVE
					} else {
						g.nextGrid[y][x] = DEAD
					}
			}
		}
	}

	// Swap buffers
	g.grid, g.nextGrid = g.nextGrid, g.grid
	g.generation++
}

// toggleCell flips a cell between ALIVE and DEAD 
func (g *Game) toggleCell(x, y int) {
	if x < 0 || y < 0 || x >= g.N || y >= g.N {
		return
	}
	if g.grid[y][x] == ALIVE {
		g.grid[y][x] = DEAD
	} else {
		g.grid[y][x] = ALIVE
	}
}

// interactAtPixel toggles the clicked cell and its neighbors
func (g *Game) interactAtPixel(px, py int) {
	// Ignore clicks in the header area
	if py < g.margin {
		return
	}

	gridX := px / g.cellSize
	gridY := (py - g.margin) / g.cellSize

	if gridX < 0 || gridY < 0 || gridX >= g.N || gridY >= g.N {
		return
	}

	// Toggle a 3x3 neighborhood around (gridX, gridY)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx := gridX + dx
			ny := gridY + dy
			g.toggleCell(nx, ny)
		}
	}
}

// Update is called every tick.
func (g *Game) Update() error {
	// ESC close
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Mouse interaction
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.interactAtPixel(x, y)
	}

	// Advance the simulation every frame
	g.step()
	return nil
}

// Draw renders the current state
func (g *Game) Draw(screen *ebiten.Image) {
	// Background.
	screen.Fill(color.Black)

	// Draw alive cells
	for y := 0; y < g.N; y++ {
		for x := 0; x < g.N; x++ {
			if g.grid[y][x] == ALIVE {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(
					float64(x*g.cellSize),
					float64(y*g.cellSize+g.margin),
				)
				screen.DrawImage(g.aliveImg, op)
			}
		}
	}

	// Generation counterr
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Generation: %d", g.generation),
		10,
		10,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.N * g.cellSize, g.N*g.cellSize + g.margin
}

func main() {

	N, cellSize := 70, 12

	game := newGame(N, cellSize)

	ebiten.SetWindowTitle("GomeOfLife")
	ebiten.SetWindowSize(N*cellSize, N*cellSize+game.margin)
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}

}