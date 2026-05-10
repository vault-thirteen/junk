package ul

import (
	"database/sql/driver"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_New(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList

	// Test #1.
	ul = New()
	aTest.MustBeDifferent(ul, (*UidList)(nil))
}

func Test_NewFromArray(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var x UidList
	var err error

	// Test #1.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	x = nil
	aTest.MustBeEqual(ul, &x)

	// Test #2.
	ul, err = NewFromArray([]int{})
	aTest.MustBeNoError(err)
	x = []int{}
	aTest.MustBeEqual(ul, &x)

	// Test #3.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	x = []int{1, 2, 3}
	aTest.MustBeEqual(ul, &x)

	// Test #4.
	ul, err = NewFromArray([]int{1, 2, 2, 3})
	aTest.MustBeAnError(err)
}

func Test_AsArray(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var x []int

	// Test #1. Case I.
	ul = nil
	aTest.MustBeEqual(ul.AsArray(), []int{})

	// Test #2. Case II.
	x = nil
	ul = (*UidList)(&x)
	aTest.MustBeEqual(ul.AsArray(), []int{})

	// Test #3. Case III.
	x = []int{}
	ul = (*UidList)(&x)
	aTest.MustBeEqual(ul.AsArray(), []int{})

	// Test #4. Case X. All Clear.
	x = []int{1, 2, 3}
	ul = (*UidList)(&x)
	aTest.MustBeEqual(ul.AsArray(), []int{1, 2, 3})
}

func Test_CheckIntegrity(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var x UidList
	var err error

	// Test #1. Golang is a garbage language.
	ul = nil
	err = ul.CheckIntegrity()
	aTest.MustBeAnError(err)

	// Test #2. Unique.
	x = []int{1, 2, 3}
	ul = &x
	err = ul.CheckIntegrity()
	aTest.MustBeNoError(err)

	// Test #3. Not unique.
	x = []int{1, 2, 3, 2}
	ul = &x
	err = ul.CheckIntegrity()
	aTest.MustBeAnError(err)
}

func Test_isUnique(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var x UidList

	// Test #1. Unique.
	x = []int{1, 2, 3}
	ul = &x
	aTest.MustBeEqual(ul.isUnique(), true)

	// Test #2. Not unique.
	x = []int{1, 2, 3, 2}
	ul = &x
	aTest.MustBeEqual(ul.isUnique(), false)
}

func Test_Size(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1. Null object.
	ul = nil
	aTest.MustBeEqual(ul.Size(), 0)

	// Test #2. Empty array (slice).
	ul, err = NewFromArray([]int{})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.Size(), 0)

	// Test #3. Null array (slice).
	var x []int = nil
	ul, err = NewFromArray(x)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.Size(), 0)

	// Test #4. Normal array (slice).
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.Size(), 3)
}

func Test_HasItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1. Empty list.
	ul = New()
	aTest.MustBeEqual(ul.HasItem(1), false)

	// Test #2. Item is not found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.HasItem(4), false)

	// Test #3. Item is found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.HasItem(3), true)

	// Test #4. Null object.
	ul = nil
	aTest.MustBeEqual(ul.HasItem(3), false)
}

func Test_AddItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1.
	ul = New()
	err = ul.AddItem(1, false)
	aTest.MustBeNoError(err)
	err = ul.AddItem(2, false)
	aTest.MustBeNoError(err)
	err = ul.AddItem(3, false)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual([]int(*ul), []int{1, 2, 3})
	err = ul.AddItem(2, false)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual([]int(*ul), []int{1, 2, 3})

	// Test #2.
	ul = New()
	err = ul.AddItem(1, false)
	aTest.MustBeNoError(err)
	err = ul.AddItem(2, false)
	aTest.MustBeNoError(err)
	err = ul.AddItem(3, true)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual([]int(*ul), []int{3, 1, 2})
}

func Test_prependItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList

	// Test.
	ul = New()
	ul.prependItem(1)
	aTest.MustBeEqual([]int(*ul), []int{1})
	ul.prependItem(2)
	aTest.MustBeEqual([]int(*ul), []int{2, 1})
	ul.prependItem(3)
	aTest.MustBeEqual([]int(*ul), []int{3, 2, 1})
}

func Test_appendItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList

	// Test.
	ul = New()
	ul.appendItem(1)
	aTest.MustBeEqual([]int(*ul), []int{1})
	ul.appendItem(2)
	aTest.MustBeEqual([]int(*ul), []int{1, 2})
	ul.appendItem(3)
	aTest.MustBeEqual([]int(*ul), []int{1, 2, 3})
}

func Test_SearchForItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var idx int
	var err error

	// Test #1.
	ul = New()
	idx, err = ul.SearchForItem(1)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(idx, IndexOnError)

	// Test #2. Non-existent item.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	idx, err = ul.SearchForItem(4)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(idx, IndexOnError)

	// Test #3. Item is found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	idx, err = ul.SearchForItem(2)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(idx, 2-1)
}

func Test_RemoveItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1. Empty list.
	ul = New()
	err = ul.RemoveItem(1)
	aTest.MustBeAnError(err)

	// Test #2. Item is not found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.RemoveItem(4)
	aTest.MustBeAnError(err)

	// Test #3. Item is found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.RemoveItem(2)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(*ul, UidList([]int{1, 3}))
}

func Test_RemoveItemAtPos(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1. Position is negative.
	ul = New()
	err = ul.RemoveItemAtPos(-1)
	aTest.MustBeAnError(err)

	// Test #2. Position is too far away.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.RemoveItemAtPos(4)
	aTest.MustBeAnError(err)

	// Test #3. Existing position.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.RemoveItemAtPos(1)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(*ul, UidList([]int{1, 3}))
}

func Test_removeItemAtPos(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error
	var lastIndex int

	// Test #1. Single item.
	ul, err = NewFromArray([]int{1})
	aTest.MustBeNoError(err)
	lastIndex = len(*ul) - 1
	ul.removeItemAtPos(0, lastIndex)
	aTest.MustBeEqual(*ul, UidList([]int{}))

	// Test #2. First item.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	lastIndex = len(*ul) - 1
	ul.removeItemAtPos(0, lastIndex)
	aTest.MustBeEqual(*ul, UidList([]int{2, 3}))

	// Test #3. Middle item.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	lastIndex = len(*ul) - 1
	ul.removeItemAtPos(1, lastIndex)
	aTest.MustBeEqual(*ul, UidList([]int{1, 3}))

	// Test #4. Last item.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	lastIndex = len(*ul) - 1
	ul.removeItemAtPos(2, lastIndex)
	aTest.MustBeEqual(*ul, UidList([]int{1, 2}))
}

func Test_removeLastItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error
	var lastIndex int

	// Test #1. Single item.
	ul, err = NewFromArray([]int{1})
	aTest.MustBeNoError(err)
	lastIndex = len(*ul) - 1
	ul.removeLastItem(lastIndex)
	aTest.MustBeEqual(*ul, UidList([]int{}))

	// Test #2. Several items.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	lastIndex = len(*ul) - 1
	ul.removeLastItem(lastIndex)
	aTest.MustBeEqual(*ul, UidList([]int{1, 2}))
}

func Test_RaiseItem(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var isAlreadyRaised bool
	var err error

	// Test #1. Item is not found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	isAlreadyRaised, err = ul.RaiseItem(4)
	aTest.MustBeAnError(err)

	// Test #2. No moving is needed.
	ul = &UidList{1, 2, 3}
	isAlreadyRaised, err = ul.RaiseItem(1)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(isAlreadyRaised, true)
	aTest.MustBeEqual([]int(*ul), []int{1, 2, 3})

	// Test #4. Middle item is moved.
	ul = &UidList{1, 2, 3}
	isAlreadyRaised, err = ul.RaiseItem(2)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(isAlreadyRaised, false)
	aTest.MustBeEqual([]int(*ul), []int{2, 1, 3})

	// Test #5. Last item is moved.
	ul = &UidList{1, 2, 3}
	isAlreadyRaised, err = ul.RaiseItem(3)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(isAlreadyRaised, false)
	aTest.MustBeEqual([]int(*ul), []int{3, 1, 2})

	// Test #6. Single item.
	ul = &UidList{1}
	isAlreadyRaised, err = ul.RaiseItem(1)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(isAlreadyRaised, true)
	aTest.MustBeEqual([]int(*ul), []int{1})
}

func Test_MoveItemUp(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1. Item is not found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.MoveItemUp(4)
	aTest.MustBeAnError(err)

	// Test #2. Item is already on top edge.
	ul = &UidList{1, 2, 3}
	err = ul.MoveItemUp(1)
	aTest.MustBeAnError(err)

	// Test #3. Middle item is moved.
	ul = &UidList{1, 2, 3}
	err = ul.MoveItemUp(2)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual([]int(*ul), []int{2, 1, 3})

	// Test #4. Last item is moved.
	ul = &UidList{1, 2, 3}
	err = ul.MoveItemUp(3)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual([]int(*ul), []int{1, 3, 2})

	// Test #5. Single item.
	ul = &UidList{1}
	err = ul.MoveItemUp(1)
	aTest.MustBeAnError(err)
}

func Test_MoveItemDown(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1. Item is not found.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	err = ul.MoveItemDown(4)
	aTest.MustBeAnError(err)

	// Test #2. Item is already on bottom edge.
	ul = &UidList{1, 2, 3}
	err = ul.MoveItemDown(3)
	aTest.MustBeAnError(err)

	// Test #3. Middle item is moved.
	ul = &UidList{1, 2, 3}
	err = ul.MoveItemDown(2)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual([]int(*ul), []int{1, 3, 2})

	// Test #4. Top item is moved.
	ul = &UidList{1, 2, 3}
	err = ul.MoveItemDown(1)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual([]int(*ul), []int{2, 1, 3})

	// Test #5. Single item.
	ul = &UidList{1}
	err = ul.MoveItemDown(1)
	aTest.MustBeAnError(err)
}

func Test_Scan(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error

	// Test #1. Null destination.
	ul = nil
	err = ul.Scan(false)
	aTest.MustBeAnError(err)

	// Test #2. Null source.
	ul = New()
	err = ul.Scan(nil)
	aTest.MustBeNoError(err)

	// Test #3. Data type is string.
	ul = New()
	err = ul.Scan("string")
	aTest.MustBeAnError(err)

	// Test #4. Data type is byte array.
	ul = New()
	err = ul.Scan([]byte("[1,2,3]"))
	aTest.MustBeNoError(err)
	tmp := UidList([]int{1, 2, 3})
	aTest.MustBeEqual(ul, &tmp)
	aTest.MustBeEqual(ul.Size(), 3)

	// Test #5.
	ul = New()
	err = ul.Scan([]byte("[1,2,3,NaN]"))
	aTest.MustBeAnError(err)

	// Test #6.
	ul = New()
	err = ul.Scan([]byte(""))
	aTest.MustBeAnError(err)

	// Test #7.
	ul = New()
	err = ul.Scan([]byte("[]"))
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.Size(), 0)
}

func Test_Value(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var dv driver.Value
	var err error

	// Test #1.
	ul = nil
	dv, err = ul.Value()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(dv, nil)

	// Test #2.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	dv, err = ul.Value()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(dv, driver.Value([]byte("[1,2,3]")))
}

func Test_ValuesString(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error
	var vs string

	// Test #1.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "")

	// Test #2.
	ul, err = NewFromArray([]int{})
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "")

	// Test #3.
	ul, err = NewFromArray([]int{1})
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "1")

	// Test #4.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "1,2,3")

	// Test #5.
	ul = nil
	vs, err = ul.ValuesString()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(vs, "")
}

func Test_OnPage(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var ulx UidList
	var ulxx *UidList
	var err error
	var nullList *UidList = nil

	// Test #1.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(0, 1), nullList)

	// Test #2.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(1, 1), nullList)

	// Test #3.
	ul, err = NewFromArray([]int{})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.OnPage(1, 1), nullList)

	// Test #4.
	ul, err = NewFromArray([]int{1, 2, 3, 4, 5})
	aTest.MustBeNoError(err)
	ulx = []int{1, 2, 3, 4, 5}
	aTest.MustBeEqual(ul.OnPage(1, 5), &ulx)

	// Test #5.
	ul, err = NewFromArray([]int{1, 2, 3})
	aTest.MustBeNoError(err)
	ulx = []int{1, 2, 3}
	aTest.MustBeEqual(ul.OnPage(1, 5), &ulx)

	// Test #6.
	ul, err = NewFromArray([]int{1, 2, 3, 4, 5, 6, 7})
	aTest.MustBeNoError(err)
	ulx = []int{1, 2, 3, 4, 5}
	aTest.MustBeEqual(ul.OnPage(1, 5), &ulx)

	// Test #7.
	ul, err = NewFromArray([]int{1, 2, 3, 4, 5})
	aTest.MustBeNoError(err)
	ulxx = nil
	aTest.MustBeEqual(ul.OnPage(2, 5), ulxx)

	// Test #8.
	ul, err = NewFromArray([]int{1, 2, 3, 4, 5, 6})
	aTest.MustBeNoError(err)
	ulx = []int{6}
	aTest.MustBeEqual(ul.OnPage(2, 5), &ulx)

	// Test #9.
	ul, err = NewFromArray([]int{1, 2, 3, 4, 5, 6, 7})
	aTest.MustBeNoError(err)
	ulx = []int{6, 7}
	aTest.MustBeEqual(ul.OnPage(2, 5), &ulx)

	// Test #10.
	ul, err = NewFromArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	aTest.MustBeNoError(err)
	ulx = []int{6, 7, 8, 9, 10}
	aTest.MustBeEqual(ul.OnPage(2, 5), &ulx)

	// Test #11.
	ul = nil
	ulxx = nil
	aTest.MustBeEqual(ul.OnPage(1, 1), ulxx)
}

func Test_LastElement(t *testing.T) {
	aTest := tester.New(t)
	var ul *UidList
	var err error
	var nullElement *int = nil

	// Test #1.
	ul, err = NewFromArray(nil)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.LastElement(), nullElement)

	// Test #2.
	ul, err = NewFromArray([]int{})
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ul.LastElement(), nullElement)

	// Test #3.
	ul, err = NewFromArray([]int{1, 2, 3, 4, 5})
	aTest.MustBeNoError(err)
	aTest.MustBeDifferent(ul.LastElement(), nullElement)
	aTest.MustBeEqual(*ul.LastElement(), 5)

	// Test #4.
	ul, err = NewFromArray([]int{1})
	aTest.MustBeNoError(err)
	aTest.MustBeDifferent(ul.LastElement(), nullElement)
	aTest.MustBeEqual(*ul.LastElement(), 1)
}
