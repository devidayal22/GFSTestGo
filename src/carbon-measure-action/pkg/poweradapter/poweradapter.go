package poweradapter

import (
	"fmt"
	EM "main/pkg/electricitymap"
	WT "main/pkg/watttime"
	"strconv"
	"strings"
	time "time"
)

// RETURN
type CarbonIntensity struct {
	LiveCarbonIntensity float64
	History             []RecentCIHistory
}

type RecentCIHistory struct {
	CarbonIntensity float64
	Datetime        string
}

// GET FROM
type TypCarbonQueryParams struct {
	IacProvider           string
	IacLocation           string
	CarbonRateProvider    string
	ElectricityMapZoneKey string
	WattTimeUser          string
	WattTimePass          string
}

func LiveCarbonIntensity(params TypCarbonQueryParams) (ci CarbonIntensity) {
	zone := GetLocation(TypCloudLocationQuery{
		Provider:      params.IacProvider,
		Location:      params.IacLocation,
		Powerprovider: params.CarbonRateProvider,
	})
	if strings.ToLower(params.CarbonRateProvider) == "electricitymap" {

		em := EM.New(params.ElectricityMapZoneKey)

		live, _ := em.LiveCarbonIntensity(EM.TypAPIParams{Zone: zone})
		ci.LiveCarbonIntensity = float64(live.CarbonIntensity)
		recent, _ := em.RecentCarbonIntensity(EM.TypAPIParams{Zone: zone})
		var historyci []RecentCIHistory
		for _, i := range recent.History {
			value := float64(i.CarbonIntensity)
			historyci = append(historyci, RecentCIHistory{value, i.Datetime})
		}
		ci.History = historyci

		return

	} else if strings.ToLower(params.CarbonRateProvider) == "watttime" {

		live, recent := Watttime(TypCarbonQueryParams{WattTimeUser: params.WattTimeUser, WattTimePass: params.WattTimePass}, zone)
		fmt.Println("devi dayal" + " h" + params.WattTimeUser)
		fmt.Println(live)
		if live != nil {
			ci.LiveCarbonIntensity, _ = strconv.ParseFloat(live.Percent, 64)
			fmt.Println("ci.LiveCarbonIntensity")
			fmt.Println(live.Freq)
			fmt.Println(live.BA)
			fmt.Println(live.Moer)
			fmt.Println(live.Percent)
			fmt.Println(recent)
		}

		if recent != nil {
			var historyci []RecentCIHistory

			for _, i := range *recent {
				historyci = append(historyci, RecentCIHistory{float64(i.Value), i.PointTime})
			}
			ci.History = historyci
		}
	}

	return
}

func GetTimeRange() (starttime, endtime string) {
	t := time.Now()
	et := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	st := et.Add((time.Hour * 23) * -1)
	starttime = st.Format(time.RFC3339)
	endtime = et.Format(time.RFC3339)
	return
}

func Watttime(params TypCarbonQueryParams, BA string) (*WT.RealTimeEmissionsIndexResp, *[]WT.GridEmissionsDataResp) {

	WT.Login(params.WattTimeUser, params.WattTimePass)
	fmt.Println("devi dayal" + " w" + params.WattTimeUser)
	starttime, endtime := GetTimeRange()

	live, _ := WT.RealTimeEmissionsIndex(BA, 0, 0, "")
	recent, _ := WT.GridEmissionsData(BA, 0, 0, starttime, endtime, "", "")
	fmt.Println(live)
	fmt.Println(recent)
	return live, recent
}
