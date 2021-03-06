module github.com/bhbosman/goSocks5

go 1.18

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20220724214237-63eea03e3695
	github.com/bhbosman/goCommsNetListener v0.0.0-20220712191421-bfbc514d3011
	github.com/bhbosman/goFxApp v0.0.0-20220707083540-2141464f4f7c
	github.com/bhbosman/gocommon v0.0.0-20220725200742-9cdc334065f3
	github.com/bhbosman/gocomms v0.0.0-20220713064628-4df3e595ffd0
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/gomessageblock v0.0.0-20220617132215-32f430d7de62
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cskr/pubsub v1.0.2
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.17.1
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
)

require (
	github.com/bhbosman/goConnectionManager v0.0.0-20220721070628-0f4b3c036d93 // indirect
	github.com/bhbosman/goFxAppManager v0.0.0-20220713064906-f13252f1e610 // indirect
	github.com/bhbosman/goUi v0.0.0-20220725200743-ddc6ed05f1d6 // indirect
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

replace github.com/cskr/pubsub => ../pubsub
