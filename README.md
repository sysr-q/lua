# sysr-q/lua

Eventual aims to be a [Lua](http://www.lua.org) lexer/parser/interpreter
(and maybe JIT compiler) in pure [Go](https://golang.org).

Mainly since I want an embeddable Lua runtime in a statically linked Go binary.

* Aims to target Lua __5.3__
* Probably not going to provide Lua's C interop; would be too complicated
* Realistically should pass all of Lua's tests (`lua-5.3.0-tests/`) with the
  `-e"_U=true"` flag equiv set
