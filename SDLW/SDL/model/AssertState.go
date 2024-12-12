package m

type AssertState = Int

const (
	ASSERTION_RETRY         = AssertState(0)
	ASSERTION_BREAK         = AssertState(1)
	ASSERTION_ABORT         = AssertState(2)
	ASSERTION_IGNORE        = AssertState(3)
	ASSERTION_ALWAYS_IGNORE = AssertState(4)
)
