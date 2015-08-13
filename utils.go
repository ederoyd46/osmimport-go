package main

import "time"

const nano float64 = 1000000000

//CalculateDegrees calcluates the real coordinate from the delta decoded one
func CalculateDegrees(coordinate float64, granularity float64) float64 {
	return (coordinate * granularity) / nano
}

//CalculateTime calculates the time
func CalculateTime(timestamp int64, granularity int64) time.Time {
	return time.Unix(0, (timestamp * granularity))
}

//ConvertStringTable converts the string table so it returns a string slice instead of bytes.
func ConvertStringTable(stringTable [][]byte) []string {
	var convertedTable []string
	for _, entry := range stringTable {
		convertedTable = append(convertedTable, string(entry))
	}
	return convertedTable
}

//DeltaDecodeInt64 Takes a seed value (normally 0) and a list. Delta decodes and returns the list (i.e. the next value is determined by the previous value plus the difference).
func DeltaDecodeInt64(seed int64, data []int64) []int64 {
	var decodedVals []int64
	for _, entry := range data {
		decodedEntry := int64(seed + entry)
		seed = int64(decodedEntry)
		decodedVals = append(decodedVals, decodedEntry)
	}
	return decodedVals
}

//DeltaDecodeInt32 Takes a seed value (normally 0) and a list. Delta decodes and returns the list (i.e. the next value is determined by the previous value plus the difference).
func DeltaDecodeInt32(seed int32, data []int32) []int32 {
	var decodedVals []int32
	for _, entry := range data {
		decodedEntry := int32(seed + entry)
		seed = int32(decodedEntry)
		decodedVals = append(decodedVals, decodedEntry)
	}
	return decodedVals
}

//DeltaDecodeInt64ToFloat Takes a seed value (normally 0) and a list. Delta decodes and returns the list (i.e. the next value is determined by the previous value plus the difference).
func DeltaDecodeInt64ToFloat(seed int64, data []int64) []float64 {
	var decodedVals []float64
	for _, entry := range data {
		decodedEntry := float64(seed + entry)
		seed = int64(decodedEntry)
		decodedVals = append(decodedVals, decodedEntry)
	}
	return decodedVals
}

//LogError Generic function loging out errors
func LogError(errReference error) {
	if errReference != nil {
		// log.Fatal(errReference)
		panic(errReference)
	}
}
