package modemmanager

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"

        "fmt"
        "github.com/maltegrosse/go-modemmanager"
        "log"
)

type ModemManager struct {
	Interface	string
	Connected	bool
	SignalRssi	int64
	SignalRsrp	int64
	SignalRsrq	int64
	SignalSnr	int64
}


func (m *ModemManager) SampleConfig() string {
	return ""
}

func (m *ModemManager) Description() string {
	return "Inserts signal statistics from ModemManager"
}

func (m *ModemManager) Gather(acc telegraf.Accumulator) error {
        mmgr, err := modemmanager.NewModemManager()
        if err != nil {
                log.Fatal(err.Error())
        }
        err = mmgr.ScanDevices()
        if err != nil {
                log.Fatal(err.Error())
        }

        modems, err := mmgr.GetModems()
        if err != nil {
                log.Fatal(err.Error())
        }
        for _, modem := range modems {
                tags := make(map[string]string)
                /* fieldsS := make(map[string]interface{}) */
                fieldsG := make(map[string]interface{})
                /* fieldsC := make(map[string]interface{}) */

                bearers, err := modem.GetBearers()
                if err != nil {
                        log.Fatal(err.Error())
                }
                for _, bearer := range bearers {
                        if err != nil {
                                fmt.Println(err)
                        } else {
				intf, err := bearer.GetInterface()
		                if err != nil {
		                        log.Fatal(err.Error())
				}
				conn, err := bearer.GetConnected()
				if err != nil {
		                        log.Fatal(err.Error())
		                }

				tags["interface"] = intf;
                                fieldsG["Connected"] = conn;

				fmt.Println("Interface", intf)
				fmt.Println("Connected", conn)
                        }
                }

                modemSignal, err := modem.GetSignal()
                if err != nil {
                        log.Fatal(err.Error())
                }

                signalRate, err := modemSignal.GetRate()
                if signalRate < 1 {
                        err := modemSignal.Setup(10)
                        if err != nil {
                                log.Fatal(err.Error())
                        }
                }

                CurrentSignals, err := modemSignal.GetCurrentSignals()
                for _, CurrentSignal := range CurrentSignals {
			tags["ConnectionType"] = fmt.Sprintf("%v", CurrentSignal.Type)
			fieldsG["SignalRssi"] = CurrentSignal.Rssi
			fieldsG["SignalRsrp"] = CurrentSignal.Rsrp
			fieldsG["SignalRsrq"] = CurrentSignal.Rsrq
			fieldsG["SignalSnr"] = CurrentSignal.Snr

                        fmt.Println("Type", CurrentSignal.Type)
                        fmt.Println("SignalRssi", CurrentSignal.Rssi)
                        fmt.Println("SignalRsrp", CurrentSignal.Rsrp)
                        fmt.Println("SignalRsrq", CurrentSignal.Rsrq)
                        fmt.Println("SignalSnr", CurrentSignal.Snr)
                }

                /* acc.AddFields("modemmanager", fieldsS, tags) */
                acc.AddGauge("modemmanager", fieldsG, tags)
                /* acc.AddCounter("modemmanager", fieldsC, tags) */
	}

	return nil
}

func init() {
	inputs.Add("modemmanager", func() telegraf.Input { return &ModemManager{} })
}
