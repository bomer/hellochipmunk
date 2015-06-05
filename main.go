package main

import (
	"fmt"

	"github.com/bomer/hellochipmunk/engine"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"log"
	"math"
	"os"
	"runtime"

	"time"
)

var (
	game *engine.Game

	jumpTick   = 2
	eSpawnTick = 2

	canJump = true
)

// key events are a way to get input from GLFW.
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//if u only want the on press, do = && && action == glfw.Press
	var speed float32
	speed = 7

	player := game.Player

	if key == glfw.KeyW && action == glfw.Press {
		fmt.Printf("W Pressed!\n")
		//Jump is controlled by a 1.5 second timer for now, should do a collision detection but that seems hard.
		if canJump {
			//Check if on floor first?
			jumpTick = 90
			canJump = false
			player.Body.AddVelocity(0, 650)
		}

	}
	if key == glfw.KeyA { //&& action == glfw.Press
		fmt.Printf("A Pressed!\n")
		player.Body.AddAngularVelocity(speed)
		player.Body.AddVelocity(-2, 0)
	}
	if key == glfw.KeyS {
		fmt.Printf("S Pressed!\n")
	}
	if key == glfw.KeyD {
		fmt.Printf("D Pressed!\n")
		player.Body.AddAngularVelocity(-speed)
		player.Body.AddVelocity(2, 0)
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

// OpenGL draw function
func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()

	player := game.Player

	//Transform screen.
	gl.PushMatrix()
	gl.Translatef((1280/2)-float32((player.Body.Position().X)), 0, 0.0)

	gl.Begin(gl.LINES)
	gl.Color3f(.2, .5, .2)
	for _, segment := range game.Level.GetChipmunkSegments() {
		x := segment.GetAsSegment().A.X
		y := segment.GetAsSegment().A.Y
		gl.Vertex3f(float32(x), float32(y), 0)
		x = segment.GetAsSegment().B.X
		y = segment.GetAsSegment().B.Y
		gl.Vertex3f(float32(x), float32(y), 0)
	}
	gl.End()

	gl.Color4f(.9, .1, 1, .9)
	// draw balls
	for _, enemy := range game.Enemies {
		gl.PushMatrix()
		pos := enemy.Body.Position()
		rot := enemy.Body.Angle() * game.DegreeConst
		gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
		gl.Rotatef(float32(rot), 0, 0, 1)
		drawCircle(float64(enemy.Radius), 60)
		gl.PopMatrix()
	}
	gl.Color4f(.3, .3, 1, .8)
	//Draw Player
	gl.PushMatrix()
	pos := player.Body.Position()
	rot := player.Body.Angle() * game.DegreeConst
	gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
	gl.Rotatef(float32(rot), 0, 0, 1)
	drawCircle(float64(player.Radius), 60)
	gl.PopMatrix()

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
	// player = createPlayer()
	// addBall()
	// set up physics

	game = engine.NewGame()

	game.Init()
	game.ReadLevelFromFile("level.json")

	//Init Controlls I think
	// glfw.KeyCallback(window)
	window.SetKeyCallback(keyCallback)

	runtime.LockOSThread()
	glfw.SwapInterval(1)

	ticker := time.NewTicker(time.Second / 60)
	for !window.ShouldClose() {

		jumpTick--
		eSpawnTick--
		if jumpTick == 0 {
			//rand.Intn(100) + 1
			// addBall()
			canJump = true
		}
		if eSpawnTick == 0 {
			game.SpawnEnemy()
			eSpawnTick = 200
		}

		game.Update(1.0 / 60.0)

		draw()
		window.SwapBuffers()
		glfw.PollEvents()

		<-ticker.C // wait up to 1/60th of a second
	}
}
