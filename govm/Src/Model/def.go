package Model

type DateStr string

type Err string

func (err Err) Exists() bool {
	return len(err) > 0
}
