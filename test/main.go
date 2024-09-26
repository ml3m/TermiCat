package main

import (
    "fmt"
    "time"
    "github.com/eiannone/keyboard"
    "github.com/inancgumus/screen"
)

func main() {
    // Clear the screen
    screen.Clear()

    barLength := 50
    currentFill := 0
    fillRate := 1

    // Goroutine to drain the loading bar
    go func() {
        for {
            if currentFill > 0 {
                currentFill -= 1 // Drain the loading bar
            }
            time.Sleep(500 * time.Millisecond)
        }
    }()

    fmt.Println("Press 'Left' and 'Right' arrow keys alternately to fill the loading bar.")
    fmt.Println("Press 'ESC' to exit.")

    for {
        char, key, err := keyboard.GetKey()
        if err != nil {
            break
        }

        // Check for left and right arrow keys
        if key == keyboard.KeyArrowLeft || key == keyboard.KeyArrowRight {
            if (key == keyboard.KeyArrowLeft && currentFill < barLength) || (key == keyboard.KeyArrowRight && currentFill < barLength) {
                currentFill += fillRate
                if currentFill > barLength {
                    currentFill = barLength
                }
            }
        } else if char == 27 { // ESC key
            break
        }

        // Draw the loading bar
        screen.Clear()
        fmt.Print("Loading: [")
        for i := 0; i < barLength; i++ {
            if i < currentFill {
                fmt.Print("=")
            } else {
                fmt.Print(" ")
            }
        }
        fmt.Printf("] %d%%\n", (currentFill*100)/barLength) // Corrected to fmt.Printf
    }

    keyboard.Close()
    screen.Clear()
}
