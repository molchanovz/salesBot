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

const (
	token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6IjBjODY1YWQxN2RkOTY4Y2Y0ODRhMDBiMGU0NjJlYWFjMDY5MzhmNGUwN2RhN2FjNTUwZDI5MmYzOGNkNmQyYWQyZTFiMmViYTgwNTZmYzI5In0.eyJhdWQiOiI0N2YzYThjOC1mYzdmLTQ3NmMtYWZlMi05MjZkYmE3M2QyMTEiLCJqdGkiOiIwYzg2NWFkMTdkZDk2OGNmNDg0YTAwYjBlNDYyZWFhYzA2OTM4ZjRlMDdkYTdhYzU1MGQyOTJmMzhjZDZkMmFkMmUxYjJlYmE4MDU2ZmMyOSIsImlhdCI6MTcyOTUxNDU4MSwibmJmIjoxNzI5NTE0NTgxLCJleHAiOjE3MzAzMzI4MDAsInN1YiI6IjExNjc1MzUwIiwiZ3JhbnRfdHlwZSI6IiIsImFjY291bnRfaWQiOjMyMDIwMjEwLCJiYXNlX2RvbWFpbiI6ImFtb2NybS5ydSIsInZlcnNpb24iOjIsInNjb3BlcyI6WyJjcm0iLCJmaWxlcyIsImZpbGVzX2RlbGV0ZSIsIm5vdGlmaWNhdGlvbnMiLCJwdXNoX25vdGlmaWNhdGlvbnMiXSwiaGFzaF91dWlkIjoiNjBkNmQ2MzYtMTljNC00ZmMxLTk3YTYtZmE0YjA4Zjg0MmYzIiwiYXBpX2RvbWFpbiI6ImFwaS1iLmFtb2NybS5ydSJ9.M_Q5LXfSnDSUki417wtg_NWKynndUFrc9qt00sPsIILeZFUfakcCh16UUP3oZ8hGZv418VARO0JIDqQDBuKOQOlCC-c2a0Wwt5T-9Vk_KVKEQefqB3a7TblTQl-0Vpx4QuILRKoPOajECD52PGp4DL-VshgUQt1aibyNEjioidmSxOhREVs0SoPxOTBENVTKaKzqc7s3Ehq0gpg-uxEHD45AUBnOaH-Oc_nZ97a6Ji39PgedtcKqlYA8Iipka0EXPNoWBZlfgZflZHxxA90ikXkLaCP56IKlGylpVObE3dSm3cP3lsonwl9bYJ6UauIbBQ8PDFuytMCABFQvEVNQsQ"
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

	fmt.Println(req)
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

	fmt.Println(req)
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
func (crm AmoCRM) AddContact(token string, contacts Contacts) string {
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
