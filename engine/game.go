package engine

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

type Game struct {
	initialized bool

	BallRadius  int32
	BallMass    int32
	DegreeConst vect.Float

	Space *chipmunk.Space

	Level   *Level
	Enemies []Enemy
	Player  *Player
}

func NewGame() *Game {

	return &Game{initialized: false, BallRadius: 25, BallMass: 1, DegreeConst: chipmunk.DegreeConst}
}

func (self *Game) Init() {
	self.Space = chipmunk.NewSpace()
	self.Space.Gravity = vect.Vect{0, -900}

	self.Player = NewPlayer(self.BallRadius, self.BallMass, self)

	self.initialized = true

}

func (self *Game) Update(dt float32) {
	self.Space.Step(vect.Float(dt))

	self.Player.Update()
	for i, enemy := range self.Enemies {
		enemy.Update(self.Player)
		if enemy.shouldRemove() {
			self.Space.RemoveBody(enemy.Body)

			// Not sure if the best way to do this ??
			self.Enemies = append(self.Enemies[:i], self.Enemies[i+1:]...)
		}
	}
}

func (self *Game) ReadLevelFromFile(filename string) {
	self.Level = NewLevelFromFile(filename)

	staticBody := chipmunk.NewBodyStatic()
	for _, segment := range self.Level.GetChipmunkSegments() {
		segment.SetElasticity(0.6)
		staticBody.AddShape(segment)
	}

	self.Space.AddBody(staticBody)
}
