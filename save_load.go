package main

import (
    "encoding/json"
    "os"
)

// Save the game state to a file
func SaveGameData(game model, filename string) error {
    data, err := json.MarshalIndent(game, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

// Load the game state from a file
func loadGameData(filename string) (model, error) {
    var game model 
    data, err := os.ReadFile(filename)
    if err != nil {
        return game, err
    }
    err = json.Unmarshal(data, &game)
    return game, err
}
