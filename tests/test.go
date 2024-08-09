package main

import "fmt"

type A interface {
	FA() string
	GA() string
}

type AS struct {
}

func (e *AS) FA() string {
	return "AS_FA"
}

func (e *AS) GA() string {
	return "AS_GA"
}

type BS struct {
}

func (e *BS) FA() string {
	return "BS_FA"
}

func (e *BS) GA() string {
	return "BS_GA"
}

type C struct {
	A
}

type D struct {
	A
}

func (e *D) FA() string {
	return "D_FA"
}

func main() {
	a := AS{}
	fmt.Printf("AS: %s\n", a.FA())
	c := C{&BS{}}
	fmt.Printf("C/BS: %s\n", c.FA())
	d := D{&BS{}}
	fmt.Printf("D: BS_FA: %s BS_GA: %s\n", d.FA(), d.GA())

}
