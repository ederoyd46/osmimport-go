package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"osm/fileformat"
	"osm/osmformat"

	"github.com/golang/protobuf/proto"
)

const nano float64 = 1000000000

func main() {
	startImport("./download/hertfordshire-latest.osm.pbf")
}

func startImport(fileName string) {
	file, err := os.Open(fileName)
	logError(err)

	getBlock(4, file)
	file.Close()
}

func getBlock(size int64, file *os.File) {

	headerSizeData := make([]byte, size)
	file.Read(headerSizeData)

	var headerSize uint32
	err := binary.Read(bytes.NewBuffer(headerSizeData), binary.BigEndian, &headerSize)
	logError(err)

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
	logError(err)

	fmt.Println("Compressed Size", len(blob.GetZlibData()))
	fmt.Println("Raw Size:", blob.GetRawSize())

	var blobUncompressed = make([]byte, blob.GetRawSize())
	io.ReadFull(zr, blobUncompressed)
	logError(err)
	zr.Close()

	// var readCount int32
	//
	// for readCount < blob.GetRawSize() {
	// 	blobSection := make([]byte, 32768)
	// 	count, err := zr.Read(blobSection)
	// 	logError(err)
	// 	readCount = readCount + int32(count)
	// 	// fmt.Println("readCount", readCount)
	// 	blobUncompressed = append(blobUncompressed, blobSection...)
	// }

	if "OSMHeader" == header.GetType() {
		osmHeader(blobUncompressed)
		getBlock(4, file)
	}

	if "OSMData" == header.GetType() {
		osmData(blobUncompressed)
	}

	// getBlock(4, file)

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
		handlePrimitiveGroupData(group, stringTable)
	}

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

func handlePrimitiveGroupData(group *osmformat.PrimitiveGroup, stringTable []string) {
	nodes := group.GetDense()
	size := len(nodes.GetId())
	uids := DeltaDecodeInt32(0, nodes.GetDenseinfo().GetUid())
	sids := DeltaDecodeInt32(0, nodes.GetDenseinfo().GetUserSid())
	timestamps := DeltaDecodeInt64(0, nodes.GetDenseinfo().GetTimestamp())
	changesets := DeltaDecodeInt64(0, nodes.GetDenseinfo().GetChangeset())
	latitudes := DeltaDecodeInt64ToFloat(0, nodes.GetLat())
	longitudes := DeltaDecodeInt64ToFloat(0, nodes.GetLon())
	keyvals := buildKeyVals(nodes.GetKeysVals(), stringTable)

	for i := 0; i < size; i++ {
		node := Node{id: nodes.GetId()[i], latitude: latitudes[i], longitude: longitudes[i],
			timestamp: timestamps[i], changeset: changesets[i], uid: uids[i],
			sid: stringTable[sids[i]], tags: keyvals[i]}
		fmt.Println(node)
	}

}

func logError(errReference error) {
	if errReference != nil {
		log.Fatal(errReference)
	}
}
