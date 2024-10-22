package amoCRM

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	//telephoneNumValue := Value{Value: "88005553535", EnumId: 273245, EnumCode: "WORK"}
	//telephoneNumValues := Values{telephoneNumValue}
	//
	//telegramNicknameValue := Value{Value: "@molchanovz"}
	//telegramNicknameValues := Values{telegramNicknameValue}
	//
	//telegramRequestValue := Value{Value: "У меня сломалась дверь, знает кто-нибудь хороших специалистов?"}
	//telegramRequestValues := Values{telegramRequestValue}
	//
	//telephoneNumField := CustomFieldsValue{Values: telephoneNumValues, FieldName: "Телефон", FieldId: 312455, FieldCode: "PHONE"}
	//telegramNicknameField := CustomFieldsValue{Values: telegramNicknameValues, FieldName: "Ник телеграм", FieldId: 337421}
	//telegramRequestField := CustomFieldsValue{Values: telegramRequestValues, FieldName: "Запрос", FieldId: 338825}
	//
	//customFieldsValues := CustomFieldsValues{telephoneNumField, telegramNicknameField, telegramRequestField}
	//contacts := Contacts{Contact{FirstName: "Сергей", LastName: "Молчанов", CustomFieldsValues: customFieldsValues}}
	//fmt.Println(addContact(token, contacts))

	//fmt.Println(getContact(token))

	//fmt.Println(getLead(token))
	//contacts := Contacts{Contact{Id: 1185541}}
	//lead := Lead{Name: "Тест", Embedded: Embedded{Contacts: contacts}}
	//leads := Leads{lead}
	//fmt.Println(addLead(token, leads))
	//fmt.Println(getLead(token, 2850715))
}

type AmoCRMConfig struct {
	Token string
}

type AmoCRM struct {
	Token string
}

func NewAmoCRM(config AmoCRMConfig) *AmoCRM {
	return &AmoCRM{Token: config.Token}
}

/*
Получение одного контакта в AmoCRM
*/
func (crm AmoCRM) GetContact(token string, id int) string {

	url := "https://molchanovtop.amocrm.ru/api/v4/contacts/" + strconv.Itoa(id)

	body := []byte(`with:{}`)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	responseBody, _ := io.ReadAll(resp.Body)

	// Если статус не 200 OK, выводим детальную ошибку и тело ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: получен статус %s", resp.Status)
		return string(responseBody) // Возвращаем тело для отладки
	}

	fmt.Println(resp.Status)

	return string(responseBody)
}

/*
Получение одной сделки в AmoCRM
*/
func (crm AmoCRM) GetLead(token string, id int) string {

	url := "https://molchanovtop.amocrm.ru/api/v4/leads/" + strconv.Itoa(id) + "?with=contacts"

	body := []byte(``)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	responseBody, _ := io.ReadAll(resp.Body)

	// Если статус не 200 OK, выводим детальную ошибку и тело ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: получен статус %s", resp.Status)
		log.Printf("Тело ответа: %s", string(responseBody))
		return string(responseBody) // Возвращаем тело для отладки
	}

	fmt.Println(resp.Status)

	return string(responseBody)
}

/*
Добавление одного контакта в AmoCRM
*/
func (crm AmoCRM) AddContact(token string) string {

	telephoneNumValue := Value{Value: "88005553535", EnumId: 273245, EnumCode: "WORK"}
	telephoneNumValues := Values{telephoneNumValue}

	telegramNicknameValue := Value{Value: "@molchanovz"}
	telegramNicknameValues := Values{telegramNicknameValue}

	telegramRequestValue := Value{Value: "У меня сломалась дверь, знает кто-нибудь хороших специалистов?"}
	telegramRequestValues := Values{telegramRequestValue}

	telephoneNumField := CustomFieldsValue{Values: telephoneNumValues, FieldName: "Телефон", FieldId: 312455, FieldCode: "PHONE"}
	telegramNicknameField := CustomFieldsValue{Values: telegramNicknameValues, FieldName: "Ник телеграм", FieldId: 337421}
	telegramRequestField := CustomFieldsValue{Values: telegramRequestValues, FieldName: "Запрос", FieldId: 338825}

	customFieldsValues := CustomFieldsValues{telephoneNumField, telegramNicknameField, telegramRequestField}
	contacts := Contacts{Contact{FirstName: "Сергей", LastName: "Молчанов", CustomFieldsValues: customFieldsValues}}

	url := "https://molchanovtop.amocrm.ru/api/v4/contacts"

	body, err := json.Marshal(contacts)
	if err != nil {
		log.Fatalf("Ошибка чтения структуры Сделка: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: получен статус %s", resp.Status)
		log.Printf("Тело ответа: %s", string(responseBody))
		return string(responseBody)
	}

	return fmt.Sprintf("Статус запроса: %+v ", resp.Status)
}

/*
Добавление одной сделки в AmoCRM
*/
func (crm AmoCRM) AddLead(token string, leads Leads) string {
	url := "https://molchanovtop.amocrm.ru/api/v4/leads"

	body, err := json.Marshal(leads)
	if err != nil {
		log.Fatalf("Ошибка чтения структуры Сделка: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: получен статус %s", resp.Status)
		log.Printf("Тело ответа: %s", string(responseBody))
		return string(responseBody)
	}

	return fmt.Sprintf("Статус запроса: %+v ", resp.Status)
}

/*
Получение информации об аккаунте в AmoCRM
*/
func (crm AmoCRM) GetParams(token string) string {

	url := "https://molchanovtop.amocrm.ru/api/v4/account"
	body := []byte(`{}`)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка: получен статус %s", resp.Status)
	}

	jsonString, _ := io.ReadAll(resp.Body)

	return string(jsonString)
}
