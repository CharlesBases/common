package file

/*
 导出Excel
*/

import (
	"common/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type excel struct {
	Path  string                `json:"path"`
	Value []map[int]interface{} `json:"value"`
}

func Excel(r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return
	}
	e := excel{}
	err = json.Unmarshal(bytes, &e)
	if err != nil {
		log.Error(err)
		return
	}

	name := fmt.Sprintf(`%s.xlsx`, time.Now().Format("2006-01-02 13:04:05"))
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("Sheet1")

	for k, v := range e.Value {
		for i := 0; i < len(v); i++ {
			xlsx.SetCellValue("Sheet1", cell(k+1, i), v[i])
		}
	}

	xlsx.SetActiveSheet(index)
	err = xlsx.SaveAs(fmt.Sprintf(`%s/%s`, e.Path, name))
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func cell(line int, row int) string {
	return fmt.Sprintf(`%s%d`, string(65+row), line)
}
