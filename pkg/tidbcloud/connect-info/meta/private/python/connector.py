import mysql.connector

connection = mysql.connector.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
)
