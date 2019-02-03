package ic

func mapRemoveNils(m map[string]interface{}) {
	for k, v := range m {
		if v == nil {
			delete(m, k)
		}
	}
}

func mapCopy(m map[string]interface{}) map[string]interface{} {
	mapCopy := make(map[string]interface{})

	for k, v := range m {
		mapCopy[k] = v
	}

	return mapCopy
}

func mapSubmap(m map[string]interface{}, fields ...string) map[string]interface{} {
	mapCopy := make(map[string]interface{})

	for _, field := range fields {
		if v, ok := m[field]; ok {
			mapCopy[field] = v
		}
	}

	return mapCopy
}

func mapGetValue(m map[string]interface{}, key string) interface{} {
	if v, ok := m[key]; ok {
		return v
	}

	return nil
}

func sliceFirstNonNil(values ...interface{}) interface{} {
	for _, v := range values {
		if v != nil {
			return v
		}
	}

	return nil
}
