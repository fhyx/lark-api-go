package lark

import "fhyx.online/lark-api-go/client"

type AuthContactResponse struct {
	client.Error

	Data struct {
		AuthedDepartments []string `json:"authed_departments"`
		AuthedEmployeeIDs []string `json:"authed_employee_ids"`
		AuthedOpenIDs     []string `json:"authed_open_ids"`
	} `json:"data"`
}
