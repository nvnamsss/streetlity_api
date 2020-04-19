package pipeline

import (
	"fmt"
	"testing"
)

func TestNextStage(t *testing.T) {

	stage := NewStage(nil)
	fn := func() (struct {
		Field string
		Meo   int
	}, error) {

		return struct {
			Field string
			Meo   int
		}{"meomeocute", 1}, nil
	}

	var p *Pipeline = NewPipeline()
	stage.Next(fn)
	p.First = stage
	p.Run()

	fmt.Println(p.GetString("Field"))
	fmt.Println(p.GetFloat("Meo"))
}