package cop

type Bar struct {}

type Usefuler interface {
	BeUseful()
}

func NewBar(foo Usefuler) *Bar {

}
