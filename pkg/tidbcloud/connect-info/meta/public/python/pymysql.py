import pymysql

connection = pymysql.connect(
  host = "${host}",
  port = ${port},
  user = "${username}",
  password = "${password}",
  database = "${database}",
  ssl_verify_cert = True,
  ssl_verify_identity = True,
  ssl_ca = "${ca_path}"
)
