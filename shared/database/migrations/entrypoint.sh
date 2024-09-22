#!/bin/bash

#DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"

goose -table _migrations mysql "root:P@ssw0r1d@tcp(mysql:3306)/ticketmasterDB?charset=utf8mb4&parseTime=True&loc=Local" up

#goose mysql "$DBSTRING" up