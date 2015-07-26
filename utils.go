package main

//ConvertStringTable converts the string table so it returns a string slice instead of bytes.
func ConvertStringTable(stringTable [][]byte) []string {
	var convertedTable []string
	for _, entry := range stringTable {
		convertedTable = append(convertedTable, string(entry))
	}
	return convertedTable
}

//DeltaDecodeInt64 for int64
func DeltaDecodeInt64(seed int64, data []int64) []int64 {
	var decodedVals []int64
	for _, entry := range data {
		decodedEntry := int64(seed + entry)
		seed = int64(decodedEntry)
		decodedVals = append(decodedVals, decodedEntry)
	}
	return decodedVals
}

//DeltaDecodeInt32 for int32
func DeltaDecodeInt32(seed int32, data []int32) []int32 {
	var decodedVals []int32
	for _, entry := range data {
		decodedEntry := int32(seed + entry)
		seed = int32(decodedEntry)
		decodedVals = append(decodedVals, decodedEntry)
	}
	return decodedVals
}

//DeltaDecodeInt64ToFloat for int64 to float64
func DeltaDecodeInt64ToFloat(seed int64, data []int64) []float64 {
	var decodedVals []float64
	for _, entry := range data {
		decodedEntry := float64(seed + entry)
		seed = int64(decodedEntry)
		decodedVals = append(decodedVals, decodedEntry)
	}
	return decodedVals
}
