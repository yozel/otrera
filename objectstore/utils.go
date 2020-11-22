package objectstore

import "encoding/json"

func DeepCopy(v *interface{}) (*interface{}, error) {
	var r interface{}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	return &r, nil
}
