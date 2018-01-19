package model

//LogFile : the configuration structure which defines a log forwarder
type LogFile struct {
	FileName      string `json:"file"`
	SplitOn       string `json:"splitOn"`
	OutputIndices []int  `json:"outputIndices"`
}
