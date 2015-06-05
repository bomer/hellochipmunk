package engine

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math"
	"math/rand"
)

type Player struct {
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

	return &Player{Body: ball.Body, ball: ball, Radius: ballRadius, Mass: ballMass}
}

func (self *Player) Update() {

}
