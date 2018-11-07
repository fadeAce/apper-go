package apper_go

type Apper struct {
}

// singleton mode
func GetApper() (*Apper, error) {
	return &Apper{}, nil
}

func (*Apper) Connect(url string) error {

	return nil
}

func (*Apper) Start(path string) (string, error) {
	return "", nil
}

func (*Apper) Stop(transactionID string) error {
	return nil
}

func (*Apper) Terminate(pass string) {

}

func (*Apper) Ls() {

}
