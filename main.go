package main

import (
	"encoding/csv"
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
	"strconv"
	"time"
)

var (
	ballRadius = 25
	ballMass   = 1

	space       *chipmunk.Space
	balls       []*chipmunk.Shape
	staticLines []*chipmunk.Shape
	deg2rad     = math.Pi / 180

	player     *chipmunk.Shape
	jumpTick   = 2
	eSpawnTick = 2

	canJump = true
)

// key events are a way to get input from GLFW.
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//if u only want the on press, do = && && action == glfw.Press
	var speed float32
	speed = 7
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

func readfile(filename string) []*chipmunk.Shape {
	csvfile, err := os.Open("level.csv")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output
	var ret []*chipmunk.Shape
	for _, each := range rawCSVdata {
		f1, _ := strconv.ParseFloat(each[0], 64)
		f2, _ := strconv.ParseFloat(each[1], 64)
		f3, _ := strconv.ParseFloat(each[2], 64)
		f4, _ := strconv.ParseFloat(each[3], 64)
		fmt.Printf("x1 : %f y1  : %f to x2: %f y2: %f \n", f1, f2, f3, f4)
		ret = append(ret, chipmunk.NewSegment(vect.Vect{vect.Float(f1), vect.Float(f2)}, vect.Vect{vect.Float(f3), vect.Float(f4)}, 0))
	}
	return ret
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

	gl.Color4f(.9, .1, 1, .9)
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
	gl.Color4f(.3, .3, 1, .8)
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

func addBall(isPlayer bool) {
	x := rand.Intn(350-115) + 115
	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(ballRadius))
	ball.SetElasticity(0.95)
	ball.SetFriction(1.5)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	if isPlayer {
		body.SetPosition(vect.Vect{vect.Float(x), 600.0})
	} else {
		x := rand.Intn(1280)
		body.SetPosition(vect.Vect{vect.Float(x), 800.0})
	}
	body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))

	body.AddShape(ball)
	if isPlayer {
		player = ball
	} else {
		balls = append(balls, ball)
	}

	space.AddBody(body)
}

// step advances the physics engine and cleans up any balls that are off-screen
func step(dt float32) {
	space.Step(vect.Float(dt))

	//Gives the velocity some torque, stops going too fast/like friction
	player.Body.SetAngularVelocity(player.Body.AngularVelocity() * 0.92)

	for i := 0; i < len(balls); i++ {
		p := balls[i].Body.Position()
		//Move Enemie TOwards player
		if p.Y < 300 { //Only move if on bottom part of screen
			if p.X < player.Body.Position().X {
				// balls[i].Body.AddVelocity(10, 0)
				balls[i].Body.AddAngularVelocity(-1)
			} else {
				// balls[i].Body.AddVelocity(-10, 0)
				balls[i].Body.AddAngularVelocity(1)
			}
		}

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
	// readfile("level.csv")
	space = chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, -900}

	staticBody := chipmunk.NewBodyStatic()
	staticLines = readfile("level.csv")
	for _, segment := range staticLines {
		segment.SetElasticity(0.6)
		staticBody.AddShape(segment)
	}
	space.AddBody(staticBody)
}

func addEnemies() {
	//Pass in Json read struct of enemies here + Do a loop
	addBall(false)
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

	addBall(true) // True for is Player

	addEnemies()

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
			addBall(false)
			eSpawnTick = 200
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
