package datareader

type DataReader interface {
	ReadRow() (*DataRow, error)
	Close()
}

type DataRow map[string]any
