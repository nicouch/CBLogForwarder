package main

type supervisor struct {
	Files []logFile `json:"files"`
}

type logFile struct {
	FileName      string `json:"file"`
	SplitOn       string `json:"splitOn"`
	OutputIndices []int  `json:"outputIndices"`
}
