package model

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
