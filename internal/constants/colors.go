package constants

type Colors struct {
	Foreground string
	Cyan       string
	Green      string
	Orange     string
	Pink       string
	Purple     string
	Red        string
	Yellow     string
	Background string // Adding background for completeness
	Comment    string // Common in Dracula theme
	Selection  string // Common in Dracula theme
}

var DraculaColors = Colors{
	Foreground: "#f8f8f2",
	Cyan:       "#8be9fd",
	Green:      "#50fa7b",
	Orange:     "#ffb86c",
	Pink:       "#ff79c6",
	Purple:     "#bd93f9",
	Red:        "#ff5555",
	Yellow:     "#f1fa8c",
	Background: "#282a36",
	Comment:    "#6272a4",
	Selection:  "#44475a",
}
