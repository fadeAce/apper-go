package apper_go

type Apper struct {
}

func NewApperInstace() (*Apper, error) {
	return &Apper{}, nil
}

func (*Apper) Connect() error {
	return nil
}

func (*Apper) Start(path string) (string, error) {
	return "", nil
}

func (*Apper) Terminate(pass string) {

}

func (*Apper) Ls() {

}
