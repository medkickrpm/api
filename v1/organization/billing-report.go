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
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Param service query string false "Service" Enums(RPM, CCM, PCM, BHI)
// @Success 200 {object} BillingReportResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organization/{id}/billing-report [get]
func getBillingReport(c echo.Context) error {
	param := struct {
		OrganizationID uint   `param:"id"`
		StartDate      string `query:"start_date" validate:"required"`
		EndDate        string `query:"end_date" validate:"required"`
		Service        string `query:"service"`
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

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to load location",
		})
	}
	startDate, err := time.ParseInLocation("2006-01-02", param.StartDate, loc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to parse start date",
		})
	}

	endDate, err := time.ParseInLocation("2006-01-02", param.EndDate, loc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Failed to parse end date",
		})
	}
	endDate = endDate.AddDate(0, 0, 1)

	rawBills, err := models.ListBillByMonth(param.Service, startDate, endDate)
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

	currentDate := time.Now().In(loc)

	res := make([]BillingRecordBody, 0)

	for _, services := range groupBill {
		for _, bills := range services {
			var record BillingRecordBody
			var once sync.Once
			var codes []string
			for _, bill := range bills {
				codes = append(codes, fmt.Sprintf("%d", bill.CPTCode))
				once.Do(func() {
					dob := bill.Patient.DOB
					if d, err2 := time.Parse("2006-01-02", dob); err2 == nil {
						dob = d.Format("02/01/2006")
					}
					record = BillingRecordBody{
						PatientID:   bill.PatientID,
						FirstName:   bill.Patient.FirstName,
						LastName:    bill.Patient.LastName,
						DOB:         dob,
						DOS:         currentDate.Format("02/01/2006"),
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
		StartDate: param.StartDate,
		EndDate:   param.EndDate,
		Service:   param.Service,
		Records:   res,
	})
}
