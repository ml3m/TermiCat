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
        return m.handleLevelUpAnimation()
    }

    m.handleCatState()
    m.decayCatFullness()

    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" || msg.String() == "ctrl+c" {
            return m.handleQuit()
        }

        if m.ShowCleanMenu {
            return m.handleCleanMenu(msg)
        }

        if m.ShowBuyMenu {
            return m.handleBuyMenu(msg)
        }

        if m.ShowFoodMenu {
            return m.handleFoodMenu(msg)
        }

        if m.ShowInventoryMenu {
            return m.handleInventoryMenu(msg)
        }

        return m.handleKeyMsg(msg)

    case time.Time:
        m.CurrentFrame = (m.CurrentFrame + 1) % len(m.Frames)
        return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
            return t
        })
    }

    return m, nil
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
    // refactor this because is boilerplate code
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
