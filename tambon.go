package main
import ("fmt"
	"net"
	"html"
	"net/http"
	"log"
	"encoding/csv"
	"strconv"
	"os"
	"runtime"
	"github.com/gorilla/mux"
	"math"
	"encoding/json"
	)

//TA_ID,TAMBON_T,TAMBON_E,AM_ID,AMPHOE_T,AMPHOE_E,CH_ID,CHANGWAT_T,CHANGWAT_E,LAT,LONG
//910106,ต. เกาะสาหร่าย,Ko Sarai,9101,อ. เมืองสตูล,Mueang Satun,91,จ. สตูล,Satun,99.706,6.546
type TamBon struct {
	TA_ID		int
	Tambon_T	string
	Tambon_E	string
	AM_ID		int
	Amphoe_T	string
	Amphoe_E	string
	CH_ID		int
	Changwat_T	string
	Changwat_E	string
	Lat		float64
	Lon		float64
}

	var tamBon [7768]TamBon

func main(){

        nCPU := runtime.NumCPU()
        runtime.GOMAXPROCS(nCPU)
        log.Println("Number of CPUs: ", nCPU)


	initData()

        rtr := mux.NewRouter()
        rtr.HandleFunc("/distance",disTance).Methods("GET").Queries("lat1","{lat1:([0-9]*.[0-9]+|[0-9]+)}","lon1","{lon1:([0-9]*.[0-9]+|[0-9]+)}")
        http.Handle("/", rtr)


        bind := ":8081"

        log.Println("Listening:" + bind + "...")
        err := http.ListenAndServe(bind, nil)
        if err != nil {
                panic(err)
        }



}

func initData() {

         csvfile, err := os.Open("TAMBON.csv")
         if err != nil {
                 fmt.Println(err)
                 return
         }

         defer csvfile.Close()

         reader := csv.NewReader(csvfile)

         reader.FieldsPerRecord = -1 // see the Reader struct information below

         rawCSVdata, err := reader.ReadAll()

         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }

         // sanity check, display to standard output
         for i, each := range rawCSVdata {
                //fmt.Printf("[%d] %s %s %s %s %s %s %s %s %s %s %s\n",i, each[0], each[1],each[2],each[3],each[4],each[5],each[6],each[7],each[8],each[9],each[10])
		tamBon[i].TA_ID,_  = strconv.Atoi(each[0])
		tamBon[i].Tambon_T = each[1]
		tamBon[i].Tambon_E = each[2]
		tamBon[i].AM_ID,_  = strconv.Atoi(each[3])
		tamBon[i].Amphoe_T = each[4]
		tamBon[i].Amphoe_E = each[5]
		tamBon[i].CH_ID,_  = strconv.Atoi(each[6])
		tamBon[i].Changwat_T = each[7]
		tamBon[i].Changwat_E = each[8]
		tamBon[i].Lat,_ = strconv.ParseFloat(each[9], 64)
		tamBon[i].Lon,_ = strconv.ParseFloat(each[10], 64)
         }
}


func disTance(w http.ResponseWriter, r *http.Request) {

	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
        params := mux.Vars(r)

        lat1,_ := strconv.ParseFloat(params["lat1"], 64)
        lon1,_ := strconv.ParseFloat(params["lon1"], 64)

	var distance [7768]float64
	var R float64 = 637100

        log.Println(lat1,lon1)

        for i, each := range tamBon {
		lat2 := each.Lat
		lon2 := each.Lon

		dLat := deg2rad(lat2-lat1)
		dLon := deg2rad(lon2-lon1)

		a:=math.Sin(dLat/2) * math.Sin(dLat/2) + math.Cos(deg2rad(lat1)) * math.Cos(deg2rad(lat2)) * math.Sin(dLon/2) * math.Sin(dLon/2)
		c:=2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		d:=R * c
		distance[i]=d
	}

	//for Debug
/*
	for i,data := range distance {
		fmt.Printf("data[%d]=%f\n",i,data)
	}
*/
	//End debug

	//find min distance
	min := distance[0]
	index_arr := 0
	for i,data := range distance {
		if data < min {
			min = data
			index_arr=i
		}
	}
	log.Printf("[%s][%s][%q]\n",ip,r.UserAgent(),html.EscapeString(r.URL.Path))
	fmt.Printf("min distance = %f at index[%d]\n",min,index_arr)
	js,err := json.Marshal(tamBon[index_arr])
	if err != nil {
		panic(err)
	 }
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(js)

}

func deg2rad(deg float64) float64 {
        return deg * (math.Pi/180)
}

