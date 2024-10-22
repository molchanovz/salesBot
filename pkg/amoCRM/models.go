package amoCRM

type Contact struct {
	Id                 int                `json:"id,omitempty"`
	FirstName          string             `json:"first_name,omitempty"`
	LastName           string             `json:"last_name,omitempty"`
	CustomFieldsValues CustomFieldsValues `json:"custom_fields_values,omitempty"`
}

type Contacts []Contact

type CustomFieldsValue struct {
	FieldId   int    `json:"field_id,omitempty"`
	FieldName string `json:"field_name,omitempty"`
	FieldCode string `json:"field_code,omitempty"`
	Values    Values `json:"values,omitempty"`
}

type CustomFieldsValues []CustomFieldsValue

type Value struct {
	Value    string `json:"value,omitempty"`
	EnumId   int    `json:"enum_id,omitempty"`
	EnumCode string `json:"enum_code,omitempty"`
}

type Values []Value

type Lead struct {
	Name               string             `json:"name"`
	CreatedBy          int                `json:"created_by,omitempty"`
	Price              int                `json:"price"`
	CustomFieldsValues CustomFieldsValues `json:"custom_fields_values,omitempty"`
	Embedded           Embedded           `json:"_embedded,omitempty"`
}

type Embedded struct {
	Contacts Contacts `json:"contacts,omitempty"`
}
type Leads []Lead
