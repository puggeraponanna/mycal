package calendar

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Calendar struct {
	rows    [][]string
	month   time.Month
	year    int
	day     int
	SelectR *rc
	NRows   int
}

type rc struct {
	Row int
	Col int
}

func New() *Calendar {
	now := time.Now()
	return &Calendar{
		month: now.Month(),
		year:  now.Year(),
		day:   now.Day(),
	}
}

func (c *Calendar) Render() string {
	return c.String()
}

func (c *Calendar) String() string {
	days := daysIn(c.month, c.year)
	rows := [][]string{}
	offset := getDaysOffset(c.year, c.month)
	var row []string
	var selected *rc
	for i := 0; i < offset; i++ {
		row = append(row, (" "))
	}
	for i := 0; i < days; i++ {
		if (i+offset)%7 == 0 && len(row) > 0 {
			rows = append(rows, row)
			row = []string{}
		}
		row = append(row, strconv.Itoa(i+1))
		if i+1 == c.day {
			selected = &rc{Row: (i + offset) / 7, Col: (i + offset) % 7}
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}
	c.rows = rows
	if c.SelectR == nil {
		c.SelectR = selected
	}
	c.NRows = len(rows)
	header := []string{"S", "M", "T", "W", "T", "F", "S"}
	t := *table.New().Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(rows...).
		Headers(header...).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)
			if col == 0 || col == 6 {
				style = style.Copy().Background(lipgloss.Color("#818285")).Foreground(lipgloss.Color("#000000"))
			}
			if col == c.SelectR.Col && row-1 == c.SelectR.Row {
				style = style.Copy().Background(lipgloss.Color("#FF0000"))
			}
			return style
		})

	re := lipgloss.NewRenderer(os.Stdout)
	labelStyle := re.NewStyle().Foreground(lipgloss.Color("241"))
	year := labelStyle.Render(strconv.Itoa(c.year))
	month := labelStyle.Render(c.month.String())
	return lipgloss.JoinVertical(lipgloss.Center, year, month, t.Render())
}

func (c *Calendar) SetMonth(m time.Month) *Calendar {
	c.month = m
	return c
}

func (c *Calendar) SetYear(y int) *Calendar {
	c.year = y
	return c
}

func (c *Calendar) SetDay(d int) *Calendar {
	c.day = d
	return c
}

func (c *Calendar) SetSelected(rc *rc) (*Calendar, bool) {
	c.SelectR = rc
	day := c.rows[rc.Row][rc.Col]
	if strings.Trim(day, " ") == "" {
		return c, false
	}
	c.day, _ = strconv.Atoi(day)
	return c, true
}

func (c *Calendar) Date() string {
	loc := time.Now().Location()
	return time.Date(c.year, c.month, c.day, 0, 0, 0, 0, loc).Format("01-02-2006")
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func getDaysOffset(y int, m time.Month) int {
	now := time.Now()
	return int(time.Date(y, m, 1, 0, 0, 0, 0, now.Location()).Weekday())
}
