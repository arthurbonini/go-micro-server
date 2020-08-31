package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "Arthur",
		Price: 1.99,
		SKU: "abs-dsa-asd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}