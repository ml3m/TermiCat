package main

import (
	"TermiCat/asciiart"
	"fmt"
	"time"
    tea "github.com/charmbracelet/bubbletea"
)

func (m model) handleLevelUpAnimation() (model, tea.Cmd) {
    if time.Since(m.LevelUpStartTime) >= 3*time.Second {
        m.Frames = asciiart.GetFrames() // Return to normal frames
        m.IsLevelUpAnimating = false    // End level-up animation
    }
    m.CurrentFrame = (m.CurrentFrame + 1) % len(m.Frames)
    return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
        return t
    })
}

func (m model) handleBuyMenu(msg tea.KeyMsg) (model, tea.Cmd) {
    switch msg.String() {
    case "left":
        m.SelectedFoodIndex = (m.SelectedFoodIndex - 1 + len(m.FoodInventory)) % len(m.FoodInventory)
    case "right":
        m.SelectedFoodIndex = (m.SelectedFoodIndex + 1) % len(m.FoodInventory)
    case "up":
        m.SelectedFoodIndex = (m.SelectedFoodIndex - 5 + len(m.FoodInventory)) % len(m.FoodInventory)
    case "down":
        m.SelectedFoodIndex = (m.SelectedFoodIndex + 5) % len(m.FoodInventory)
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
        m.ShowBuyMenu = false
        m.ShowFoodMenu = true
    }
    return m, nil
}


func (m model) handleCleanMenu(msg tea.KeyMsg) (model, tea.Cmd) {
    switch msg.String() {
    case "enter":
        m.MyCat.Dirtiness = 0
        m.ActionMessage = "Your cat has been cleaned!"
        m.ShowCleanMenu = false
    case "esc":
        m.ShowCleanMenu = false
    }
    return m, nil
}


func (m model) handleFoodMenu(msg tea.KeyMsg) (model, tea.Cmd) {
    switch msg.String() {
    case "left":
        m.FocusedFoodButton = (m.FocusedFoodButton - 1 + len(m.FoodButtonLabels)) % len(m.FoodButtonLabels)
    case "right":
        m.FocusedFoodButton = (m.FocusedFoodButton + 1) % len(m.FoodButtonLabels)
    case "enter":
        if m.FoodButtonLabels[m.FocusedFoodButton] == "Feed" {
            m.ShowFoodMenu = false
            m.ShowInventoryMenu = true
            m.SelectedFoodIndex = 0
        } else {
            m.ActionMessage = "Opening buy food menu..."
            m.ShowFoodMenu = false
            m.ShowBuyMenu = true
        }
    case "esc":
        m.ShowFoodMenu = false
        m.FocusedFoodButton = 0
    }
    return m, nil
}

func (m model) handleInventoryMenu(msg tea.KeyMsg) (model, tea.Cmd) {
    switch msg.String() {
    case "left":
        m.SelectedFoodIndex = (m.SelectedFoodIndex - 1 + len(m.FoodInventory)) % len(m.FoodInventory)
    case "right":
        m.SelectedFoodIndex = (m.SelectedFoodIndex + 1) % len(m.FoodInventory)
    case "up":
        m.SelectedFoodIndex = (m.SelectedFoodIndex - 5 + len(m.FoodInventory)) % len(m.FoodInventory)
    case "down":
        m.SelectedFoodIndex = (m.SelectedFoodIndex + 5) % len(m.FoodInventory)
    case "enter":
        selectedFood := &m.FoodInventory[m.SelectedFoodIndex]
        if selectedFood.Quantity > 0 && m.MyCat.Fullness <= 99 {
            selectedFood.Quantity--
            m.ActionMessage = fmt.Sprintf("You fed the cat: %s", selectedFood.Name)
            m.MyCat.GainXP(XP_FEEDING, &m)
            m.MyCat.Fullness += float64(selectedFood.FeedingPower)
            m.MyCat.Hunger += HUNGER_SATIATED_FACTOR
            m.MyCat.Boredom += BOREDOM_FEEDING_FACTOR
            m.MyCat.Health += HEALTH_FEEDING_REGEN_FACTOR
            m.MyCat.Coins += COINS_FEEDING_FACTOR
            m.MyCat.LastFed = time.Now()
        } else if selectedFood.Quantity <= 0 {
            m.ActionMessage = "Out of stock!"
        } else if m.MyCat.Fullness >= 100 {
            m.ActionMessage = "Termi is full, can't have any more!"
        }
    case "esc":
        m.ShowInventoryMenu = false
        m.ShowFoodMenu = true
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

