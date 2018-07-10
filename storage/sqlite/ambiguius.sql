select
    group_concat(distinct curse_addon.id) as ids,
    group_concat(distinct curse_addon.name) as names,
    t.dirs
from(
        select
            release_id,
            group_concat(directory, "-") as dirs
        from
            curse_release_directory
        group by
            release_id
        order by
            directory
    ) as t
    INNER JOIN curse_release on curse_release.id = t.release_id