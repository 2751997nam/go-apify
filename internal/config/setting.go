package config

func GetSetting() map[string]any {
	return map[string]any{
		"top": map[uint64]string{
			5: "type",
			7: "style",
			2: "color",
		},
		"bottom": map[uint64]string{
			1: "size",
		},
	}

}
