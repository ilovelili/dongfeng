package controller

import (
	"fmt"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/core/repository"
)

var (
	pmasters    []*model.AgeHeightWeightPMaster
	sdmasters   []*model.AgeHeightWeightSDMaster
	hwpmasters  []*model.HeightToWeightPMaster
	hwsdmasters []*model.HeightToWeightSDMaster
	bmimasters  []*model.BMIMaster
)

// Physique controller
type Physique struct {
	physiqueRepo *repository.Physique
}

// NewPhysiqueController constructor
func NewPhysiqueController() *Physique {
	physique := &Physique{
		physiqueRepo: repository.NewPhysiqueRepository(),
	}

	physique.LoadMasters()
	return physique
}

// ResolvePhysique resolve physique based on the following steps:
// 1 测量出身高体重后，按照男女性别及年龄，分别比对《0-6岁男、女童体格发育五项指标评价参考值》。
// 2 身高小于P3,怀疑是生长迟缓，根据性别及年龄比对《5岁以下儿童低体重/生长迟缓标准表》。身高小于-2SD，为生长迟缓。
// 3 五项指标评价参考值核对出来后体重小于P10，按照性别及年龄比对《0-6岁按身高测体重》。
// 如身高测体重也小于P10，为营养不良。（年龄测身高和身高测体重两项都小于P3的为重度营养不良；一项小于P10，一项小于P3的或者两项都小于P10的为轻度营养不良。）
// 4 五项指标评价参考值核对出来后体重较重的幼儿
// 五岁以下，按照性别及年龄，核对《5岁以下男/女童身高别体重标准》表，
// 根据身高，大于+1SD为超重，大于+2SD为轻度肥胖，大于+3SD的为中重度肥胖。
// 大于5岁的幼儿，计算BMI指数【体重/身高(米)的平方】，然后比对《5-19岁BMI指数》表，
// 按照性别与年龄比对BMI指数，大于+1SD为超重，大于+2SD为轻度肥胖，大于+3SD的为中重度肥胖。
// 5 5岁以下超重或肥胖的幼儿在计算肥胖度时，根据《5岁以下男/女童身高别体重标准》表，对应相应的身高后，计算公式为实测体重（kg）-中位数/中位数。
func (c *Physique) ResolvePhysique(physique *model.Physique) (err error) {
	physique.ResolveAge()
	physique.ResolveBMI()

	if found := physique.ResolveAgeHeightP(pmasters); !found {
		err = fmt.Errorf("P height master data not found")
		return
	}

	if found := physique.ResolveAgeWeightP(pmasters); !found {
		err = fmt.Errorf("P weight master data not found")
		return
	}

	if found := physique.ResolveAgeHeightSD(sdmasters); !found {
		// if sd not found... then how about we set it as unknown
		physique.HeightSD = "Unknown"
	}

	if found := physique.ResolveAgeWeightSD(sdmasters); !found {
		// if sd not found... then how about we set it as unknown
		physique.WeightSD = "Unknown"
	}

	if found := physique.ResolveHeightToWeightP(hwpmasters); !found {
		// if hwp not found...
		physique.HeightToWeightP = "Unknown"
	}

	if found := physique.ResolveHeightToWeightSD(hwsdmasters); !found {
		// if hwp not found...
		physique.HeightToWeightSD = "Unknown"
	}

	if found := physique.ResolveFatCofficient(hwsdmasters); !found {
		// if sd not found... then how about we set it as unknown
		physique.FatCofficient = 0.0
	}

	if found := physique.ResolveBMISD(bmimasters); !found {
		// if bmisd not found...
		physique.BMISD = "Unknown"
	}

	physique.ResolveConclusion()
	return nil
}

// LoadMasters load master data
func (c *Physique) LoadMasters() (err error) {
	pmasters, err = c.physiqueRepo.SelectAgeHeightWeightPMasters()
	if err != nil {
		return
	}

	sdmasters, err = c.physiqueRepo.SelectAgeHeightWeightSDMasters()
	if err != nil {
		return
	}

	hwpmasters, err = c.physiqueRepo.SelectHeightToWeightPMasters()
	if err != nil {
		return
	}

	hwsdmasters, err = c.physiqueRepo.SelectHeightToWeightSDMasters()
	if err != nil {
		return
	}

	bmimasters, err = c.physiqueRepo.SelectBMIMasters()
	if err != nil {
		return
	}

	return
}
