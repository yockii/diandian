module diandian

go 1.24.0

toolchain go1.24.3

replace diandian/background/automation/core => ./background/automation/core

replace diandian/background/automation/hybrid => ./background/automation/hybrid

require (
	diandian/background/automation/core v0.0.0-00010101000000-000000000000
	diandian/background/automation/hybrid v0.0.0-00010101000000-000000000000
	github.com/glebarez/sqlite v1.11.0
	github.com/go-vgo/robotgo v0.110.8
	github.com/kbinani/screenshot v0.0.0-20250624051815-089614a94018
	github.com/sashabaranov/go-openai v1.41.2
	github.com/wailsapp/wails/v3 v3.0.0-alpha.28
	github.com/yockii/snowflake_ext v0.1.0
	golang.org/x/sys v0.33.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gorm.io/gorm v1.31.0
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.1.6 // indirect
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/bep/debounce v1.2.1 // indirect
	github.com/cloudflare/circl v1.6.0 // indirect
	github.com/cyphar/filepath-securejoin v0.4.1 // indirect
	github.com/dblohm7/wingoes v0.0.0-20240820181039-f2b84150679e // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/ebitengine/purego v0.8.3 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gen2brain/shm v0.1.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.2 // indirect
	github.com/go-git/go-git/v5 v5.13.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/leaanthony/go-ansi-parser v1.6.1 // indirect
	github.com/leaanthony/u v1.1.1 // indirect
	github.com/lmittmann/tint v1.0.7 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/micmonay/keybd_event v1.1.2 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/otiai10/gosseract v2.2.1+incompatible // indirect
	github.com/otiai10/mint v1.6.3 // indirect
	github.com/pjbgf/sha1cd v0.3.2 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robotn/xgb v0.10.0 // indirect
	github.com/robotn/xgbutil v0.10.0 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/shirou/gopsutil/v4 v4.25.4 // indirect
	github.com/skeema/knownhosts v1.3.1 // indirect
	github.com/tailscale/win v0.0.0-20250213223159-5992cb43ca35 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/vcaesar/gops v0.41.0 // indirect
	github.com/vcaesar/imgo v0.41.0 // indirect
	github.com/vcaesar/keycode v0.10.1 // indirect
	github.com/vcaesar/screenshot v0.11.1 // indirect
	github.com/vcaesar/tt v0.20.1 // indirect
	github.com/wailsapp/go-webview2 v1.0.21 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/image v0.27.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	modernc.org/libc v1.61.13 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.8.2 // indirect
	modernc.org/sqlite v1.36.0 // indirect
)
