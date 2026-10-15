package stat

type (
	// Response usage statistics
	// see docs https://dadata.ru/api/stat/
	Stat struct {
		Date     string `json:"date"`
		Services *struct {
			Brand       int `json:"brand"`
			Clean       int `json:"clean"`
			Company     int `json:"company"`
			Docs        int `json:"docs"`
			Merging     int `json:"merging"`
			Suggestions int `json:"suggestions"`
		} `json:"services"`
		Remaining struct {
			Brand       int `json:"brand"`
			Clean       int `json:"clean"`
			Company     int `json:"company"`
			Docs        int `json:"docs"`
			Merging     int `json:"merging"`
			Suggestions int `json:"suggestions"`
		} `json:"remaining"`
	}
)
