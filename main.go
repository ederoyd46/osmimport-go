package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/ederoyd46/osm/fileformat"
	"github.com/ederoyd46/osm/osmformat"

	"github.com/golang/protobuf/proto"
	"github.com/jeffail/tunny"
)

var (
	wg       sync.WaitGroup
	nodePool *tunny.WorkPool
)

func createWorkerPool() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	nodePool, _ = tunny.CreatePool(numCPUs, func(nodes interface{}) interface{} {
		input, _ := nodes.([]Node)
		SaveNodes(input)
		return 1
	}).Open()

}

func main() {
	if len(os.Args) != 4 {
		help()
		return //Same as os.Exit(0)
	}

	dbconnection := os.Args[1]
	dbname := os.Args[2]
	filename := os.Args[3]

	createWorkerPool()
	defer nodePool.Close()

	InitDB(dbconnection, dbname)
	startImport(filename)

	fmt.Println("Waiting for all go routines to finish")
	wg.Wait()
	fmt.Println("Import Done")
}

func help() {
	fmt.Println("usage: dbconnection dbname filename")
	fmt.Println("example: osmimport-go '127.0.0.1:28015' 'geo' './download/england-latest.osm.pbf'")
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

	blobData := make([]byte, header.GetDatasize())
	file.Read(blobData)

	var blob fileformat.Blob
	proto.Unmarshal(blobData, &blob)

	zr, err := zlib.NewReader(bytes.NewBuffer(blob.GetZlibData()))
	LogError(err)

	var blobUncompressed = make([]byte, blob.GetRawSize())
	io.ReadFull(zr, blobUncompressed)
	LogError(err)
	zr.Close()

	if "OSMHeader" == header.GetType() {
		osmHeader(blobUncompressed)
	}

	if "OSMData" == header.GetType() {
		osmData(blobUncompressed)
	}

	getBlock(4, file)

}

func osmHeader(blobUncompressed []byte) {
	var headerBlock osmformat.HeaderBlock
	proto.Unmarshal(blobUncompressed, &headerBlock)

	fmt.Println(" Min Lat:", float64(headerBlock.GetBbox().GetBottom())/nano)
	fmt.Println(" Min Lon:", float64(headerBlock.GetBbox().GetLeft())/nano)
	fmt.Println(" Max Lat:", float64(headerBlock.GetBbox().GetTop())/nano)
	fmt.Println(" Max Lon:", float64(headerBlock.GetBbox().GetRight())/nano)
}

func osmData(blobUncompressed []byte) {
	var primitiveBlock osmformat.PrimitiveBlock
	proto.Unmarshal(blobUncompressed, &primitiveBlock)

	// fmt.Println(" DateGranularity:", primitiveBlock.GetDateGranularity())
	// fmt.Println(" LatOffset:", primitiveBlock.GetLatOffset())
	// fmt.Println(" LonOffset:", primitiveBlock.GetLonOffset())
	// fmt.Println(" Granularity:", primitiveBlock.GetGranularity())
	// fmt.Println(" PrimitiveGroups:", len(primitiveBlock.GetPrimitivegroup()))

	var stringTable = ConvertStringTable(primitiveBlock.GetStringtable().GetS())
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
	handleNodes(group.GetDense(), stringTable, granularity, dateGranularity)
	handleWays(group.GetWays(), stringTable, granularity, dateGranularity)
	handleRelations(group.GetRelations(), stringTable, granularity, dateGranularity)
}

func handleRelations(pbRelations []*osmformat.Relation, stringTable []string, granularity float64, dateGranularity int64) {
	var relations []Relation
	for _, pbRelation := range pbRelations {
		relation := Relation{
			ID:        pbRelation.GetId(),
			Version:   pbRelation.GetInfo().GetVersion(),
			Timestamp: CalculateTime(int64(pbRelation.GetInfo().GetTimestamp()), dateGranularity),
			Changeset: pbRelation.GetInfo().GetChangeset(),
			UID:       pbRelation.GetInfo().GetUid(),
			User:      stringTable[pbRelation.GetInfo().GetUserSid()],
			Tags:      BuildTags(pbRelation.GetKeys(), pbRelation.GetVals(), stringTable),
			MemIds:    DeltaDecodeInt64(0, pbRelation.GetMemids()),
			Roles:     BuildStringList(pbRelation.GetRolesSid(), stringTable),
			Types:     ParseMemberTypes(pbRelation.GetTypes()),
		}
		relations = append(relations, relation)
	}

	if len(relations) > 0 {
		go func(data []Relation) {
			wg.Add(1)
			defer wg.Done()

			SaveRelations(data)
		}(relations)

		fmt.Println("Relations: ", len(relations))
	}
}

func handleWays(pbWays []*osmformat.Way, stringTable []string, granularity float64, dateGranularity int64) {
	var ways []Way
	for _, pbWay := range pbWays {
		way := Way{
			ID:        pbWay.GetId(),
			Version:   pbWay.GetInfo().GetVersion(),
			Timestamp: CalculateTime(int64(pbWay.GetInfo().GetTimestamp()), dateGranularity),
			Changeset: pbWay.GetInfo().GetChangeset(),
			UID:       pbWay.GetInfo().GetUid(),
			User:      stringTable[pbWay.GetInfo().GetUserSid()],
			Tags:      BuildTags(pbWay.GetKeys(), pbWay.GetVals(), stringTable),
			Refs:      DeltaDecodeInt64(0, pbWay.GetRefs()),
		}
		ways = append(ways, way)
	}

	if len(ways) > 0 {
		go func(data []Way) {
			wg.Add(1)
			defer wg.Done()

			SaveWays(data)
		}(ways)

		fmt.Println("Ways: ", len(ways))
	}
}

func handleNodes(denseNodes *osmformat.DenseNodes, stringTable []string, granularity float64, dateGranularity int64) {
	size := len(denseNodes.GetId())
	if size == 0 {
		return
	}
	ids := DeltaDecodeInt64(0, denseNodes.GetId())
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
			ID:        ids[i],
			Latitude:  CalculateDegrees(latitudes[i], granularity),
			Longitude: CalculateDegrees(longitudes[i], granularity),
			Timestamp: CalculateTime(timestamps[i], dateGranularity),
			Changeset: changesets[i],
			UID:       uids[i],
			User:      stringTable[sids[i]],
			Tags:      keyvals[i],
		}
		nodes = append(nodes, node)
	}

	if len(nodes) > 0 {
		go func(data []Node) {
			wg.Add(1)
			defer wg.Done()
			nodePool.SendWork(data)
		}(nodes)
		fmt.Println("Nodes: ", len(nodes))
	}
}
