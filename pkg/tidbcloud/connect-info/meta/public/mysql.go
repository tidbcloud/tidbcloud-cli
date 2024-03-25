mysql.RegisterTLSConfig("tidb", &tls.Config{
  MinVersion: tls.VersionTLS12,
  ServerName: "${host}",
})

db, err := sql.Open("mysql", "${username}:${password}@tcp(${host}:${port})/${database}?tls=tidb")
