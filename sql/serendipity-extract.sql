copy (
select	e.timestamp, 
	e.title, 
	'' as description, 
	'["' || (
		select 	array_to_string(array_agg(tag), '","')
		from 	serendipity_entrytags where entryid = e.id
	) || '"]' as tags, 
	'["' || (
		select	array_to_string(array_agg(category_name), '","') 
		from  	serendipity_category c, serendipity_entrycat ec 
		where  	c.categoryid = ec.categoryid and 
			ec.entryid = e.id
	) || '"]' as topics, 
	-- '' as slug,	-- slug is a relative path to the web root
	p.permalink as url,	-- URL gives a hardcoded absolute path
	'http://diaspora.gen.nz/~rodgerd/' as project_url,
	e.body as body 
from	serendipity_entries e, serendipity_permalinks p
where 	e.id = p.entry_id and
	e.id = 1334 
order by e.timestamp 
limit 1
) to stdout with csv;
