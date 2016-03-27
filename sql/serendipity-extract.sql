copy (
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
from	serendipity_entries e
where 	e.id = 1334
order by e.timestamp
limit 1
) to stdout with delimiter = '|';
