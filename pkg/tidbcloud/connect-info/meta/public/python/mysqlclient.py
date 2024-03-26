import MySQLdb

connection = MySQLdb.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
  ssl_mode = "VERIFY_IDENTITY",
  ssl = {
    "ca": "${ca_path}"
  }
)
