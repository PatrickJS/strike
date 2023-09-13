package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	. "github.com/JLarky/strike/pkg/h"
	"github.com/JLarky/strike/pkg/strike"
)

func main() {
	http.Handle("/favicon.ico", http.FileServer(http.Dir("public")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flush := func() {
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
				time.Sleep(1000 * time.Millisecond)
			}
		}

		page := Page(r.URL.Path)

		rsc := r.Header.Get("RSC")
		if rsc == "1" {
			jsonData, err := json.Marshal(page)
			if err != nil {
				fmt.Printf("Error serializing data: %v", err)
				return
			}

			w.Header().Set("Content-Type", "text/x-component; charset=utf-8")
			w.Write(jsonData)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<!doctype html>"))
		// flush()
		strike.RenderToString(w, page)
		flush()
		// w.Write([]byte("Hello, World!"))
		jsonData, err := json.Marshal(page)

		if err != nil {
			fmt.Printf("Error serializing data: %v", err)
			return
		}

		const tpl = `<script>self.__rsc=self.__rsc||[];__rsc.push({{.}})</script>`

		t, err := template.New("webpage").Parse(tpl)

		if err != nil {
			fmt.Printf("Error parsing template: %v", err)
			return
		}

		err = t.Execute(w, string(jsonData))

		if err != nil {
			fmt.Printf("Error parsing template: %v", err)
			return
		}
	})
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Page(url string) Component {
	// time.Sleep(1000 * time.Millisecond)
	fmt.Println("Page", url)
	nav := H("nav",
		H("a", Props{"href": "/"}, "Home"), " ",
		H("a", Props{"href": "/about"}, "About"),
	)

	return H("html", Props{"lang": "en"},
		H("head",
			H("meta", Props{"charset": "utf-8"}),
			H("meta", Props{"name": "description", "content": "React with Server Components demo"}),
			H("meta", Props{"name": "viewport", "content": "width=device-width, initial-scale=1"}),
			H("link", Props{"rel": "stylesheet", "href": "/static/style.css"}),
			H("title", "React Notes"),
		),
		H("body",
			App(),
			H("div", Props{"id": "root"}, nav, "Loading... "+url),
			H("script", Props{"src": "/static/strike/bootstrap.js", "type": "module"}),
		),
	)
}

func App() Component {
	return H("div", Props{"class": "main"},
		H("section", Props{"class": "col sidebar"},
			H("section", Props{"class": "sidebar-header"},
				H("img", Props{"class": "logo", "src": "/static/logo.svg", "width": "22px", "height": "20px", "alt": "", "role": "presentation"}),
				H("strong", "React Notes"),
			),
			H("section", Props{"class": "sidebar-menu", "role": "menubar"},
				SearchField(),
				EditButton(nil, "New"),
			),
			H("nav"), // <Suspense fallback={<NoteListSkeleton />}>
			// 	<NoteList searchText={searchText} />
			// </Suspense>

		),
		H("section", Props{"class": "col note-viewer"}), // 	<Suspense fallback={<NoteSkeleton isEditing={isEditing} />}>
	// 	<Note selectedId={selectedId} isEditing={isEditing} />
	// </Suspense>

	)
}

func SearchField() Component {
	return H("input")
}

func EditButton(rest ...any) Component {
	fmt.Println("EditButton", rest)
	return H("button")
}