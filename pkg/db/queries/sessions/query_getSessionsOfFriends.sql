WITH {{.DB_Tables.SessionUsers.TableName}} AS (
    SELECT DISTINCT {{.DB_Tables.SessionUsers.SessionID}}, {{.DB_Tables.SessionUsers.UserID}}
    FROM {{.DB_Tables.SessionUsers}}
),
non_friend_users AS (
    SELECT su.{{.DB_Tables.SessionUsers.SessionID}}
    FROM {{.DB_Tables.SessionUsers}} su
    WHERE su.{{.DB_Tables.SessionUsers.UserID}} != $1
    AND NOT EXISTS (
        SELECT 1 FROM {{.DB_Tables.Friendships}} f
        WHERE (f.{{.DB_Tables.Friendships.UserID1}} = $1 AND f.{{.DB_Tables.Friendships.UserID2}} = su.{{.DB_Tables.SessionUsers.UserID}})
        OR (f.{{.DB_Tables.Friendships.UserID2}} = $1 AND f.{{.DB_Tables.Friendships.UserID1}} = su.{{.DB_Tables.SessionUsers.UserID}})
    )
)
SELECT DISTINCT s.*
FROM {{.DB_Tables.Sessions}} s
WHERE NOT EXISTS (
    SELECT 1 FROM non_friend_users nfu
    WHERE nfu.{{.DB_Tables.SessionUsers.SessionID}} = s.{{.DB_Tables.Sessions.ID}}
)
AND EXISTS (
    SELECT 1 FROM {{.DB_Tables.SessionUsers.TableName}} su
    WHERE su.{{.DB_Tables.SessionUsers.SessionID}} = s.{{.DB_Tables.Sessions.ID}}
);
