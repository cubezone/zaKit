package main
 
import (
//  "encoding/json"
//   iconv "iconv-go"
    "os"
    "bufio"
    "io"
    "fmt"
    "strconv"   
    "strings"
)
 
//对应json天气数据源的结构，头字母大写
type dayvalue struct {    
    stockid string
    day     string
    c       float64
    high    float64
    low     float64    
    value   float64

    EMA_12  float64
    EMA_26  float64
    DIFF    float64
    DEA     float64
    BAR     float64

    High_9  float64
    Low_9   float64
    RSV     float64
    Kt      float64
    Dt      float64
    Jt      float64
}

//http://blog.sina.com.cn/s/blog_85c43c210102v98k.html
/* MACD 算法
    具体计算公式及例子如下： 
    EMA（12）= 前一日EMA（12）×11/13＋今日收盘价×2/13
    EMA（26）= 前一日EMA（26）×25/27＋今日收盘价×2/27
    DIFF=今日EMA（12）- 今日EMA（26）
    DEA（MACD）= 前一日DEA×8/10＋今日DIF×2/10 
    BAR=2×(DIFF－DEA)
    对理工检测：
    20091218日：
        新股上市，DIFF=0, DEA=0, MACD=0，收盘价55.01
    20091219日：
        收盘价53.7
        EMA（12）=55.01+(53.7-55.01)×2/13=54.8085   
        EMA（26）=55.01+(53.7-55.01)×2/27=54.913  
    DIFF=EMA（12）- EMA（26）= 54.8085 - 54.913 = -0.1045  （-0.104？）
*/

func cal_macd(arr [] dayvalue , a int , b int){

    var ee dayvalue 

    arr[0].EMA_12,arr[0].EMA_26 = arr[0].value,arr[0].value
    fmt.Print(arr[0])

    for index := 1; index < len(arr); index++ {   
     ee = arr[index]
     ee.EMA_12 = arr[index-1].EMA_12*float64(a-1)/float64(a+1) + ee.value*2/float64(a+1);
     ee.EMA_26 = arr[index-1].EMA_26*float64(b-1)/float64(b+1) + ee.value*2/float64(b+1);
     ee.DIFF =  ee.EMA_12  -  ee.EMA_26 ;
     ee.DEA = arr[index-1].DEA *8/10 + ee.DIFF*2/10;
     ee.BAR = 2*(ee.DIFF-ee.DEA)

     arr[index] = ee
    }
}

/*  KDJ计算公式：
    　　以9日周期的KDJ为例，首先算出最近9天的“未成熟随机值”即RSV值，
    RSV的计算公式如下：　　
    RSVt＝(Ct－L9)／(H9－L9)＊100　　
    式中： Ct-------当日收盘价　　L9-------9天内最低价　　H9-------9天内最高价　　
    从计算公式可以看出，RSV指标和%R计算很类似。
    事实上，同周期的RSV值与%R值之和等于100，因而RSV值也介于与100之间。得出RSV值后，便可求出K值
    与D值：K值为RSV值3日平滑移动平均线，而D值为K值的3日平滑移动平均线三倍K值减二倍D值所得的J线，
    其计算公式为：　　
    Kt＝RSVt／3＋2＊Kt-1／3　　
    Dt＝Kt／3＋2＊Dt-1／3　　
    Jt＝3＊Dt－2＊Kt  ;  Jt＝Kt+ (Kt－Dt) ＊2  ??
    KD线中的RSV，随着9日中高低价、收盘价的变动而变动。
    如果没有KD的数值，就可以用当日的RSV值或50代替前一日的KD之值。经过平滑运算之后，起算基期不同
    的KD值将趋于一致，不会有任何差异，K值与K值永远介于0至100之间。根据快、慢移动平均线的交叉原理
    ，K线向上突破K线为买进信号，K线跌破D线为卖出信号，即行情是一个明显的涨势，会带动K线(快速平均
    值)与D线(慢速平均值)向上升，如果涨势开始迟缓，便会慢慢反应到K值与D值，使K线跌破D线，此时中短
    期调整跌势确立，这是一个常用的简单应用原则。
*/
func cal_kdj(arr [] dayvalue,cnt int){
     
    var ee dayvalue
    if (cnt > len(arr)){
        return
    }

    arr[0].EMA_12,arr[0].EMA_26 = arr[0].value,arr[0].value

    var tt_high ,tt_low float64 = 0 , 0

    for index := 0; index < len(arr); index++ {
        ee = arr[index] 
        tt_high = arr[index].high
        tt_low = arr[index].low
        var a int
        if index > cnt {
            a = cnt
        }else{
            a  = index
        }
        if (index > a){
            for i:= 0; i < a ; i++ {
                if arr[index-i].high > tt_high{
                    tt_high = arr[index-i].high
                }
                if arr[index-i].low < tt_low{
                    tt_low = arr[index-i].low
                }
            }       
        }
        ee.High_9 = tt_high
        ee.Low_9 = tt_low
        
        arr[index] = ee;
    }

    for index := 1; index < len(arr); index++ {
    //fmt.Printf("arr[%d]=%f \n", index, arr[index].a)
     ee = arr[index]
     ee.RSV = (ee.value - ee.Low_9)/(ee.High_9 - ee.Low_9)*100
     ee.Kt = ee.RSV/3 + 2* arr[index-1].Kt/3;
     ee.Dt = ee.Kt/3 + 2 *arr[index-1].Dt/3;
     ee.Jt =  3*ee.Dt - 2*ee.Kt ;
     arr[index] = ee
    }
}

/*
func arr_init() {    
    for index := 0; index < len(arr); index++ {
    //fmt.Printf("arr[%d]=%f \n", index, arr[index].a)
     arr[index].c = 12.3 + float64(index)   
     arr[index].value = 7.8+ float64(index)
     arr[index].high = 9+ float64(index)
     arr[index].low = 7+ float64(index)
    }
}
*/

func arr_init(fn string)(arr [] dayvalue, err error) {
    f, err := os.Open(fn)//打开文件     
    defer f.Close() //打开文件出错处理     
    arr = make( [] dayvalue, 0, 30)
    if nil != err {  
      return arr,err       
    }
    buff := bufio.NewReader(f) //读入缓存         
    for { 
        line, err := buff.ReadString('\n') //以'\n'为结束符读入一行    
        if err != nil || io.EOF == err { 
          break             
        }
        //fmt.Print(line)  //可以对一行进行处理 
         rets := strings.Split(line," ")   
         var t dayvalue 
         t.stockid = rets[0]
         t.day = rets[1]
         // open,higt, close,low
         t.high,_ = strconv.ParseFloat(rets[3],3)
         t.value,_ = strconv.ParseFloat(rets[4],3)         
         t.low,_ = strconv.ParseFloat(rets[5],3)

         //fmt.Println(t.value,t.high);
         arr = append(arr,t)
    }
    //根据时间 正向排序， 
    for i := 0 ; i < len(arr)/2 ;i++ {
      arr[i],arr[len(arr)-1-i] = arr[len(arr)-1-i],arr[i]
    }
    return  arr,err
}

func main(){

     var fname string

     arg_num := len(os.Args)
     if (arg_num == 1 ){
         fmt.Println("require data file, C:\\php\\600887.txt " )
         return        
     }else{
        fname = os.Args[1]
     }
  
    //fname = "C:\\php\\600887.txt"
  
    arr, err := arr_init(fname) 
    if (err !=  nil){
        fmt.Println("can't open ", fname)
        return 
    } 

    cal_macd(arr,12,26)
    cal_kdj(arr,9)

    for index, cur := range arr {
    //    value.a = 12.3;
        //fmt.Println(cur)        
        if (cur.BAR < -0.5 && cur.Kt < 15){
        fmt.Printf("\n %s %s %6f %6f  \n",cur.stockid,cur.day ,cur.high,cur.low);
        fmt.Printf("arr[%4d]=%6f VALUE:%6f ,DEA :%6f ", index,  cur.EMA_12, cur.value,cur.DEA)
        fmt.Printf("DIFF:%6f  MACD:%6f\n ", cur.DIFF , cur.BAR)
        fmt.Printf("RSV:%6f,value:%6f ,high9:%6f, low9: %6f ,K:%6f ,D:%6f  J:%6f\n",cur.RSV,cur.value,cur.High_9, cur.Low_9, cur.Kt, cur.Dt,cur.Jt)
        }
        
    }
}