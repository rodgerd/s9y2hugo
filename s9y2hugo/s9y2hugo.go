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
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"net/url"
	"os"
	"strconv"
	"text/template"
	"time"
)

type Post struct {
	Date       string
	Title      string
	Tags       string
	Categories string
	Permalinks string
	isDraft    string
	Body       string
	Extended   string
}

const templ = `
+++
Title	= "{{.Title}}"
Date	= "{{.Date}}"
Categories = {{.Categories}}
Tags	= {{.Tags}}
Aliases = {{.Permalinks}}
+++
{{.Body}}
{{.Extended}}
`

func main() {

	userPtr := flag.String("user", "serendipity", "Database username")
	passwordPtr := flag.String("password", "", "Database password")
	dbnamePtr := flag.String("dbname", "serendipity", "Database name")
	hostnamePtr := flag.String("host", "localhost", "Host name")
	sslPtr := flag.String("sslmode", "disable", "SSL Mode")
	flag.Parse()

	connStr := "user=" + *userPtr + " password=" + *passwordPtr + " dbname=" + *dbnamePtr + " host=" + *hostnamePtr + " sslmode=" + *sslPtr
	db, err := sql.Open("postgres", connStr)
	checkError(err)
	defer db.Close()

	rows, err := db.Query(`
		select	e.id,
						e.timestamp,
						e.title,
						(
							select coalesce(json_agg(tag), '["no-tag"]')
							from		serendipity_entrytags where entryid = e.id
							) as tags,
						(
							select	coalesce(json_agg(category_name), '["no-cat"]')
							from  	serendipity_category c, serendipity_entrycat ec
							where  	c.categoryid = ec.categoryid and
											ec.entryid = e.id
						) as categories,
						(
							select	coalesce(json_agg(permalink), '["no-link"]')
							from		serendipity_permalinks p
							where		p.entry_id = e.id
						) as url,
						e.isdraft as isdraft,
						e.body as body,
						e.extended as extended
		from		serendipity_entries e
		order by e.id
	`)
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		var (
			id         string
			timestamp  string
			title      string
			tags       string
			categories string
			permalinks string
			isDraft    string
			body       string
			extended   string
		)
		err := rows.Scan(&id, &timestamp, &title, &tags, &categories, &permalinks, &isDraft, &body, &extended)
		// Transform the record into a Post
		post := Post{
			Title:      title,
			Date:       makeDate(timestamp),
			Tags:       tags,
			Categories: categories,
			Permalinks: permalinks,
			isDraft:    isDraft,
			Body:       body,
			Extended:   extended,
		}

		// Process the entries through the blog template.
		// Output one entry per file.
		t := template.New("Post template")
		t, err = t.Parse(templ)
		checkError(err)

		filename := makeFilename(id, title)
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
func makeDate(old string) string {
	i, err := strconv.ParseInt(old, 10, 64)
	checkError(err)
	t := time.Unix(i, 0)
	// fmt.Println(t.Format(time.RFC3339))
	return t.Format(time.RFC3339)
}

/*
 Accepts a permalink and turns it into a file.

 permalinks are assumed to be in the format 'archives/entry_id-slug.html; these will be transformed into entry_id-slug.md as the output file.'
*/
func makeFilename(id string, title string) string {
	name := id + "-" + url.QueryEscape(title) + ".md"
	return name
}
