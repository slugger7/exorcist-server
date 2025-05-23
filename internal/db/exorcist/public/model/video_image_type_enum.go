//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type VideoImageTypeEnum string

const (
	VideoImageTypeEnum_Thumbnail VideoImageTypeEnum = "thumbnail"
	VideoImageTypeEnum_Chapter   VideoImageTypeEnum = "chapter"
)

var VideoImageTypeEnumAllValues = []VideoImageTypeEnum{
	VideoImageTypeEnum_Thumbnail,
	VideoImageTypeEnum_Chapter,
}

func (e *VideoImageTypeEnum) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "thumbnail":
		*e = VideoImageTypeEnum_Thumbnail
	case "chapter":
		*e = VideoImageTypeEnum_Chapter
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for VideoImageTypeEnum enum")
	}

	return nil
}

func (e VideoImageTypeEnum) String() string {
	return string(e)
}
