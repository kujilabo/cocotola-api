migrate -database 'mysql://user:password@tcp(127.0.0.1:3306)/development' -source file://sqls/msyql/ drop
migrate -database 'mysql://user:password@tcp(127.0.0.1:3306)/development' -source file://sqls/mysql/ up