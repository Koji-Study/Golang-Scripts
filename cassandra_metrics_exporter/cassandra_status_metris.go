package main
import (
        "fmt"
        "os/exec"
        "strings"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "net/http"
        "github.com/robfig/cron/v3"
)

var(
        MyGauge *prometheus.GaugeVec= prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
        Name: "Cassandra_node_status",
        Help: "Status_of_Cassandia_node",
        },
   // 标签集合
        []string{"instance"},
        )

)

func init(){
        prometheus.MustRegister(MyGauge)
}

func main() {
        crontab := cron.New()
        crontask := func(){
                go_metrics()
        }
        cronfreq := "*/1 * * * *"
        crontab.AddFunc(cronfreq,crontask)
        crontab.Start()
        //select {}
        //go_metrics()
        http.Handle("/metrics", promhttp.Handler())
        fmt.Println("metrics url: http://0.0.0.0:18888/metrics")
        http.ListenAndServe(":18888", nil)
}

func go_metrics(){
        shell_cmd := "sshpass -p passwd ssh -o StrictHostKeychecking=no ubuntu@ip \"/cassandra/apache-cassandra-4.1.0/bin/nodetool -u nodetool -pw passwd status\""
        fmt.Println(string(shell_cmd))
        status,ip := shell_command(shell_cmd)
        //fmt.Println(status)
        //fmt.Println(ip)
        make_metrics(status,ip)
}

func shell_command(cmd string) ([5]int,[5]string) {
        var result_list [5]string
        var result_status [5]int
        var result_ip [5]string
        c := exec.Command("bash", "-c", cmd)
    // 此处是windows版本
    // c := exec.Command("cmd", "/C", cmd)
        output,err := c.CombinedOutput()
        if err != nil{
                panic(err)
        }
        fmt.Println(string(output))
        result := strings.Split(string(output),"\n")
        //fmt.Println(len(result))
        i := 0
        for index,value := range result{
                if index >= 6 && index <=10{
                result_list[i] = value
                i++
        }
        }
        i = 0
        for _,value1 := range result_list{
                if strings.Split(value1,"  ")[0] != "UN"{
                        result_status[i] = 0
                }else{
                        result_status[i] = 1
                }
                result_ip[i] = strings.Split(value1,"  ")[1]
                i++
        }
        return result_status,result_ip
}


func make_metrics(status [5]int,ip [5]string){
        for index2,value2 := range status{
                MyGauge.With(prometheus.Labels{"instance":ip[index2]}).Set(float64(value2))
        }
}
