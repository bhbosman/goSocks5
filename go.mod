module github.com/bhbosman/goSocks5

go 1.18

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20220628093721-7191a8828f7f
	github.com/bhbosman/goCommsNetListener v0.0.0-20220611182354-46d10d89b8e1
	github.com/bhbosman/goCommsSshListener v0.0.0-00010101000000-000000000000
	github.com/bhbosman/goCommsStacks v0.0.0-20220701193851-2f07d2023eb2
	github.com/bhbosman/goFxApp v0.0.0-20220621055337-d5d1bfa131ee
	github.com/bhbosman/goMessages v0.0.0-20220701083921-8724d3c12872
	github.com/bhbosman/gocommon v0.0.0-20220628091717-f03fe56d9feb
	github.com/bhbosman/gocomms v0.0.0-20220628074707-e93417aaaed2
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cskr/pubsub v1.0.2
	github.com/golang/protobuf v1.4.1
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.17.1
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	google.golang.org/protobuf v1.25.0
)

require (
	github.com/bhbosman/goConnectionManager v0.0.0-20220628055237-4a995b1ae47c // indirect
	github.com/bhbosman/goFxAppManager v0.0.0-20220625173841-bbd050c3bfe2 // indirect
	github.com/bhbosman/goUi v0.0.0-20220625174028-03193a90ee79 // indirect
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc // indirect
	github.com/bhbosman/gomessageblock v0.0.0-20210901070622-be36a3f8d303 // indirect
	github.com/cenkalti/backoff/v4 v4.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.5.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/tview v0.0.0-20220610163003-691f46d6f500 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/bhbosman/gocomms => ../gocomms

replace github.com/bhbosman/gocommon => ../gocommon

replace github.com/golang/mock => ../gomock

replace github.com/bhbosman/goCommsNetListener => ../goCommsNetListener

replace github.com/bhbosman/goCommsSshListener => ../goCommsSshListener

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

replace github.com/bhbosman/goCommsStacks => ../goCommsStacks

replace github.com/bhbosman/goCommsSSHProtocols => ../goCommsSSHProtocols

replace github.com/bhbosman/goCommsSSH => ../goCommsSSH

replace github.com/bhbosman/goFxApp => ../goFxApp

replace github.com/bhbosman/goUi => ../goUi

replace github.com/bhbosman/goerrors => ../goerrors

replace github.com/bhbosman/goFxAppManager => ../goFxAppManager

replace github.com/bhbosman/goConnectionManager => ../goConnectionManager

replace github.com/rivo/tview => ../tview

replace github.com/bhbosman/goprotoextra => ../goprotoextra
