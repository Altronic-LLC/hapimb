package hapimb

import (
	"fmt"
	"log"
	"math"

	modbus "github.com/Altronic-LLC/altronic_modbus"
)

func SliceToJson(slice []uint16, start uint16, inc int) string {
	sliceLen := len(slice)
	if sliceLen == 1 && inc == 8 {
		//fmt.Println("breaking up slice",sliceLen)
		newSlice := make([]uint16, 8)
		byteVal := slice[0]
		for i := 0; i < 8; i += 1 {
			newSlice[i] = uint16((byteVal >> i) & 1)
		}
		slice = newSlice
		sliceLen = len(slice)
		inc = 1
	}
	jsonStr := "{"
	index := int(start)
	for i := 0; i < sliceLen; i++ {
		commaStr := ","
		val := slice[i]
		if i == sliceLen-1 {
			commaStr = ""
		}
		jsonStr = jsonStr + fmt.Sprintf(`"%d":%d%s`, index, val, commaStr)
		index += inc
	}
	jsonStr += "}"
	return jsonStr
}

// HmbDoFc handles the Modbus requests for the HMB devices and returns the result as a slice of uint16 or json string
func HmbDoFc(node uint8, fc uint8, start uint16, count uint16, data []byte, localClient *modbus.ModbusClient, returnJson bool) interface{} {

	typ := modbus.HOLDING_REGISTER
	rv := []uint16{}
	var err error = nil
	//inc := 1
	if fc == 17 {
		rvB, err := localClient.ReadFunction17()
		if err != nil {
			fmt.Println("err fc17: ", err, start)
		}
		rv = make([]uint16,len(rvB))
		for i,v := range(rvB) {
			rv[i] = uint16(v)
		}
		start = 0
	}
	if fc == 16 {
		err = localClient.WriteBytes(start, data)
		if err != nil {
			log.Println("error in hmbdofc ", err)
			log.Println("parameters: ", node, fc, start, count, data)
		}
		//log.Println(rv)
	}
	if fc == 6 {
		data1 := uint16(data[0])*256 + uint16(data[1])
		err = localClient.WriteRegister(start, data1)
		if err != nil {
			log.Println("error in hmbdofc", err)
			log.Println("parameters: ", node, fc, start, count, data)
		}
	}
	if fc == 5 {
		coilVal := false
		if data[1] == 1 {
			coilVal = true
		}
		err = localClient.WriteCoil(start, coilVal)
		if err != nil {
			log.Println("error in hmbdofc", err)
			log.Println("parameters: ", node, fc, start, count, data)
		}
		//log.Println(rv)
	}
	if fc == 4 {
		typ = modbus.INPUT_REGISTER
		rv, err = localClient.ReadRegisters(start, count, typ)
		if err != nil {
			log.Println("error in hmbdofc", err)
			log.Println("parameters: ", node, fc, start, count, data)
		}
		//log.Println(rv)
	}
	if fc == 3 {
		typ = modbus.HOLDING_REGISTER
		//log.Println("pre-read FC3", start, count, typ)
		rv, err = localClient.ReadRegisters(start, count, typ)
		if err != nil {
			//log.Println("in err fc3")
			log.Println("error in hmbdofc", err)
			log.Println("parameters: ", node, fc, start, count, data)
		}
		//log.Println(rv)
	}
	if fc == 2 {
		//inc = 8
		rv2, err := localClient.ReadDiscreteInputs(start, count)
		if err != nil {
			log.Println("error in hmbdofc", err)
			log.Println("parameters: ", node, fc, start, count, data)
		}

		_ = rv2
		//log.Println(rv2)
		cnt := uint16(0)
		registerValue := uint16(0)
		for i, v := range rv2 {
			_ = i
			//log.Println("cnt: ",cnt,"v: ",v)
			if v {
				registerValue = registerValue + uint16(math.Pow(2, float64(cnt)))
				//log.Println(registerValue,uint16(math.Pow(2,float64(cnt))))
			}

			if cnt == 7 {
				rv = append(rv, registerValue)
				//log.Println("final registerValue: ",registerValue)
				registerValue = 0
				cnt = 0
			} else {
				cnt += 1
			}
		}
		//log.Println("registerValue: ",registerValue)
		if cnt > 0 {
			rv = append(rv, registerValue)
		}
		//log.Println("rv: ",rv)

	}

	if fc == 1 {
		//inc = 8
		rv2, err := localClient.ReadCoils(start, count)
		if err != nil {
			log.Println("error in hmbdofc", err)
			log.Println("parameters: ", node, fc, start, count, data)
		}
		_ = rv2
		//log.Println(rv2)
		cnt := uint16(0)
		registerValue := uint16(0)
		for i, v := range rv2 {
			_ = i
			//log.Println("cnt: ",cnt,"v: ",v)
			if v {
				registerValue = registerValue + uint16(math.Pow(2, float64(cnt)))
				//log.Println(registerValue,uint16(math.Pow(2,float64(cnt))))
			}

			if cnt == 7 {
				rv = append(rv, registerValue)
				//log.Println("final registerValue: ",registerValue)
				registerValue = 0
				cnt = 0
			} else {
				cnt += 1
			}
		}
		//log.Println("registerValue: ",registerValue)
		if cnt > 0 {
			rv = append(rv, registerValue)
		}
		//log.Println("rv: ",rv)

	}
	//log.Println("rv: ", rv)
	if returnJson {
		return SliceToJson(rv, start, 1)
	}
	return rv
}
