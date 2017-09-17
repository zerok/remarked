package commandchain

type Command struct {
	Type       string `json:"type"`
	SlideIndex int    `json:"slideIndex"`
	Token      string `json:"token"`
}
