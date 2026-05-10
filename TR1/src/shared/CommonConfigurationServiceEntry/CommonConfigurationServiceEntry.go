package ccse

import (
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

type CommonConfigurationServiceEntry struct {
	Type       string
	Protocol   string
	Parameters []ccp.CommonConfigurationParameter
}

func (se *CommonConfigurationServiceEntry) getParameter(parameterName string) *ccp.CommonConfigurationParameter {
	for _, p := range se.Parameters {
		if p.Name == parameterName {
			return &p
		}
	}

	return nil
}

func (se *CommonConfigurationServiceEntry) GetParameterAsString(parameterName string) string {
	return getParameterValue[string](se.getParameter(parameterName))
}
func (se *CommonConfigurationServiceEntry) GetParameterAsStrings(parameterName string) []string {
	return getParameterValue[[]string](se.getParameter(parameterName))
}
func (se *CommonConfigurationServiceEntry) GetParameterAsInt(parameterName string) int {
	return getParameterValue[int](se.getParameter(parameterName))
}
func (se *CommonConfigurationServiceEntry) GetParameterAsInts(parameterName string) []int {
	return getParameterValue[[]int](se.getParameter(parameterName))
}
func (se *CommonConfigurationServiceEntry) GetParameterAsBool(parameterName string) bool {
	return getParameterValue[bool](se.getParameter(parameterName))
}
func (se *CommonConfigurationServiceEntry) GetParameterAsMap(parameterName string) map[string]string {
	return getParameterValue[map[string]string](se.getParameter(parameterName))
}

func getParameterValue[T any](param *ccp.CommonConfigurationParameter) T {
	return (param.Value).(T)
}
