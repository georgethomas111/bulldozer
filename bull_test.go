package bulldozer

import (
	"testing"
)

type NumPrinter struct {
}

func (n *NumPrinter) DoWork(i int) {
	fmt.Println(i)
}

func TestBull(t *testing.T) {

}
