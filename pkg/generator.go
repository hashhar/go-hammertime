package pkg

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hashhar/go-hammertime/pkg/constants"
	"github.com/hashhar/go-hammertime/pkg/database"
	"github.com/hashhar/go-hammertime/pkg/models"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// SendWork generates Work with the following distribution of values and pushes
// them onto the passed channel.
//  status |       proportion
// --------+------------------------
//  ACC    |    36.8047308399355773
//  STK    |    23.7282207291480258
//  RCV    |    15.7014104670646993
//  LST    |     4.7528022255120467
//  PAK    |     4.1329775984642665
//  PIK    |     3.3407083245213034
//  ALC    |     2.7908376580064748
//  REJ    |     2.1905350664562625
//  SHP    |     2.0254111828726675
//  TBQ    |     1.4218549187395272
//  TBP    |     1.1615611121052888
//  HLD    | 1.01595926401926174169
//  RTS    | 0.54743041207763262783
//  CAN    | 0.34570271193609787047
//  RPO    | 0.02358912622622785469
//  RWR    | 0.01138785404024792985
//  RTD    | 0.00488050887439196994
func SendWork(workChan chan models.Work, count int) {
	managementKeys, managementTable := database.LoadManagementTable()
	for count > 0 {
		random := rand.Float32()
		chance := random * 100
		// RTD
		if chance < 0.004 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "RTD",
			}
		}
		// RWR
		if chance < 0.01 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "RWR",
			}
		}
		// RPO
		if chance < 0.02 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "RPO",
			}
		}
		// CAN
		if chance < 0.34 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "CAN",
			}
		}
		// RTS
		if chance < 0.54 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "RTS",
			}
		}
		// HLD
		if chance < 1.01 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "HLD",
			}
		}
		// TBP
		if chance < 1.16 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "TBP",
			}
		}
		// TBQ
		if chance < 1.42 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "TBQ",
			}
		}
		// SHP
		if chance < 2.02 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "SHP",
			}
		}
		// REJ
		if chance < 2.19 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "REJ",
			}
		}
		// ALC
		if chance < 2.79 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "ALC",
			}
		}
		// PIK
		if chance < 3.34 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "PIK",
			}
		}
		// PAK
		if chance < 4.13 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "PAK",
			}
		}
		// LST
		if chance < 4.75 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "LST",
			}
		}
		// RCV
		if chance < 15.70 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "RCV",
			}
		}
		// STK
		if chance < 23.72 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "STK",
			}
		}
		// ACC
		if chance < 36.80 && count > 0 {
			count--
			key := managementKeys[int(random*float32(len(managementTable)))]
			workChan <- models.Work{
				Serial:       getRandomSerial(),
				FcSk:         managementTable[key],
				ManagementID: key,
				Status:       "ACC",
			}
		}
	}
	close(workChan)
}

func getRandomSerial() string {
	// 10 bytes in hex would take up 20 chars
	data := make([]byte, constants.SerialLenHexBytes)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", data)
}
