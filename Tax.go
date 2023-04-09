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

	calculateButton := widget.NewButton("Calculate", func() {
		handleCalculateButtonClick(w, salaryEntry)
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
			calculateButton,
			layout.NewSpacer(),
			layout.NewSpacer(),
		),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(350, 450))
	w.SetFixedSize(true)
	w.ShowAndRun()
}

func handleCalculateButtonClick(w fyne.Window, salaryEntry *widget.Entry) {
	salary, err := strconv.ParseFloat(salaryEntry.Text, 64)
	if err != nil {
		fmt.Println("Invalid Input:", err)
		return
	}

	sssCon, phCon, piCon, totalCon := computeTotalContributions(salary)
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

	w.SetContent(container.NewPadded(
		container.NewVBox(
			layout.NewSpacer(),
			container.NewMax(title),
			salaryEntry,
			layout.NewSpacer(),
			widget.NewButton("Calculate", func() {
				handleCalculateButtonClick(w, salaryEntry)
			}),
			layout.NewSpacer(),
			resultTable,
			layout.NewSpacer(),
		),
	))
}

func computeTotalContributions(salary float64) (float64, float64, float64, float64) {
	sssContribution := computeSSSContribution(salary)
	phContribution := computePHContribution(salary)
	piContribution := computePAGIBIGContribution(salary)
	return sssContribution, phContribution, piContribution, (sssContribution + phContribution + piContribution)
}

func computeSSSContribution(salary float64) float64 {
	incomeRange := getSSSIncomeRange()
	contributionRange := getSSSContributionRange()

	for i := 0; i < 45; i++ {
		switch {
		case i != 44 && (salary >= incomeRange[i] && salary < incomeRange[i+1]):
			return contributionRange[i]

		case i == 44 && salary >= incomeRange[i]:
			return contributionRange[i]
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

func getSSSContributionRange() []float64 {
	contributionRange := make([]float64, 45)
	contributionRange[0] = 135.00
	for i := 0; i < 44; i++ {
		contributionRange[i+1] = contributionRange[i] + 22.5
	}

	return contributionRange
}

func computePHContribution(salary float64) float64 {
	return salary * 0.04
}

func computePAGIBIGContribution(salary float64) float64 {
	switch {
	case salary >= 1000 && salary <= 1500:
		return salary * 0.01
	case salary > 1500:
		return salary * 0.02
	default:
		return 0
	}
}

func computeWithholdingTax(salary float64) float64 {
	incomeRange := []float64{20833, 33333, 66667, 166667, 666667}
	switch {
	case salary < incomeRange[0]:
		return 0
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
