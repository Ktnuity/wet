package types

type WetFlag uint8
const (
	WetFlagNone WetFlag = 0
	WetFlagVerboseTokenize WetFlag = 0x1
	WetFlagVerboseRuntime WetFlag = 0x2
	WetFlagVerboseCompile WetFlag = 0x4
	WetFlagVerbose WetFlag = 0x7
	WetFlagDev WetFlag = 0x8
	WetFlagHelp WetFlag = 0x10
	WetFlagVersion WetFlag = 0x20
	WetFlagLicense WetFlag = 0x40
)

func (flag WetFlag) Is(other WetFlag) bool {
	return flag & other == other
}

type WetBin struct {
	Path			string
	Name			string
}

type WetArgs struct {
	Bin				WetBin
	Flags			WetFlag
	Path			*string
}
