# Conway's Game of Life

It has been some time since I:

- actively used GitHub
- programmed in Go

I have been thinking more and more of wanting to do _something_ graphical, so why not return to my Golang/GitHub roots and once again tackle Conway's Game of Life. (For anyone interested, my first pass is [here](https://github.com/PJSoftware/game-of-life).)

This code will be adapted from an [OpenGL & Go Tutorial](https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl) I found online; just something to play with to familiarise myself with `OpenGL` while refreshing my `Go` memory.

## Notes in Installing OpenGL

The first steps are to init our module, and download the relevant packages:

```sh
go mod init github.com/PJSoftware/go-life
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw
```
