package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/index", index)
	e.POST("/text", readText)

	e.Logger.Fatal(e.Start(":1323"))
}

func parse(textarea string) (res string, err error) {
	analyzes := []analyze{}
	parts := []Part{}
	lines := strings.Split(textarea, "\n")

	//分区
	for i := 0; i < len(lines); {
		if area, b := containsAreas(lines[i]); b {
			ana := analyze{area: area}
			for i++; i < len(lines); i++ {
				if _, b2 := containsAreas(lines[i]); b2 {
					break
				}
				ana.lines = append(ana.lines, lines[i])
			}
			analyzes = append(analyzes, ana)

		} else {
			i++
		}

	}

	fmt.Println(textarea, analyzes)

	for _, analyze := range analyzes {
		part := Part{area: analyze.area}

		daysold, err := parseOneFloat(analyze.lines, []string{"日销售", "日销"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.DaySold = daysold

		monthsold, err := parseOneFloat(analyze.lines, []string{"月累售", "月累计"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.MonthSold = monthsold

		turnOverRate, err := parsethreeFloat(analyze.lines, []string{"成交率"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.TurnoverRate = turnOverRate

		bookingRate, err := parseTwoFloat(analyze.lines, []string{"预约金额占比", "约销售占比"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.BookingRate = bookingRate

		comeOfBookRate, err := parsethreeFloat(analyze.lines, []string{"预约到店率"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.ComeOfBookRate = comeOfBookRate

		turnOfbookRate, err := parsethreeFloat(analyze.lines, []string{"到店成交率"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.TurnOfbookRate = turnOfbookRate

		DayFansAddition, err := parseOneFloat(analyze.lines, []string{"今日新增微粉", "今日新增潜客"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.DayFansAddition = DayFansAddition

		MonthFansAddition, err := parseOneFloat(analyze.lines, []string{"本月新增微粉", "月累新增潜客"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.MonthFansAddition = MonthFansAddition

		MonthFansAndTurnOfRate, err := parseTwoFloat(analyze.lines, []string{"本月微粉成交单数"})
		if err != nil {
			fmt.Printf("parseOneFloat error:%s", err.Error())
			return "", err
		}
		part.Detail.MonthFansAndTurnOfRate = MonthFansAndTurnOfRate

		parts = append(parts, part)
	}

	//解析具体数据
	var daySold, monthSold, CmlRate, diff, turnoverCml, turnoverGoal, booking1, booking2, bookingcome1, bookingcome2, monthfans1, monthfas2, dayfansAddtion, monthfansAddtion float64
	var  bookingcometurnover1 ,bookingcometurnover2 float64
	for _, part := range parts {
		detail := part.Detail
		daySold += detail.DaySold
		monthSold += detail.MonthSold
		turnoverCml += detail.TurnoverRate.Num1
		turnoverGoal += detail.TurnoverRate.Num2
		booking1 += detail.BookingRate.Num1
		booking2 += detail.BookingRate.Num2
		bookingcome1 += detail.ComeOfBookRate.Num1
		bookingcome2 += detail.ComeOfBookRate.Num2
		bookingcometurnover1+=detail.TurnOfbookRate.Num1
		bookingcometurnover2+=detail.TurnOfbookRate.Num2
		monthfans1 += detail.MonthFansAndTurnOfRate.Num1
		monthfas2 += detail.MonthFansAndTurnOfRate.Num2
		dayfansAddtion += detail.DayFansAddition
		monthfansAddtion += detail.MonthFansAddition
	}
	CmlRate = monthSold / goal

	now := time.Now()
	end := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	end = end.AddDate(0, 0, -1)

	diff = CmlRate-(float64(now.Day()) / float64(end.Day()))

	//统计数据

	report := Report{
		Date:                    time.Now().Format("2006-01-02"),                                                         //日期
		MonthGoal:               fmt.Sprintf("%+v亿", goal/10000),                                                         //月目标
		DaySell:                 fmt.Sprintf("%+v万(%.2f%%)", (daySold), daySold/goal*100),                                    //日销售
		MonthSell:               fmt.Sprintf("%.2f万", monthSold),                                                         //月销售
		CmlRate:                 fmt.Sprintf("%.2f%%", CmlRate*100),                                                      // string //完成率
		Diff:                    fmt.Sprintf("%.2f%%", diff),                                                             //string //差异
		TurnOverRate:            fmt.Sprintf("%.2f%% (%d/%d)", turnoverCml/turnoverGoal*100, int(turnoverCml),int(turnoverGoal) ), //string //成交率
		BookingMonyRate:         fmt.Sprintf("%.2f%% /%.0f万", booking2/monthSold*100, booking2),                               //string //预约金额占比
		BookingComeRate:         fmt.Sprintf("%.2f%%（%.0f/%.0f）", bookingcome1/bookingcome2*100,bookingcome1,bookingcome2),                                        //string //预约到店率
		BookingComeTurnOverRate: fmt.Sprintf("%.2f%%(%.0f/%.0f)", bookingcometurnover1/bookingcometurnover2*100,bookingcometurnover1,bookingcometurnover2),      //预约到店成交率                                   //预约到店成交率
		DayFansAddition:         fmt.Sprintf("%d", int(dayfansAddtion)),                                                  // string //今日新增微粉
		MonthFansAddition:       fmt.Sprintf("%d", int(monthfansAddtion)),                                                //string //本月新增微粉
		MonthFansAndMony:        fmt.Sprintf("%+v/%+v  (占比 %.2f%%)", monthfans1, monthfas2,monthfas2/monthSold*100),                                          //string //本月微粉成交单数/金额
	}

	res = printPart(parts) + "\n\n" + printReport(report)

	return res, nil

}

func index(c echo.Context) error {

	//dat, err := ioutil.ReadFile("./index.html")
	//if err!=nil{
	//	return c.HTML(http.StatusOK, "read index.html error")
	//}
	return c.HTML(http.StatusOK, htmlmsg)
}

func readText(c echo.Context) error {

	texts := c.FormValue("area")

	res, err := parse(texts)
	if err != nil {
		fmt.Println("error", err.Error())
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, res)

}

const goal = 19400.0 //万

type Rate struct {
	Num1 float64
	Num2 float64
}

type Detail struct {
	DaySold                float64 //日销售
	MonthSold              float64 //月销售
	TurnoverRate           Rate    //成交率
	BookingRate            Rate    //预约金占比
	ComeOfBookRate         Rate    //预约到店率
	TurnOfbookRate         Rate    //预约到店成交率
	DayFansAddition        float64 //日增加粉丝
	MonthFansAddition      float64 //月增加粉丝
	MonthFansAndTurnOfRate Rate    //本月微粉成交单数/金额

}

type Report struct {
	Date                    string //日期
	MonthGoal               string //全月目标
	DaySell                 string //日销售
	MonthSell               string //本月累计
	CmlRate                 string //完成率
	Diff                    string //差异
	TurnOverRate            string //成交率
	BookingMonyRate         string //预约金额占比
	BookingComeRate         string //预约到店率
	BookingComeTurnOverRate string //  预约到店成交率
	DayFansAddition         string //今日新增微粉
	MonthFansAddition       string //本月新增微粉
	MonthFansAndMony        string //本月微粉成交单数/金额
}

type analyze struct {
	area  string
	lines []string
	//Detail Detail
}

type Part struct {
	area   string
	Detail Detail
}

var areas = []string{"东北区", "华东二区", "华南区", "西南区", "华北区", "西北区", "华东一区", "华中区"}

func parseTwoFloat(lines []string, contains []string) (Rate, error) {
	var res Rate

	for _, line := range lines {

		for _, contain := range contains {
			if strings.Contains(line, contain) {
				re := regexp.MustCompile(`[\d.]+`)
				nstrs := re.FindAllString(line, 2)

				number, err := strconv.ParseFloat(nstrs[0], 64)
				if err != nil {
					return res, err
				}
				res.Num1 = number

				number, err = strconv.ParseFloat(nstrs[1], 64)
				if err != nil {
					return res, err
				}
				res.Num2 = number

				return res, nil
			}
		}

	}

	return res, nil
}

func parsethreeFloat(lines []string, contains []string) (Rate, error) {
	var res Rate

	for _, line := range lines {

		for _, contain := range contains {
			if strings.Contains(line, contain) {
				re := regexp.MustCompile(`[\d.]+`)
				nstrs := re.FindAllString(line, 3)

				number, err := strconv.ParseFloat(nstrs[1], 64)
				if err != nil {
					return res, err
				}
				res.Num1 = number

				number, err = strconv.ParseFloat(nstrs[2], 64)
				if err != nil {
					return res, err
				}
				res.Num2 = number

				return res, nil
			}
		}

	}

	return res, nil
}

func parseOneFloat(lines []string, contains []string) (float64, error) {
	var res float64

	for _, line := range lines {

		for _, contain := range contains {
			if strings.Contains(line, contain) {
				re := regexp.MustCompile(`[\d.]+`)
				nstrs := re.FindStringSubmatch(line)

				number, err := strconv.ParseFloat(nstrs[0], 64)
				if err != nil {
					return res, err
				}

				return number, nil
			}
		}

	}

	return res, nil
}

func printPart(parts []Part) string {
	res := ""

	for _, part := range parts {

		detail := part.Detail
		partStr := fmt.Sprintf("区域:%s     \n\t日销售:%+v   \n\t月销售:%+v  \n\t成交率:%+v      \n\t预约金占比:%+v    \n\t预约到店率:%+v      \n\t预约到店成交率:%+v  \n\t日增加粉丝:%+v        \n\t月增加粉丝:%+v        \n\t本月微粉成交单数/金额:%+v\n\t ",
			part.area, detail.DaySold, detail.MonthSold, detail.TurnoverRate, detail.BookingRate, detail.ComeOfBookRate, detail.TurnOfbookRate, detail.DayFansAddition, detail.MonthFansAddition, detail.MonthFansAndTurnOfRate)
		res += "\n" + partStr
	}

	return res

}

func printReport(report Report) string {
	res := fmt.Sprintf("全国%s       \n\t全月目标：%s  \n\t日销售：%s      \n\t本月累计：%s        \n\t完成率：%s       \n\t差异：%s          \n\t成交率：%s       \n\t预约金额占比：%s      \n\t预约到店率:%s  \n\t预约到店成交率：%s  \n\t今日新增微粉：%s     \n\t本月新增微粉：%s       \n\t本月微粉成交单数/金额：%s\n\t",
		report.Date, report.MonthGoal, report.DaySell, report.MonthSell, report.CmlRate, report.Diff, report.TurnOverRate, report.BookingMonyRate, report.BookingComeRate,report.BookingComeTurnOverRate, report.DayFansAddition, report.MonthFansAddition, report.MonthFansAndMony)
	return res
}

func containsAreas(line string) (string, bool) {
	for _, area := range areas {
		if strings.Contains(line, area) {
			return area, true
		}
	}

	return "", false
}

var htmlmsg = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Single file upload</title>
</head>
<body>
<h1>请佳佳把信息粘贴到这里</h1>

<form action="/text" method="post" enctype="multipart/form-data">
   <textarea rows="10" cols="30" name="area">

</textarea>
    <input type="submit" value="Submit">
</form>
</body>
</html>`
