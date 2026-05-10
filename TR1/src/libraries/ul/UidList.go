package ul

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	Err_DestinationIsNotInitialised = "destination is not initialised"
	Err_ItemsAreNotUnique           = "items are not unique"
	Err_EdgePosition                = "edge position"
	Err_Position                    = "position error"

	ErrF_UidIsNotFound       = "uid is not found: %v"
	ErrF_DuplicateUid        = "duplicate uid: %v"
	ErrF_UnsupportedDataType = "unsupported data type: %s"
)

const (
	IndexOnError      = -1
	ListItemSeparator = ","
	StringOnError     = ""
)

// UidList is a list unique identifiers.
//
// The main purpose of this list is to store a chronological order of all added
// identifiers. The order of items in the list is important and, thus, the list
// may not be sorted. New items are added to the end of the list, deleted items
// shift existing items. All operations on the list assume that the list is
// unique before the operation, thus, every operation must ensure that its
// results do not break the uniqueness of items in the list.
type UidList []int

func New() (ul *UidList) {
	return new(UidList)
}

func NewFromArray(uids []int) (ul *UidList, err error) {
	tmp := UidList(uids)
	ul = &tmp

	err = ul.CheckIntegrity()
	if err != nil {
		return nil, err
	}

	return ul, nil
}

// AsArray returns the list as an array.
func (ul *UidList) AsArray() (arr []int) {
	if ul == nil {
		return []int{}
	}

	if *ul == nil {
		return []int{}
	}

	if len(*ul) == 0 {
		return []int{}
	}

	return *ul
}

// CheckIntegrity verifies integrity of the list.
func (ul *UidList) CheckIntegrity() (err error) {
	if ul == nil {
		return errors.New(Err_DestinationIsNotInitialised)
	}

	if !ul.isUnique() {
		return errors.New(Err_ItemsAreNotUnique)
	}
	return nil
}

// isUnique checks uniqueness of all items.
func (ul *UidList) isUnique() bool {
	m := make(map[int]bool)
	var isDuplicate bool

	for _, uid := range *ul {
		_, isDuplicate = m[uid]
		if isDuplicate {
			return false
		}

		m[uid] = true
	}

	return true
}

// Size returns list's size, i.e. it counts the items.
func (ul *UidList) Size() (n int) {
	if ul == nil {
		return 0
	}

	return int(len(*ul))
}

// HasItem checks whether an item is contained in the list or not.
func (ul *UidList) HasItem(uid int) bool {
	if ul == nil {
		return false
	}

	for _, x := range *ul {
		if x == uid {
			return true
		}
	}

	return false
}

// AddItem add a new identifier to the end of the list.
// If 'addToTop' is set to 'True', then the item is added to the beginning
// (top) of the list; otherwise – to the end (bottom) of the list.
func (ul *UidList) AddItem(uid int, addToTop bool) (err error) {
	if ul.HasItem(uid) {
		return fmt.Errorf(ErrF_DuplicateUid, uid)
	}

	if addToTop {
		ul.prependItem(uid)
	} else {
		ul.appendItem(uid)
	}

	return nil
}

// prependItem adds an item to the beginning of the list.
func (ul *UidList) prependItem(uid int) {
	// Add an empty item.
	*ul = append(*ul, 0)

	// Shift elements.
	for i := len(*ul) - 1; i > 0; i-- {
		(*ul)[i] = (*ul)[i-1]
	}

	// Set the new item.
	(*ul)[0] = uid
}

// appendItem adds an item to the end of the list.
func (ul *UidList) appendItem(uid int) {
	*ul = append(*ul, uid)
}

// SearchForItem searches for an item in the list.
// If an item is found, its index is returned without error.
// If an item is not found, an error is returned.
func (ul *UidList) SearchForItem(uid int) (idx int, err error) {
	for pos, x := range *ul {
		if x == uid {
			return int(pos), nil
		}
	}

	return IndexOnError, fmt.Errorf(ErrF_UidIsNotFound, uid)
}

// RemoveItem deletes an identifier from the list shifting its items.
func (ul *UidList) RemoveItem(uid int) (err error) {
	var pos int
	pos, err = ul.SearchForItem(uid)
	if err != nil {
		return err
	}

	return ul.RemoveItemAtPos(pos)
}

// RemoveItemAtPos removes the item at position.
func (ul *UidList) RemoveItemAtPos(pos int) (err error) {
	if pos < 0 {
		return errors.New(Err_Position)
	}

	lastIndex := int(len(*ul) - 1)
	if pos > lastIndex {
		return errors.New(Err_Position)
	}

	ul.removeItemAtPos(pos, lastIndex)
	return nil
}

// removeItemAtPos removes the existing item at position.
func (ul *UidList) removeItemAtPos(pos int, lastIndex int) {
	if pos != lastIndex {
		copy((*ul)[pos:], (*ul)[pos+1:])
	}
	ul.removeLastItem(lastIndex)
}

// removeLastItem removes the last existing item.
func (ul *UidList) removeLastItem(lastIndex int) {
	[]int(*ul)[lastIndex] = 0
	*ul = (*ul)[:lastIndex]
}

// RaiseItem moves an existing identifier to the top of the list.
func (ul *UidList) RaiseItem(uid int) (isAlreadyRaised bool, err error) {
	var pos int
	pos, err = ul.SearchForItem(uid)
	if err != nil {
		return false, err
	}

	if pos == 0 {
		return true, nil
	}

	var movedItem = (*ul)[pos]
	for i := pos; i > 0; i-- {
		(*ul)[i] = (*ul)[i-1]
	}
	(*ul)[0] = movedItem

	return false, nil
}

// MoveItemUp moves an existing identifier one position upwards if possible.
func (ul *UidList) MoveItemUp(uid int) (err error) {
	var pos int
	pos, err = ul.SearchForItem(uid)
	if err != nil {
		return err
	}

	// Check for top edge position.
	if pos == 0 {
		return errors.New(Err_EdgePosition)
	}

	// Move the item one position upwards.
	(*ul)[pos-1], (*ul)[pos] = (*ul)[pos], (*ul)[pos-1]
	return nil
}

// MoveItemDown moves an existing identifier one position downwards if possible.
func (ul *UidList) MoveItemDown(uid int) (err error) {
	var pos int
	pos, err = ul.SearchForItem(uid)
	if err != nil {
		return err
	}

	// Check for bottom edge position.
	if pos == len(*ul)-1 {
		return errors.New(Err_EdgePosition)
	}

	// Move the item one position upwards.
	(*ul)[pos+1], (*ul)[pos] = (*ul)[pos], (*ul)[pos+1]
	return nil
}

// Scan method provides compatibility with SQL JSON data type.
func (ul *UidList) Scan(src any) (err error) {
	if ul == nil {
		return errors.New(Err_DestinationIsNotInitialised)
	}

	switch src.(type) {
	case []byte:
		{
			data := new(UidList)

			err = json.Unmarshal(src.([]byte), data)
			if err != nil {
				return err
			}

			if data != nil {
				*ul = *data
			}

			return nil
		}

	case nil:
		return nil

	default:
		return fmt.Errorf(ErrF_UnsupportedDataType, reflect.TypeOf(src).String())
	}
}

// Value method provides compatibility with SQL JSON data type.
func (ul *UidList) Value() (dv driver.Value, err error) {
	if ul == nil {
		return nil, nil
	}

	var buf []byte
	buf, err = json.Marshal(ul)
	if err != nil {
		return nil, err
	}

	return driver.Value(buf), nil
}

// ValuesString lists items as a simple plain text with a comma as separator.
func (ul *UidList) ValuesString() (values string, err error) {
	if ul == nil {
		return StringOnError, nil
	}

	if len(*ul) == 0 {
		return StringOnError, nil
	}

	var sb = strings.Builder{}
	iLast := len(*ul) - 1
	for i, uid := range *ul {
		if i < iLast {
			_, err = sb.WriteString(strconv.Itoa(uid) + ListItemSeparator)
		} else {
			_, err = sb.WriteString(strconv.Itoa(uid))
		}
		if err != nil {
			return StringOnError, err
		}
	}

	return sb.String(), nil
}

// OnPage returns paginated items.
func (ul *UidList) OnPage(pageNumber int, pageSize int) (ulop *UidList) {
	if pageNumber < 1 {
		return nil
	}

	if ul == nil {
		return nil
	}

	if *ul == nil {
		return nil
	}

	if len(*ul) == 0 {
		return nil
	}

	// Last index in array.
	iLast := len(*ul) - 1

	// Left index of a virtual page.
	ipL := pageSize * (pageNumber - 1)
	if iLast < ipL {
		return nil
	}

	// Right index of a virtual page.
	ipR := ipL + pageSize - 1
	var x UidList
	if iLast < ipR {
		x = (*ul)[ipL : iLast+1]
	} else {
		x = (*ul)[ipL : ipR+1]
	}

	return &x
}

// LastElement returns the last item of the list.
func (ul *UidList) LastElement() (lastElement *int) {
	if ul == nil {
		return nil
	}

	if *ul == nil {
		return nil
	}

	if len(*ul) == 0 {
		return nil
	}

	iLast := len(*ul) - 1
	x := (*ul)[iLast]
	return &x
}
