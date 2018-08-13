# merge-config

A small library to merge multiple unmarshalled `map[string]interface{}`
variables together into one `map[string]interface{}` that can then be parsed.

The library assumes that the `map[string]interface{}` was created by a call to
an `Unmarshal` method or a previous call to this library's `Merge` function.
That is, it assumes all maps are of type `map[string]interface{}`, all slices
are `[]interface{}`, and all other types present are not pointer types.