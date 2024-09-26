package main

import (
	"TermiCat/asciiart"
	"time"
)


func (c *cat) GainXP(amount int, m *model) {
    c.Xp += amount
    if c.Xp >= c.XPThreshold() {
        c.LevelUp(m)
    }
}

func (c *cat) LevelUp(m *model) {
    c.Level++
    c.Xp = 0 
    c.Health += 10 
    c.Happiness += 5 
    c.Energy += 10
    c.Coins += 100

    m.IsLevelUpAnimating = true
    m.Frames = asciiart.GetLevelUpFrames()
    m.LevelUpStartTime = time.Now()
    
}

func (c *cat) XPThreshold() int {
    // level up hardness grows by 100 each time. 
    return 100 * c.Level 
}


