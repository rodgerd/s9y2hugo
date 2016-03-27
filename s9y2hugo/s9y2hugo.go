package main

/*

s9y2hugo is designed to aid in migrating from a Serendipity blog to a Hugo one.

It connects directly to postgresql and writes a set of files in Hugo format.

It will carry over the title, frontmatter, and body of the entries,
and use Hugo's "alias" feature to carry over the permalinks as
aliases.

The mapping is straightforward:

    Column    Source Description     				Output Field
    ------    ------------------	      	  ------------
    [0]       Unix epoch										Date in RFC3339 format
    [1]       Title			      							Title
    [2]       Array of tags  								Tags in JSON format
    [4]       Array of categories						Categories in JSON format
    [5]       Array of permalinks 		      					Aliases in JSON format
		[6]				isdraft												Whether the entry is draft or not.
    [7]       Body						      				Body text.
		[8]				extended											Any extended text.

*/

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
	"github.com/lib/pq"
)

type Post struct {
     Date 			string
     Title			string
     Tags	[]		string
     Categories	[]string
     Permalinks	[]string
		 isDraft		string
     Body				string
		 Extended		string
}

const templ = `
+++
Title	= "{{.Title}}"
Date	= "{{.Date}}"
Categories = [{{range $index, $elmt := .Categories}}{{if $index}},"{{$elmt}}"{{else}}"{{$elmt}}"{{end}}{{end}}]
Tags	= [{{range $index, $elmt := .Tags}}{{if $index}},"{{$elmt}}"{{else}}"{{$elmt}}"{{end}}{{end}}]
Aliases = [{{range $index, $elmt := .Permalinks}}{{if $index}},"{{$elmt}}"{{else}}"{{$elmt}}"{{end}}{{end}}]
+++
{{.Body}}
{{.Extended}}
`

func main() {

	db, err := sql.Open("postgres", "user=foo dbname=bar")
	checkError(err)

	rows, err := db.Query(`
		select	e.timestamp,
						e.title,
						(
							select 	array_agg(tag)
							from		serendipity_entrytags where entryid = e.id
							) as tags,
						(
							select	array_agg(category_name)
							from  	serendipity_category c, serendipity_entrycat ec
							where  	c.categoryid = ec.categoryid and
											ec.entryid = e.id
						) as categories,
						(
							select	array_agg(permalink)
							from		serendipity_permalinks p
							where		p.entry_id = e.id
						) as url,
						e.isdraft as isdraft,
						e.body as body,
						e.extended as extended
		from		serendipity_entries e
		order by e.timestamp
		limit 10
	`)

	for row in rows  {

		// Transform the record into a Post
		post := Post{
			Title: 	row[1],
			Date:		makeDate(record[0]),
			Tags:		strings.Split(record[2], ", "),
			Categories: strings.Split(record[3], ", "),
			Permalinks: strings.Split(record[4], ", "),
			isDraft: record[5],
			Body:	record[6],
			Extended: record[7],
		}

		// Process the entries through the blog template.
		// Output one entry per file.
		t := template.New("Post template")
		t, err = t.Parse(templ)
		checkError(err)

		filename := makeFilename(post.Permalinks[0])
		file, err := os.Create(filename)
		checkError(err)
		defer file.Close()

		err = t.Execute(file, post)
		checkError(err)

		file.Sync()
	}
}

func checkError(err error) {
     if err != nil {
     	fmt.Println("Fatal error ", err.Error())
			os.Exit(1)
     }
}

/*
	Convert the date as extracted from postgresql into RFC3339 format so hugo will parse it correctly
*/
func makeDate(old string) (string) {
	i, err := strconv.ParseInt(old, 10, 64)
	checkError(err)
	t := time.Unix(i,0)
	// fmt.Println(t.Format(time.RFC3339))
	return t.Format(time.RFC3339)
}

/*
 Accepts a permalink and turns it into a file.

 permalinks are assumed to be in the format 'archives/entry_id-slug.html; these will be transformed into entry_id-slug.md as the output file.'
*/
func makeFilename(permalink string) (string) {
	i, j := strings.LastIndex(permalink, "/") + 1, strings.LastIndex(permalink, path.Ext(permalink))
	name := permalink[i:j] + ".md"
	return name
}
