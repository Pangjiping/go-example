package gen_id

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sony/sonyflake"
	"time"
)

func GenBySnowflake() error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}

	for i := 0; i < 20; i++ {
		id := node.Generate()

		fmt.Printf("int64 ID: %d\n", id)
		fmt.Printf("string ID: %s\n", id)
		fmt.Printf("base2 ID: %s\n", id.Base2())
		fmt.Printf("base64 ID: %s\n", id.Base64())
		fmt.Printf("ID time: %d\n", id.Time())
		fmt.Printf("ID node: %d\n", id.Node())
		fmt.Printf("ID step: %d\n", id.Step())
		fmt.Println("--------------------------------")
	}
	return nil
}

type Settings struct {
	StartTime      time.Time              // 起始时间，默认2014-09-01 00:00:00 +0000 UTC
	MachineID      func() (uint16, error) // 返回实例ID的函数，如果不定义此函外，默认用本机ip
	CheckMachineID func(uint16) bool      // 验证实例ID/计算机ID的唯一性，返回true时才创建
}

// 我们需要自己来实现这两个函数 MachineID CheckMachineID
func getMachineID() (uint16, error) {
	var machineID uint16 = 6
	return machineID, nil
}

func checkMachineID(machineID uint16) bool {
	existsMachines := []uint16{1, 2, 3, 4, 5}
	for _, v := range existsMachines {
		if v == machineID {
			return false
		}
	}
	return true
}

func GenbySonyflake() error {
	t, _ := time.Parse("2006-01-02", "2021-01-01")
	settings := sonyflake.Settings{
		StartTime:      t,
		MachineID:      getMachineID,
		CheckMachineID: checkMachineID,
	}

	sf := sonyflake.NewSonyflake(settings)

	for i := 0; i < 20; i++ {
		id, err := sf.NextID()
		if err != nil {
			return err
		}
		fmt.Println(id)
	}

	return nil
}
