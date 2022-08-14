package cv

import (
	"githab.com/kbats183/argotech/backend/models"
	rtfdoc "github.com/therox/rtf-doc"
	"strconv"
	"strings"
)

const (
	defaultFont  = rtfdoc.FontTimesNewRoman
	defaultColor = rtfdoc.ColorBlack
)

func GenerateResume(details models.UserCV) []byte {
	d := rtfdoc.NewDocument()
	d.SetOrientation(rtfdoc.OrientationPortrait)
	d.SetFormat(rtfdoc.FormatA4)

	p := d.AddParagraph()
	p.SetAlign(rtfdoc.AlignCenter)
	p.AddText("Резюме", 20, defaultFont, defaultColor).SetBold()

	addVeriticalMargin(d)

	p = d.AddParagraph()
	p.SetAlign(rtfdoc.AlignLeft)
	p.AddText("Фамилия, имя, отчество:   ", 12, defaultFont, defaultColor)
	p.AddText(" "+details.LastName+" "+details.FirstName+" "+details.Patronymic,
		12, defaultFont, defaultColor)
	if details.DateOfBirth != nil {
		p.AddNewLine()
		p.AddText("Дата рождения:", 12, defaultFont, defaultColor)
		p.AddText(" ", 6, defaultFont, defaultColor)
		p.AddText("            ", 12, defaultFont, defaultColor)
		p.AddText(" "+*details.DateOfBirth, 12, defaultFont, defaultColor)
	}
	if details.Address != nil {
		p.AddNewLine()
		p.AddText("Адрес:                      ", 12, defaultFont, defaultColor)
		p.AddText(" "+*details.Address, 12, defaultFont, defaultColor)
	}
	if details.ContactData != nil {
		p.AddNewLine()
		p.AddText("Контактные данные:       ", 12, defaultFont, defaultColor)
		p.AddText(" "+*details.ContactData,
			12, defaultFont, defaultColor)
	}
	addVeriticalMargin(d)

	p = d.AddParagraph()
	p.SetAlign(rtfdoc.AlignLeft)
	p.AddText("Образование   ", 12, defaultFont, defaultColor).SetBold()

	p = d.AddParagraph()
	p.SetAlign(rtfdoc.AlignLeft)
	if details.SchoolName != nil && details.SchoolBeginYear != nil && details.SchoolEndYear != nil {
		schoolName := *details.SchoolName
		p.AddText(strconv.Itoa(*details.SchoolBeginYear)+" - "+strconv.Itoa(*details.SchoolEndYear)+
			"    "+schoolName, 12, defaultFont, defaultColor)
		p.AddNewLine()
	}
	if details.UniversityName != nil && details.UniversityBeginYear != nil && details.UniversityEndYear != nil {
		p.AddText(strconv.Itoa(*details.UniversityBeginYear)+" - "+strconv.Itoa(*details.UniversityEndYear)+
			"    "+*details.UniversityName, 12, defaultFont, defaultColor)
		if details.UniversityStudyProgram != nil {
			p.AddText(", образовательная программа "+*details.UniversityStudyProgram+"", 12, defaultFont, defaultColor)
		}
	}

	if details.Skills != nil {
		addVeriticalMargin(d)

		p = d.AddParagraph()
		p.SetAlign(rtfdoc.AlignLeft)
		p.AddText("Навыки", 12, defaultFont, defaultColor).SetBold()

		p = d.AddParagraph()
		p.SetAlign(rtfdoc.AlignLeft)
		for _, skill := range strings.Split(*details.Skills, "\n") {
			p.AddText(skill, 12, defaultFont, defaultColor)
			p.AddNewLine()
		}
	}

	return d.Export()
}

func addVeriticalMargin(d *rtfdoc.Document) {
	d.AddParagraph().
		AddText("", 24, defaultFont, defaultColor)
}
