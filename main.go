package main

import (
	"TermiCat/asciiart"
	"fmt"
	"log"
	"time"
    "strings"
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)


func (m model) Init() tea.Cmd {
	// Start the animation loop with a delay
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
    m.decayCatFullness()

    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" || msg.String() == "ctrl+c" {
            return m.handleQuit()
        }

        if m.ShowCleanMenu {
            switch msg.String() {
            case "enter":
                // Perform clean action
                m.MyCat.Dirtiness = 0 // Reset dirtiness
                m.ActionMessage = "Your cat has been cleaned!"
                m.ShowCleanMenu = false // Hide the clean menu after action
            case "esc":
                m.ShowCleanMenu = false // Hide the clean menu
            }
            return m, nil // Early return to avoid further processing
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

        return m.handleKeyMsg(msg)

    case time.Time:
        // Handle animation frame updates
        m.CurrentFrame = (m.CurrentFrame + 1) % len(m.Frames)
        // Continue the animation loop
        return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
            return t})
        }

        return m, nil
    }

// decayCatFullness reduces the cat's fullness based on elapsed time.
func (m *model) decayCatFullness() {
    currentTime := time.Now()
    elapsed := currentTime.Sub(m.MyCat.LastFed)
    seconds := elapsed.Seconds()
    m.MyCat.Fullness -= FULLNESS_DECAY_RATE_PER_SECOND * seconds
    m.MyCat.LastFed = currentTime
}

// handleQuit handles the quit logic and saves the game data.
func (m *model) handleQuit() (tea.Model, tea.Cmd) {
    err := SaveGameData(m, "game_data.json")
    if err != nil {
        log.Fatalf("Error saving game data: %v", err)
    }
    return m, tea.Quit
}


// handleKeyMsg processes key messages for navigation and actions.
func (m *model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Handle main menu navigation if no submenus are open
    switch msg.String() {
    case "left":
        m.FocusedButton = (m.FocusedButton - 1 + len(m.ButtonLabels)) % len(m.ButtonLabels) // Navigate left
    case "right":
        m.FocusedButton = (m.FocusedButton + 1) % len(m.ButtonLabels) // Navigate right
    case "enter":
        return m.handleMenuAction()
    }
    return m, nil
}

// handleMenuAction performs actions based on the currently focused button.
func (m *model) handleMenuAction() (tea.Model, tea.Cmd) {
    switch m.FocusedButton {
    case 0:
        m.feedMenu()
    case 1:
        m.gameMenu()
    case 2:
        m.cleanMenu()
    case 3:
        m.sleepMenu()
    }
    return m, nil
}

func (m *model) feedMenu() {
    print("feed_main clicked")
    m.ShowFoodMenu = true // Show food menu
    m.FocusedFoodButton = 0 // Reset food button focus when entering
    m.FocusedButton = -1 // Reset focused button
}
func (m *model) gameMenu() {
    print("game_main clicked")
}
func (m *model) cleanMenu() {
    print("clean_main clicked")
    m.ShowCleanMenu = true // show clean menu
    m.FocusedButton = -1 // Reset focused button
}
func (m *model) sleepMenu() {
    print("sleep_main clicked")
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

    if m.ShowCleanMenu {
        print("showcleanmenu View() reached ")
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
