package model

import (
	"fmt"
	"math"

	"github.com/ilovelili/dongfeng/util"
)

// PhysiqueMasterType physique master type
type PhysiqueMasterType int

const (
	// AgeHeightWeightP 1
	AgeHeightWeightP PhysiqueMasterType = 1
	// AgeHeightWeightSD 2
	AgeHeightWeightSD PhysiqueMasterType = 2
	// BMI 3
	BMI PhysiqueMasterType = 3
	// HeightToWeightP 4
	HeightToWeightP PhysiqueMasterType = 4
	// HeightToWeightSD 5
	HeightToWeightSD PhysiqueMasterType = 5
)

// Physique physique entity
type Physique struct {
	BaseModel
	Pupil            *Pupil  `json:"pupil" csv:"-"`
	PupilID          uint    `json:"pupil_id" csv:"园儿ID"`
	Gender           string  `json:"gender" csv:"性别"`
	BirthDate        string  `json:"birth_date" csv:"出生日期"`
	ExamDate         string  `json:"exam_date" csv:"体检日期"`
	Age              string  `json:"age" csv:"-"`
	AgeComparison    float64 `json:"age_cmp" csv:"-"`
	Height           float64 `json:"height" csv:"身高"`
	HeightP          string  `json:"height_p" csv:"-"`  // Height P Zone
	HeightSD         string  `json:"height_sd" csv:"-"` // Height SD Zone
	Weight           float64 `json:"weight" csv:"体重"`
	WeightP          string  `json:"weight_p" csv:"-"`         // Weight P Zone
	WeightSD         string  `json:"weight_sd" csv:"-"`        // Weight SD Zone
	HeightToWeightP  string  `json:"height_weight_p" csv:"-"`  // Height to Weight P Zone
	HeightToWeightSD string  `json:"height_weight_sd" csv:"-"` // Height to Weight P Zone
	BMI              float64 `json:"bmi" csv:"-"`              // BMI for age > 5
	BMISD            string  `json:"bmi_sd" csv:"-"`           // BMI SD Zone
	FatCofficient    float64 `json:"fat_cofficient" csv:"-"`   // FatCofficient for age < 5
	Conclusion       string  `json:"conclusion" csv:"-"`
	CreatedBy        string  `json:"created_by" csv:"-"`
}

// AgeHeightWeightPMaster age height / weight p standard master
type AgeHeightWeightPMaster struct {
	BaseModel
	HeightOrWeight string  `json:"height_or_weight"`
	Gender         string  `json:"gender"`
	AgeMin         float64 `json:"age_min"`
	AgeMax         float64 `json:"age_max"`
	P3             float64 `json:"p3"`
	P10            float64 `json:"p10"`
	P20            float64 `json:"p20"`
	P50            float64 `json:"p50"`
	P80            float64 `json:"p80"`
	P97            float64 `json:"p97"`
}

// AgeHeightWeightSDMaster age height / weight sd standard master
type AgeHeightWeightSDMaster struct {
	BaseModel
	HeightOrWeight string  `json:"height_or_weight"`
	Gender         string  `json:"gender"`
	Age            string  `json:"age"` // x年x月
	SDM2           float64 `json:"sdm2"`
	SDM1           float64 `json:"sdm1"`
	Average        float64 `json:"avg"`
	SD1            float64 `json:"sd1"`
	SD2            float64 `json:"sd2"`
}

// HeightToWeightPMaster height to weight p master
type HeightToWeightPMaster struct {
	BaseModel
	Gender string  `json:"gender"`
	Height float64 `json:"height"`
	P3     float64 `json:"p3"`
	P10    float64 `json:"p10"`
	P20    float64 `json:"p20"`
	P50    float64 `json:"p50"`
	P80    float64 `json:"p80"`
	P97    float64 `json:"p97"`
}

// HeightToWeightSDMaster height to weight sd master
type HeightToWeightSDMaster struct {
	BaseModel
	Gender string  `json:"gender"`
	Height float64 `json:"height"`
	SDM3   float64 `json:"sdm3"`
	SDM2   float64 `json:"sdm2"`
	SDM1   float64 `json:"sdm1"`
	SD0    float64 `json:"sd0"`
	SD1    float64 `json:"sd1"`
	SD2    float64 `json:"sd2"`
	SD3    float64 `json:"sd3"`
}

// BMIMaster bmi standard master
type BMIMaster struct {
	BaseModel
	Gender  string  `json:"gender"`
	Age     string  `json:"age"` // x年x月
	Average float64 `json:"avg"`
	SD1     float64 `json:"sd1"`
	SD2     float64 `json:"sd2"`
	SD3     float64 `json:"sd3"`
}

// ResolveAge diff by birth date and exam date
func (p *Physique) ResolveAge() {
	birthdate, _ := util.StringToTime(p.BirthDate)
	examdate, _ := util.StringToTime(p.ExamDate)

	year, month, _, _, _, _ := util.Diff(birthdate, examdate)
	p.Age = fmt.Sprintf("%d岁%d月", year, month)
	cmp := float64(year) + float64(month)/12.0
	p.AgeComparison = math.Round(cmp*100) / 100
}

// ResolveBMI kg/m^2
func (p *Physique) ResolveBMI() {
	bmi := p.Weight / (p.Height / 100 * p.Height / 100)
	p.BMI = math.Round(bmi*100) / 100
}

// ResolveAgeHeightP get the corresponding p zone
func (p *Physique) ResolveAgeHeightP(pmasters []*AgeHeightWeightPMaster) (found bool) {
	found = true
	for _, m := range pmasters {
		if m.AgeMin <= p.AgeComparison && m.AgeMax >= p.AgeComparison && m.HeightOrWeight == "h" && m.Gender == p.Gender {
			if p.Height >= m.P97 {
				p.HeightP = ">P97"
				return
			}

			if p.Height >= m.P80 {
				p.HeightP = "P80~P97"
				return
			}

			if p.Height >= m.P50 {
				p.HeightP = "P50~P80"
				return
			}

			if p.Height >= m.P20 {
				p.HeightP = "P20~P50"
				return
			}

			if p.Height >= m.P10 {
				p.HeightP = "P10~P20"
				return
			}

			if p.Height >= m.P3 {
				p.HeightP = "P3~P10"
				return
			}

			p.HeightP = "<P3"
			return
		}
	}

	// not found
	found = false
	return
}

// ResolveAgeWeightP get the corresponding p zone
func (p *Physique) ResolveAgeWeightP(pmasters []*AgeHeightWeightPMaster) (found bool) {
	found = true
	for _, m := range pmasters {
		if m.AgeMin <= p.AgeComparison && m.AgeMax >= p.AgeComparison && m.HeightOrWeight == "w" && m.Gender == p.Gender {
			if p.Weight >= m.P97 {
				p.WeightP = ">P97"
				return
			}

			if p.Weight >= m.P80 {
				p.WeightP = "P80~P97"
				return
			}

			if p.Weight >= m.P50 {
				p.WeightP = "P50~P80"
				return
			}

			if p.Weight >= m.P20 {
				p.WeightP = "P20~P50"
				return
			}

			if p.Weight >= m.P10 {
				p.WeightP = "P10~P20"
				return
			}

			if p.Weight >= m.P3 {
				p.WeightP = "P3~P10"
				return
			}

			p.WeightP = "<P3"
			return
		}
	}

	found = false
	return
}

// ResolveAgeHeightSD get the corresponding sd zone
func (p *Physique) ResolveAgeHeightSD(sdmasters []*AgeHeightWeightSDMaster) (found bool) {
	found = true
	for _, sd := range sdmasters {
		if sd.Age == p.Age && sd.HeightOrWeight == "h" && sd.Gender == p.Gender {
			if p.Height >= sd.SD2 {
				p.HeightSD = ">2SD"
				return
			}

			if p.Height >= sd.SD1 {
				p.HeightSD = "1SD~2SD"
				return
			}

			if p.Height >= sd.Average {
				p.HeightSD = "AVG~1SD"
				return
			}

			if p.Height >= sd.SDM1 {
				p.HeightSD = "-1SD~AVG"
				return
			}

			if p.Height >= sd.SDM2 {
				p.HeightSD = "-2SD~-1SD"
				return
			}

			p.HeightSD = "<-2SD"
			return
		}
	}

	found = false
	return
}

// ResolveAgeWeightSD get the corresponding sd zone
func (p *Physique) ResolveAgeWeightSD(sdmasters []*AgeHeightWeightSDMaster) (found bool) {
	found = true
	for _, sd := range sdmasters {
		if sd.Age == p.Age && sd.HeightOrWeight == "w" && sd.Gender == p.Gender {
			if p.Weight >= sd.SD2 {
				p.WeightSD = ">2SD"
				return
			}

			if p.Weight >= sd.SD1 {
				p.WeightSD = "1SD~2SD"
				return
			}

			if p.Weight >= sd.Average {
				p.WeightSD = "AVG~1SD"
				return
			}

			if p.Weight >= sd.SDM1 {
				p.WeightSD = "-1SD~AVG"
				return
			}

			if p.Weight >= sd.SDM2 {
				p.WeightSD = "-2SD~-1SD"
				return
			}

			p.WeightSD = "<-2SD"
			return
		}
	}

	found = false
	return
}

// ResolveBMISD get the corresponding bmi sd zone
func (p *Physique) ResolveBMISD(bmimasters []*BMIMaster) (found bool) {
	found = true
	for _, bmi := range bmimasters {
		if bmi.Age == p.Age && bmi.Gender == p.Gender {
			if p.BMI >= bmi.SD3 {
				p.BMISD = ">3SD"
				return
			}

			if p.BMI >= bmi.SD2 {
				p.BMISD = "2SD~3SD"
				return
			}

			if p.BMI >= bmi.SD1 {
				p.BMISD = "1SD~2SD"
				return
			}

			if p.Weight >= bmi.Average {
				p.BMISD = "AVG~1SD"
				return
			}

			p.BMISD = "<AVG"
			return
		}
	}

	found = false
	return
}

// ResolveHeightToWeightP get the corresponding height to weight p zone
func (p *Physique) ResolveHeightToWeightP(hwpmasters []*HeightToWeightPMaster) (found bool) {
	found = true
	for _, hwp := range hwpmasters {
		if hwp.Gender == p.Gender && math.Abs(hwp.Height-p.Height) <= 0.5 /* 1 一个区间,选择最近点 */ {
			if p.Weight >= hwp.P97 {
				p.HeightToWeightP = ">P97"
				return
			}

			if p.Weight >= hwp.P80 {
				p.HeightToWeightP = "P80~P97"
				return
			}

			if p.Weight >= hwp.P50 {
				p.HeightToWeightP = "P50~P80"
				return
			}

			if p.Weight >= hwp.P20 {
				p.HeightToWeightP = "P20~P50"
				return
			}

			if p.Weight >= hwp.P10 {
				p.HeightToWeightP = "P10~P20"
				return
			}

			if p.Weight >= hwp.P3 {
				p.HeightToWeightP = "P3~P10"
				return
			}

			p.HeightToWeightP = "<P3"
			return
		}
	}

	found = false
	return
}

// ResolveHeightToWeightSD get the corresponding height to weight sd zone
func (p *Physique) ResolveHeightToWeightSD(hwsdmasters []*HeightToWeightSDMaster) (found bool) {
	found = true
	for _, hwsd := range hwsdmasters {
		if hwsd.Gender == p.Gender && math.Abs(hwsd.Height-p.Height) <= 0.25 /* 0.5 一个区间,选择最近点 */ {
			if p.Weight >= hwsd.SD3 {
				p.HeightToWeightSD = ">3SD"
				return
			}

			if p.Weight >= hwsd.SD2 {
				p.HeightToWeightSD = "2SD~3SD"
				return
			}

			if p.Weight >= hwsd.SD1 {
				p.HeightToWeightSD = "1SD~2SD"
				return
			}

			if p.Weight >= hwsd.SD0 {
				p.HeightToWeightSD = "0SD~1SD"
				return
			}

			if p.Weight >= hwsd.SDM1 {
				p.HeightToWeightSD = "-1SD~0SD"
				return
			}

			if p.Weight >= hwsd.SDM2 {
				p.HeightToWeightSD = "-2SD~-1SD"
				return
			}

			if p.Weight >= hwsd.SDM3 {
				p.HeightToWeightSD = "-3SD~-2SD"
				return
			}

			p.HeightToWeightSD = "<-3SD"
			return
		}
	}

	found = false
	return
}

// ResolveFatCofficient (体重（kg）-中位数)/中位数
func (p *Physique) ResolveFatCofficient(hwsdmasters []*HeightToWeightSDMaster) (found bool) {
	found = true
	for _, hwsd := range hwsdmasters {
		if hwsd.Gender == p.Gender && math.Abs(hwsd.Height-p.Height) <= 0.25 /* 0.5 一个区间,选择最近点 */ {
			median := hwsd.SD0
			fatcoff := 100 * (p.Weight - median) / median
			p.FatCofficient = math.Round(fatcoff*100) / 100
			return
		}
	}

	found = false
	return
}

// ResolveConclusion resolve conclusion
func (p *Physique) ResolveConclusion() {
	p.Conclusion = "正常"

	// 身高小于P3,疑似是生长迟缓，根据性别及年龄比对《5岁以下儿童低体重/生长迟缓标准表》。身高小于-2SD，为生长迟缓。
	if p.HeightP == "<P3" {
		if p.HeightSD == "<-2SD" {
			p.Conclusion = "生长迟缓"
			// do not return since conclusion can be overwritten
			// return
		} else if p.HeightSD == "Unknown" {
			p.Conclusion = "疑似生长迟缓(数据不足)"
		}
	}

	// 五项指标评价参考值核对出来后体重小于P10，按照性别及年龄比对《0-6岁按身高测体重》。如身高测体重也小于P10，为营养不良。
	if p.WeightP == "P3~P10" || p.WeightP == "<P3" {
		if p.HeightToWeightP == "P3~P10" || p.HeightToWeightP == "<P3" {
			p.Conclusion = "营养不良"
		} else if p.HeightToWeightP == "Unknown" {
			p.Conclusion = "疑似营养不良(数据不足)"
		}
	}

	// 年龄测身高和身高测体重两项都小于P3的为重度营养不良；一项小于P10，一项小于P3的或者两项都小于P10的为轻度营养不良
	if p.HeightP == "<P3" && p.HeightToWeightP == "<P3" {
		p.Conclusion = "重度营养不良"
	}

	if (p.HeightP == "<P3" && p.HeightToWeightP == "P3~P10") ||
		(p.HeightP == "P3~P10" && p.HeightToWeightP == "P3~P10") ||
		(p.HeightP == "P3~P10" && p.HeightToWeightP == "P3") {
		p.Conclusion = "轻度营养不良"
	}

	// 五项指标评价参考值核对出来后体重较重的幼儿
	if p.WeightP == "P50~P80" || p.WeightP == "P80~P97" || p.WeightP == ">P97" {
		// 五岁以下，按照性别及年龄，核对《5岁以下男/女童身高别体重标准》表，
		// 根据身高，大于+1SD为超重，大于+2SD为轻度肥胖，大于+3SD的为中重度肥胖。
		if p.AgeComparison < 5 {
			if p.HeightToWeightSD == "1SD~2SD" {
				p.Conclusion = "超重"
			} else if p.HeightToWeightSD == "2SD~3SD" {
				p.Conclusion = "轻度肥胖"
			} else if p.HeightToWeightSD == ">3SD" {
				p.Conclusion = "中重度肥胖"
			} else if p.HeightToWeightSD == "Unknown" {
				p.Conclusion = "疑似肥胖(数据不足)"
			}
		} else {
			// 大于5岁的幼儿，计算BMI指数【体重/身高(米)的平方】，
			// 然后比对《5-19岁BMI指数》表，按照性别与年龄比对BMI指数，大于+1SD为超重，大于+2SD为轻度肥胖，大于+3SD的为中重度肥胖。
			if p.BMISD == "1SD~2SD" {
				p.Conclusion = "超重"
			} else if p.BMISD == "2SD~3SD" {
				p.Conclusion = "轻度肥胖"
			} else if p.BMISD == ">3SD" {
				p.Conclusion = "中重度肥胖"
			} else if p.BMISD == "Unknown" {
				p.Conclusion = "疑似肥胖(数据不足)"
			}
		}
	}
}
