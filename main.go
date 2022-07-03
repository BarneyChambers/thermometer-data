package main

import (
    "bufio"
    "fmt"
    "os"
	"strings"
	"strconv"
	"math"
)

type Sensor struct
{
	sensor_type string; //thermometer, humidity
	name string; //temp-1, hum-1 
	data map[string]float64; //2007-04-05T22:01 69.5  
}

func getMean (arr map[string]float64) float64 {
	sum := 0.0
	for  _, value := range arr {
		sum += value
	}
	avg := (float64(sum)) / (float64(len(arr)))
	return avg;
} 

func getThermometerPrecision(mean float64, reference float64, stddev float64) string {
	if(mean > (reference + 0.5) || mean < (reference - 0.5)){
		return "precise";
	}
	if(stddev < 3){
		return "ultra precise";
	}
	if(stddev < 5){
		return "very precise";
	}	
	return "precise";
}

func getHumidityQuality(mean float64, reference float64) string {
	if((mean / reference) >= 0.99){
		return "OK";
	}
	return "discard";
}

func getStdDeviation(arr map[string]float64, mean float64) float64{	
	var sd float64
	for  _, value := range arr {
		sd += math.Pow(value - mean, 2);
	}
	sd = math.Sqrt(sd/float64(len(arr)));
	return sd;
}

func getSensorQuality(sensors []Sensor, temperature float64, humidity float64){
	for _, element := range sensors {
		
		if(element.sensor_type == "thermometer"){
			var mean float64 = getMean(element.data);
			var stddev float64 = getStdDeviation(element.data, mean);
			
			fmt.Print(element.name+": ");
			fmt.Println(getThermometerPrecision(mean, temperature, stddev));
		}
		if(element.sensor_type == "humidity"){
			fmt.Print(element.name+": ");

			var mean float64 = getMean(element.data);
			var result string = getHumidityQuality(mean, humidity);
			fmt.Println(result);
		}		
	}
}

func parseThermometerInput(input *os.File) (sensors []Sensor, temperature float64, humidity float64){
	scanner := bufio.NewScanner(input)	
	var currSensor *Sensor;
	var err error;

	for scanner.Scan() {
		
		var currentLine string = scanner.Text(); 
		words := strings.Fields(currentLine);

		if(strings.Contains(currentLine, "reference")){
			humidity,err = strconv.ParseFloat(words[2], 64);
			temperature,err = strconv.ParseFloat(words[1], 64);
			if err != nil {
				throwException("Error: data was malformed, temperature or humidity value could not be converted to float");
			}
		}else if(strings.Contains(currentLine, "thermometer") || strings.Contains(currentLine, "humidity")){
			s := Sensor{ 
				sensor_type: words[0],
				name: words[1],
				data: make(map[string]float64)}

			currSensor = &s;
			sensors = append(sensors, s);
		}else{
			var f float64;
			f,err = strconv.ParseFloat(words[2], 64);
			if err != nil {
				throwException("Error: data was malformed, temperature or humidity value could not be converted to float");
			}
			currSensor.data[words[0]] = f;
		}
	}
	return sensors, temperature, humidity;
}

func throwException(errorMessage string){
	fmt.Println(errorMessage);
	os.Exit(1);
}

func main() {
	var sensors, temperature, humidity = parseThermometerInput(os.Stdin);
	getSensorQuality(sensors, temperature, humidity);
}

