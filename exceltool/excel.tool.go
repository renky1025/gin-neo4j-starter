package exceltool

import (
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	excelize "github.com/xuri/excelize/v2"
)

type Style struct {
	Border        []Border    `json:"border"`
	Fill          Fill        `json:"fill"`
	Font          *Font       `json:"font"`
	Alignment     *Alignment  `json:"alignment"`
	Protection    *Protection `json:"protection"`
	NumFmt        int         `json:"number_format"`
	DecimalPlaces int         `json:"decimal_places"`
	CustomNumFmt  *string     `json:"custom_number_format"`
	Lang          string      `json:"lang"`
	NegRed        bool        `json:"negred"`
}

// Border 边框
type Border struct {
	Type  string `json:"type"`
	Color string `json:"color"`
	Style int    `json:"style"`
}

// Fill 填充
type Fill struct {
	Type    string   `json:"type"`
	Pattern int      `json:"pattern"`
	Color   []string `json:"color"`
	Shading int      `json:"shading"`
}

// Font 字体
type Font struct {
	Bold      bool    `json:"bold"`      // 是否加粗
	Italic    bool    `json:"italic"`    // 是否倾斜
	Underline string  `json:"underline"` // single    double
	Family    string  `json:"family"`    // 字体样式
	Size      float64 `json:"size"`      // 字体大小
	Strike    bool    `json:"strike"`    // 删除线
	Color     string  `json:"color"`     // 字体颜色
}

// Protection 保护
type Protection struct {
	Hidden bool `json:"hidden"`
	Locked bool `json:"locked"`
}

// Alignment 对齐
type Alignment struct {
	Horizontal      string `json:"horizontal"`        // 水平对齐方式
	Indent          int    `json:"indent"`            // 缩进  只要设置了值，就变成了左对齐
	JustifyLastLine bool   `json:"justify_last_line"` // 两端分散对齐，只有在水平对齐选择 distributed 时起作用
	ReadingOrder    uint64 `json:"reading_order"`     // 文字方向 不知道值范围和具体的含义
	RelativeIndent  int    `json:"relative_indent"`   // 不知道具体的含义
	ShrinkToFit     bool   `json:"shrink_to_fit"`     // 缩小字体填充
	TextRotation    int    `json:"text_rotation"`     // 文本旋转
	Vertical        string `json:"vertical"`          // 垂直对齐
	WrapText        bool   `json:"wrap_text"`         // 自动换行
}

var (
	defaultSheetName = "Sheet1" //默认Sheet名称
	defaultHeight    = 25.0     //默认行高度
)

type lzExcelExport struct {
	file      *excelize.File
	sheetName string //可定义默认sheet名称
}

func NewMyExcel() *lzExcelExport {
	return &lzExcelExport{file: createFile(), sheetName: defaultSheetName}
}

// 导出基本的表格
func (l *lzExcelExport) ExportToPath(params []map[string]string, data []map[string]interface{}, path string) (string, error) {
	l.export(params, data)
	name := createFileName()
	filePath := path + "/" + name
	err := l.file.SaveAs(filePath)
	return filePath, err
}

// 导出到浏览器。此处使用的gin框架 其他框架可自行修改ctx
func (l *lzExcelExport) ExportToWeb(params []map[string]string, data []map[string]interface{}, c *gin.Context) {
	l.export(params, data)
	buffer, _ := l.file.WriteToBuffer()
	//设置文件类型
	c.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	//设置文件名称
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName()))
	_, _ = c.Writer.Write(buffer.Bytes())
}

// 设置首行
func (l *lzExcelExport) writeTop(params []map[string]string) {
	topStyle, _ := l.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			// Italic: false,
			// Underline: "single",
			Size:   14,
			Family: "宋体",
			// Strike:    true, // 删除线
			Color: "#0000FF",
		}, Alignment: &excelize.Alignment{
			// 水平对齐方式 center left right fill(填充) justify(两端对齐)  centerContinuous(跨列居中) distributed(分散对齐)
			Horizontal: "center",
			// 垂直对齐方式 center top  justify distributed
			Vertical: "center",
			// Indent:     1,        // 缩进  只要有值就变成了左对齐 + 缩进
			// TextRotation: 30, // 旋转
			// RelativeIndent:  10,   // 好像没啥用
			// ReadingOrder:    0,    // 不知道怎么设置
			// JustifyLastLine: true, // 两端分散对齐，只有 水平对齐 为 distributed 时 设置true 才有效
			// WrapText:        true, // 自动换行
			// ShrinkToFit:     true, // 缩小字体以填充单元格
		},
	})
	var word = 'A'
	//首行写入
	for _, conf := range params {
		title := conf["title"]
		width, _ := strconv.ParseFloat(conf["width"], 64)
		line := fmt.Sprintf("%c1", word)
		//设置标题
		_ = l.file.SetCellValue(l.sheetName, line, title)
		//列宽
		_ = l.file.SetColWidth(l.sheetName, fmt.Sprintf("%c", word), fmt.Sprintf("%c", word), width)
		//设置样式
		_ = l.file.SetCellStyle(l.sheetName, line, line, topStyle)
		word++
	}
}

// 写入数据
func (l *lzExcelExport) writeData(params []map[string]string, data []map[string]interface{}) {
	lineStyle, _ := l.file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			// 水平对齐方式 center left right fill(填充) justify(两端对齐)  centerContinuous(跨列居中) distributed(分散对齐)
			Horizontal: "center",
			// 垂直对齐方式 center top  justify distributed
			Vertical: "center",
			// Indent:     1,        // 缩进  只要有值就变成了左对齐 + 缩进
			// TextRotation: 30, // 旋转
			// RelativeIndent:  10,   // 好像没啥用
			// ReadingOrder:    0,    // 不知道怎么设置
			// JustifyLastLine: true, // 两端分散对齐，只有 水平对齐 为 distributed 时 设置true 才有效
			// WrapText:        true, // 自动换行
			// ShrinkToFit:     true, // 缩小字体以填充单元格
		},
	})
	//数据写入
	var j = 2 //数据开始行数
	for i, val := range data {
		//设置行高
		_ = l.file.SetRowHeight(l.sheetName, i+1, defaultHeight)
		//逐列写入
		var word = 'A'
		for _, conf := range params {
			valKey := conf["key"]
			line := fmt.Sprintf("%c%v", word, j)
			isNum := conf["is_num"]

			//设置值
			if isNum != "0" {
				valNum := fmt.Sprintf("'%v", val[valKey])
				_ = l.file.SetCellValue(l.sheetName, line, valNum)
			} else {
				_ = l.file.SetCellValue(l.sheetName, line, val[valKey])
			}

			//设置样式
			_ = l.file.SetCellStyle(l.sheetName, line, line, lineStyle)
			word++
		}
		j++
	}
	//设置行高 尾行
	_ = l.file.SetRowHeight(l.sheetName, len(data)+1, defaultHeight)
}

func (l *lzExcelExport) export(params []map[string]string, data []map[string]interface{}) {
	l.writeTop(params)
	l.writeData(params, data)
}

func createFile() *excelize.File {
	f := excelize.NewFile()
	// 创建一个默认工作表
	sheetName := defaultSheetName
	index, _ := f.NewSheet(sheetName)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	return f
}

func createFileName() string {
	name := time.Now().Format("2006-01-02-15-04-05")
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("excle-%v-%v.xlsx", name, rand.Int63n(time.Now().Unix()))
}

// excel导出(数据源为Struct) []interface{}
func (l *lzExcelExport) ExportExcelByStruct(titleList []string, data []interface{}, fileName string, sheetName string, c *gin.Context) error {
	l.file.SetSheetName("Sheet1", sheetName)
	header := make([]string, 0)
	header = append(header, titleList...)
	rowStyleID, _ := l.file.NewStyle(&excelize.Style{Font: &excelize.Font{
		Bold: true,
		// Italic: false,
		// Underline: "single",
		Size:   13,
		Family: "宋体",
		// Strike:    true, // 删除线
		Color: "#666666",
	}, Alignment: &excelize.Alignment{
		// 水平对齐方式 center left right fill(填充) justify(两端对齐)  centerContinuous(跨列居中) distributed(分散对齐)
		Horizontal: "center",
		// 垂直对齐方式 center top  justify distributed
		Vertical: "center",
		// Indent:     1,        // 缩进  只要有值就变成了左对齐 + 缩进
		// TextRotation: 30, // 旋转
		// RelativeIndent:  10,   // 好像没啥用
		// ReadingOrder:    0,    // 不知道怎么设置
		// JustifyLastLine: true, // 两端分散对齐，只有 水平对齐 为 distributed 时 设置true 才有效
		// WrapText:        true, // 自动换行
		// ShrinkToFit:     true, // 缩小字体以填充单元格
	}})
	_ = l.file.SetSheetRow(sheetName, "A1", &header)
	_ = l.file.SetRowHeight("Sheet1", 1, 30)
	length := len(titleList)
	headStyle := Letter(length)
	var lastRow string
	var widthRow string
	for k, v := range headStyle {

		if k == length-1 {

			lastRow = fmt.Sprintf("%s1", v)
			widthRow = v
		}
	}
	if err := l.file.SetColWidth(sheetName, "A", widthRow, 30); err != nil {
		fmt.Print("错误--", err.Error())
	}
	rowNum := 1
	for _, v := range data {

		t := reflect.TypeOf(v)
		fmt.Print("--ttt--", t.NumField())
		value := reflect.ValueOf(v)
		row := make([]interface {
		}, 0)
		for l := 0; l < t.NumField(); l++ {

			val := value.Field(l).Interface()
			row = append(row, val)
		}
		rowNum++
		err := l.file.SetSheetRow(sheetName, "A"+strconv.Itoa(rowNum), &row)
		_ = l.file.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), lastRow, rowStyleID)
		if err != nil {
			return err
		}
	}
	disposition := fmt.Sprintf("attachment; filename=%s.xlsx", url.QueryEscape(fileName))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", disposition)
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	return l.file.Write(c.Writer)
}

// Letter 遍历a-z
func Letter(length int) []string {
	var str []string
	for i := 0; i < length; i++ {
		str = append(str, string(rune('A'+i)))
	}
	return str
}
