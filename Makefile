.PHONY: run

setup:
	sudo npm install -g typescript
	sudo npm install -g less
	sudo apt-get install inotify-tools

compile: ts-compile less-compile

clean: ts-clean less-clean

run: clean compile
	@go run main.go

dev-run: ts-compile-watch run

# TS RELATED TARGETS AND VARIABLES
ts_src_root_path = presentation/web/assets/js/ts
ts_out_root_path = presentation/web/assets/js/ts_compiled

ts_files = $(shell find $(ts_src_root_path) -name '*.ts')
ts_flags = --pretty --noImplicitAny --noImplicitReturns --noFallthroughCasesInSwitch --rootDir $(ts_src_root_path)
ts_unbundled_out_path = --outDir $(ts_out_root_path)
#ts_bundled_out_path = --outFile $(ts_out_root_path)/ts-bundle.js

ts-compile: $(ts_files)
	@echo 'Compiling typescript files at "$(ts_src_root_path)"...'
	@tsc $(ts_flags) $(ts_unbundled_out_path) $?
	@echo 'Done compiling typescript'

#ts-bundle: $(ts_files)
#	@tsc $(ts_flags) $(ts_bundled_out_path) $?

ts-compile-watch: $(ts_files)
# TODO(vgiordano): The above cmd only watchs modifications over existing files. Not new ones. Would be nice to polish in order to compile new ones.
	@tsc -w $(ts_flags) $(ts_unbundled_out_path) $?

ts-clean:
	rm -rfv $(ts_out_root_path)

# LESS RELATED TARGETS AND VARIABLES
less_src_root_path = presentation/web/assets/css/less
less_out_root_path = presentation/web/assets/css/less_compiled

less_files = $(shell find $(less_src_root_path) -name '*.less')

less-compile:
	@echo 'Compiling less files at "$(less_src_root_path)"...'
	@echo $(shell for file in $(less_files);do lessc --strict-imports $$file  `dirname $$file | sed -e "s/less/less_compiled/"`/`basename $$file | sed -e "s/less/css/"`; done)
	@echo 'Done compiling less'

less-clean:
	rm -rfv $(less_out_root_path)

less-compile-watch: $(less_files)
	chmod +x scripts/less_compiler_deamon.sh
	scripts/less_compiler_deamon.sh $(less_src_root_path)

all-compile-watch: ts-compile-watch less-compile-watch

