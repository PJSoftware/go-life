package shader

import (
	"log"
	"os"
)

func Import(glslFile string) string {
	content, err := os.ReadFile("shader/shaders/" + glslFile + ".glsl")
	if err != nil {
		log.Fatal(err)
	}
	return string(content) + "\x00"	// shader string must be null-terminated to compile
}
