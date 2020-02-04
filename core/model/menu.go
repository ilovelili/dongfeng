package model

import "strconv"

// BreakfastOrLunch breakfast or lunch emum
type BreakfastOrLunch uint

const (
	// AllMenuType all
	AllMenuType BreakfastOrLunch = 0
	// Breakfast breakfast
	Breakfast BreakfastOrLunch = 1
	// Lunch lunch
	Lunch BreakfastOrLunch = 2
	// Snack snack
	Snack BreakfastOrLunch = 3
)

// String to string
func (b BreakfastOrLunch) String() string {
	return strconv.Itoa(int(b))
}

// JuniorOrSenior juinor class menu or senior class menu
type JuniorOrSenior uint

const (
	// AllClass all
	AllClass JuniorOrSenior = 0
	// Junior junior class
	Junior JuniorOrSenior = 1
	// Senior senior class
	Senior JuniorOrSenior = 2
)

// String to string
func (j JuniorOrSenior) String() string {
	return strconv.Itoa(int(j))
}

// Menu entity
type Menu struct {
	BaseModel
	Date             string           `json:"date"`
	Recipe           *Recipe          `json:"recipe"`
	RecipeID         uint             `json:"recipe_id"`
	BreakfastOrLunch BreakfastOrLunch `json:"breakfast_or_lunch"`
	JuniorOrSenior   JuniorOrSenior   `json:"junior_or_senior"`
}

// Recipe entity
type Recipe struct {
	BaseModel
	Name              string           `json:"name" csv:"菜品名称"`
	Ingredients       []*Ingredient    `json:"ingredients" gorm:"many2many:recipe_ingredients" csv:"-"`
	CSVIngredient     string           `gorm:"-" json:"-" csv:"原料名称"`
	RecipeNutrition   *RecipeNutrition `json:"nutrition" csv:"-"`
	RecipeNutritionID *uint            `csv:"nutrition_id"`
}

// RecipeNutrition entity
type RecipeNutrition struct {
	BaseModel
	Recipe       string  `json:"recipe" csv:"recipe"`
	Carbohydrate float64 `json:"carbohydrate" csv:"carbohydrate"`
	Dietaryfiber float64 `json:"dietaryfiber" csv:"dietaryfiber"`
	Protein      float64 `json:"protein" csv:"protein"`
	Fat          float64 `json:"fat" csv:"fat"`
	Heat         float64 `json:"heat" csv:"heat"`
}

// Ingredient ingredient entity
type Ingredient struct {
	BaseModel
	Ingredient            string             `gorm:"unique_index" json:"ingredient" csv:"食材"`
	Alias                 string             `json:"alias" csv:"-"`
	IngredientCategory    IngredientCategory `json:"category"`
	IngredientCategoryID  uint               `json:"category_id"`
	CSVIngredientCategory string             `gorm:"-" json:"-" csv:"类别"` // for csv upload
	Protein100g           float64            `gorm:"column:protein_100g" json:"protein_100g" csv:"蛋白质(100g)"`
	ProteinDaily          float64            `json:"protein_daily" csv:"蛋白质(每人日)"`
	Fat100g               float64            `gorm:"column:fat_100g" json:"fat_100g" csv:"脂肪(100g)"`
	FatDaily              float64            `json:"fat_daily" csv:"脂肪(每人日)"`
	Carbohydrate100g      float64            `gorm:"column:carbohydrate_100g" json:"carbohydrate_100g" csv:"碳水化合物(100g)"`
	CarbohydrateDaily     float64            `json:"carbohydrate_daily" csv:"碳水化合物(每人日)"`
	Heat100g              float64            `gorm:"column:heat_100g" json:"heat_100g" csv:"热量(100g)"`
	HeatDaily             float64            `json:"heat_daily" csv:"热量(每人日)"`
	Calcium100g           float64            `gorm:"column:calcium_100g" json:"calcium_100g" csv:"钙(100g)"`
	CalciumDaily          float64            `json:"calcium_daily" csv:"钙(每人日)"`
	Iron100g              float64            `gorm:"column:iron_100g" json:"iron_100g" csv:"铁(100g)"`
	IronDaily             float64            `json:"iron_daily" csv:"铁(每人日)"`
	Zinc100g              float64            `gorm:"column:zinc_100g" json:"zinc_100g" csv:"锌(100g)"`
	ZincDaily             float64            `json:"zinc_daily" csv:"锌(每人日)"`
	VA100g                float64            `gorm:"column:va_100g" json:"va_100g" csv:"VA(100g)"`
	VADaily               float64            `json:"va_daily" csv:"VA(每人日)"`
	VB1100g               float64            `gorm:"column:vb1_100g" json:"vb1_100g" csv:"VB1(100g)"`
	VB1Daily              float64            `json:"vb1_daily" csv:"VB1(每人日)"`
	VB2100g               float64            `gorm:"column:vb2_100g" json:"vb2_100g" csv:"VB2(100g)"`
	VB2Daily              float64            `json:"vb2_daily" csv:"VB2(每人日)"`
	VC100g                float64            `gorm:"column:vc_100g" json:"vc_100g" csv:"VC(100g)"`
	VCDaily               float64            `json:"vc_daily" csv:"VC(每人日)"`
	CreatedBy             string             `json:"-" csv:"-"`
}

// IngredientCategory ingredient category
type IngredientCategory struct {
	BaseModel
	Category  string `json:"category"`
	CreatedBy string `json:"-"`
}
