package main

import (
	"TermiCat/asciiart"
    "time"
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

func loadDefaultSettings() model {
	frames := asciiart.GetFrames()

	m := model{
        // cat init
        MyCat: cat{
            Wellness:          100,
            Fullness:           50,
            Hunger:             50,
            Dirtiness:          10,
            Happiness:          70,
            Energy:             80,
            Coins:               0,
            Level:               1,
            Xp:                  0,
            Health:            100,
            Boredom:            20,
            Age:                 1,
            Breed:       "Siamese",
            Name:          "Termi",
            LastFed:     time.Now(),
            LastCleaned: time.Now(),
        },

		Frames:       frames,

		ButtonLabels: []string{
            "Feed", 
            "Games", 
            "Clean", 
            "Sleep"}, // Define button labels

        ShowFoodMenu: false,
        

        FoodButtonLabels: []string{ "Feed", 
            "Buy Food"},

        FoodInventory: []food{
            {Name: "Apple", FeedingPower: 10, Cost: 5,    Quantity: 3},
            {Name: "Fish", FeedingPower: 20,  Cost: 10,   Quantity: 2},
            {Name: "Meat", FeedingPower: 30,  Cost: 15,   Quantity: 5},
            {Name: "Milk", FeedingPower: 5,   Cost: 3,    Quantity: 8},
            {Name: "Bread", FeedingPower: 8,  Cost: 4,    Quantity: 6},
            {Name: "Cheese", FeedingPower: 12,Cost: 6,    Quantity: 4},
            {Name: "Cake", FeedingPower: 25,  Cost: 20,   Quantity: 1},
            {Name: "Steak", FeedingPower: 35, Cost: 25,   Quantity: 2},
            {Name: "Carrot", FeedingPower: 7, Cost: 2,    Quantity: 0},
            {Name: "Berry", FeedingPower: 4,  Cost: 1,    Quantity: 5},
            {Name: "Nut", FeedingPower: 3,    Cost: 1,    Quantity: 0},
            {Name: "Egg", FeedingPower: 10,   Cost: 5,    Quantity: 7},
            {Name: "Soup", FeedingPower: 15,  Cost: 12,   Quantity: 3},
            {Name: "Pudding", FeedingPower: 18, Cost: 14, Quantity: 2},
            {Name: "Pie", FeedingPower: 22,     Cost: 18, Quantity: 1},
        },

        ShowInventoryMenu: false,
        SelectedFoodIndex: 0, // Reset selected food index

        // ANIMATION
        IsLevelUpAnimating: false,

        // <> Cleaning <>
        ShowCleanMenu: false,
	}
    
    return m
}
