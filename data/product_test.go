package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "york",
		Price: 1.22,
		SKU:   "abd-fdg-gfsd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
