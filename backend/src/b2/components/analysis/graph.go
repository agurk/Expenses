package analysis

import (
	"b2/errors"
	"b2/moneyutils"
	"database/sql"
	"fmt"
	"math"
	"time"
)

type graphParams struct {
	ccy            string
	canvasMaxX     int64
	canvasMaxY     int64
	padding        int64
	xIncrement     float64
	amountMaximum  float64
	maxX           int64
	incNonDayBox   bool
	nonDayBoxStyle string
	axisStyle      string
	areas          []*area
	lines          []*line
	from           string
	to             string
	periodLength   int
	periodType     int
	lookbackPeriod int
	minYIncrement  int64
}

const (
	monthperiod = 1
	yearperiod  = 2
	rangeperiod = 3
)

type area struct {
	e0     *[]int64
	e1     *[]int64
	colour string
}

type line struct {
	points            []float64
	colour            string
	stroke            int64
	extrapolate       bool
	extrapolateStroke string
}

func gInitialise(p *totalsParams) *graphParams {
	params := new(graphParams)
	params.ccy = p.CCY
	params.from = p.From
	params.to = p.To
	params.canvasMaxX = 4750
	params.canvasMaxY = 2500
	params.padding = 100
	params.amountMaximum = 2000
	params.maxX = 0
	params.incNonDayBox = true
	params.nonDayBoxStyle = `style="fill:rgb(20, 20, 20)" fill-opacity="0.3"`
	params.axisStyle = `style="stroke:rgb(0,0,0);stroke-width:10"`
	params.minYIncrement = 100
	// todo: 31 is period days, and 12 is months
	switch periodType(params) {
	case monthperiod:
		params.periodLength = 31
		params.lookbackPeriod = 12
		params.periodType = monthperiod
	case yearperiod:
		params.periodLength = 366
		params.lookbackPeriod = 1
		params.periodType = yearperiod
	}
	params.xIncrement = float64(params.canvasMaxX-params.padding) / float64(params.periodLength)
	return params
}

// periodType will decide what type of graph to be rendered - a single month view, a year view
// or an arbitary range
func periodType(params *graphParams) int {
	// todo deal with errors
	from, _ := time.Parse("2006-01-02", params.from)
	to, _ := time.Parse("2006-01-02", params.to)
	if from.Month() == to.Month() && from.Year() == to.Year() {
		return monthperiod
	}
	return yearperiod
}

func graph(params *graphParams, fx *moneyutils.FxValues, db *sql.DB) (string, error) {
	cumulative, sdData, err := averageSpend(params, fx, db)
	if err != nil {
		return "", errors.Wrap(err, "analysis.graph")
	}
	cs, err := cumulativeSpend(params, fx, db)
	if err != nil {
		return "", errors.Wrap(err, "analysis.graph")
	}
	addLine(cs, "rgb(165,0,0)", 20, true, "20, 20", params)
	if params.periodType == monthperiod {
		sd(cumulative, sdData, params)
	}
	addLine(cumulative, "rgb(165, 165, 165)", 4, false, "", params)
	svg := fmt.Sprintf("<svg viewBox=\"%d %d %d %d\">", params.padding*-2,
		params.padding*-1,
		params.canvasMaxX+3*params.padding,
		params.canvasMaxY+2*params.padding)
	svg += makeAreas(params)
	for _, line := range params.lines {
		svg += buildLine(line, params)
	}
	if params.incNonDayBox {
		svg += buildNonDayBox(params)
	}
	svg += axis(params)
	svg += "</svg>"
	return svg, nil
}

func buildNonDayBox(params *graphParams) string {
	tStart, _ := time.Parse("2006-01-02", params.from)
	tStop, _ := time.Parse("2006-01-02", params.to)
	days := (tStop.Sub(tStart).Hours() / 24) + 1
	height := params.canvasMaxY + params.padding
	width := int(params.xIncrement) * (params.periodLength - int(days))
	xPos := params.canvasMaxX - int64(width)
	yPos := 0 - params.padding
	return fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" %s/>`, xPos, yPos, width, height, params.nonDayBoxStyle)
}

func axis(params *graphParams) string {
	svg := fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="0" %s/>`, params.padding, params.canvasMaxY+params.padding, params.padding, params.axisStyle)
	svg += fmt.Sprintf(`<line x1="0" y1="%d" x2="%d" y2="%d" %s />`, params.canvasMaxY, params.canvasMaxX, params.canvasMaxY, params.axisStyle)
	switch params.periodType {
	case monthperiod:
		tStart, _ := time.Parse("2006-01-02", params.from)
		tStop, _ := time.Parse("2006-01-02", params.to)
		days := (tStop.Sub(tStart).Hours() / 24) + 1
		fill := "black"
		for i := 1; i <= params.periodLength; i++ {
			xPos := float64(i)*params.xIncrement + float64(params.padding)
			yPos := params.canvasMaxY + (params.padding / 3)
			svg += fmt.Sprintf(`<line x1="%f" y1="%d" x2="%f" y2="%d" %s />`, xPos, yPos, xPos, params.canvasMaxY, params.axisStyle)
			if i > int(days) {
				fill = "grey"
			}
			xPos -= params.xIncrement * 0.5
			svg += fmt.Sprintf(`<text x="%f" y=%d font-size="80" text-anchor="middle" fill="%s">%d</text>`, xPos, yPos+60, fill, i)
		}
	case yearperiod:
		months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		for i, month := range months {
			xPos := float64(i*31)*params.xIncrement + float64(params.padding)
			yPos := params.canvasMaxY + params.padding/3
			svg += fmt.Sprintf(`<line x1="%f" y1="%d" x2="%f" y2="%d" %s />`, xPos, yPos, xPos, params.canvasMaxY, params.axisStyle)
			xPos += float64(15) * params.xIncrement
			// todo: Magic number makes tails of chars visible. Should be dealt with better
			yPos -= 12
			svg += fmt.Sprintf(`<text x="%f" y=%d font-size="80" text-anchor="middle">%s</text>`, xPos, yPos+60, month)
		}
	}
	increment := int64(math.Pow(10, math.Round(math.Log10(params.amountMaximum)-1)))
	if increment < params.minYIncrement {
		increment = params.minYIncrement
	}
	var amount int64 = 0
	yFactor := float64(params.canvasMaxY) / params.amountMaximum
	for amount <= int64(params.amountMaximum) {
		amount += increment
		yPos := float64(params.canvasMaxY) - (float64(amount) * yFactor)
		xPos := (1 / 3) * params.padding
		svg += fmt.Sprintf(`<line x1="%d" y1="%f" x2="%d" y2="%f" %s />`, xPos, yPos, params.padding, yPos, params.axisStyle)
		svg += fmt.Sprintf(`<text x="%d" y="%f" font-size="80" text-anchor="end" dominant-baseline="middle">%d</text>`, xPos, yPos, amount)
	}
	return svg
}

func buildLine(l *line, params *graphParams) string {
	yFactor := float64(params.canvasMaxY) / params.amountMaximum
	maxRange := len(l.points)
	points := make([]int64, maxRange)
	points[0] = params.canvasMaxY
	for i := 1; i < maxRange; i++ {
		points[i] = int64((params.amountMaximum - math.Abs(l.points[i])) * yFactor)
	}
	xPos := float64(params.padding)
	line := `<polyline points="`
	for _, yPos := range points {
		line += fmt.Sprintf(" %f %d", xPos, yPos)
		xPos += params.xIncrement
	}
	line += fmt.Sprintf(`" stroke="%s" stroke-width="%d" stroke-linecap="square" fill="none" stroke-linejoin="round"/>`, l.colour, l.stroke)
	if l.extrapolate {
		line += extrapolateLine(l, params)
	}
	return line
}

func extrapolateLine(l *line, params *graphParams) string {
	tStart, _ := time.Parse("2006-01-02", params.from)
	tStop, _ := time.Parse("2006-01-02", params.to)
	tNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, tStart.Location())

	complete := false

	// if all of the line is in the future
	if tStart.Sub(tNow).Hours()/24 >= 0 {
		return ""
	} else if tStop.Sub(tNow).Hours()/24 < 0 {
		complete = true
	}

	val := tStop.Sub(tStart).Hours()/24 + 2
	yFactor := float64(params.canvasMaxY) / params.amountMaximum
	y := int64((params.amountMaximum - math.Abs(l.points[len(l.points)-1])) * yFactor)
	x1 := params.padding + int64(params.xIncrement)*int64(len(l.points)-1)

	var line string
	if complete {
		x2 := x1 + int64(params.xIncrement)*int64(int(val)-len(l.points))
		line = fmt.Sprintf(`<line stroke="%s" stroke-width="%d" x1="%d" y1="%d" x2="%d" y2="%d" />`,
			l.colour, l.stroke, x1, y, x2, y)
	} else {
		val -= tStop.Sub(tNow).Hours() / 24
		x2 := x1 + int64(params.xIncrement)*int64(int(val)-len(l.points))
		line = fmt.Sprintf(`<line stroke="%s" stroke-width="%d" stroke-dasharray="%s" x1="%d" y1="%d" x2="%d" y2="%d" />`,
			l.colour, l.stroke, l.extrapolateStroke, x1, y, x2, y)
	}

	return line
}

func addLine(points []float64, colour string, stroke int64, extrapolate bool, extrapolateStroke string, params *graphParams) {
	l := new(line)
	l.extrapolate = extrapolate
	l.extrapolateStroke = extrapolateStroke
	l.stroke = stroke
	l.colour = colour
	l.points = points
	for _, value := range points {
		if math.Abs(value) > params.amountMaximum {
			params.amountMaximum = math.Abs(value)
		}
	}
	params.lines = append(params.lines, l)
}

func addArea(up, down *[]int64, colour string, params *graphParams) {
	a := new(area)
	a.e0 = up
	a.e1 = down
	a.colour = colour
	params.areas = append(params.areas, a)
}

func makeAreas(params *graphParams) string {
	areas := ""
	for _, area := range params.areas {
		xPos := float64(params.padding)
		areas += `<polygon points="`
		for _, yPos := range *area.e0 {
			areas += fmt.Sprintf("%f, %d ", xPos, yPos)
			xPos += params.xIncrement
		}
		for i := range *area.e1 {
			i = len(*area.e1) - 1 - i
			yPos := (*area.e1)[i]
			xPos -= params.xIncrement
			areas += fmt.Sprintf("%f, %d ", xPos, yPos)
		}
		areas += fmt.Sprintf("\" fill=\"%s\" stroke-width=\"0\" />", area.colour)
	}
	return areas
}

func makeSdSlice(size int, val int64) *[]int64 {
	slice := make([]int64, size)
	for i := range slice {
		slice[i] = val
	}
	return &slice
}

func sd(average, sd []float64, params *graphParams) {
	means := makeSdSlice(params.periodLength+1, params.canvasMaxY)
	sdUp := makeSdSlice(params.periodLength+1, params.canvasMaxY)
	sdDown := makeSdSlice(params.periodLength+1, params.canvasMaxY)
	twosdUp := makeSdSlice(params.periodLength+1, params.canvasMaxY)
	twosdDown := makeSdSlice(params.periodLength+1, params.canvasMaxY)
	yFactor := float64(params.canvasMaxY) / params.amountMaximum
	for i := range *means {
		(*means)[i] = int64((params.amountMaximum - math.Abs(average[i])) * yFactor)
	}
	for i := range *means {
		sdi := int64(sd[i] * yFactor)
		(*sdUp)[i] = (*means)[i] - sdi
		(*twosdUp)[i] = (*means)[i] - 2*sdi
		(*sdDown)[i] = (*means)[i] + sdi
		(*twosdDown)[i] = (*means)[i] + 2*sdi
		if (*sdDown)[i] > params.canvasMaxY {
			(*sdDown)[i] = params.canvasMaxY
		}
		if (*twosdDown)[i] > params.canvasMaxY {
			(*twosdDown)[i] = params.canvasMaxY
		}
	}
	addArea(twosdUp, twosdDown, "rgb(240, 240, 240)", params)
	addArea(sdUp, sdDown, "rgb(225, 225, 225)", params)
}

func cumulativeSpend(params *graphParams, fx *moneyutils.FxValues, db *sql.DB) (points []float64, err error) {
	rows, err := cumulativeData(params, db)
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "analysis.cumulativeSpend (cumulative data db err)")
	}

	points = make([]float64, params.periodLength+1)
	var localMaxX int64
	for rows.Next() {
		var amount, day int64
		var ccy, date string
		err = rows.Scan(&amount, &ccy, &date, &day)
		if err != nil {
			return nil, errors.Wrap(err, "analysis.cumulativeSpend (rows scan)")
		}
		// todo: better date handling
		date = date[:10]
		rate, err := fx.Rate(date, params.ccy, ccy)
		if err != nil {
			return nil, errors.Wrap(err, "analysis.cumulativeSpend (rates err)")
		}
		ccyAmt, err := moneyutils.CurrencyAmount(amount, ccy)
		if err != nil {
			return nil, errors.Wrap(err, "analysis.cumulativeSpend (ccy amount)")
		}
		day = getDay(day, date, params)
		points[day] += ccyAmt / rate
		if day > localMaxX {
			localMaxX = day
		}
	}
	if localMaxX > params.maxX {
		params.maxX = localMaxX
	}
	var i int64
	for i = 1; i <= localMaxX; i++ {
		points[i] += points[i-1]
	}
	return points[:localMaxX+1], nil
}

func getDay(day int64, date string, params *graphParams) int64 {
	switch params.periodType {
	case monthperiod:
		return day
	case yearperiod:
		// todo deal with err
		t, _ := time.Parse("2006-01-02", date)
		return int64(t.YearDay())
	}
	return 0
}

// cumulativeData get that data used for drawing the main spend line
func cumulativeData(params *graphParams, db *sql.DB) (*sql.Rows, error) {
	query := `
		select
			amount,
			ccy,
			date,
			strftime('%d', e.date) day
		from
			expenses e,
			classifications c,
			classificationdef cd
		where
			e.eid = c.eid
			and c.cid = cd.cid
			and cd.isexpense `
	switch params.periodType {
	case monthperiod:
		query += `and strftime(date) >= date($1,'start of month') and strftime(date) < date($2,'start of month','+1 month')`
	case yearperiod:
		query += `and strftime(date) >= date($1,'start of year') and strftime(date) < date($2,'start of year','+12 month')`
	default:
		query += `and strftime(date) >= date($1) and strftime(date) < date($2)`
	}
	return db.Query(query, params.from, params.to)
}

func averageSpend(params *graphParams, fx *moneyutils.FxValues, db *sql.DB) (cumulative, sd []float64, err error) {
	averageSpend := make([][]float64, params.lookbackPeriod+1)
	spends := make([]float64, (params.periodLength+1)*(params.lookbackPeriod+1))
	for i := range averageSpend {
		averageSpend[i], spends = spends[:params.periodLength+1], spends[params.periodLength+1:]
	}

	cumulative = make([]float64, params.periodLength+1)
	sd = make([]float64, params.periodLength+1)

	rows, err := data(params, db)
	defer rows.Close()
	if err != nil {
		return nil, nil, errors.Wrap(err, "analysis.averageSpend")
	}

	for rows.Next() {
		var amount int64
		var day, month, year int
		var ccy string
		err = rows.Scan(&amount, &day, &month, &year, &ccy)
		if err != nil {
			return nil, nil, errors.Wrap(err, "analysis.averageSpend")
		}
		date := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		rate, err := fx.Rate(date, params.ccy, ccy)
		if err != nil {
			return nil, nil, errors.Wrap(err, "analysis.averageSpend")
		}
		ccyAmt, err := moneyutils.CurrencyAmount(amount, ccy)
		if err != nil {
			errors.Wrap(err, "analysis.averageSpend")
		}
		lookback, day := lookbackDay(day, month, year, params)
		averageSpend[lookback][day] += ccyAmt / rate
	}
	for day := 1; day <= params.periodLength; day++ {
		for i := 1; i <= params.lookbackPeriod; i++ {
			cumulative[day] += math.Abs(averageSpend[i][day])
			averageSpend[i][day] += averageSpend[i][day-1]
		}
		cumulative[day] = cumulative[day] / float64(params.lookbackPeriod)
		cumulative[day] += cumulative[day-1]
		for year := 1; year <= params.lookbackPeriod; year++ {
			sd[day] += math.Pow((math.Abs(averageSpend[year][day]) - cumulative[day]), 2)
		}
		sd[day] = math.Sqrt(sd[day] / float64(params.lookbackPeriod))
		if math.Abs(cumulative[day]) > params.amountMaximum {
			params.amountMaximum = math.Abs(cumulative[day])
		}
	}
	return
}

func lookbackDay(day, month, year int, params *graphParams) (int, int) {
	switch params.periodType {
	case monthperiod:
		return month, day
	case yearperiod:
		// todo deal with err
		t, _ := time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-%02d", year, month, day))
		return 1, t.YearDay()
	}
	return 0, 0
}

func data(params *graphParams, db *sql.DB) (*sql.Rows, error) {
	query := `
		select
			sum (e.amount),
			strftime('%d', e.date) day,
			strftime('%m', e.date) month,
			strftime('%Y', e.date) year,
			ccy
		from
			expenses e,
			classifications c,
			classificationdef cd
		where `
	switch params.periodType {
	case monthperiod:
		query += `date(e.date) < date($1,'start of month') and date(e.date) > date($1,'start of month','-12 months') `
	case yearperiod:
		query += `date(e.date) < date($1,'start of year') and date(e.date) >= date($2,'start of year','-12 months') `
	default:
		query += `date(e.date) =< date($2) and date(e.date) >= date($1) `
	}
	query += `
			and e.eid = c.eid
			and c.cid = cd.cid
			and cd.isexpense
		group by
			day,
			month,
			ccy`
	return db.Query(query, params.from, params.to)
}
