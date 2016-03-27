package main

/*

s9y2hugo is designed to aid in migrating from a Serendipity blog to a Hugo one.

It takes input as CSV, generated by an accompanying SQL script, and
writes a set of files in Hugo format.

It will carry over the title, frontmatter, and body of the entries,
and use Hugo's "alias" feature to carry over the permalinks as
aliases.

The input spec is straightforward and derived from the serendipity-extract:
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
	//"encoding/csv"
	"fmt"
	//"io"
	//"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Post struct {
     Date 			string
     Title			string
     Tags	[]		string
     Categories	[]string
     Permalink	[]string
     Body				string
		 Extended		string
}

const templ = `
+++
Title	= "{{.Title}}"
Date	= "{{.Date}}"
Categories = [{{range $index, $elmt := .Categories}}{{if $index}},"{{$elmt}}"{{else}}"{{$elmt}}"{{end}}{{end}}]
Tags	= [{{range $index, $elmt := .Tags}}{{if $index}},"{{$elmt}}"{{else}}"{{$elmt}}"{{end}}{{end}}]
Aliases = [{{range $index, $elmt := .Permalink}}{{if $index}},"{{$elmt}}"{{else}}"{{$elmt}}"{{end}}{{end}}]
+++
{{.Body}}
{{.Extended}}
`

func main() {
/*
	 in := `timestamp,title,description, tags, categories, permalink, project_url, body
2005-11-26T19:07:00ZNZDT,French Toast,Test,,"[""Food""]",archives/775-French-Toast.html,http://diaspora.gen.nz/~rodgerd/,"<p>I'm just not that fond of french toast as it appears in most Wellington cafés. Mostly it's done with thick bread and incredibly sweet, usually with fruit or maple syrup.</p><p>I grew up with french toast being a savoury treat: egg, milk, pepper, and salt whisked together, with plain toast slice bread dunked into it and then pan-fried in butter until golden brown. Dee-licious.</p>"
`
*/

	// r := csv.NewReader(strings.NewReader(in))

	posts := []Post{
		{
			Title:	"French Toast",
			Date:	"1132985220",
			Tags:	[]string{"food", "wellington"},
			Categories: []string{"Food"},
			Permalink: []string{"archives/775-French-Toast"},
			Body: 	"<p>I'm just not that fond of french toast as it appears in most Wellington cafés. Mostly it's done with thick bread and incredibly sweet, usually with fruit or maple syrup.</p><p>I grew up with french toast being a savoury treat: egg, milk, pepper, and salt whisked together, with plain toast slice bread dunked into it and then pan-fried in butter until golden brown. Dee-licious.</p>",
			Extended: "",
		},
		{
			Title:	"Bubba Ho-Tep",
			Date:	"1132985220",
			Tags:	[]string{"food", "wellington"},
		 	Categories: []string{"Movies"},
			Permalink: []string{"archives/776-Bubba-Ho-Tep"},
			Body: 	`<p>The premise of <a href="http://www.imdb.com/title/tt0281686/">Bubba Ho-Tep</a> is bizarre and amusing: Elvis lives.  In a rest home.  With a black guy who thinks he&#8217;s JFK.  A mummy is stalking the corridors and, well, &#8220;ask not what your rest home can do for you, but what you can do for your rest home.&#8221;</p>

					 	<p>The thing is, it&#8217;s creepy, but not like you&#8217;d expect&#8230;</p>

					 	<p>As a horror movie, it&#8217;s pretty much what you&#8217;d expect from anything with <a href="http://www.imdb.com/name/nm0132257/">Bruce Campbell</a>: a good, solid, slightly scary, tongue in cheek horror.  The thing is it&#8217;s actually creepier at the start than once the movie gets into full swing&#151;and mummies have nothing to do with it.  Rather, it&#8217;s Campbell&#8217;s confused old man in a shitty rest home that sent chills up and down my spine.</p>

					 	<p>It&#8217;s enough to make you hope you don&#8217;t ever grow old.</p>`,
			Extended: "",
		},
	}

	for _, post := range posts {
		/*
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		*/

		// We call some routines that tidy up the input data in a variety of ways.
		filename := makeFilename(post.Permalink[0])
		post.Date = (makeDate(post.Date))

		// Process the entries through the blog template.
		// Output one entry per file.
		t := template.New("Post template")
		t, err := t.Parse(templ)
		checkError(err)

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
