package html_preprocessor

type MacroWrapperResult struct {
	Before []byte
	After  []byte
	Call   []byte
}

type MacroWrapper = func(macroName string, uniqueName string, attributes map[string]any) (*MacroWrapperResult, error)
