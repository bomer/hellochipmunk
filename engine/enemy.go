package engine

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math"
	"math/rand"
)

type Enemy struct {
	alive bool

	Mass   int32
	Radius int32

	ball *chipmunk.Shape
	Body *chipmunk.Body
}

func NewEnemy(ballRadius, ballMass int32, game *Game) *Enemy {
	x := rand.Intn(1280)
	ball := chipmunk.NewCircle(vect.Vector_Zero, float32(ballRadius))
	ball.SetElasticity(0.95)
	ball.SetFriction(1.5)

	body := chipmunk.NewBody(vect.Float(ballMass), ball.Moment(float32(ballMass)))
	body.SetPosition(vect.Vect{vect.Float(x), 800.0})
	body.SetPosition(vect.Vect{vect.Float(x), 600.0})
	body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))

	body.AddShape(ball)

	game.Space.AddBody(body)

	return &Enemy{Body: ball.Body, ball: ball, alive: true, Radius: ballRadius, Mass: ballMass}
}

func (self *Enemy) shouldRemove() bool {
	return !self.alive
}

func (self *Enemy) Update(player *Player) {

	if !self.alive {
		return
	}

	p := self.Body.Position()
	//Move Enemie TOwards player
	if p.Y < 300 { //Only move if on bottom part of screen
		if p.X < player.Body.Position().X {
			// balls[i].Body.AddVelocity(10, 0)
			self.Body.AddAngularVelocity(-1)
		} else {
			// balls[i].Body.AddVelocity(-10, 0)
			self.Body.AddAngularVelocity(1)
		}
	}

	if p.Y < -100 {
		self.alive = false
	}
}
