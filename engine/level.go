package engine

import (
	"encoding/json"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"io/ioutil"
)

type LevelPoint struct {
	Ax float32
	Ay float32
	Bx float32
	By float32
}

func (self LevelPoint) getChipmunkSegment() *chipmunk.Shape {
	return chipmunk.NewSegment(vect.Vect{vect.Float(self.Ax), vect.Float(self.Ay)}, vect.Vect{vect.Float(self.Bx), vect.Float(self.By)}, 0)
}

type Level struct {
	Points []LevelPoint
}

func NewLevelFromFile(filename string) *Level {
	var level Level

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &level)
	if err != nil {
		panic(err)
	}

	return &level
}

func (self Level) GetChipmunkSegments() []*chipmunk.Shape {
	segments := make([]*chipmunk.Shape, len(self.Points))

	for i, lp := range self.Points {
		segments[i] = lp.getChipmunkSegment()
	}

	return segments
}
