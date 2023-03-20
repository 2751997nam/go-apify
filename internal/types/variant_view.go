package types

type VariantView struct {
	ID     uint64       `json:"id"`
	Name   string       `json:"name,omitempty"`
	Slug   string       `json:"slug,omitempty"`
	Type   string       `json:"type"`
	Values []OptionView `json:"values"`
}
