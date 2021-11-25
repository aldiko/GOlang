package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"v0/office"
)

func main() {
	employeeActions,err := readCsv("employeeActions")
	if err!=nil {
		fmt.Printf("Error running program: %s \n",err.Error())		
	}
	
	Loop:
		for i, v := range employeeActions {
			fmt.Printf("%d: %s ->",i,v[0])
			employee,err := office.NewEmployeeFactory(v[0])
			if err!=nil {
				fmt.Printf("Error running program: %s \n",err.Error())
				continue		
			}
			if employee.GetCurrentLocation()==nil {
				loc,err := office.NewLocationFactory("office")
				if err!=nil {
					fmt.Printf("Error running program: %s \n",err.Error())
				}
				employee.MoveToLocation(loc)
			}
			places:=v[1:]
			for _, v2 := range places {
				loc,err := office.NewLocationFactory(v2)
				if err!=nil{
					fmt.Printf("Error running program: %s \n",err.Error())
					continue Loop
				}
				if employee.GetCurrentLocation().GetLocationTitle()==loc.GetLocationTitle() {
					fmt.Printf("%s ->",loc.GetLocationTitle())
					continue
				}
				if employee.GetCurrentLocation().CheckMoveToArea(loc) {
					err:=employee.MoveToLocation(loc)
					if err!=nil {
						fmt.Printf("Error running program: %s \n",err.Error())
						continue Loop	
					}
					fmt.Printf("%s->",loc.GetLocationTitle())
				}else {
					fmt.Printf("Error can not move to %s from %s",loc.GetLocationTitle(),employee.GetCurrentLocation().GetLocationTitle())
					continue
				}
			}
			fmt.Printf("\n")
		}

}


func readCsv(name string)([][]string,error){
	f, err := os.Open(name)
    if err != nil {
		return nil, fmt.Errorf("readcsv error open %s: %w",name,err)
    }
    defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	employeeActions , err := reader.ReadAll()
	if err!=nil{
		return nil, fmt.Errorf("readcsv error parse %s: %w",name,err)
	}
	return employeeActions,nil
}