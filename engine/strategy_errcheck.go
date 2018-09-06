package engine

import (
	"strings"
	"strconv"
	"github.com/cyhon/goreporter/linters/errorcheck"
	"github.com/cyhon/goreporter/utils"
)

type StrategyErrorCheck struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyErrorCheck) GetName() string {
	return "ErrorCheck"
}

func (s *StrategyErrorCheck) GetDescription() string {
	return "Query all duplicate code in the project and give duplicate code locations and rows."
}

func (s *StrategyErrorCheck) GetWeight() float64 {
	return 0.05
}

// linterCopy provides a function that scans all duplicate code in the project and give
// duplicate code locations and rows.It will extract from the linter need to convert the
// data.The result will be saved in the r's attributes.
func (s *StrategyErrorCheck) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	errCodeList := errorcheck.ErrorCheck(parameters.ProjectPath)
	sumProcessNumber := int64(7)
	processUnit := utils.GetProcessUnit(sumProcessNumber, len(errCodeList))

	for _, errTip := range errCodeList {
		simpleTips := strings.Split(errTip, ":")
		if len(simpleTips) == 4 {
			packageName := utils.PackageNameFromGoPath(simpleTips[0])
			line, _ := strconv.Atoi(simpleTips[1])
			erroru := Error{
				LineNumber:  line,
				ErrorString: utils.AbsPath(simpleTips[0]) + ":" + strings.Join(simpleTips[1:], ":"),
			}
			summaries.Lock()
			if summarie, ok := summaries.Summaries[packageName]; ok {
				summarie.Errors = append(summarie.Errors, erroru)
				summaries.Summaries[packageName] = summarie
			} else {
				summarie := Summary{
					Name:   packageName,
					Errors: make([]Error, 0),
				}
				summarie.Errors = append(summarie.Errors, erroru)
				summaries.Summaries[packageName] = summarie
			}
			summaries.Unlock()
		}
		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}

	return
}

func (s *StrategyErrorCheck) Percentage(summaries *Summaries) float64 {
	summaries.RLock()
	defer summaries.RUnlock()
	return utils.CountPercentage(len(summaries.Summaries))
}
