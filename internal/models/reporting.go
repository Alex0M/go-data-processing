package models

type ViewerCount struct {
	ViewerCount int
}

type StreamCount struct {
	StreamCount int
}

type ViewerCountPerState struct {
	State       string
	ViewerCount int
}

type ViewerCountPerDevice struct {
	Device      string
	ViewerCount int
}
