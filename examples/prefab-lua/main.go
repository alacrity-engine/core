package main

import (
	"fmt"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/scripting"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

const script = `
f = loadstring("a = 6")
f()
print(a)

prefab = {
	Name = "zergon321",
	TransformRoot = {
		Position = Vec(32.12, 74.18),
		Angle = 16.24,
		Scale = Vec(1.5, 0.5),
		Gmob = {
			Name = "playah",
			Components = {
				{
					TypeName = "Health",
					Active = true,
					Data = {
						MaxHP = 100,
						CurrentHP = 50,
						Frequency = 10.01,
						Caption = "Life",
						Pos = Vec(10.25, 13.15)
					}
				}
			}
		}
	}
}

registerPrefab(prefab)
`

func registerPrefab(prefab *scripting.Prefab) {
	fmt.Println(prefab)
}

func Vec(x, y float64) geometry.Vec {
	return geometry.V(x, y)
}

func main() {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("registerPrefab", luar.New(L, registerPrefab))
	L.SetGlobal("Vec", luar.New(L, Vec))

	err := L.DoString(script)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
