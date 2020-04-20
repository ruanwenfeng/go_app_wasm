package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"
)

type E_VR_TYPE uint16

const (
	/// application entity title
	EVR_AE E_VR_TYPE = 0x4145

	/// age string
	EVR_AS E_VR_TYPE = 0x4153

	/// attribute tag
	EVR_AT E_VR_TYPE = 0x4154

	/// code string
	EVR_CS E_VR_TYPE = 0x4353

	/// date string
	EVR_DA E_VR_TYPE = 0x4441

	/// decimal string
	EVR_DS

	/// date time string
	EVR_DT

	/// float single-precision
	EVR_FL

	/// float double-precision
	EVR_FD

	/// integer string
	EVR_IS

	/// long string
	EVR_LO

	/// long text
	EVR_LT

	/// other byte
	EVR_OB

	/// other double
	EVR_OD

	/// other float
	EVR_OF

	/// other long
	EVR_OL

	/// other word
	EVR_OW

	/// person name
	EVR_PN

	/// short string
	EVR_SH

	/// signed long
	EVR_SL

	/// sequence of items
	EVR_SQ

	/// signed short
	EVR_SS

	/// short text
	EVR_ST

	/// time string
	EVR_TM

	/// unlimited characters
	EVR_UC

	/// unique identifier
	EVR_UI

	/// unsigned long
	EVR_UL E_VR_TYPE = 0x554c

	/// universal resource identifier or universal resource locator (URI/URL)
	EVR_UR

	/// unsigned short
	EVR_US

	/// unlimited text
	EVR_UT

	/// OB or OW depending on context
	EVR_ox

	/// SS or US depending on context
	EVR_xs

	/// US, SS or OW depending on context, used for LUT Data (thus the name)
	EVR_lt

	/// na="not applicable", for data which has no VR
	EVR_na

	/// up="unsigned pointer", used internally for DICOMDIR support
	EVR_up

	/// used internally for items
	EVR_item

	/// used internally for meta info datasets
	EVR_metainfo

	/// used internally for datasets
	EVR_dataset

	/// used internally for DICOM files
	EVR_fileFormat

	/// used internally for DICOMDIR objects
	EVR_dicomDir

	/// used internally for DICOMDIR records
	EVR_dirRecord

	/// used internally for pixel sequences in a compressed image
	EVR_pixelSQ

	/// used internally for pixel items in a compressed image
	EVR_pixelItem

	/// used internally for elements with unknown VR (encoded with 4-byte length field in explicit VR)
	EVR_UNKNOWN

	/// unknown value representation
	EVR_UN

	/// used internally for uncompressed pixel data
	EVR_PixelData

	/// used internally for overlay data
	EVR_OverlayData

	/// used internally for elements with unknown VR with 2-byte length field in explicit VR
	EVR_UNKNOWN2B
)

type Tag struct {
	groupid   uint16
	elementid uint16
}
type Element struct {
	vr     E_VR_TYPE
	vm     uint8
	length uint32
	value  interface{} //VR or DicomElement(嵌套Tag)
}
type DicomElement struct {
	Tag
	Element
}

type TagDecoder interface {
	//
}

// 文件头Dicom Tag 小端显示VR
type HeaderTagDecoder interface {
}

// 标准Dicom Tag
type StandardTagDecoder interface {
}

// 私有Dicom Tag
type PrivateTagDecoder interface {
}

func main() {

	dic, err := os.Open("D:/4thMR/AI_ZHAO_DI/SDY00000/SRS00001/IMG00000.DCM")
	if err != nil {
		fmt.Println(err)
	}
	defer dic.Close()
	ret, err := dic.Seek(128, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ret:%ld", ret)

	prefix := make([]byte, 4)
	dic.Read(prefix)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*(*string)(unsafe.Pointer(&prefix)))

	// tag := make([]byte, 2)

	// dic.Read(tag)
	// fmt.Println()
	// fmt.Println(tag[0])
	// fmt.Println(tag[1])

	type Tag struct {
		group uint16
		ele   uint16
		// vr     [2]byte
		// _      [2]byte
		// length int32
		// // data     []byte
	}
	type Element struct {
		vr     uint16
		length uint16
		data   []byte
	}
	tmp_element := make([]byte, 10)
	for i := 0; i < 1; i++ {
		var offset_size int64 = 0
		dic.Read(tmp_element)
		fmt.Println(tmp_element)
		tag := Tag{}
		ele := Element{}
		tag.group = binary.LittleEndian.Uint16(tmp_element[:2])
		tag.ele = binary.LittleEndian.Uint16(tmp_element[2:4])
		ele.vr = binary.BigEndian.Uint16(tmp_element[4:6])
		offset_size = -2
		ele.length = binary.LittleEndian.Uint16(tmp_element[6:8])

		ele.data = make([]byte, ele.length)

		dic.Seek(offset_size, io.SeekCurrent)
		dic.Read(ele.data)

		fmt.Println(tag)
		fmt.Println(ele)
	}
}
