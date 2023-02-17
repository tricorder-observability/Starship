  (module
    (import "" "hello" (func $hello))
    (func (export "run")
      (call $hello))
  )
