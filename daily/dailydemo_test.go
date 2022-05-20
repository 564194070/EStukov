package daily

import "testing"

func TestShowStr(t *testing.T)  {
	//SendFile()
}

func TestJobSort(t *testing.T) {
	A :=Job2{nil,"A",nil}
	B :=Job2{nil,"B",nil}
	C :=Job2{nil,"C",nil}
	D :=Job2{nil,"D",nil}
	E :=Job2{nil,"E",nil}
	F := Job2{nil,"F",[]*Job2{&A,&B,&C}}
	G := Job2{nil,"G",[]*Job2{&C,&D,&E}}
	H := Job2{nil,"H",[]*Job2{&F,&C,&G}}
	I := Job2{nil,"I",[]*Job2{&H}}
	//jobList1 := []Job{A,B,C,D,E,F,G,H,I}
	jobList2 := []Job2{I,F,B,C,D,E,G,H,A}

	//sortJob(jobList2)
	sortJob2(jobList2)
}