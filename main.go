package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	ballRadius = 25
	ballMass   = 1

	space       *chipmunk.Space
	balls       []*chipmunk.Shape
	staticLines []*chipmunk.Shape
	deg2rad     = math.Pi / 180

	player *chipmunk.Shape
)

// key events are a way to get input from GLFW.
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//if u only want the on press, do = && && action == glfw.Press
	var speed float32
	speed = 7
	if key == glfw.KeyW && action == glfw.Press {
		fmt.Printf("W Pressed!\n")

		//Check if on floor first?
		player.Body.AddVelocity(0, 750)
	}
	if key == glfw.KeyA { //&& action == glfw.Press
		fmt.Printf("A Pressed!\n")
		player.Body.AddAngularVelocity(speed)
	}
	if key == glfw.KeyS {
		fmt.Printf("S Pressed!\n")
	}
	if key == glfw.KeyD {
		fmt.Printf("D Pressed!\n")
		player.Body.AddAngularVelocity(-speed)
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

	//Transform screen.
	gl.PushMatrix()
	gl.Translatef((1280/2)-float32((player.Body.Position().X)), 0, 0.0)

	gl.Begin(gl.LINES)
	gl.Color3f(.2, .5, .2)
	for i := range staticLines {
		x := staticLines[i].GetAsSegment().A.X
		y := staticLines[i].GetAsSegment().A.Y
		gl.Vertex3f(float32(x), float32(y), 0)
		x = staticLines[i].GetAsSegment().B.X
		y = staticLines[i].GetAsSegment().B.Y
		gl.Vertex3f(float32(x), float32(y), 0)
	}
	gl.End()

	gl.Color4f(.3, .3, 1, .8)
	// draw balls
	for _, ball := range balls {
		gl.PushMatrix()
		pos := ball.Body.Position()
		rot := ball.Body.Angle() * chipmunk.DegreeConst
		gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
		gl.Rotatef(float32(rot), 0, 0, 1)
		drawCircle(float64(ballRadius), 60)
		gl.PopMatrix()
	}

	//Draw Player
	gl.PushMatrix()
	pos := player.Body.Position()
	rot := player.Body.Angle() * chipmunk.DegreeConst
	gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
	gl.Rotatef(float32(rot), 0, 0, 1)
	drawCircle(float64(ballRadius), 60)
	gl.PopMatrix()

	gl.PopMatrix()
}

func addBall() {
	x := rand.Intn(350-115) + 115
	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(ballRadius))
	ball.SetElasticity(0.95)
	// ball.SetFriction(0.9)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(x), 600.0})
	body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))

	body.AddShape(ball)
	// space.AddBody(body)
	// balls = append(balls, ball)
	player = ball
	space.AddBody(body)
}

// step advances the physics engine and cleans up any balls that are off-screen
func step(dt float32) {
	space.Step(vect.Float(dt))

	//Gives the velocity some torque, stops going too fast/like friction
	player.Body.SetAngularVelocity(player.Body.AngularVelocity() * 0.92)

	for i := 0; i < len(balls); i++ {
		p := balls[i].Body.Position()
		if p.Y < -100 {
			space.RemoveBody(balls[i].Body)
			balls[i] = nil
			balls = append(balls[:i], balls[i+1:]...)
			i-- // consider same index again
		}
	}
}

// createBodies sets up the chipmunk space and static bodies
func createBodies() {
	space = chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, -900}

	staticBody := chipmunk.NewBodyStatic()
	staticLines = []*chipmunk.Shape{
		//TODO Load from a CSV instead
		chipmunk.NewSegment(vect.Vect{100.0, 200.0}, vect.Vect{407.0, 200.0}, 0),
		chipmunk.NewSegment(vect.Vect{407.0, 200.0}, vect.Vect{407.0, 343.0}, 0),
		chipmunk.NewSegment(vect.Vect{0.0, 100.0}, vect.Vect{500.0, 10.0}, 0),
		// chipmunk.NewSegment(vect.Vect{0.0, 100.0}, vect.Vect{500.0, 100.0}, 0),
		chipmunk.NewSegment(vect.Vect{5.0, 100.0}, vect.Vect{5.0, 500.0}, 0),
		chipmunk.NewSegment(vect.Vect{1280.0, 100.0}, vect.Vect{500.0, 10.0}, 0),
	}
	for _, segment := range staticLines {
		segment.SetElasticity(0.6)
		staticBody.AddShape(segment)
	}
	space.AddBody(staticBody)
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
	createBodies()

	addBall()

	//Init Controlls I think
	// glfw.KeyCallback(window)
	window.SetKeyCallback(keyCallback)

	runtime.LockOSThread()
	glfw.SwapInterval(1)

	ticksToNextBall := 2
	ticker := time.NewTicker(time.Second / 60)
	for !window.ShouldClose() {
		ticksToNextBall--
		if ticksToNextBall == 0 {
			ticksToNextBall = rand.Intn(100) + 1
			// addBall()
		}

		//Input Handling

		//Output
		draw()
		step(1.0 / 60.0)
		window.SwapBuffers()
		glfw.PollEvents()

		<-ticker.C // wait up to 1/60th of a second
	}
}
