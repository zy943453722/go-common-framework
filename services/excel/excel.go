package excel

import (
	"path/filepath"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go-common-framework/util"
)

type IExcel interface {
	GetSheetData() ([][]string, error)
	SetSheetData(string, []string, [][]string) error
	SetNewSheet(string) int
	SaveFile(string) (string, error)
	SetSheetHeader(string, []string) error
	SetSheetDataRandom(string, string, [][]string, int) error
}

type Excel struct {
	File *excelize.File
}

func NewExcelInstance() *Excel {
	return &Excel{
		File: excelize.NewFile(),
	}
}

func LoadExcel(fileName string) (*Excel, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Excel{
		File: f,
	}, nil
}

func (e *Excel) GetSheetData() ([][]string, error) {
	index := e.File.GetActiveSheetIndex()
	name := e.File.GetSheetName(index)
	return e.File.GetRows(name)
}

func (e *Excel) SetNewSheet(sheetName string) int {
	return e.File.NewSheet(sheetName)
}

func (e *Excel) SetSheetHeader(sheetName string, sheetHeaders []string) error {
	colNum, rowNum := 1, 1
	for _, header := range sheetHeaders {
		//从第一列开始获取每列第一个元素位置
		axis, err := excelize.ColumnNumberToName(colNum)
		if err != nil {
			return err
		}
		axis = axis + strconv.Itoa(rowNum)
		if err = e.File.SetCellValue(sheetName, axis, header); err != nil {
			return err
		}
		colNum++
	}
	return nil
}

func (e *Excel) SetSheetDataRandom(requestId, sheetName string, sheetData [][]string, rowNum int) error {
	for _, row := range sheetData {
		colNum := 1
		for _, value := range row {
			axis, err := excelize.ColumnNumberToName(colNum)
			if err != nil {
				return err
			}
			axis = axis + strconv.Itoa(rowNum)
			if err := e.File.SetCellValue(sheetName, axis, value); err != nil {
				return err
			}
			colNum++
		}
		rowNum++
	}
	return nil
}

func (e *Excel) SetSheetData(sheetName string, sheetHeaders []string, sheetData [][]string) error {
	if err := e.SetSheetHeader(sheetName, sheetHeaders); err != nil {
		return err
	}
	rowNum := 2
	for _, row := range sheetData {
		colNum := 1
		for _, value := range row {
			axis, err := excelize.ColumnNumberToName(colNum)
			if err != nil {
				return err
			}
			axis = axis + strconv.Itoa(rowNum)
			if err := e.File.SetCellValue(sheetName, axis, value); err != nil {
				return err
			}
			colNum++
		}
		rowNum++
	}
	return nil
}

func (e *Excel) SaveFile(fileName string) (string, error) {
	res, err := util.GetBaseDir()
	if err != nil {
		return "", err
	}
	fileName = filepath.Join(res, "./logs/batch/", util.GenUuid()+fileName)
	return fileName, e.File.SaveAs(fileName)
}
