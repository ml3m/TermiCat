package main

import (
	"TermiCat/asciiart"
)


func (m *model) handleCatState() {
    if m.MyCat.Health <= 0 {
        // death
        m.Frames = asciiart.GetDeadCat() // Load the death frames if the cat is dead
    }
    
    if m.MyCat.Hunger > 100 {
        m.MyCat.Hunger = 100 // Cap hunger at 100
    }
    
    if m.MyCat.Fullness > 100 {
        m.MyCat.Fullness = 100 // Cap fullness at 100
    }

    // Ensure fullness does not drop below 0
    if m.MyCat.Fullness < 0 {
        m.MyCat.Fullness = 0
    }
    
    if m.MyCat.Wellness < 60 {
        m.MyCat.Happiness -= 10 // Decrease happiness
        if m.MyCat.Happiness < 0 {
            m.MyCat.Happiness = 0 // Ensure happiness doesn't go below 0
        }
    }
    
    if m.MyCat.Dirtiness > 80 {
        m.MyCat.Wellness -= 5 // Decrease wellness
        if m.MyCat.Wellness < 0 {
            m.MyCat.Wellness = 0 // Ensure wellness doesn't go below 0
        }
    }

    if m.MyCat.Energy < 50 {
        m.MyCat.Happiness -= 5 // Decrease happiness
        if m.MyCat.Happiness < 0 {
            m.MyCat.Happiness = 0 // Ensure happiness doesn't go below 0
        }
    }

    if m.MyCat.Age > 15 {
        m.MyCat.Health -= 10 // Decrease health due to age
        if m.MyCat.Health < 0 {
            m.MyCat.Health = 0 // Ensure health doesn't go below 0
        }
    }
}
