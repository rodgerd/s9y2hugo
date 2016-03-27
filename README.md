# s9y2hugo

s9y2hugo is designed to aid in migrating from a Serendipity blog to a Hugo one.

It connects to a postgresql database, pulls the data out, and will dump it out as a set of files, with the names based on the filename of the original.

It will carry over the title, front matter, and body of the entries,
and use Hugo's "alias" feature to carry over the permalinks as
aliases.

## Requirements

It has been written against postgresql and uses the to_json feature which was merged in postgresql 9.3.

It relies on the github.com/lib/pg driver for postgresql access.

## Usage

s9y2hugo takes the following parameters:

-user
-password
-dbname
-host

## Limitations and Gotchas

This was a go learning/quick hack exercise.  It bears all the hallmarks of someone who can write the same code in BASIC, Perl, TCL, Python, and Java.

It doesn't support databases other than postgresql.  If you want to add support for other DBs, feel free.

lib/pg does not always fail gracefully or informatively, and this code will tend to leave you staring at a cryptic error.  Sorry.  It's probably a null column returned in a row.
