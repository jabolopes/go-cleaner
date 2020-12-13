package cleaner

// Cleaner is useful to aggregate cleanup functions when a function
// has multiple success / error exit conditions. Cleaner is especially
// useful in functions / methods that open or acquire multiple
// resources which must all be destroyed if an error is
// reached. Cleaner is also useful when the 'defer' idiom would lead
// to complicated code and error checking.
//
// Example:
//
//   func LoadFileAndTexture() (file *os.File, tex *sdl.Texture, clean func(), err error) {
//     c, cleanup := cleaner.New()
//     defer cleanup()
//
//     if file, err = os.Open(...); err != nil {
//       return
//     }
//     c.Add(func() { file.Close() })
//
//     if tex, err = sdl.LoadTexture(...); err != nil {
//       // file is automatically closed by the Cleaner when the function returns.
//       return
//     }
//     c.Add(func() { tex.Destroy() })
//
//     clean = c.Ok()  // file and tex remain in existence until the caller runs this 'clean' function.
//     return
//   }
type Cleaner struct {
  closures []func()
}

// Add adds a new closure to this cleaner.
func (d *Cleaner) Add(closure func()) {
  d.closures = append(d.closures, closure)
}

// Ok removes all closures from this cleaner and returns a function
// that runs all of those closures instead. This is typically used
// when a function has completed successfully and the cleanup closures
// should be returned to the caller.
func (d *Cleaner) Ok() func() {
  closures := d.closures
  d.closures = nil

  return func() {
    for i := len(closures) - 1; i >= 0; i-- {
      closures[i]()
    }
  }
}

// New returns a new cleaner and a cleanup function that must
// eventually be called to call then cleanup closures (unless Ok is
// called).
func New() (*Cleaner, func()) {
  d := &Cleaner{nil /* closures */}
  return d, func() { d.Ok()() }
}
