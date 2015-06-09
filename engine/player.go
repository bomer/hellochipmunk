package engine

import (
	// "fmt"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math"
	"math/rand"
)

type Player struct {
	jumpTick     int32
	jumpTickStop int32
	jumping      bool

	speed float32

	Mass   int32
	Radius int32

	ball *chipmunk.Shape
	Body *chipmunk.Body
}

func NewPlayer(ballRadius, ballMass int32, game *Game) *Player {
	x := rand.Intn(350-115) + 115
	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(ballRadius))
	ball.SetElasticity(0.95)
	ball.SetFriction(1.5)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(x), 600.0})
	body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))

	body.AddShape(ball)

	game.Space.AddBody(body)

	return &Player{Body: ball.Body, ball: ball,
		Radius: ballRadius, Mass: ballMass,
		jumpTickStop: 2, jumpTick: 0, jumping: false, speed: 5}
}

func (self *Player) Update() {
	if self.jumping {

		if self.jumpTick == self.jumpTickStop {
			self.jumping = false
			self.jumpTick = 0
		}

		self.jumpTick += 1
	}

}

func (self *Player) MoveRight() {
	self.Body.AddAngularVelocity(-self.speed)
	self.Body.AddVelocity(2, 0)
}

func (self *Player) MoveLeft() {
	self.Body.AddAngularVelocity(self.speed)
	self.Body.AddVelocity(-2, 0)

}

func (self *Player) Jump() {
	if !self.jumping {
		self.jumpTickStop = 90
		self.jumpTick = 0
		self.Body.AddVelocity(0, 650)

		self.jumping = true
	}

}
