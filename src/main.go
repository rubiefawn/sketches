package main

import (
	. "fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	fpath "path/filepath"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	IsDraft uint = iota
	IsWip
	IsRelease
)

type Song struct {
	Name   string
	Emoji  string
	Type   uint
	Drafts []Draft
}

type Draft struct {
	Path     template.URL
	Modified time.Time
}

func main() {
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".mjs", "text/javascript")

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /", serve_template)
	http.ListenAndServe(":8080", mux)
}

func serve_template(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("template/*.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	songs, err := scan_all_songs("static/audio")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if err = tmpl.ExecuteTemplate(w, "main", songs); err != nil {
		Println(err)
	}
}

func song_from_path(path string) (song Song, err error) {
	path = fpath.Clean(path)
	var files []os.DirEntry
	if files, err = os.ReadDir(path); err != nil {
		return
	}
	for _, file := range files {
		ext := fpath.Ext(file.Name())
		if file.IsDir() || !(ext == ".mp3" || ext == ".wav" || ext == ".flac" || file.Name() == "info.toml") {
			continue
		}
		dpath := fpath.ToSlash(fpath.Join(path, file.Name()))
		if file.Name() == "info.toml" {
			var data []byte
			if data, err = os.ReadFile(dpath); err != nil {
				return
			}
			var info struct {
				Emoji string
				Type  string
			}
			if err = toml.Unmarshal(data, &info); err != nil {
				return
			}
			song.Emoji = info.Emoji
			switch info.Type {
			case "draft":
				song.Type = IsDraft
			case "wip":
				song.Type = IsWip
			case "rel", "release":
				song.Type = IsRelease
			default:
				song.Type = IsDraft
			}
		}
		var i os.FileInfo
		if i, err = file.Info(); err != nil {
			return
		}
		song.Drafts = append(song.Drafts, Draft{template.URL(dpath), i.ModTime()})
	}
	sort.Slice(song.Drafts, func(a int, b int) bool { return song.Drafts[a].Modified.After(song.Drafts[b].Modified) })
	return
}

func scan_all_songs(in_path string) (songs []Song, err error) {
	in_path = fpath.Clean(in_path)
	var dirs []os.DirEntry
	if dirs, err = os.ReadDir(in_path); err != nil {
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		var song Song
		if song, err = song_from_path(fpath.Join(in_path, dir.Name())); err != nil {
			return
		}
		songs = append(songs, song)
	}
	sort.Slice(songs, func(a int, b int) bool { return songs[a].Drafts[0].Modified.After(songs[b].Drafts[0].Modified) })
	return
}
