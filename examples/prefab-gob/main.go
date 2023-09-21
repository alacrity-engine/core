package main

import (
	"encoding/gob"
	"os"

	"github.com/alacrity-engine/core/definitions"
	"github.com/alacrity-engine/core/geometry"
)

func main() {
	gob.Register(definitions.GameObjectPointer{})
	gob.Register(definitions.ComponentPointer{})

	prefab := &definitions.Prefab{
		Name: "player",
		TransformRoot: &definitions.TransformDefinition{
			Position: geometry.V(16.34, 12.32),
			Angle:    128.92,
			Scale:    geometry.V(1.5, 0.5),
			Gmob: &definitions.GameObjectDefinition{
				Name:    "player",
				ZUpdate: 14,
				Components: []*definitions.ComponentDefinition{
					{
						TypeName: "danmaku__player__Health",
						Active:   true,
						Data: map[string]interface{}{
							"maxHp": int64(100),
							"healthGUI": definitions.ComponentPointer{
								GmobName: "player-health",
								CompType: "health-gui",
							},
						},
					},
				},
			},
		},
	}

	file, err := os.Create("encoded.bin")
	handleError(err)
	enc := gob.NewEncoder(file)

	err = enc.Encode(prefab)
	handleError(err)
	err = file.Close()
	handleError(err)

	file, err = os.Open("encoded.bin")
	handleError(err)
	dec := gob.NewDecoder(file)

	var restoredPrefab definitions.Prefab
	err = dec.Decode(&restoredPrefab)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
