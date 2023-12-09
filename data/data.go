package data

import "fmt"

type Word struct {
	OriginFile string
	Word       string
}

func (w Word) String() string {
	return fmt.Sprintf("%s\t%s", w.OriginFile, w.Word)
}
