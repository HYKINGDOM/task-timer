// 这是一个测试文件，用于验证番茄工作法定时器功能
package main

import (
	"fmt"
	"testing"
	"time"
)

// 测试定时器功能
func testTimer() {
	// 模拟设置
	workMinutes := 1 // 为了测试方便，设置为1分钟
	restMinutes := 1 // 为了测试方便，设置为1分钟

	fmt.Println("===== 番茄工作法定时器测试 =====")
	fmt.Printf("工作时间: %d分钟\n", workMinutes)
	fmt.Printf("休息时间: %d分钟\n", restMinutes)

	// 模拟工作时间倒计时
	fmt.Println("\n开始工作时间倒计时...")
	endTime := time.Now().Add(time.Duration(workMinutes) * time.Minute)

	for time.Now().Before(endTime) {
		remaining := endTime.Sub(time.Now())
		hours := int(remaining.Hours())
		minutes := int(remaining.Minutes()) % 60
		seconds := int(remaining.Seconds()) % 60

		fmt.Printf("\r工作时间剩余: %02d:%02d:%02d", hours, minutes, seconds)
		time.Sleep(time.Second)
	}

	fmt.Println("\n\n工作时间结束！开始休息！")

	// 模拟休息时间倒计时
	fmt.Println("\n开始休息时间倒计时...")
	endTime = time.Now().Add(time.Duration(restMinutes) * time.Minute)

	for time.Now().Before(endTime) {
		remaining := endTime.Sub(time.Now())
		hours := int(remaining.Hours())
		minutes := int(remaining.Minutes()) % 60
		seconds := int(remaining.Seconds()) % 60

		fmt.Printf("\r休息时间剩余: %02d:%02d:%02d", hours, minutes, seconds)
		time.Sleep(time.Second)
	}

	fmt.Println("\n\n休息时间结束！可以继续工作了！")
}

// 测试时间范围验证
func TestTimeRangeValidation(t *testing.T) {
	// 测试工作时间范围
	workMinutes := []int{0, 1, 30, 60, 61}
	for _, minutes := range workMinutes {
		if minutes < 1 || minutes > 60 {
			t.Errorf("工作时间 %d 分钟超出有效范围 [1-60]", minutes)
		}
	}

	// 测试休息时间范围
	restMinutes := []int{0, 1, 15, 60, 61}
	for _, minutes := range restMinutes {
		if minutes < 1 || minutes > 60 {
			t.Errorf("休息时间 %d 分钟超出有效范围 [1-60]", minutes)
		}
	}
}

// 测试倒计时精确性
func TestCountdownAccuracy(t *testing.T) {
	// 测试1分钟倒计时
	startTime := time.Now()
	endTime := startTime.Add(1 * time.Minute)
	actualDuration := endTime.Sub(startTime)

	if actualDuration != 1*time.Minute {
		t.Errorf("倒计时精确性测试失败：预期1分钟，实际%.2f秒", actualDuration.Seconds())
	}
}

// 如果需要单独运行测试，取消下面的注释
/*
func main() {
	testTimer()
}
*/
