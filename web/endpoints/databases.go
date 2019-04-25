package endpoints

import (
	// "encoding/json"
	// "go-contacts/models"
	// u "go-contacts/utils"
	// "strconv"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DbInfo struct {
	Name           string
	Parameters     Parameters
	Samples        []string
	NumSamples     uint64
	BlockSize      uint64
	KeepEmptyBlock bool
	NumRegisters   uint64
	NumSNPS        uint64
	NumBlocks      uint64
	CounterBits    int
}

type ChromInfo struct {
	Name        string
	Pos         int
	MinPosition uint64
	MaxPosition uint64
	NumBlocks   uint64
	NumSNPS     uint64
}

func Databases(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Databases", DATABASE_DIR)

	files := listDatabases()

	resp := Message(true, "success")
	resp["data"] = files
	Respond(w, resp)
}

func Database(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	database := params["database"]

	listDatabases()

	ib, ok := databases.Get(database)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such database: " + database
		Respond(w, resp)
		return
	}

	infos := make([]ChromInfo, 0, 0)

	for _, chromNamePosPair := range ib.ChromosomesNames {
		chromName := chromNamePosPair.Name
		chromPos := chromNamePosPair.Pos
		chromosome := ib.Chromosomes[chromName]

		infos = append(infos, ChromInfo{
			Name:        chromName,
			Pos:         chromPos,
			MinPosition: chromosome.MinPosition,
			MaxPosition: chromosome.MaxPosition,
			NumBlocks:   chromosome.NumBlocks,
			NumSNPS:     chromosome.NumSNPS,
		})
	}

	resp := Message(true, "success")
	resp["data"] = infos

	Respond(w, resp)

	// params := mux.Vars(r)
	// category := vars["category"]
	// id, err := strconv.Atoi(params["id"])
	// if err != nil {
	// 	//The passed path parameter is not an integer
	// 	Respond(w, Message(false, "There was an error in your request"))
	// 	return
	// }

	// data := models.GetContacts(uint(id))

	// params := mux.Vars(r)
	// id, err := strconv.Atoi(params["id"])
	// if err != nil {
	// 	//The passed path parameter is not an integer
	// 	Respond(w, Message(false, "There was an error in your request"))
	// 	return
	// }

	// data := models.GetContacts(uint(id))
	// resp := Message(true, "success")
	// resp["data"] = data
	// Respond(w, resp)
}

func listDatabases() (files []DbInfo) {
	err := filepath.Walk(DATABASE_DIR, func(path string, info os.FileInfo, err error) error {
		found, _, _, prefix := GuessFormat(path)

		if found {
			fi, err := os.Stat(path)

			if err != nil {
				fmt.Println(err)
				return nil
			}

			if fi.Mode().IsRegular() {
				fn := strings.TrimPrefix(prefix, DATABASE_DIR)

				ib, ok := databases.Register(fn, path)

				if ok {
					files = append(files, DbInfo{
						Name:           fn,
						Parameters:     ib.Parameters,
						Samples:        ib.Samples,
						NumSamples:     ib.NumSamples,
						BlockSize:      ib.BlockSize,
						KeepEmptyBlock: ib.KeepEmptyBlock,
						NumRegisters:   ib.NumRegisters,
						NumSNPS:        ib.NumSNPS,
						NumBlocks:      ib.NumBlocks,
						CounterBits:    ib.CounterBits,
					})
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}
