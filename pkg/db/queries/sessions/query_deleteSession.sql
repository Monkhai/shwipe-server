DELETE FROM {{.SessionsTable.TableName}} WHERE {{.SessionsTable.ID}} = $1;