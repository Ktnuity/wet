package types

type WetFlag uint8
const (
	WetFlagNone WetFlag = 0
	WetFlagVerboseTokenize WetFlag = 0x1
	WetFlagVerboseRuntime WetFlag = 0x2
	WetFlagVerboseTypeCheck WetFlag = 0x4
	WetFlagVerboseCompile WetFlag = 0x8
	WetFlagVerbose WetFlag = 0xf
	WetFlagDev WetFlag = 0x10
	WetFlagHelp WetFlag = 0x20
	WetFlagVersion WetFlag = 0x40
	WetFlagLicense WetFlag = 0x80
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
