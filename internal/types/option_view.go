package types

type OptionView struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// Slug        string `json:"slug,omitempty"`
	// VariantSlug string `json:"variant_slug,omitempty"`
	// VariantType string `json:"variant_type,omitempty"`
	VariantId uint64 `json:"variant_id"`
	// Variant     string `json:"variant,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
	// Options     []OptionView `json:"options"`
}
