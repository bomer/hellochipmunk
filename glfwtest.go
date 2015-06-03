package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"log"
	"math"
	// "math/rand"
	"os"
	"runtime"
	"time"
)

const (
	BALL_RADIUS = 25
	GRAVITY     = -5
)

var (
	player *Ball
)

type Ball struct {
	x   float32
	y   float32
	rot float32

	max_y float32

	color_r float32
	color_g float32
	color_b float32
	color_a float32

	velocity_x float32
	velocity_y float32
}

func NewBall(x, y float32) *Ball {
	var velocity_x float32 = 0
	var velocity_y float32 = 0

	return &Ball{x: x, y: y, rot: 0,
		color_r: 0.1, color_g: 0.2, color_b: 0.5, color_a: 0.8,
		velocity_x: velocity_x, velocity_y: velocity_y}
}

func (b *Ball) update() {
	b.velocity_y += b.velocity_y * 0.01

	b.x += b.velocity_x

	if b.y >= b.max_y {
		b.velocity_y = -3
	}

	b.y += b.velocity_y + GRAVITY

	if b.y < 100 {
		b.y = 100
	}

	b.rot += 5 * (-b.velocity_x)

}

func (b *Ball) moveRight() {
	b.velocity_x = 3
	b.velocity_y = 0
}

func (b *Ball) moveLeft() {
	b.velocity_x = -3
	b.velocity_y = 0
}

func (b *Ball) moveUp() {
	b.velocity_y = 20
	b.max_y = b.y + 300
	// b.velocity_x = 0
}

func (b *Ball) moveDown() {
	// b.velocity_y = -3
	// b.velocity_x = 0

}

// key events are a way to get input from GLFW.
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//if u only want the on press, do = && && action == glfw.Press
	if key == glfw.KeyW { // && action == glfw.Press {
		fmt.Printf("W Pressed!\n")
		player.moveUp()
	}
	if key == glfw.KeyA { //&& action == glfw.Press
		fmt.Printf("A Pressed!\n")
		player.moveLeft()
	}
	if key == glfw.KeyS {
		fmt.Printf("S Pressed!\n")
		player.moveDown()
	}
	if key == glfw.KeyD {
		fmt.Printf("D Pressed!\n")
		player.moveRight()
	}

	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

// drawCircle draws a circle for the specified radius, rotation angle, and the specified number of sides
func drawCircle(radius float64, sides int) {
	gl.Begin(gl.LINE_LOOP)
	for a := 0.0; a < 2*math.Pi; a += (2 * math.Pi / float64(sides)) {
		gl.Vertex2d(math.Sin(a)*radius, math.Cos(a)*radius)
	}
	gl.Vertex3f(0, 0, 0)
	gl.End()
}

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()

	// gl.Begin(gl.LINES)
	// gl.Color3f(.2, .5, .2)
	// for i := range staticLines {
	// 	x := staticLines[i].GetAsSegment().A.X
	// 	y := staticLines[i].GetAsSegment().A.Y
	// 	gl.Vertex3f(float32(x), float32(y), 0)
	// 	x = staticLines[i].GetAsSegment().B.X
	// 	y = staticLines[i].GetAsSegment().B.Y
	// 	gl.Vertex3f(float32(x), float32(y), 0)
	// }
	// gl.End()

	gl.Color4f(player.color_r, player.color_g, player.color_b, player.color_a)

	//Draw Player
	gl.PushMatrix()

	rot := player.rot
	pos_x := player.x
	pos_y := player.y

	gl.Translatef(pos_x, pos_y, 0.0)
	gl.Rotatef(float32(rot), 0, 0, 1)
	drawCircle(float64(BALL_RADIUS), 60)
	gl.PopMatrix()
}

// onResize sets up a simple 2d ortho context based on the window size
func onResize(window *glfw.Window, w, h int) {
	w, h = window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(1, 1, 1, 1)
}

func main() {
	runtime.LockOSThread()

	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize GLFW: ", err)
	}
	defer glfw.Terminate()

	// create window
	window, err := glfw.CreateWindow(1280, 720, os.Args[0], nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.SetFramebufferSizeCallback(onResize)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	// set up opengl context
	onResize(window, 600, 600)

	// glfw.KeyCallback(window)
	window.SetKeyCallback(keyCallback)

	runtime.LockOSThread()
	glfw.SwapInterval(1)

	player = NewBall(600, 600)

	ticker := time.NewTicker(time.Second / 60)
	for !window.ShouldClose() {

		player.update()
		//Output
		draw()

		window.SwapBuffers()
		glfw.PollEvents()

		<-ticker.C // wait up to 1/60th of a second
	}
}
