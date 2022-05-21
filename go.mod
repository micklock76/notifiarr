module github.com/Notifiarr/notifiarr

go 1.17

// home grown goodness.
require (
	golift.io/cnfg v0.2.0
	golift.io/cnfgfile v0.0.0-20220509075834-08755d9ef3f5
	golift.io/deluge v0.9.4-0.20220317021225-51ab09943349
	golift.io/qbit v0.0.0-20220317021235-151a7d2d428e
	golift.io/rotatorr v0.0.0-20220126065426-d6b5acaac41c
	golift.io/starr v0.14.1-0.20220517052405-2529a7c6ece7
	golift.io/version v0.0.2
	golift.io/xtractr v0.1.1
)

// menu-ui
require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/gen2brain/beeep v0.0.0-20220402123239-6a3042f4b71a
	github.com/gen2brain/dlgs v0.0.0-20211108104213-bade24837f0b
	github.com/getlantern/context v0.0.0-20220418194847-3d5e7a086201 // indirect
	github.com/getlantern/errors v1.0.1 // indirect
	github.com/getlantern/golog v0.0.0-20211223150227-d4d95a44d873 // indirect
	github.com/getlantern/hex v0.0.0-20220104173244-ad7e4b9194dc // indirect
	github.com/getlantern/hidden v0.0.0-20220104173330-f221c5a24770 // indirect
	github.com/getlantern/ops v0.0.0-20220418195917-45286e0140f6 // indirect
	github.com/getlantern/systray v1.2.1
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-toast/toast v0.0.0-20190211030409-01e6764cf0a4 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gonutz/w32 v1.0.0
	github.com/iamacarpet/go-win64api v0.0.0-20220420140031-e6e5e4bf90e4
	github.com/tadvi/systray v0.0.0-20190226123456-11a2b8fa57af // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	gopkg.in/toast.v1 v1.0.0-20180812000517-0a84660828b2 // indirect
)

// snapshot and other stuff.
require (
	github.com/BurntSushi/toml v1.1.0
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/google/cabbie v1.0.3 // indirect
	github.com/google/glazier v0.0.0-20220513203518-91541a8d376d // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gopherjs/gopherwasm v1.1.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/jaypipes/ghw v0.9.0
	github.com/jaypipes/pcidb v1.0.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/lestrrat-go/apache-logformat v0.0.0-20200929122403-cd9b7dc018c7
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/lufia/plan9stats v0.0.0-20220326011226-f1430873d8db // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/scjalliance/comshim v0.0.0-20190308082608-cf06d2532c4e // indirect
	github.com/shirou/gopsutil/v3 v3.22.4
	github.com/spf13/pflag v1.0.6-0.20201009195203-85dd5c8bc61c
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.5.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3
	golang.org/x/net v0.0.0-20220513224357-95641704303c // indirect
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20220512140231-539c8e751b99 // indirect
	howett.net/plist v1.0.0 // indirect
)

// zip extraction for databsae corruption checks.
require (
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/saracen/go7z v0.0.0-20191010121135-9c09b6bd7fda // indirect
	github.com/saracen/go7z-fixtures v0.0.0-20190623165746-aa6b8fba1d2f // indirect
	github.com/saracen/solidblock v0.0.0-20190426153529-45df20abab6f // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
)

// sqlite3 abstraction for databsae corruption checks.
require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	lukechampine.com/uint128 v1.2.0 // indirect
	modernc.org/cc/v3 v3.36.0 // indirect
	modernc.org/ccgo/v3 v3.16.6 // indirect
	modernc.org/libc v1.16.8 // indirect
	modernc.org/mathutil v1.4.1 // indirect
	modernc.org/memory v1.1.1 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/sqlite v1.17.2
	modernc.org/strutil v1.1.1 // indirect
	modernc.org/token v1.0.0 // indirect
)

require (
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/fsnotify/fsnotify v1.5.4
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/schema v1.2.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/nxadm/tail v1.4.8
	golang.org/x/crypto v0.0.0-20220513210258-46612604a0f9
	golift.io/datacounter v1.0.3
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
)
