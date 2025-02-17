package src

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileLink struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Slug     string             `json:"slug" bson:"slug"`
	Filename string             `json:"filename" bson:"filename"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi stinkyhead u shouldnt be here :3"))
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	db := GetDB(r)

	// Get the file ID from the URL path parameter
	id := chi.URLParam(r, "id")

	// Get the file from the database
	file := db.Collection("docs").FindOne(r.Context(), bson.M{"slug": id})
	if file.Err() != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Decode the file
	var f FileLink
	file.Decode(&f)

	// Serve the file
	http.ServeFile(w, r, "./files/html/"+f.Filename+".html")
}

func PullFiles(w http.ResponseWriter, r *http.Request) {
	// Get the repository URL from the environment
	repoURL := os.Getenv("REPO_URL")

	// Get the database from the context
	db := GetDB(r)

	// Pull the files
	err := CloneOrPullRepo(repoURL, "./files/raw")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all linked files
	cursor, err := db.Collection("docs").Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode the files
	var files []FileLink
	cursor.All(r.Context(), &files)

	// Convert all files to HTML
	err = ConvAllToHtml(files, "./files/html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Files pulled and converted"))
}
