package models

type LoaderType struct {
	FileSize  int64
	ChankSize int64
}

type Prog struct {
	Done  int64
	Total int64
	Err   error
}
