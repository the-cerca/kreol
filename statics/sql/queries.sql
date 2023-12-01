WITH UnseenWord AS (
    SELECT
        w.id,
        w.word,
        w.language_id,
        w.theme_id
    FROM
        words w
        LEFT JOIN words_views wv ON w.id = wv.word_id AND wv.user_id = $1    WHERE
        wv.id IS NULL
        AND w.theme_id = $2 
    ORDER BY
        w.created_at DESC
    LIMIT 1
),
WordTranslation AS (
    SELECT
        t.translation
    FROM
        translations t
        JOIN UnseenWord uw ON t.word_id = uw.id
),
RandomWord AS (
    SELECT
        t.id,
        t.translation
    FROM
        words w
        JOIN translations t ON w.id = t.word_id
    WHERE
        w.theme_id = (SELECT theme_id FROM UnseenWord)
        AND w.language_id = (SELECT language_id FROM UnseenWord)
        AND t.translation != (SELECT translation FROM WordTranslation)
    ORDER BY
        RANDOM()
    LIMIT 1
)
SELECT
    uw.id AS word_id,
    uw.word,
    wt.translation AS translation,
    rw.translation AS random_translation
FROM
    UnseenWord uw
    CROSS JOIN WordTranslation wt
    CROSS JOIN RandomWord rw;


-------------------------------------
SELECT
    l.id,
    l.name
FROM
    user_languages ul
    JOIN languages l ON ul.language_id = l.id
WHERE
    ul.user_id = $1
