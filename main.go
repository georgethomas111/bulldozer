package bulldozer

type BullInput interface {
	DoWork(interface{})
	GetData()
}

type BullDozer struct {
	NumOfParallel int
	Task          BullInput
}

func main() {

}
