package endpoints

import (
	// "encoding/json"
	// "go-contacts/models"
	// u "go-contacts/utils"
	// "github.com/gorilla/mux"
	// "strconv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Databases(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Databases", DATABASE_DIR)
	var files []string
	err := filepath.Walk(DATABASE_DIR, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yaml") {
			fi, err := os.Stat(path)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			if fi.Mode().IsRegular() {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	// params := mux.Vars(r)
	// category := vars["category"]
	// id, err := strconv.Atoi(params["id"])
	// if err != nil {
	// 	//The passed path parameter is not an integer
	// 	Respond(w, Message(false, "There was an error in your request"))
	// 	return
	// }

	// data := models.GetContacts(uint(id))
	resp := Message(true, "success")
	resp["data"] = files
	Respond(w, resp)
}

func Database(w http.ResponseWriter, r *http.Request) {
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
