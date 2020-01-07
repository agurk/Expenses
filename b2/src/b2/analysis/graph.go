package analysis

import (
	"b2/fxrates"
	"database/sql"
	"fmt"
	"math"
)

type graphParams struct {
	ccy            string
	canvasMaxX     int64
	canvasMaxY     int64
	padding        int64
	xIncrement     int64
	amountMaximum  float64
	maxX           int64
	incNonDayBox   bool
	nonDayBoxStyle string
	axisStyle      string
	areas          []*area
	lines          []*line
	from           string
	to             string
	periodDays     int
	lookbackPeriod int
	minYIncrement  int64
}

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
	params.xIncrement = (params.canvasMaxX - params.padding) / 31
	// todo: 31 is period days, and 12 is months
	params.periodDays = 31
	params.lookbackPeriod = 12
	params.minYIncrement = 100
	return params
}

func graph(params *graphParams, fx *fxrates.FxValues, db *sql.DB) (string, error) {
	cumulative, sdData, err := averageSpend(params, fx, db)
	sd(cumulative, sdData, params)
	if err != nil {
		return "", err
	}
	addLine(cumulative, "rgb(165, 165, 165)", 4, false, params)
	cs, err := cumulativeSpend(params, fx, db)
	if err != nil {
		return "", err
	}
	addLine(cs, "rgb(165,0,0)", 20, false, params)
	svg := fmt.Sprintf("<svg viewBox=\"%d %d %d %d\">", params.padding*-2,
		params.padding*-1,
		params.canvasMaxX+3*params.padding,
		params.canvasMaxY+2*params.padding)
	svg += makeAreas(params)
	for _, line := range params.lines {
		svg += buildLine(line, params)
	}
	svg += axis(params)
	svg += "</svg>"
	return svg, nil
}

func axis(params *graphParams) string {
	svg := fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="0" %s/>`, params.padding, params.canvasMaxY+params.padding, params.padding, params.axisStyle)
	svg += fmt.Sprintf(`<line x1="0" y1="%d" x2="%d" y2="%d" %s />`, params.canvasMaxY, params.canvasMaxX, params.canvasMaxY, params.axisStyle)
	for i := 1; i <= params.periodDays; i++ {
		xPos := (int64(i) * params.xIncrement) + params.padding
		yPos := params.canvasMaxY + (params.padding / 3)
		svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" %s />`, xPos, yPos, xPos, params.canvasMaxY, params.axisStyle)
		svg += fmt.Sprintf(`<text x="%d" y=%d font-size="80" text-anchor="middle">%d</text>`, xPos, yPos+60, i)
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
		svg += fmt.Sprintf(`<line x1="%f" y1="%f" x2="%d" y2="%f" %s />`, xPos, yPos, params.padding, yPos, params.axisStyle)
		svg += fmt.Sprintf(`<text x="%f" y="%f" font-size="80" text-anchor="end" dominant-baseline="middle">%d</text>`, xPos, yPos, amount)
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
	xPos := params.padding
	line := `<polyline points="`
	for _, yPos := range points {
		line += fmt.Sprintf(" %d %d", xPos, yPos)
		xPos += params.xIncrement
	}
	line += fmt.Sprintf(`" stroke="%s" stroke-width="%d" stroke-linecap="square" fill="none" stroke-linejoin="round"/>`, l.colour, l.stroke)
	return line
}

func addLine(points []float64, colour string, stroke int64, extrapolate bool, params *graphParams) {
	l := new(line)
	l.extrapolate = extrapolate
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
		xPos := params.padding
		areas += `<polygon points="`
		for _, yPos := range *area.e0 {
			areas += fmt.Sprintf("%d, %d ", xPos, yPos)
			xPos += params.xIncrement
		}
		// todo: could be wrong here
		for i := range *area.e1 {
			i = len(*area.e1) - 1 - i
			yPos := (*area.e1)[i]
			xPos -= params.xIncrement
			areas += fmt.Sprintf("%d, %d ", xPos, yPos)
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
	means := makeSdSlice(params.periodDays+1, params.canvasMaxY)
	sdUp := makeSdSlice(params.periodDays+1, params.canvasMaxY)
	sdDown := makeSdSlice(params.periodDays+1, params.canvasMaxY)
	twosdUp := makeSdSlice(params.periodDays+1, params.canvasMaxY)
	twosdDown := makeSdSlice(params.periodDays+1, params.canvasMaxY)
	yFactor := float64(params.canvasMaxY) / params.amountMaximum
	for i := range *means {
		(*means)[i] = int64((params.amountMaximum - math.Abs(average[i])) * yFactor)
	}
	for i := range *means {
		sdi := int64(sd[i] * yFactor)
		(*sdUp)[i] = (*means)[i] - sdi
		(*twosdUp)[i] = (*means)[i] - 2*sdi
		(*sdDown)[i] = (*means)[i] + sdi
		(*twosdDown)[i] = (*means)[i] + sdi*2
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

func cumulativeSpend(params *graphParams, fx *fxrates.FxValues, db *sql.DB) (points []float64, err error) {
	rows, err := getCumulativeData(params, db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points = make([]float64, params.periodDays+1)
	var localMaxX int64 = 0
	for rows.Next() {
		var amount float64
		var ccy, date string
		var day int64
		err = rows.Scan(&amount, &ccy, &date, &day)
		if err != nil {
			return nil, err
		}
		// todo: better date handling
		date = date[:10]
		rate, err := fx.Get(date, params.ccy, ccy)
		if err != nil {
			return nil, err
		}
		points[day] += amount / rate
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

func getCumulativeData(params *graphParams, db *sql.DB) (*sql.Rows, error) {
	return db.Query(`
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
			and cd.isexpense
			and strftime(date) >= date($1,'start of month')
			and strftime(date) < date($2,'start of month','+1 month')`,
		params.from, params.to)
}

func averageSpend(params *graphParams, fx *fxrates.FxValues, db *sql.DB) (cumulative, sd []float64, err error) {
	averageSpend := make([][]float64, params.lookbackPeriod+1)
	spends := make([]float64, (params.periodDays+1)*(params.lookbackPeriod+1))
	for i := range averageSpend {
		averageSpend[i], spends = spends[:params.periodDays+1], spends[params.periodDays+1:]
	}

	cumulative = make([]float64, params.periodDays+1)
	sd = make([]float64, params.periodDays+1)

	rows, err := getData(params, db)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var amount float64
		var day, month, year int
		var ccy string
		err = rows.Scan(&amount, &day, &month, &year, &ccy)
		if err != nil {
			return nil, nil, err
		}
		date := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		rate, err := fx.Get(date, params.ccy, ccy)
		if err != nil {
			return nil, nil, err
		}
		averageSpend[month][day] += amount / rate
	}
	for day := 1; day <= params.periodDays; day++ {
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

func getData(params *graphParams, db *sql.DB) (*sql.Rows, error) {
	return db.Query(`
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
		where
			date(e.date) < date($1,'start of month')
			and date(e.date) > date($1,'start of month','-12 months')
			and e.eid = c.eid
			and c.cid = cd.cid
			and cd.isexpense
		group by
			day,
			month,
			ccy`,
		params.from, params.to)
}
