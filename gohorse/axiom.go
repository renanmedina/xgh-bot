package gohorse

import (
	"fmt"
)

type Axiom struct {
	Number      uint   `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (axiom Axiom) Url() string {
	return fmt.Sprintf("https://gohorse.com.br/extreme-go-horse-xgh/#%d", axiom.Number)
}

func (axiom Axiom) AudioUrl() string {
	return fmt.Sprintf("/assets/static/audios/axioms/axiom_%d.mp3", axiom.Number)
}

func (axiom Axiom) ToQuote() string {
	return fmt.Sprintf("%d - %s\r\n%s", axiom.Number, axiom.Title, axiom.Description)
}

func (axiom Axiom) ImageUrl() string {
	return fmt.Sprintf("frontend/static/images/axioms/%d.jpg", axiom.Number)
}
