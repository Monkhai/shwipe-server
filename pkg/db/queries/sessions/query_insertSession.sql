INSERT INTO
 {{.SessionsTable.TableName}} ({{.SessionsTable.ID}}, {{.SessionsTable.CreatedAt}}) 
VALUES ($1, $2);