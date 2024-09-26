package main

import (
	"TermiCat/asciiart" // my ascii
	"fmt"
	"log"
	"os"
	"time"
    "strings"
	"golang.org/x/term"
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

func splitLines(content string) []string {return strings.Split(content, "\n")}


func CenterEngine(content string) {
	// Get terminal size
	width, height, _ := term.GetSize(int(os.Stdout.Fd())) // _ is err ignored for now 

	// Split the content into lines
	lines := splitLines(content)
	contentHeight := len(lines)

	// Find the maximum width of the content
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Calculate center position
	x := (width - maxWidth) / 2
	y := (height - contentHeight) / 2

	// Print each line at the center position
	for i, line := range lines {
		fmt.Printf("\033[%d;%dH%s", y+i+1, x+1, line) // Add 1 because terminal positions are 1-based
	}
	fmt.Println("\033[0m") // Reset terminal color
}

// Define styles using Lipgloss
var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).MarginBottom(1)
	coinsStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Align(lipgloss.Center)
	levelStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Align(lipgloss.Center)
	catStyle          = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Align(lipgloss.Center).Height(8).Width(18)
	actionsStyle      = lipgloss.NewStyle().Align(lipgloss.Center).MarginTop(1)
	buttonStyle       = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.NormalBorder()).Align(lipgloss.Center)
	focusedButtonStyle = lipgloss.NewStyle().Padding(0, 2).Border(lipgloss.DoubleBorder()).Align(lipgloss.Center)
	centerStyle       = lipgloss.NewStyle().Align(lipgloss.Center).MarginBottom(2)
	actionMessageStyle = lipgloss.NewStyle().Align(lipgloss.Center).MarginTop(1)
)


func (m model) Init() tea.Cmd {
	// Start the animation loop with a delay
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle all possible outcomes of the poor cat
    if m.MyCat.Health <= 0 {
        m.Frames = asciiart.GetDeadCat() // Load the death frames if the cat is dead
    }

    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle quit (q or ctrl+c)
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
                    // Feed the cat
                    m.MyCat.Xp += XP_FEEDING // XP rewarded for feeding
                    m.MyCat.Fullness += selectedFood.FeedingPower
                    selectedFood.Quantity-- // Decrease the quantity of the food
                    m.MyCat.Hunger += HUNGER_SATIATED_FACTOR
                    m.ActionMessage = fmt.Sprintf("You fed the cat: %s", selectedFood.Name)
                    m.MyCat.Boredom += BOREDOM_FEEDING_FACTOR
                    m.MyCat.Health += HEALTH_FEEDING_REGEN_FACTOR
                    m.MyCat.Coins += COINS_FEEDING_FACTOR 

                    // Cap hunger and fullness at 100
                    if m.MyCat.Hunger > 100 {
                        m.MyCat.Hunger = 100
                    }
                    if m.MyCat.Fullness > 100 {
                        m.MyCat.Fullness = 100
                    }
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
    coins := coinsStyle.Render("Coins: 100")
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
        "Xp: %d\nCoins: %d\nName: %s\nBreed: %s\nAge: %d days\nWellness: %d\nFullness: %d\nHunger: %d\nDirtiness: %d\nHappiness: %d\nEnergy: %d\nHealth: %d\nBoredom: %d\nLast Fed: %s\nLast Cleaned: %s\n",
        
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

    // program 
	p := tea.NewProgram(m)
	p.Run() 
}
