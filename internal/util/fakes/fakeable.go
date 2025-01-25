package fakes

import "encoding/json"

func Fakeable[T any](
	concrete func() (T, error),
	fake func(data json.RawMessage) (T, error),
	faked bool,
	data json.RawMessage,
) (T, error) {
	if faked {
		return fake(data)
	}

	return concrete()
}
