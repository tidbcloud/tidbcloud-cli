mysql --comments -u '${username}' -h ${host} -P ${port} -D '${database}' --ssl-mode=VERIFY_IDENTITY --ssl-ca=${ca_path} -p'${password}'
