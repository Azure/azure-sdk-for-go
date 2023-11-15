package nonexport

type Public struct {
}

func (p *Public) PublicMethod() {

}

func (p *Public) privateMethod() {

}

func (p *Public) removePrivateMethod() {

}

type private struct {
}
