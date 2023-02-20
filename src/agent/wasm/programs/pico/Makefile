# Copyright (C) 2023  tricorder-observability
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

all: wasm wat go-wrapper

native:
	@rm -f ./pico
	clang -Wall -Wextra -Werror picohttpparser.c main.c -o pico

wasm:
	WASI_SDK_PATH=/opt/wasi-sdk /opt/wasi-sdk/bin/clang --sysroot=/opt/wasi-sdk/share/wasi-sysroot \
	  -Wall -Wextra -Werror \
	  -Wl,--export=get_result_buf \
	  -Wl,--export=get_result_count \
	  -Wl,--export=get_result_struct_size \
	  -Wl,--export=phr_parse_request \
	  picohttpparser.c main.c -o pico.wasm

wat:
	/opt/wabt-1.0.31/bin/wasm2wat pico.wasm > pico.wat

go-wrapper:
	go build main-wasmedge.go
	go build main-wasmtime.go

clean:
	@rm -f *.wasm *.wat pico main-wasmedge main-wasmtime
