module github.com/kujilabo/cocotola-api

go 1.16

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.4.0
	github.com/casbin/casbin/v2 v2.34.1
	github.com/casbin/gorm-adapter/v3 v3.3.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.3.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/onrik/gorm-logrus v0.3.0
	github.com/onrik/logrus v0.9.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.6
	github.com/vektra/mockery v1.1.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.31.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.31.0
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/exporters/jaeger v1.6.3
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.6.3
	go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace v1.6.3
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/tools v0.1.8 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)
