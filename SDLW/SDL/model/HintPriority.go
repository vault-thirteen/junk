package m

type HintPriority = Int

const (
	HINT_DEFAULT  = HintPriority(0)
	HINT_NORMAL   = HintPriority(1)
	HINT_OVERRIDE = HintPriority(2)
)
