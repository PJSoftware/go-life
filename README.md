# Conway's Game of Life

It has been some time since I:

- actively used GitHub
- programmed in Go

I have been thinking more and more of wanting to do _something_ graphical, so why not return to my Golang/GitHub roots and once again tackle Conway's Game of Life. (For anyone interested, my first pass is [here](https://github.com/PJSoftware/game-of-life).)

This code will be adapted from an [OpenGL & Go Tutorial](https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl) I found online; just something to play with to familiarise myself with `OpenGL` while refreshing my `Go` memory.

## Notes on Installing OpenGL

The first steps are to init our module, and download the relevant packages:

```sh
go mod init github.com/PJSoftware/go-life
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw
```

## OpenGL Shaders

See the [Version 4.10 Reference Manual](https://www.khronos.org/registry/OpenGL/specs/gl/GLSLangSpec.4.10.pdf).

The code worked -- produced a white triangle on a black screen -- before we added the shader code. The tutorial said it wouldn't.

Clearly OpenGL shaders are an important topic that I'll need to wrap my head around, but for now we'll simply keep them in here as presented.

## Next Steps

The tutorial suggests the following challenges:

- [ ] Give each cell a unique color.
- [ ] Allow the user to specify, via command-line arguments, the grid size, frame rate, seed and threshold.
  - You can see this one implemented [on GitHub](https://github.com/KyleBanks/conways-gol).
- [ ] Change the shape of the cells into something more interesting, like a hexagon.
- [ ] Use color to indicate the cell’s state - for example, make cells green on the first frame that they’re alive, and make them yellow if they’ve been alive more than three frames.
- [ ] Automatically close the window if the simulation completes, meaning all cells are dead or no cells have changed state in the last two frames.
- [x] Move the shader source code out into their own files, rather than having them as string constants in the Go source code.

I've already implemented #6 because it seemed a better approach at the time.

The other thing I've modified from the original tutorial is changing the size of the cells to allow a pixel border around each cell. I much prefer how this looks.
