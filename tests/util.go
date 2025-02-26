package tests

func WrapItemsResponse(input []byte) []byte {
	result := append([]byte(`{"items":[`), input...)
	result = append(result, "]}"...)
	
	return result
}