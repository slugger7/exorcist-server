//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type MediaRelationTypeEnum string

const (
	MediaRelationTypeEnum_Thumbnail MediaRelationTypeEnum = "thumbnail"
	MediaRelationTypeEnum_Chapter   MediaRelationTypeEnum = "chapter"
	MediaRelationTypeEnum_Media     MediaRelationTypeEnum = "media"
)

var MediaRelationTypeEnumAllValues = []MediaRelationTypeEnum{
	MediaRelationTypeEnum_Thumbnail,
	MediaRelationTypeEnum_Chapter,
	MediaRelationTypeEnum_Media,
}

func (e *MediaRelationTypeEnum) Scan(value interface{}) error {
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
		*e = MediaRelationTypeEnum_Thumbnail
	case "chapter":
		*e = MediaRelationTypeEnum_Chapter
	case "media":
		*e = MediaRelationTypeEnum_Media
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for MediaRelationTypeEnum enum")
	}

	return nil
}

func (e MediaRelationTypeEnum) String() string {
	return string(e)
}
