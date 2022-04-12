package pigeoncache

type PeerPicker interface {
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}