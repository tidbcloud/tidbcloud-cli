import pymysql

connection = pymysql.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
)
