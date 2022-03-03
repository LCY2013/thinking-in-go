package packet

import "bytes"

type Submit struct {
	ID      string
	Payload []byte
}

func (sub *Submit) Decode(pkgBody []byte) error {
	sub.ID = string(pkgBody[:8])
	sub.Payload = pkgBody[8:]
	return nil
}

func (sub *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(sub.ID[:8]), sub.Payload}, nil), nil
}
