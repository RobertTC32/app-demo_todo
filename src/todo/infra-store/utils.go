package todo_store

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// only errors but no sentinel errors are returned from mysql driver;
// example:
// Error - Number: 1062; Message: "Duplicate entry '%s' for key %d"
var messageFormats = map[uint16]string{
	1062: "Duplicate entry '%s' for key '%d'",       // Symbol: ER_DUP_ENTRY
	1406: "Data too long for column '%s' at row %d", // Symbol: ER_DATA_TOO_LONG
}

func ParseMysqlError(err error) (uint16, []string, error) {
	slog.Debug("store::ParseMysqlError() - Executing")
	driverErr, ok := err.(*mysql.MySQLError)
	if !ok {
		slog.Error("store::ParseMysqlError() - Not a MySQL error")
		return 0, nil, errors.New("Not a MySQL error")
	}
	number := driverErr.Number
	message := driverErr.Message
	messageFormat := messageFormats[number]
	slog.Debug("store::ParseMysqlError() - Found messageFormat", "Number", number, "Format", messageFormat, "Message", message)
	valuesLength := strings.Count(messageFormat, "%")
	values := make([]string, valuesLength)
	//
	ivalues := 0
	im := 0
	for i := 0; i < len(messageFormat); i++ {
		cf := messageFormat[i : i+1]
		cm := message[im : im+1]
		if cf == "%" {
			i++
			cfn := ""
			if i+2 <= len(messageFormat) {
				cfn = messageFormat[i+1 : i+2]
			}
			v := ""
			for len(cm) != 0 && (len(cfn) == 0 || cm != cfn) {
				v = v + cm
				im++
				cm = ""
				if im+1 <= len(message) {
					cm = message[im : im+1]
				}
			}
			values[ivalues] = v
			ivalues++
		} else if cm != cf {
			slog.Error("store::ParseMysqlError() - Incorrect format of MySQL error", "message", im, "format", i)
			return 0, nil, errors.New("Incorrect format of MySQL error")
		} else {
			im++
		}
	}
	slog.Debug("store::ParseMysqlError() - Returning", "number", number, "values", values)
	return number, values, nil
}
