# go-cleaner

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jabolopes/go-cleaner)](https://pkg.go.dev/github.com/jabolopes/go-cleaner)

Cleaner is useful to aggregate cleanup functions when a function has multiple success /
error exit conditions. Cleaner is especially useful in functions / methods that open or
acquire multiple resources which must all be destroyed if an error is reached. Cleaner is
also useful when the 'defer' idiom would lead to complicated code and error checking.

Example:

```go
func LoadFileAndTexture() (file *os.File, tex *sdl.Texture, clean func(), err error) {
  c, cleanup := cleaner.New()
  defer cleanup()

  if file, err = os.Open(...); err != nil {
    return
  }
  c.Add(func() { file.Close() })

  if tex, err = sdl.LoadTexture(...); err != nil {
    // file is automatically closed by the Cleaner when the function returns.
    return
  }
  c.Add(func() { tex.Destroy() })

  clean = c.Ok()  // file and tex remain in existence until the caller runs this 'clean' function.
  return
}
```

Example:


```go
// LoadGame loads this game and returns a cleanup function, and either
// nil (in case of success) or an error.
func LoadGame() (func(), error) {
  c, cleanup := cleaner.New()
  defer cleanup()

  tex, texCleanup, err := game.LoadTexture(...)
  if err != nil {
    return func() {}, err
  }
  c.Add(texCleanup)

  spr, sprCleanup, err := game.LoadSprite(...)
  if err != nil {
    return func() {}, err
  }
  c.Add(sprCleanup)

  lvl, lvlCleanup, err := game.LoadLevel(...)
  if err != nil {
    return func() {}, err
  }
  c.Add(lvlCleanup)

  log.Printf("All assets loaded")

  return c.Ok(), nil
}
```
