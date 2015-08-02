package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/ederoyd46/osm/fileformat"
	"github.com/ederoyd46/osm/osmformat"

	"github.com/golang/protobuf/proto"
)

func main() {
	startImport("./download/hertfordshire-latest.osm.pbf")
}

func startImport(fileName string) {
	file, err := os.Open(fileName)
	LogError(err)

	getBlock(4, file)
	file.Close()
}

func getBlock(size int64, file *os.File) {

	headerSizeData := make([]byte, size)
	_, err := file.Read(headerSizeData)
	if err != nil {
		//Probably end of file
		return
	}

	var headerSize uint32
	err = binary.Read(bytes.NewBuffer(headerSizeData), binary.BigEndian, &headerSize)
	LogError(err)

	headerData := make([]byte, headerSize)
	file.Read(headerData)

	var header fileformat.BlockHeader
	proto.Unmarshal(headerData, &header)

	fmt.Println(header.GetType())

	blobData := make([]byte, header.GetDatasize())
	file.Read(blobData)

	var blob fileformat.Blob
	proto.Unmarshal(blobData, &blob)

	zr, err := zlib.NewReader(bytes.NewBuffer(blob.GetZlibData()))
	LogError(err)

	fmt.Println("Compressed Size", len(blob.GetZlibData()))
	fmt.Println("Raw Size:", blob.GetRawSize())

	var blobUncompressed = make([]byte, blob.GetRawSize())
	io.ReadFull(zr, blobUncompressed)
	LogError(err)
	zr.Close()

	if "OSMHeader" == header.GetType() {
		osmHeader(blobUncompressed)
		// getBlock(4, file)
	}

	if "OSMData" == header.GetType() {
		osmData(blobUncompressed)
	}

	getBlock(4, file)

}

func osmHeader(blobUncompressed []byte) {
	var headerBlock osmformat.HeaderBlock
	proto.Unmarshal(blobUncompressed, &headerBlock)

	minlat := float64(headerBlock.GetBbox().GetBottom()) / nano
	minlon := float64(headerBlock.GetBbox().GetLeft()) / nano
	maxlat := float64(headerBlock.GetBbox().GetTop()) / nano
	maxlon := float64(headerBlock.GetBbox().GetRight()) / nano

	fmt.Println(" Min Lat:", minlat)
	fmt.Println(" Min Lon:", minlon)
	fmt.Println(" Max Lat:", maxlat)
	fmt.Println(" Max Lon:", maxlon)

}

func osmData(blobUncompressed []byte) {
	var primitiveBlock osmformat.PrimitiveBlock
	proto.Unmarshal(blobUncompressed, &primitiveBlock)

	fmt.Println(" DateGranularity:", primitiveBlock.GetDateGranularity())
	fmt.Println(" LatOffset:", primitiveBlock.GetLatOffset())
	fmt.Println(" LonOffset:", primitiveBlock.GetLonOffset())
	fmt.Println(" Granularity:", primitiveBlock.GetGranularity())
	fmt.Println(" PrimitiveGroups:", len(primitiveBlock.GetPrimitivegroup()))

	var stringTable = ConvertStringTable(primitiveBlock.GetStringtable().GetS())
	fmt.Println(" StringTable:", len(stringTable))

	var primitiveGroup = primitiveBlock.GetPrimitivegroup()

	for _, group := range primitiveGroup {
		handlePrimitiveGroupData(group,
			stringTable,
			float64(primitiveBlock.GetGranularity()),
			int64(primitiveBlock.GetDateGranularity()),
		)
	}

}

func buildKeyVals(mixedKeyVals []int32, stringTable []string) []map[string]string {
	var keyvals []map[string]string
	keyvalEntry := make(map[string]string)

	count := len(mixedKeyVals)
	for i := 0; i < count; {
		if mixedKeyVals[i] == 0 {
			keyvals = append(keyvals, keyvalEntry)
			keyvalEntry = map[string]string{}
			i = i + 1
		} else {
			key := stringTable[mixedKeyVals[i]]
			val := stringTable[mixedKeyVals[i+1]]
			keyvalEntry[key] = val
			i = i + 2
		}
	}
	return keyvals
}

func handlePrimitiveGroupData(group *osmformat.PrimitiveGroup, stringTable []string, granularity float64, dateGranularity int64) {
	denseNodes := group.GetDense()
	size := len(denseNodes.GetId())
	uids := DeltaDecodeInt32(0, denseNodes.GetDenseinfo().GetUid())
	sids := DeltaDecodeInt32(0, denseNodes.GetDenseinfo().GetUserSid())
	timestamps := DeltaDecodeInt64(0, denseNodes.GetDenseinfo().GetTimestamp())
	changesets := DeltaDecodeInt64(0, denseNodes.GetDenseinfo().GetChangeset())
	latitudes := DeltaDecodeInt64ToFloat(0, denseNodes.GetLat())
	longitudes := DeltaDecodeInt64ToFloat(0, denseNodes.GetLon())
	keyvals := buildKeyVals(denseNodes.GetKeysVals(), stringTable)

	var nodes []Node
	for i := 0; i < size; i++ {
		node := Node{
			ID:        denseNodes.GetId()[i],
			Latitude:  CalculateDegrees(latitudes[i], granularity),
			Longitude: CalculateDegrees(longitudes[i], granularity),
			Timestamp: CalculateTime(timestamps[i], dateGranularity),
			Changeset: changesets[i],
			UID:       uids[i],
			SID:       stringTable[sids[i]],
			Tags:      keyvals[i]}
		nodes = append(nodes, node)
	}
	SaveNodes(nodes)
}
