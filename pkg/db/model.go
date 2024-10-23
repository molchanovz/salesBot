// Code generated by mfd-generator v0.4.5; DO NOT EDIT.

//nolint:all
//lint:file-ignore U1000 ignore unused code, it's generated
package db

var Columns = struct {
	Gigachatmessage struct {
		ID, Message, Tgid string
	}
}{
	Gigachatmessage: struct {
		ID, Message, Tgid string
	}{
		ID:      "messageid",
		Message: "message",
		Tgid:    "tgid",
	},
}

var Tables = struct {
	Gigachatmessage struct {
		Name, Alias string
	}
}{
	Gigachatmessage: struct {
		Name, Alias string
	}{
		Name:  "studup.gigachatmessages",
		Alias: "t",
	},
}

type Gigachatmessage struct {
	tableName struct{} `pg:"studup.gigachatmessages,alias:t,discard_unknown_columns"`

	ID      int    `pg:"messageid,pk"`
	Message string `pg:"message,use_zero"`
	Tgid    *int64 `pg:"tgid"`
}
