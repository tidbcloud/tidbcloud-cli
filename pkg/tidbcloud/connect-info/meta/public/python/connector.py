import mysql.connector

connection = mysql.connector.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
  ssl_ca = "${ca_path}",
  ssl_verify_cert = True,
  ssl_verify_identity = True
)
