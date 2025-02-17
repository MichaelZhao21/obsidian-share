package src

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi stinkyhead u shouldnt be here :3"))
}
