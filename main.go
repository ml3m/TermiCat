package main

import (
	"TermiCat/asciiart" // my ascii
	"fmt"
	"log"
	//"os"
	"time"
    "strings"
	//"golang.org/x/term"
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)


func (m model) Init() tea.Cmd {
	// Start the animation loop with a delay
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return t
	})
}

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
    return 100 * c.Level 
    // level up hardness grows by 100 each time. 
}


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


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {


    // Handle level-up animation timing
    if m.IsLevelUpAnimating {
        if time.Since(m.LevelUpStartTime) >= 3*time.Second {
            // End the animation after 3 seconds
            m.Frames = asciiart.GetFrames() // Return to normal frames
            m.IsLevelUpAnimating = false    // End level-up animation
        }
        // Continue the animation loop even during the level-up animation
        m.CurrentFrame = (m.CurrentFrame + 1) % len(m.Frames)
        return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
            return t
        })
    }

    m.handleCatState()

    currentTime := time.Now()

    elapsed := currentTime.Sub(m.MyCat.LastFed)

    // Convert elapsed time to seconds
    seconds := elapsed.Seconds()

    // Smoothly decrease fullness by a proportional amount
    m.MyCat.Fullness -= FULLNESS_DECAY_RATE_PER_SECOND * seconds

    // Ensure fullness does not drop below 0
    if m.MyCat.Fullness < 0 {
        m.MyCat.Fullness = 0
    }

    // Reset the last update time
    m.MyCat.LastFed = currentTime


    switch msg := msg.(type) {
    case tea.KeyMsg:

        if msg.String() == "q" || msg.String() == "ctrl+c" {
            /// SAVING THIS EVIL JSON
            err := SaveGameData(m, "game_data.json")
            if err != nil {
                log.Fatalf("Error saving game data: %v", err)
            }

            return m, tea.Quit
        }

        // If the buy menu is open
        if m.ShowBuyMenu {
            switch msg.String() {
            case "left":
                m.SelectedFoodIndex = (m.SelectedFoodIndex - 1 + len(m.FoodInventory)) % len(m.FoodInventory) // Adjust to number of items
            case "right":
                m.SelectedFoodIndex = (m.SelectedFoodIndex + 1) % len(m.FoodInventory)
            case "up":
                m.SelectedFoodIndex = (m.SelectedFoodIndex - 5 + len(m.FoodInventory)) % len(m.FoodInventory) // Move up a row
            case "down":
                m.SelectedFoodIndex = (m.SelectedFoodIndex + 5) % len(m.FoodInventory) // Move down a row
            case "enter":
                selectedFood := &m.FoodInventory[m.SelectedFoodIndex]
                if m.MyCat.Coins >= selectedFood.Cost && selectedFood.Quantity < MAX_FOOD_INVENTORY_QUANTITY {
                    m.MyCat.Coins -= selectedFood.Cost // Deduct coins
                    selectedFood.Quantity++            // Increase quantity
                    m.ActionMessage = fmt.Sprintf("Bought %s for %d coins", selectedFood.Name, selectedFood.Cost)
                } else if selectedFood.Quantity >= MAX_FOOD_INVENTORY_QUANTITY {
                    m.ActionMessage = "Already stocked up!"
                } else {
                    m.ActionMessage = "Not enough coins!"
                }
            case "esc":
                m.ShowBuyMenu = false // Hide the buy food menu
                m.ShowFoodMenu = true // Return to the food menu
            }
            return m, nil // Early return to avoid further processing
        }

        // If the food menu is open
        if m.ShowFoodMenu {
            switch msg.String() {
            case "left":
                m.FocusedFoodButton = (m.FocusedFoodButton - 1 + len(m.FoodButtonLabels)) % len(m.FoodButtonLabels)
            case "right":
                m.FocusedFoodButton = (m.FocusedFoodButton + 1) % len(m.FoodButtonLabels)
            case "enter":
                if m.FoodButtonLabels[m.FocusedFoodButton] == "Feed" {
                    m.ShowFoodMenu = false  // Hide the food menu
                    m.ShowInventoryMenu = true // Show inventory menu
                    m.SelectedFoodIndex = 0 // Reset selected food index
                } else {
                    m.ActionMessage = "Opening buy food menu..."
                    m.ShowFoodMenu = false
                    m.ShowBuyMenu = true // Show buy menu
                }
            case "esc":
                m.ShowFoodMenu = false // Hide the food menu
                m.FocusedFoodButton = 0 // Reset focused food button
            }
            return m, nil // Early return to avoid further processing
        }

        // If the inventory menu is open
        if m.ShowInventoryMenu {
            switch msg.String() {
            case "left":
                m.SelectedFoodIndex = (m.SelectedFoodIndex - 1 + len(m.FoodInventory)) % len(m.FoodInventory) // Adjust to number of food items
            case "right":
                m.SelectedFoodIndex = (m.SelectedFoodIndex + 1) % len(m.FoodInventory)
            case "up":
                m.SelectedFoodIndex = (m.SelectedFoodIndex - 5 + len(m.FoodInventory)) % len(m.FoodInventory) // Move up a row
            case "down":
                m.SelectedFoodIndex = (m.SelectedFoodIndex + 5) % len(m.FoodInventory) // Move down a row
            case "enter":
                selectedFood := &m.FoodInventory[m.SelectedFoodIndex]
                if selectedFood.Quantity > 0 && m.MyCat.Fullness < 100 {
                    // Feed food to cat
                    selectedFood.Quantity-- 

                    m.ActionMessage = fmt.Sprintf("You fed the cat: %s", selectedFood.Name)

                    m.MyCat.GainXP(XP_FEEDING, &m)
                    m.MyCat.Fullness += float64(selectedFood.FeedingPower)
                    m.MyCat.Hunger += HUNGER_SATIATED_FACTOR
                    m.MyCat.Boredom += BOREDOM_FEEDING_FACTOR
                    m.MyCat.Health += HEALTH_FEEDING_REGEN_FACTOR
                    m.MyCat.Coins += COINS_FEEDING_FACTOR 
                    m.MyCat.LastFed = time.Now()

                    // Cap hunger and fullness at 100
                } else if selectedFood.Quantity <= 0 {
                    m.ActionMessage = "Out of stock!"
                } else if m.MyCat.Fullness >= 100 {
                    m.ActionMessage = "Termi is full, can't have any more!" 
                }
            case "esc":
                m.ShowInventoryMenu = false // Hide the inventory menu
                m.ShowFoodMenu = true // Return to the food menu
            }
            return m, nil // Early return to avoid further processing
        }

        // Handle main menu navigation if no submenus are open
        switch msg.String() {
        case "left":
            m.FocusedButton = (m.FocusedButton - 1 + len(m.ButtonLabels)) % len(m.ButtonLabels) // Navigate left
        case "right":
            m.FocusedButton = (m.FocusedButton + 1) % len(m.ButtonLabels) // Navigate right
        case "enter":
            if m.FocusedButton == 0 { // "Feed" button
                m.ShowFoodMenu = true // Show food menu
                m.FocusedFoodButton = 0 // Reset food button focus when entering
                m.FocusedButton = -1 // Reset focused button
            }
        }

    case time.Time:
        // Handle animation frame updates
        m.CurrentFrame = (m.CurrentFrame + 1) % len(m.Frames)
        // Continue the animation loop
        return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
            return t
        })
    }

    return m, nil
}

func (m model) View() string {
    // layout using Lipgloss
    header := titleStyle.Render("Termicat")
    
    // fix coins display in this render !!!
    coins := coinsStyle.Render(fmt.Sprintf("Coins: %d", m.MyCat.Coins))
    level := levelStyle.Render(fmt.Sprintf("Lv: %d", m.MyCat.Level))

    // Check if frames are loaded before accessing them
    var cat string
    if len(m.Frames) > 0 {
        cat = catStyle.Render(m.Frames[m.CurrentFrame])
    } else {
        cat = catStyle.Render("No animation frames loaded")
    }

    buttons := make([]string, len(m.ButtonLabels))
    for i, label := range m.ButtonLabels {
        // Check if no submenu is open, otherwise don't render focus on main buttons
        if i == m.FocusedButton && !m.ShowFoodMenu && !m.ShowInventoryMenu && !m.ShowBuyMenu {
            buttons[i] = focusedButtonStyle.Render(label)
        } else {
            buttons[i] = buttonStyle.Render(label)
        }
    }

    // Food Menu button
    var foodMenu string
    if m.ShowFoodMenu {
        foodButtons := make([]string, len(m.FoodButtonLabels))
        for i, label := range m.FoodButtonLabels {
            if i == m.FocusedFoodButton {
                foodButtons[i] = focusedButtonStyle.Render(label) // Highlight focused food button
            } else {
                foodButtons[i] = buttonStyle.Render(label)
            }
        }
        foodMenu = lipgloss.JoinHorizontal(lipgloss.Top, foodButtons...)
    }

    // Render inventory menu if it is visible
    var inventoryView string
    if m.ShowInventoryMenu {
        inventoryView += "Food Inventory:\n"
        for i := 0; i < 3; i++ { // 3 rows
            row := []string{}
            for j := 0; j < 5; j++ { // 5 columns
                index := i*5 + j
                if index < len(m.FoodInventory) {
                    food := m.FoodInventory[index]
                    // Display only the food name and quantity like "Apple (2)"
                    foodDetails := fmt.Sprintf("%s (%d)", food.Name, food.Quantity)
                    if index == m.SelectedFoodIndex {
                        row = append(row, fmt.Sprintf("[%s]", foodDetails)) // Highlight selected item
                    } else {
                        row = append(row, fmt.Sprintf(" %s ", foodDetails))
                    }
                }
            }
            inventoryView += strings.Join(row, " | ") + "\n"
        }
    } 

    var buyMenuView string
    if m.ShowBuyMenu {
        buyMenuView += "Buy Menu:\n"
        for i := 0; i < 3; i++ { // 3 rows
            row := []string{}
            for j := 0; j < 5; j++ { // 5 columns
                index := i*5 + j
                if index < len(m.FoodInventory) {
                    food := m.FoodInventory[index]
                    // Display food name, quantity, and cost like "Apple (2) - 10 coins"
                    foodDetails := fmt.Sprintf("%s (%d) - %d coins", food.Name, food.Quantity, food.Cost)
                    if index == m.SelectedFoodIndex {
                        row = append(row, fmt.Sprintf("[%s]", foodDetails)) // Highlight selected item
                    } else {
                        row = append(row, fmt.Sprintf(" %s ", foodDetails))
                    }
                }
            }
            buyMenuView += strings.Join(row, " | ") + "\n"
        }
    }

    // Display cat attributes for debugging purposes
    catAttributes := fmt.Sprintf(
        "Xp: %d\nCoins: %d\nName: %s\nBreed: %s\nAge: %d days\nWellness: %d\nFullness: %.0f\nHunger: %d\nDirtiness: %d\nHappiness: %d\nEnergy: %d\nHealth: %d\nBoredom: %d\nLast Fed: %s\nLast Cleaned: %s\n",
        
        m.MyCat.Xp, m.MyCat.Coins, m.MyCat.Name, m.MyCat.Breed, m.MyCat.Age, m.MyCat.Wellness, m.MyCat.Fullness, m.MyCat.Hunger,
        m.MyCat.Dirtiness, m.MyCat.Happiness, m.MyCat.Energy, m.MyCat.Health, m.MyCat.Boredom,
        m.MyCat.LastFed.Format(time.RFC822), m.MyCat.LastCleaned.Format(time.RFC822))

    catStatus := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("33")).Render(catAttributes)

    actionText := actionsStyle.Render("More actions if needed")
    actionMessage := actionMessageStyle.Render(m.ActionMessage)

    // Combine everything into the final layout
    layout := lipgloss.JoinVertical(
        lipgloss.Center,
        header,
        lipgloss.JoinHorizontal(lipgloss.Center, coins, level),
        cat,
        lipgloss.JoinHorizontal(lipgloss.Top, buttons...),
        actionText,
        actionMessage,
        foodMenu,        // Render the food menu if it is visible
        inventoryView,   // Render the food inventory if it is visible
        buyMenuView,     // Render the food buy inventory offers.
        catStatus,       // Render the cat's attributes for debugging
    )

    return layout
}

// Move cursor to the home position and clear the screen
func clearTerminal() {fmt.Print("\033[H\033[2J") }

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
	}
    
    return m
}

/********************************************************************************/
//                                    main 
/********************************************************************************/
func main() {
    clearTerminal()

    m := loadDefaultSettings()

    loaded_m_data, err := loadGameData("game_data.json")

    if err != nil {
        log.Fatalf("Error loading game data: %v", err)
    }

    m = loaded_m_data

    m.Frames = asciiart.GetFrames()
    // program 
	p := tea.NewProgram(m)
	p.Run() 
}
