Name = BBoltViewer

Version = 0.0.1

BuildTime = $(shell date +'%Y-%m-%d %H:%M:%S')

LDFlags = -ldflags "-s -w -X '${Name}/version.version=$(Version)' -X '${Name}/version.buildTime=${BuildTime}'"

os-archs=darwin:arm64 linux:amd64 windows:amd64

build: app

win:
	go build -trimpath $(LDFlags) -o ./release/${Name}_windows_amd64.exe ./main.go

linux:
	go build -trimpath $(LDFlags) -o ./release/${Name}_linux_amd64 ./main.go

mac:
	go build -trimpath $(LDFlags) -o ./release/${Name}_darwin_arm64 ./main.go

app:
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} go build -trimpath $(LDFlags) -o ./release/${Name}_$${target_suffix} ./main.go;\
		echo "Build $${os}-$${arch} done";\
	)
	@mv ./release/${Name}_windows_amd64 ./release/${Name}_windows_amd64.exe

clean:
	rm -rf ./release