package amoCRM

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	fieldReqText = 1433145
)

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
func (crm AmoCRM) GetContact(token string, id int) []byte {

	url := "https://sindoor.amocrm.ru/api/v4/contacts/" + strconv.Itoa(id)

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
		return responseBody // Возвращаем тело для отладки
	}

	fmt.Println(resp.Status)

	return responseBody
}

/*
Получение одной сделки в AmoCRM
*/
func (crm AmoCRM) GetLead(token string, id int) []byte {

	url := "https://sindoor.amocrm.ru/api/v4/leads/" + strconv.Itoa(id) + "?with=contacts"

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
		return responseBody // Возвращаем тело для отладки
	}

	fmt.Println(resp.Status)

	return responseBody
}

/*
Добавление одного контакта в AmoCRM
*/
func (crm AmoCRM) AddContact(nickName, number, request string) string {

	telephoneNumValue := Value{Value: number, EnumId: 273245, EnumCode: "WORK"}
	telephoneNumValues := Values{telephoneNumValue}

	telegramNicknameValue := Value{Value: "t.me/" + nickName}
	telegramNicknameValues := Values{telegramNicknameValue}

	telegramRequestValue := Value{Value: request}
	telegramRequestValues := Values{telegramRequestValue}

	telephoneNumField := CustomFieldsValue{Values: telephoneNumValues, FieldName: "Телефон", FieldId: 312455, FieldCode: "PHONE"}
	telegramNicknameField := CustomFieldsValue{Values: telegramNicknameValues, FieldName: "Ссылка на телеграм", FieldId: 387109}
	telegramRequestField := CustomFieldsValue{Values: telegramRequestValues, FieldName: "Запрос", FieldId: 338825}

	customFieldsValues := CustomFieldsValues{telephoneNumField, telegramNicknameField, telegramRequestField}
	contacts := Contacts{Contact{FirstName: "Сергей", LastName: "Молчанов", CustomFieldsValues: customFieldsValues}}

	url := "https://sindoor.amocrm.ru/api/v4/contacts"

	body, err := json.Marshal(contacts)
	if err != nil {
		log.Fatalf("Ошибка чтения структуры Сделка: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+crm.Token)

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

func (crm AmoCRM) EditContact(contactId int, request string) (string, error) {

	telegramRequestValue := Value{Value: request}
	telegramRequestValues := Values{telegramRequestValue}

	telegramRequestField := CustomFieldsValue{Values: telegramRequestValues, FieldName: "Текст запроса", FieldId: fieldReqText}

	customFieldsValues := CustomFieldsValues{telegramRequestField}
	contacts := Contacts{Contact{Id: contactId, CustomFieldsValues: customFieldsValues}}

	url := "https://sindoor.amocrm.ru/api/v4/contacts"

	body, err := json.Marshal(contacts)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка чтения структуры Сделка: %v", err))
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))

	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка создания запроса: %v", err))
	}

	req.Header.Set("Authorization", "Bearer "+crm.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка выполнения запроса: %v", err))
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return string(responseBody), errors.New(fmt.Sprintf("Ошибка: получен статус %s", resp.Status))
	}

	return fmt.Sprintf("Статус запроса: %+v ", resp.Status), nil
}

/*
Добавление одной сделки в AmoCRM
*/
func (crm AmoCRM) AddLead(token string, leads Leads) string {
	url := "https://sindoor.amocrm.ru/api/v4/leads"

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

	url := "https://sindoor.amocrm.ru/api/v4/account"
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
