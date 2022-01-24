module velvet

go 1.17

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/df-mc/dragonfly v0.5.2-0.20220122125012-15d7cef7ff2a
	github.com/emperials/df-worldmanager v0.0.0-20210810032941-29fc34569fb4
	github.com/go-gl/mathgl v1.0.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/pelletier/go-toml v1.9.3
	github.com/sandertv/gophertunnel v1.18.2
	github.com/sirupsen/logrus v1.8.1
	go.uber.org/atomic v1.9.0
)

require (
	github.com/brentp/intintmap v0.0.0-20190211203843-30dc0ade9af9 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/df-mc/goleveldb v1.1.9 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/muhammadmuzzammil1998/jsonc v0.0.0-20201229145248-615b0916ca38 // indirect
	github.com/sandertv/go-raknet v1.10.2 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

//replace github.com/df-mc/dragonfly => ../dragonfly

replace github.com/emperials/df-worldmanager => ../df-worldmanager
