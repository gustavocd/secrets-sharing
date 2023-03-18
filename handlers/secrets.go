package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gustavocd/secrets-sharing/filestore"
	"github.com/gustavocd/secrets-sharing/types"
)

func secretHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSecret(w, r)
	case http.MethodPost:
		saveSecret(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.String(), "/")
	if len(paths) < 2 {
		http.Error(w, "No secret ID specified", http.StatusBadRequest)
		return
	}

	var gs types.GetSecretResponse
	id := paths[1]
	v, err := filestore.FileStoreConfig.Fs.Read(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gs.Data = v
	jd, err := json.Marshal(gs)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(gs.Data) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write(jd)
}

func saveSecret(w http.ResponseWriter, r *http.Request) {
	var sp types.CreateSecretPayload
	err := json.NewDecoder(r.Body).Decode(&sp)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	sd := types.SecretData{
		Id:     fmt.Sprintf("%x", md5.Sum([]byte(sp.PlainText))),
		Secret: sp.PlainText,
	}

	err = filestore.FileStoreConfig.Fs.Write(sd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := types.CreateSecretResponse{Id: sd.Id}

	jd, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jd)
}
