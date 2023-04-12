********************
Name: DATARIO, AUDREY; DE GRACIA, SHANKY; EDRALIN, PHILIPPE; MENDOZA, ANTONIO
Language: Go
Paradigm: multi-paradigm
********************

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.example.salarytaxcalculator")
	customTheme := MyTheme{}
	a.Settings().SetTheme(customTheme)
	w := a.NewWindow("Salary Tax Calculator")

	salaryEntry := widget.NewEntry()
	salaryEntry.SetPlaceHolder("Please Type in Your Salary Here")

	empTypeRadio := widget.NewRadioGroup([]string{"Employee", "Self-Employed"}, func(selected string) {
	})
	empTypeRadio.SetSelected("Employee")

	calculateButton := widget.NewButton("Calculate", func() {
		handleCalculateButtonClick(w, salaryEntry, empTypeRadio)
	})

	title := canvas.NewText("Personal Tax Calculator", color.White)
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	content := container.NewPadded(
		container.NewVBox(
			layout.NewSpacer(),
			container.NewMax(title),
			widget.NewLabelWithStyle("MONTHLY SALARY**", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
			salaryEntry,
			empTypeRadio,
			calculateButton,
			layout.NewSpacer(),
			layout.NewSpacer(),
		),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(350, 550))
	w.SetFixedSize(true)
	w.ShowAndRun()
}

func handleCalculateButtonClick(w fyne.Window, salaryEntry *widget.Entry, empTypeRadio *widget.RadioGroup) {
	empType := 0
	if empTypeRadio.Selected == "Self-Employed" {
		empType = 1
	}

	salary, err := strconv.ParseFloat(salaryEntry.Text, 64)
	if err != nil {
		salaryEntry.SetText(fmt.Sprintf("Invalid Input (Only numbers are allowed)"))
		return
	}

	sssCon, phCon, piCon, totalCon := computeTotalContributions(salary, empType)

	taxableIncome := salary - totalCon
	incomeTax := computeWithholdingTax(taxableIncome)

	netIncome := computeNetSalary(salary, totalCon, incomeTax)

	title := canvas.NewText("Personal Tax Calculator", color.White)
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	resultTable := container.NewGridWithColumns(2,
		widget.NewLabel("Salary"), widget.NewLabel(fmt.Sprintf("%.2f", salary)),
		widget.NewLabel("SSS"), widget.NewLabel(fmt.Sprintf("%.2f", sssCon)),
		widget.NewLabel("PhilHealth"), widget.NewLabel(fmt.Sprintf("%.2f", phCon)),
		widget.NewLabel("Pag-IBIG"), widget.NewLabel(fmt.Sprintf("%.2f", piCon)),
		widget.NewLabel("Total Contributions"), widget.NewLabel(fmt.Sprintf("%.2f", totalCon)),
		widget.NewLabel("Taxable Income"), widget.NewLabel(fmt.Sprintf("%.2f", taxableIncome)),
		widget.NewLabel("Income Tax"), widget.NewLabel(fmt.Sprintf("%.2f", incomeTax)),
		widget.NewLabel("Net Income"), widget.NewLabel(fmt.Sprintf("%.2f", netIncome)),
	)

	content := container.NewPadded(
		container.NewVBox(
			layout.NewSpacer(),
			container.NewMax(title),
			salaryEntry,
			layout.NewSpacer(),
			widget.NewButton("Calculate", func() {
				handleCalculateButtonClick(w, salaryEntry, empTypeRadio)
			}),
			layout.NewSpacer(),
			empTypeRadio,
			resultTable,
			layout.NewSpacer(),
		),
	)

	w.SetContent(content)
}

func computeTotalContributions(salary float64, empType int) (float64, float64, float64, float64) {
	sssContribution := computeSSSContribution(salary, empType)
	phContribution := computePHContribution(salary, empType)
	piContribution := computePAGIBIGContribution(salary, empType)
	return sssContribution, phContribution, piContribution, (sssContribution + phContribution + piContribution)
}

func computeSSSContribution(salary float64, empType int) float64 {
	incomeRange := getSSSIncomeRange()
	contributionRange := getSSSContributionRange(empType)
	mPFRange := getSSSMandatoryProvidentFundRange(empType)

	for i := 0; i < 45; i++ {
		switch {
		case i != 44 && (salary >= incomeRange[i] && salary < incomeRange[i+1]):
			if empType == 1 {
				if i >= 25 && i < 35 {
					return contributionRange[i] + 30
				} else if i >= 35 {
					return contributionRange[i] + mPFRange[i-35] + 30
				} else {
					return contributionRange[i] + 10
				}
			} else {
				if i < 35 {
					return contributionRange[i]
				} else if i >= 35 {
					return contributionRange[i] + mPFRange[i-35]
				}

			}

		case i == 44 && salary >= incomeRange[i]:
			if empType == 1 {
				return contributionRange[i] + mPFRange[i-35] + 30
			} else {
				return contributionRange[i] + mPFRange[i-35]
			}
		}
	}

	return 0
}

func getSSSIncomeRange() []float64 {
	incomeRange := make([]float64, 45)
	incomeRange[0] = 1000.00
	incomeRange[1] = 3250.00
	for i := 1; i < 44; i++ {
		incomeRange[i+1] = incomeRange[i] + 500
	}
	return incomeRange
}

func getSSSContributionRange(empType int) []float64 {
	contributionRange := make([]float64, 45)
	var interval, fixed float64

	if empType == 1 {
		contributionRange[0] = 390.00
		interval = 65
		fixed = 2600
	} else {
		contributionRange[0] = 135.00
		interval = 22.5
		fixed = 900
	}

	for i := 0; i < 44; i++ {
		if i < 34 {
			contributionRange[i+1] = contributionRange[i] + interval
		} else {
			contributionRange[i+1] = fixed
		}

	}
	return contributionRange
}

func getSSSMandatoryProvidentFundRange(empType int) []float64 {
	mPFRange := make([]float64, 10)
	var interval float64
	if empType == 1 {
		interval = 65
	} else {
		interval = 22.5
	}
	mPFRange[0] = interval

	for i := 1; i < 10; i++ {
		mPFRange[i] += mPFRange[i-1] + interval
	}
	return mPFRange
}

func computePHContribution(salary float64, empType int) float64 {
  	divisor := 1.0
	if empType == 0 {
      		divisor = 2
	} else {
		divisor = 1
	}
  	if salary > 0 && salary <= 10000 {
		return 400.00 / divisor
    	} else if salary > 10000 && salary < 80000 {
        	return (salary * 0.04) / divisor
    	} else if salary >= 80000 {
        	return 3200.00 / divisor
    	} else {
      		return 0
    	}
}

func computePAGIBIGContribution(salary float64, empType int) float64 {
	contri := 0.0
	switch {
	case salary > 0 && salary <= 1500:
		contri = salary * 0.01
	case salary > 1500:
		if salary*0.02 < 100 {
			contri = salary * 0.02
		} else {
			contri = 100
		}
	default:
		contri = 0
	}
	if empType == 1 {
		if salary*0.02 < 100 {
			contri += salary * 0.02
		} else {
			contri += 100
		}
	}
	return contri
}

func computeWithholdingTax(salary float64) float64 {
	incomeRange := []float64{20833, 33333, 66667, 166667, 666667}
	switch {
	case salary >= incomeRange[0] && salary < incomeRange[1]:
		return (salary - 20833) * 0.2
	case salary >= incomeRange[1] && salary < incomeRange[2]:
		return 2500 + ((salary - 33333) * 0.25)
	case salary >= incomeRange[2] && salary < incomeRange[3]:
		return 10833.33 + ((salary - 66667) * 0.3)
	case salary >= incomeRange[3] && salary < incomeRange[4]:
		return 40833.33 + ((salary - 166667) * 0.32)
	case salary >= incomeRange[4]:
		return 200833.33 + ((salary - 666667) * 0.35)
	default:
		return 0
	}
}

func computeNetSalary(salary, totalCon, incomeTax float64) float64 {
	return salary - totalCon - incomeTax
}
