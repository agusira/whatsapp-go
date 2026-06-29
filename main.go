package main

import (
	"agus/configs"
	"agus/handler"
	"agus/lib"
	_ "agus/plugins"
	"strings"

	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	meowcaller "github.com/purpshell/meowcaller"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	// "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	// "google.golang.org/protobuf/proto"
)

var dbAddress = flag.String("db-address", "./tes.db", "Database Address")
var deviceName = flag.String("device-name", "Agus", "Device name")

func init() {
	configs.LoadEnv()
	flag.Parse()
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	clearTerminal()
	fmt.Println("<< Welcome To WhatsGO >>")

	logger := os.Getenv("LOGGER")
	// dbFile := os.Getenv("DB_FILE")
	ctx := context.Background()
	dbLog := waLog.Stdout("Database", logger, true)
	container, _ := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", *dbAddress), dbLog)

	deviceStore := loadSession(ctx, *container)
	// deviceStore, _ := container.GetFirstDevice(ctx)
	clientLog := waLog.Stdout("Client", logger, true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	// call event handler
	call := meowcaller.NewClient(client)
	call.OnIncomingCall(func(c *meowcaller.Call) {
		if err := c.Answer(); err != nil {
			fmt.Println(err)
			return
		}
		aud, err := meowcaller.MP3File("./call.mp3")
		if err != nil {
			fmt.Println(err)
			return
		}
		c.Play(aud)
	})

	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_CHROME.Enum()
	store.SetOSInfo(*deviceName, [3]uint32{2, 3000, 1031080782})

	client.AddEventHandler(func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			conn := lib.SerializeClient(client)
			m := lib.Serialize(v, conn)

			go handler.Handler(conn, m)
		case *events.Connected:
			fmt.Println("[!] Connected to Whatsapp")
			// client.SendPresence(ctx, types.PresenceAvailable)
		case *events.CallOffer:
			if configs.CONFIG.AntiCall {
				client.RejectCall(ctx, v.From, v.CallID)
				text := "Halo, saat ini saya sedang dalam kondisi yang tidak memungkinkan untuk menerima telepon. Mohon untuk meninggalkan pesan!"
				client.SendMessage(ctx, v.From, &waE2E.Message{
					ExtendedTextMessage: &waE2E.ExtendedTextMessage{
						Text: &text,
						ContextInfo: &waE2E.ContextInfo{
							QuotedMessage: &waE2E.Message{
								Call: &waE2E.Call{
									ContextInfo: &waE2E.ContextInfo{
										StanzaID:    &v.CallID,
										Participant: &v.From.User,
									},
								},
							},
						},
					},
				})
			}
		}
	})

	if client.Store.ID == nil {
		var number string
		fmt.Println("[!] No Session Available")
		fmt.Println("[!] Creating New Session")
		fmt.Print("[?] Type Your Number (62): ")
		fmt.Scanln(&number)
		fmt.Printf("[!] Connecting With %s\n", number)

		// No ID stored, new login
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		code, err := client.PairPhone(ctx, number, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if err != nil {
			fmt.Println("[!] Connection Failure!")
			return
		}
		fmt.Println("[!] Kode : ", code)
	} else {
		// Already logged in, just connect
		err := client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func loadSession(ctx context.Context, container sqlstore.Container) *store.Device {
	var pilih int
	var deviceStore *store.Device

	// container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", *dbAddress), dbLog)

	clearTerminal()
	fmt.Printf("\n[1] New Session\n[2] Load Session\n[?] Your Choice: ")
	fmt.Scanln(&pilih)
	switch pilih {
	case 1:
		deviceStore = container.NewDevice()
	case 2:
		var str strings.Builder
		var id int
		all, _ := container.GetAllDevices(ctx)
		fmt.Println()

		if len(all) <= 1 {
			device, err := container.GetFirstDevice(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(device.GetJID())
			deviceStore = device
		}

		if len(all) > 1 {
			str.WriteString("List Device\n")
			for i, a := range all {
				fmt.Fprintf(&str, "%d. %s\n", i, a.ID.User)
			}
			fmt.Println(str.String())
			fmt.Print("[?] Your Choice: ")
			fmt.Scanln(&id)
			if all[id] != nil {
				device, err := container.GetDevice(ctx, all[id].GetJID())
				if err != nil {
					panic(err)
				}
				deviceStore = device
			}
			if all[id] == nil {
				device, err := container.GetFirstDevice(ctx)
				if err != nil {
					panic(err)
				}
				deviceStore = device
			}
		}
	default:
		if pilih > 2 || pilih == 0 {
			os.Exit(0)
		}
	}
	clearTerminal()
	fmt.Println(deviceStore.GetJID())
	return deviceStore
}
