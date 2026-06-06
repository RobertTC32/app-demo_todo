package todo_store

import (
	todo_ports "RobertTC32/example-demo_hello/src/todo/ports"
	"encoding/json"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
)

// examples

func (this *MysqlStore) GetDbmapping1() ([]todo_ports.Dbmapping1DTO, error) {
	slog.Debug("store::GetDbmapping1() - Executing")
	var tt []todo_ports.Dbmapping1DTO
	sql := `SELECT id, anint, abigint, anintunsighed, abigintunsigned, adecimal, afloat, adouble, aboolean, avarchar, adatetime, ablob
		FROM dbmapping ORDER BY id`
	if err := this.Db.Select(&tt, sql); err != nil {
		slog.Error("store::GetDbmapping1() - Failed to find dbmapping", "error", err)
		return []todo_ports.Dbmapping1DTO{}, err
	}
	for i := range tt {
		t := tt[i]
		j, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println("GetDbmapping1() -", "index", i, "value \n", string(j))
	}
	return tt, nil
}

func (this *MysqlStore) GetDbmapping2() ([]todo_ports.Dbmapping2DTO, error) {
	slog.Debug("store::GetDbmapping2() - Executing")
	var tt []todo_ports.Dbmapping2DTO
	sql := `SELECT id, anint, abigint, anintunsighed, abigintunsigned, adecimal, afloat, adouble, aboolean, avarchar, adatetime, ablob
		FROM dbmapping ORDER BY id`
	if err := this.Db.Select(&tt, sql); err != nil {
		slog.Error("store::GetDbmapping2() - Failed to find dbmapping", "error", err)
		return []todo_ports.Dbmapping2DTO{}, err
	}
	for i := range tt {
		t := tt[i]
		j, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println("GetDbmapping2() -", "index", i, "value \n", string(j))
	}
	return tt, nil
}
