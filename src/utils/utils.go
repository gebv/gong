package utils

//

type M map[string]interface{}

func NewM() M {
	return make(M)
}

func (m M) Set(key string, v interface{}) M {
	m[key] = v

	return m
}

func (m M) Get(key string) interface{} {
	return m[key]
}
