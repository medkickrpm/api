package organization

import (
	"MedKick-backend/pkg/database/models"
	"MedKick-backend/pkg/echo/dto"
	"MedKick-backend/pkg/validator"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// getBillingReport godoc
// @Summary Get Billing Report
// @Description Get Billing Report
// @Tags Organization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} BillingReportResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organization/{id}/billing-report [get]
func getBillingReport(c echo.Context) error {
	param := struct {
		OrganizationID uint `param:"id"`
		Year           int  `query:"year" validate:"required"`
		Month          int  `query:"month" validate:"required"`
	}{}

	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := validator.Validate.Struct(param); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	rawBills, err := models.ListBillByMonth(param.Year, param.Month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list bills",
		})
	}

	groupBill := make(map[uint]map[string][]models.Bill)

	for _, bill := range rawBills {
		patientID := bill.PatientID
		if _, ok := groupBill[patientID]; !ok {
			groupBill[patientID] = make(map[string][]models.Bill)
		}

		serviceCode := bill.ServiceCode
		if _, ok := groupBill[patientID][serviceCode]; !ok {
			groupBill[patientID][serviceCode] = []models.Bill{}
		}

		groupBill[patientID][serviceCode] = append(groupBill[patientID][serviceCode], bill)
	}

	loc, _ := time.LoadLocation("EST")
	lastDayOfMonth := time.Now().In(loc)

	res := make([]BillingRecordBody, 0)

	for _, services := range groupBill {
		for _, bills := range services {
			var record BillingRecordBody
			var once sync.Once
			var codes []string
			for _, bill := range bills {
				codes = append(codes, fmt.Sprintf("%d", bill.CPTCode))
				once.Do(func() {
					record = BillingRecordBody{
						PatientID:   bill.PatientID,
						FirstName:   bill.Patient.FirstName,
						LastName:    bill.Patient.LastName,
						DOB:         bill.Patient.DOB,
						DOS:         lastDayOfMonth.Format("02/01/2006"),
						ServiceCode: bill.ServiceCode,
					}
				})
			}
			record.CPTCodes = strings.Join(codes, ", ")
			res = append(res, record)
		}
	}

	sort.SliceStable(res, func(i, j int) bool {
		if res[i].PatientID == res[j].PatientID {
			return res[i].ServiceCode < res[j].ServiceCode
		}
		return res[i].PatientID < res[j].PatientID
	})

	return c.JSON(http.StatusOK, BillingReportResponse{
		Year:    int64(param.Year),
		Month:   int64(param.Month),
		Records: res,
	})
}
