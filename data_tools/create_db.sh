#!/bin/sh

DATA_FILE="../data.db"
if [ ! -f "$DATA_FILE" ]; then
	echo "CREATE TABLE users ( id integer primary key, first_name TEXT, last_name TEXT, email text unique);"  | sqlite3 ../data.db
fi
