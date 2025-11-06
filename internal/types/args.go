package types

const (
	WetFlagNone uint8 = 0
	WetFlagVerboseTokenize uint8 = 0x1
	WetFlagVerboseRuntime uint8 = 0x2
	WetFlagVerboseCompile uint8 = 0x4
	WetFlagVerbose uint8 = 0x7
	WetFlagHelp uint8 = 0x8
	WetFlagDev uint8 = 0x10
)

type WetBin struct {
	Path			string
	Name			string
}

type WetArgs struct {
	Bin				WetBin
	Flags			uint8
	Path			*string
}
