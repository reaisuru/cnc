package archs

type Arch int

const (
	EmNone  Arch = 0
	EmM32        = 1
	EmSparc      = 2
	Em386        = 3
	Em68k        = 4
	Em88k        = 5
	Em486        = 6
	Em860        = 7
	EmMips       = 8

	EmMipsRs3Le = 10
	EmMipsRs4Be = 10

	EmParisc      Arch = 15
	EmSparc32Plus      = 18
	EmPowerPC          = 20
	EmPowerPC64        = 21
	EmSPU              = 23
	EmARM              = 40
	EmSH               = 42
	EmSPARCv9          = 43
	Emx8664            = 62
	EmOpenRisc         = 92
	EmAArch64          = 183
	EmArc              = 93
	EmCSky             = 39
)

func (c Arch) String() string {
	var mappings = map[Arch]string{
		EmNone: "none",

		EmSparc:     "sparc",
		Em68k:       "m68k",
		EmARM:       "arm",
		EmMips:      "mips",
		EmMipsRs3Le: "mips",
		EmPowerPC:   "powerpc",
		EmPowerPC64: "powerpc",
		EmSH:        "superh",
		EmAArch64:   "aarch64",

		Em386:   "i386",
		Em486:   "i486",
		Em860:   "i860",
		Emx8664: "x86_64",

		// exotic arch
		EmM32:         "M32",
		Em88k:         "88K",
		EmParisc:      "pa-risc",
		EmSparc32Plus: "spc32+",
		EmSPU:         "spu",
		EmSPARCv9:     "spcv9",
		EmOpenRisc:    "openrisc",
		EmArc:         "arc",
		EmCSky:        "csky",
	}

	archName, ok := mappings[c]
	if !ok {
		return "unknown"
	}

	return archName
}
