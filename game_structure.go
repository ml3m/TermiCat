package main
import (
    "time"
    "github.com/charmbracelet/lipgloss"
)


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

const (
    // Core game settings.
    XP_FEEDING = 15
    XP_LEVELUP = 100
    HUNGER_SATIATED_FACTOR = 20
    BOREDOM_FEEDING_FACTOR = 2
    HEALTH_FEEDING_REGEN_FACTOR = 5
    COINS_FEEDING_FACTOR = 30
    MAX_FOOD_INVENTORY_QUANTITY = 10
    FULLNESS_DECAY_RATE_PER_SECOND = 100.0 / 600.0
)

type cat struct {
    Wellness      int         `json:"wellness"`
    Fullness      float64     `json:"fullness"`
    Hunger        int         `json:"hunger"`
    Dirtiness     int         `json:"dirtiness"`
    Happiness     int         `json:"happiness"`
    Energy        int         `json:"energy"`
    Health        int         `json:"health"`
    Xp            int         `json:"xp"`
    Boredom       int         `json:"boredom"`
    Age           int         `json:"age"`
    Coins         int         `json:"coins"`
    Level         int         `json:"level"`
    Breed         string      `json:"breed"`
    Name          string      `json:"name_cat"`  
    LastFed       time.Time   `json:"lastFed"`
    LastCleaned   time.Time   `json:"lastCleaned"`
}

type food struct {
    FeedingPower int          `json:"FeedingPower"`
    Cost         int          `json:"cost"`
    Quantity     int          `json:"quantity"`
    Name         string       `json:"name_food"` 
}

type model struct {

    MyCat cat                          `json:"cat"`
    CurrentFrame          int          `json:"currentFrame"`
    FocusedButton         int          `json:"focusedButton"`
    FocusedFoodButton     int          `json:"focusedFoodButton"`
    SelectedFoodIndex     int          `json:"selectedFoodIndex"`
    ShowFoodMenu          bool         `json:"showFoodMenu"`
    ShowBuyMenu           bool         `json:"showBuyMenu"`
    ShowInventoryMenu     bool         `json:"showInventoryMenu"`
    ActionMessage         string       `json:"actionMessage"`
    Frames              []string       `json:"frames"`
    ButtonLabels        []string       `json:"buttonLabels"`
    FoodButtonLabels    []string       `json:"foodButtonLabels"`
    FoodInventory       []food         `json:"foodInventory"`

    //level
    IsLevelUpAnimating    bool         `json:"isLevelUpAnimating"`
    LevelUpStartTime      time.Time

    ShowCleanMenu         bool         `json:"showCleanMenu"`
}
