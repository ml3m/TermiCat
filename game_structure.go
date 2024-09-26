package main
import ("time")

const (
    // Core game settings.
    XP_FEEDING = 15
    XP_LEVELUP = 100
    HUNGER_SATIATED_FACTOR = 20
    BOREDOM_FEEDING_FACTOR = 2
    HEALTH_FEEDING_REGEN_FACTOR = 5
    COINS_FEEDING_FACTOR = 30
    MAX_FOOD_INVENTORY_QUANTITY = 10
)

type cat struct {
    Wellness      int         `json:"wellness"`
    Fullness      int         `json:"fullness"`
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
}
