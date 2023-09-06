.PHONY: run

setup:
	sudo npm install -g typescript@4.2.3
	sudo npm install -g less@3.12.2
	sudo apt-get install inotify-tools

compile: ts-compile less-compile

clean: ts-clean less-clean

run: clean compile
	@go run main.go

dev-run: less-compile-watch ts-compile-watch run

# TS RELATED TARGETS AND VARIABLES
ts_src_root_path = presentation/web/assets/js/ts
ts_out_root_path = presentation/web/assets/js/ts_compiled

ts_files = $(shell find $(ts_src_root_path) -name '*.ts')
ts_flags = --project presentation/web/tsconfig.json
#ts_bundled_out_path = --outFile $(ts_out_root_path)/ts-bundle.js

ts-compile:
	@echo 'Compiling typescript files'
	@tsc $(ts_flags) $?
	@echo 'Done compiling typescript'

#ts-bundle:
#	@tsc $(ts_flags) $(ts_bundled_out_path) $?

ts-compile-watch:
	@tsc -w $(ts_flags) $?

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
	chmod +x ./less_compiler_deamon.sh
	./less_compiler_deamon.sh $(less_src_root_path)

all-compile-watch: ts-compile-watch less-compile-watch

