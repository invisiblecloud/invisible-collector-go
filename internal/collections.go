package internal

func MapRemoveNils(m map[string]interface{}) {
	for k, v := range m {
		if v == nil {
			delete(m, k)
		}
	}
}

func MapCopy(m map[string]interface{}) map[string]interface{} {
	mapCopy := make(map[string]interface{})

	for k, v := range m {
		mapCopy[k] = v
	}

	return mapCopy
}

func MapSubmap(m map[string]interface{}, fields ...string) map[string]interface{} {
	mapCopy := make(map[string]interface{})

	for _, field := range fields {
		if v, ok := m[field]; ok {
			mapCopy[field] = v
		}
	}

	return mapCopy
}

func MapGetValue(m map[string]interface{}, key string) interface{} {
	if v, ok := m[key]; ok {
		return v
	}

	return nil
}

func SliceFirstNonNil(values ...interface{}) interface{} {
	for _, v := range values {
		if v != nil {
			return v
		}
	}

	return nil
}
