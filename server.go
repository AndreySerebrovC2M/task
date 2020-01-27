package main

import (
	"2/example"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DBConn   = "user=dev password=dev host=example_db port=5432 dbname=example_db sslmode=disable connect_timeout=5"
	DBDriver = "postgres"
)

func main() {
	store, err := example.NewStore(DBDriver, DBConn)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect to db")
	}
	http.HandleFunc("/your_url", func(w http.ResponseWriter, httpReq *http.Request) {
		if httpReq.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(httpReq.Body)
		defer func() {
			_ = httpReq.Body.Close()
		}()
		tr := example.Transaction{}
		if err := decoder.Decode(&tr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tr.Source = httpReq.Header.Get("Source-Type")
		err := store.InsertTransaction(tr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("ok"))
		w.WriteHeader(http.StatusOK)
		return

	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("failed to start")
	}
}
