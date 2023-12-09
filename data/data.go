package data

type Word struct {
	OriginFile string
	Word       string
}

func (w Word) String() string {
	return w.Word
}
