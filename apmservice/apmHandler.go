package apmservice

import (
	"Golangcode/config"
	"Golangcode/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// To Test whether the GoCode is in running state/not
func test(w http.ResponseWriter, r *http.Request) {
	body := config.Body{ResponseCode: 200, Message: "OK"}
	jsonBody, err := json.Marshal(body)
	//if there is error in converting to json marshaling the below response will be sent
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//if no error then sent the response back to the port
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
func compareResults(writer http.ResponseWriter, request *http.Request) {
	//get the values from the parameters in the POST request to variables
	//http://localhost:6666/compareResults?oldBuildNum=2&oldRelease=202301&
	//newBuildNum=3&newRelease=202302&email=asdfg&metric=95th&index=jmeter-aggregate-jtl&
	//Scenario=Dummy&Iteration=1
	oldBuildNum := request.URL.Query().Get("oldBuildNum")
	oldReleaseNumber := request.URL.Query().Get("oldReleaseNumber")
	oldreleaseIteration := request.URL.Query().Get("oldRelease_Iteration")
	newReleaseNumber := request.URL.Query().Get("newReleaseNumber")
	newBuildNum := request.URL.Query().Get("newBuildNum")
	newreleaseIteration := request.URL.Query().Get("newRelease_Iteration")
	emailID := request.URL.Query().Get("email")
	metric := request.URL.Query().Get("metric")     //Metric options may be 95th,99th,90th,Average
	selectIndex := request.URL.Query().Get("index") //Mention the ES Index from where results needs to be obtained

	log.Println("OldBuildNUm is: ", oldBuildNum)
	log.Println("OldRelease is: ", oldReleaseNumber)
	log.Println("OldRelease_Iteration is: ", oldreleaseIteration)
	log.Println("newBuildNUm is: ", newBuildNum)
	log.Println("newRelease is: ", newReleaseNumber)
	log.Println("newRelease_Iteration is: ", newreleaseIteration)
	log.Println("emailID is: ", emailID)
	log.Println("metric is: ", metric)
	oldReleaseData := utils.GetReleaseData(oldreleaseIteration, oldBuildNum, metric, selectIndex, oldReleaseNumber)
	newReleaseData := utils.GetReleaseData(newreleaseIteration, newBuildNum, metric, selectIndex, newReleaseNumber)

	if len(oldReleaseData) == 0 || len(newReleaseData) == 0 {
		utils.RespondWithJSON("Error Select correct Data", writer, request)
	}
	//log.Println(oldReleaseData)
	///log.Println(newReleaseData)
	if len(oldReleaseData) != 0 && len(newReleaseData) != 0 {
		pnewapis := fmt.Sprintf("")
		if len(newReleaseData) > len(oldReleaseData) {
			pnewapis = fmt.Sprintf("<table style='backgound:#99ebff;border-collapse: collapse;' border = '2' cellpadding = '6'><tbody><tr><td colspan=2 style='text-align:center;background-color:Lavender;color:Black;'><b>%s Response Time of Newly added usecases in Regression Suite</b></td></tr><tr><th>API</th><th>Release: %s (in ms) </th></tr> ", metric, newReleaseNumber)
		}
		subject := fmt.Sprintf("Release Comparison Report for %s (%s) & %s (%s)", oldReleaseNumber, oldBuildNum, newReleaseNumber, newBuildNum)
		//declaring the header to be used in html report
		p := fmt.Sprintf("<body style='background:White'><h3 style='background:#0790bd;color:#fff;padding:5px;text-align:center;border-radius:5px;'>%s Response Time Comparison for %s (%s) & %s (%s) </h3> <br/> <br/>", metric, oldReleaseNumber, oldBuildNum, newReleaseNumber, newBuildNum)
		p = p + fmt.Sprintf("<div style='background:#80bfff;text-align:center'><p><b>IICS Platform Performance Regression</p> </b></div>")
		p = p + fmt.Sprintf("<div style='text-align:left'><p><b>API Labeling : <i>{Servicename}_{API_Name}_{Concurrency}</i></p></b></div>")

		countapis := fmt.Sprintf("<table style='backgound:#99ebff;border-collapse: collapse;' border = '2'cellpadding = '6'><tbody><tr><td colspan=4 style='text-align:center;background-color:Lavender;color:Black;'><b>Performance Summary </b></td></tr><tr><th>Label</th><th>Range</th><th>Use case Count</th><th>Color Code</th></tr> ")
		p10T := fmt.Sprintf("<table style='backgound:#99ebff;;border-collapse: collapse;' border = '2' cellpadding = '6'><tbody><tr><td colspan=5 style='text-align:center;background-color:Lavender;color:Black;'><b> %s Response Time (ms) for 10 user Concurrency </b></td></tr><tr><th>API</th><th>Release: %s (in ms) </th><th>Release: %s (in ms)</th><th>Time Difference</th><th> %% Time Difference</th></tr> ", metric, oldReleaseNumber, newReleaseNumber)
		p100T := fmt.Sprintf("<table style='backgound:#99ebff;border-collapse: collapse;' border = '2' cellpadding = '6'><tbody><tr><td colspan=5 style='text-align:center;background-color:Lavender;color:Black;'><b>%s Response Time (ms) for 100 user Concurrency </b></td></tr><tr><th>API</th><th>Release: %s (in ms) </th><th>Release: %s (in ms)</th><th>Time Difference</th><th> %% Time Difference</th></tr> ", metric, oldReleaseNumber, newReleaseNumber)

		newReleaseDataSorted := utils.SortingMap(newReleaseData) //[all the labels]
		log.Println(newReleaseDataSorted)
		var green int
		var yellow int
		var red int
		var total int
		var newApiCount int
		total = len(oldReleaseData)
		yellow = 0
		red = 0
		green = 0
		newApiCount = 0

		//iterating over the map values of new release data and storing the key/value pair as Label/_
		for _, Label := range newReleaseDataSorted {
			//log.Println(Label)
			_, isNewApi := oldReleaseData[Label]
			if isNewApi {
				if newReleaseData[Label] != 0 {
					if strings.Contains(Label, "10T") {
						//log.Println(Label)
						var timeOld int
						var timeNew int
						timeOld = oldReleaseData[Label]
						timeNew = newReleaseData[Label]
						diff := timeOld - timeNew
						percDiff := utils.CalcPerc(float64(diff), float64(timeOld))
						if percDiff < 0 && percDiff > -20 {
							yellow = yellow + 1
							p10T = p10T + "<tr style='background:Yellow'><td>" + Label + "</td><td>" + strconv.FormatInt(int64(timeOld), 10) + "</td><td>" + strconv.FormatInt(int64(timeNew), 10) + "</td><td>" + strconv.FormatInt(int64(diff), 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}
						if percDiff <= -20 {
							red = red + 1
							p10T = p10T + "<tr style='background:Red'><td>" + Label + "</td><td>" + strconv.FormatInt(int64(timeOld), 10) + "</td><td>" + strconv.FormatInt(int64(timeNew), 10) + "</td><td>" + strconv.FormatInt(int64(diff), 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}
						if percDiff >= 0 {
							green = green + 1
							p10T = p10T + "<tr style='background:Green'><td>" + Label + "</td><td>" + strconv.FormatInt(int64(timeOld), 10) + "</td><td>" + strconv.FormatInt(int64(timeNew), 10) + "</td><td>" + strconv.FormatInt(int64(diff), 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}

					}
					if strings.Contains(Label, "100T") {
						//log.Println(Label)
						var timeOld int
						var timeNew int
						timeOld = oldReleaseData[Label]
						timeNew = newReleaseData[Label]
						diff := timeOld - timeNew
						percDiff := utils.CalcPerc(float64(diff), float64(timeOld))
						if percDiff < 0 && percDiff > -20 {
							yellow = yellow + 1
							p100T = p100T + "<tr style='background:Yellow'><td>" + Label + "</td><td>" + strconv.FormatInt(int64(timeOld), 10) + "</td><td>" + strconv.FormatInt(int64(timeNew), 10) + "</td><td>" + strconv.FormatInt(int64(diff), 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}
						if percDiff <= -20 {
							red = red + 1
							p100T = p100T + "<tr style='background:Red'><td>" + Label + "</td><td>" + strconv.FormatInt(int64(timeOld), 10) + "</td><td>" + strconv.FormatInt(int64(timeNew), 10) + "</td><td>" + strconv.FormatInt(int64(diff), 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}
						if percDiff >= 0 {
							green = green + 1
							p100T = p100T + "<tr style='background:Green'><td>" + Label + "</td><td>" + strconv.FormatInt(int64(timeOld), 10) + "</td><td>" + strconv.FormatInt(int64(timeNew), 10) + "</td><td>" + strconv.FormatInt(int64(diff), 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}

					}

				}
			}
			if !isNewApi {
				log.Println("the " + Label + "is Newly added api")
				var timenewapi int
				newApiCount = newApiCount + 1
				timenewapi = newReleaseData[Label]
				pnewapis = pnewapis + "<tr style='background:White'><td>" + Label + "</td><td>" + strconv.FormatInt(int64(timenewapi), 10) + "</td></tr>"
			}
		}
		countapis = countapis + "<tr><td colspan=2 style='text-align:center;color:Black;'>Total API Count</td><td>" + strconv.FormatInt(int64(total), 10) + "</td><td style='text-align:center;color:Black;'>-</td></tr>"
		countapis = countapis + "<tr><td>% Improvement</td><td> > 0 %</td><td>" + strconv.FormatInt(int64(green), 10) + "</td><td style='background-color: Green;'></td>"
		countapis = countapis + "<tr><td>% Degradation</td><td> 0 to 20 %</td><td>" + strconv.FormatInt(int64(yellow), 10) + "</td><td style='background-color: Yellow;'></td>"
		countapis = countapis + "<tr><td>% Degradation</td><td> > 20 %</td><td>" + strconv.FormatInt(int64(red), 10) + "</td><td style='background-color: Red;'></td>"
		conf := utils.ReadConfig()
		p = p + fmt.Sprintf("<b>Dashboard URL : </b><a href='%s'> %s </a><br><br> %s </tbody></table><br><br> %s </tbody></table><br><br> %s </tbody></table><br><br>", conf.DashboardURL, conf.DashboardURL, countapis, p10T, p100T)
		if newApiCount > 0 {
			p = p + pnewapis + "</tbody></table>"
		}

		utils.SendMail(p, subject, emailID)
		//write to file
		fileName := conf.HtmlFolderPath + oldreleaseIteration + "_" + oldBuildNum + "vs" + newreleaseIteration + "_" + newBuildNum + "_" + ".html"
		f, err := os.Create(fileName)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()

		_, err = f.WriteString(p)
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(p)
		utils.RespondWithJSON("Email Sent Successfully", writer, request)
		log.Println("email sent ")
		log.Println("totalAPIs are ", total, "improvement >0%", green, "degraded 0-20%", yellow, "degraded >20%", red, "NelyAdded", newApiCount)
	}
}
 func htmlReport(w http.ResponseWriter, r *http.Request) {
    log.Println("html report inside ")
	p := "./" + r.URL.Path
	http.ServeFile(w, r, p)

}