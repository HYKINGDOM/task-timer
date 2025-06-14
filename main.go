package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// 应用状态
type AppState int

const (
	StateInit AppState = iota
	StateWorking
	StateResting
)

// 主应用结构体
type EyeProtectionApp struct {
	// 窗口
	mainWindow *walk.MainWindow

	// 当前状态
	state AppState

	// 设置界面控件
	workTimeEdit  *walk.NumberEdit
	restTimeEdit  *walk.NumberEdit
	startTimeEdit *walk.DateEdit
	endTimeEdit   *walk.DateEdit
	startButton   *walk.PushButton

	// 倒计时界面控件
	countdownLabel *walk.Label
	endButton      *walk.PushButton

	// 设置值
	workMinutes int
	restMinutes int
	startTime   time.Time
	endTime     time.Time

	// 倒计时相关
	countdownTimer *time.Timer
	countdownEnd   time.Time
}

func main() {
	app := &EyeProtectionApp{
		state: StateInit,
	}

	// 创建主窗口
	if err := app.createMainWindow(); err != nil {
		fmt.Println("创建主窗口失败:", err)
		return
	}

	// 运行应用
	app.mainWindow.Run()
}

// 创建主窗口
func (app *EyeProtectionApp) createMainWindow() error {
	// 创建主窗口
	return MainWindow{
		AssignTo: &app.mainWindow,
		Title:    "番茄工作法",
		MinSize:  Size{Width: 600, Height: 400},
		Size:     Size{Width: 600, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			// 初始化设置界面
			app.createInitUI(),
		},
		OnSizeChanged: func() {
			// 窗口大小变化时，更新布局
			app.mainWindow.Synchronize(func() {
				app.mainWindow.SetSize(walk.Size{Width: 600, Height: 400})
			})
		},
	}.Create()
}

// 创建初始化设置界面
func (app *EyeProtectionApp) createInitUI() Widget {
	return Composite{
		Layout: VBox{MarginsZero: false, Margins: Margins{Left: 50, Top: 30, Right: 50, Bottom: 30}},
		Children: []Widget{
			// 标题
			Label{
				Text:      "番茄工作法",
				Font:      Font{Family: "微软雅黑", PointSize: 16, Bold: true},
				Alignment: AlignHNearVNear,
			},
			VSpacer{},

			// 工作时间设置
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text:    "工作时间:",
						MinSize: Size{Width: 100},
					},
					NumberEdit{
						AssignTo: &app.workTimeEdit,
						Value:    25.0,
						MinValue: 1.0,
						MaxValue: 60.0,
						Decimals: 0,
						MinSize:  Size{Width: 80},
					},
					Label{
						Text:    "分钟",
						MinSize: Size{Width: 40},
					},
				},
			},
			VSpacer{},

			// 护眼时间设置
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text:    "休息时间:",
						MinSize: Size{Width: 100},
					},
					NumberEdit{
						AssignTo: &app.restTimeEdit,
						Value:    5.0,
						MinValue: 1.0,
						MaxValue: 60.0,
						Decimals: 0,
						MinSize:  Size{Width: 80},
					},
					Label{
						Text:    "分钟",
						MinSize: Size{Width: 40},
					},
				},
			},
			VSpacer{},

			// 起始时间设置
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text:    "起始时间:",
						MinSize: Size{Width: 100},
					},
					DateEdit{
						AssignTo: &app.startTimeEdit,
						Format:   "yyyy-MM-dd HH:mm",
						Date:     time.Now(),
						MinSize:  Size{Width: 150},
					},
					HSpacer{},
					Label{
						Text:    "结束时间:",
						MinSize: Size{Width: 100},
					},
					DateEdit{
						AssignTo: &app.endTimeEdit,
						Format:   "yyyy-MM-dd HH:mm",
						Date:     time.Now().Add(8 * time.Hour),
						MinSize:  Size{Width: 150},
					},
				},
			},
			VSpacer{},

			// 开始按钮
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &app.startButton,
						Text:     "开始",
						MinSize:  Size{Width: 100, Height: 30},
						OnClicked: func() {
							app.onStartButtonClicked()
						},
					},
				},
			},
		},
	}
}

// 创建工作时间倒计时界面
func (app *EyeProtectionApp) createWorkingUI() Widget {
	return Composite{
		Layout: VBox{MarginsZero: false, Margins: Margins{Left: 50, Top: 30, Right: 50, Bottom: 30}},
		Children: []Widget{
			// 标题
			Label{
				Text:      "番茄工作法-工作时间",
				Font:      Font{Family: "微软雅黑", PointSize: 16, Bold: true},
				Alignment: AlignHNearVNear,
			},
			VSpacer{},

			// 倒计时显示
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "工作时间倒计时: ",
						Font: Font{Family: "微软雅黑", PointSize: 14},
					},
					Label{
						AssignTo: &app.countdownLabel,
						Text:     "00:00:00",
						Font:     Font{Family: "微软雅黑", PointSize: 14},
					},
					Label{
						Text: "分钟",
						Font: Font{Family: "微软雅黑", PointSize: 14},
					},
				},
			},
			VSpacer{},

			// 结束按钮
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &app.endButton,
						Text:     "结束",
						MinSize:  Size{Width: 100, Height: 30},
						OnClicked: func() {
							app.onEndButtonClicked()
						},
					},
				},
			},
		},
	}
}

// 创建护眼时间倒计时界面
func (app *EyeProtectionApp) createRestingUI() Widget {
	return Composite{
		Layout: VBox{MarginsZero: false, Margins: Margins{Left: 50, Top: 30, Right: 50, Bottom: 30}},
		Children: []Widget{
			// 标题
			Label{
				Text:      "番茄工作法-休息时间",
				Font:      Font{Family: "微软雅黑", PointSize: 16, Bold: true},
				Alignment: AlignHNearVNear,
			},
			VSpacer{},

			// 倒计时显示
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					Label{
						Text:      "休息时间倒计时: ",
						TextColor: walk.RGB(47, 158, 68), // 绿色
						Font:      Font{Family: "微软雅黑", PointSize: 14},
					},
					Label{
						AssignTo:  &app.countdownLabel,
						Text:      "00:00:00",
						TextColor: walk.RGB(47, 158, 68), // 绿色
						Font:      Font{Family: "微软雅黑", PointSize: 14},
					},
					Label{
						Text:      "分钟",
						TextColor: walk.RGB(47, 158, 68), // 绿色
						Font:      Font{Family: "微软雅黑", PointSize: 14},
					},
				},
			},
			VSpacer{},

			// 结束按钮
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &app.endButton,
						Text:     "结束",
						MinSize:  Size{Width: 100, Height: 30},
						OnClicked: func() {
							app.onEndButtonClicked()
						},
					},
				},
			},
		},
	}
}

// 开始按钮点击事件
func (app *EyeProtectionApp) onStartButtonClicked() {
	// 获取设置的时间
	app.workMinutes = int(app.workTimeEdit.Value())
	app.restMinutes = int(app.restTimeEdit.Value())

	// 获取起始和结束时间
	now := time.Now()

	// 直接使用DateEdit的完整日期时间
	app.startTime = app.startTimeEdit.Date()
	app.endTime = app.endTimeEdit.Date()

	// 如果结束时间早于起始时间，认为是第二天
	if app.endTime.Before(app.startTime) {
		app.endTime = app.endTime.Add(24 * time.Hour)
	}

	// 验证当前时间是否在起始和结束时间范围内
	if now.Before(app.startTime) || now.After(app.endTime) {
		walk.MsgBox(app.mainWindow, "时间范围错误",
			"当前时间不在设定的起始和结束时间范围内，请调整时间设置。",
			walk.MsgBoxIconWarning)
		return
	}

	// 切换到工作时间界面
	app.switchToWorkingState()
}

// 结束按钮点击事件
func (app *EyeProtectionApp) onEndButtonClicked() {
	// 停止倒计时
	if app.countdownTimer != nil {
		app.countdownTimer.Stop()
		app.countdownTimer = nil
	}

	// 切换回初始界面
	app.switchToInitState()
}

// 切换到初始状态
func (app *EyeProtectionApp) switchToInitState() {
	// 更新状态
	app.state = StateInit

	// 清除主窗口内容
	app.clearMainWindow()

	// 设置初始界面
	initUI := app.createInitUI()
	if err := initUI.(Widget).Create(NewBuilder(app.mainWindow)); err != nil {
		walk.MsgBox(app.mainWindow, "错误", "创建初始界面失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	// 更新窗口标题
	app.mainWindow.SetTitle("番茄工作法")
}

// 切换到工作状态
func (app *EyeProtectionApp) switchToWorkingState() {
	// 更新状态
	app.state = StateWorking

	// 清除主窗口内容
	app.clearMainWindow()

	// 设置工作界面
	workingUI := app.createWorkingUI()
	if err := workingUI.(Widget).Create(NewBuilder(app.mainWindow)); err != nil {
		walk.MsgBox(app.mainWindow, "错误", "创建工作界面失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	// 更新窗口标题
	app.mainWindow.SetTitle("番茄工作法-工作时间")

	// 开始工作时间倒计时
	app.startWorkingCountdown()
}

// 切换到休息状态
func (app *EyeProtectionApp) switchToRestingState() {
	// 更新状态
	app.state = StateResting

	// 清除主窗口内容
	app.clearMainWindow()

	// 设置休息界面
	restingUI := app.createRestingUI()
	if err := restingUI.(Widget).Create(NewBuilder(app.mainWindow)); err != nil {
		walk.MsgBox(app.mainWindow, "错误", "创建休息界面失败: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	// 更新窗口标题
	app.mainWindow.SetTitle("番茄工作法-休息时间")

	// 开始休息时间倒计时
	app.startRestingCountdown()
}

// 清除主窗口内容
func (app *EyeProtectionApp) clearMainWindow() {
	// 移除所有子控件
	children := app.mainWindow.Children()
	for i := children.Len() - 1; i >= 0; i-- {
		if widget := children.At(i); widget != nil {
			widget.SetParent(nil)
			widget.Dispose()
		}
	}
}

// 切换到初始状态

// 开始工作时间倒计时
func (app *EyeProtectionApp) startWorkingCountdown() {
	// 设置倒计时结束时间
	app.countdownEnd = time.Now().Add(time.Duration(app.workMinutes) * time.Minute)

	// 更新倒计时显示
	app.updateCountdown()

	// 创建定时器，每秒更新一次
	app.countdownTimer = time.AfterFunc(time.Second, func() {
		app.mainWindow.Synchronize(func() {
			// 检查是否倒计时结束
			if time.Now().After(app.countdownEnd) {
				// 停止定时器
				app.countdownTimer.Stop()
				app.countdownTimer = nil

				// 弹出提示
				app.showWorkFinishedNotification()

				// 切换到休息状态
				app.switchToRestingState()
			} else {
				// 更新倒计时显示
				app.updateCountdown()

				// 继续倒计时
				app.countdownTimer.Reset(time.Second)
			}
		})
	})
}

// 开始休息时间倒计时
func (app *EyeProtectionApp) startRestingCountdown() {
	// 设置倒计时结束时间
	app.countdownEnd = time.Now().Add(time.Duration(app.restMinutes) * time.Minute)

	// 更新倒计时显示
	app.updateCountdown()

	// 创建定时器，每秒更新一次
	app.countdownTimer = time.AfterFunc(time.Second, func() {
		app.mainWindow.Synchronize(func() {
			// 检查是否倒计时结束
			if time.Now().After(app.countdownEnd) {
				// 停止定时器
				app.countdownTimer.Stop()
				app.countdownTimer = nil

				// 弹出提示
				app.showRestFinishedNotification()

				// 切换到工作状态
				app.switchToWorkingState()
			} else {
				// 更新倒计时显示
				app.updateCountdown()

				// 继续倒计时
				app.countdownTimer.Reset(time.Second)
			}
		})
	})
}

// 更新倒计时显示
func (app *EyeProtectionApp) updateCountdown() {
	// 计算剩余时间
	remaining := app.countdownEnd.Sub(time.Now())
	if remaining < 0 {
		remaining = 0
	}

	// 格式化剩余时间
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	seconds := int(remaining.Seconds()) % 60

	// 更新倒计时标签
	countdownText := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	app.countdownLabel.SetText(countdownText)
}

// 显示工作时间结束通知
func (app *EyeProtectionApp) showWorkFinishedNotification() {
	// 将窗口置顶
	app.bringWindowToFront()

	// 显示消息框
	walk.MsgBox(
		app.mainWindow,
		"工作时间结束",
		"您已经工作了 "+strconv.Itoa(app.workMinutes)+" 分钟，请休息一下！",
		walk.MsgBoxIconInformation,
	)
}

// 显示休息时间结束通知
func (app *EyeProtectionApp) showRestFinishedNotification() {
	// 将窗口置顶
	app.bringWindowToFront()

	// 显示消息框
	walk.MsgBox(
		app.mainWindow,
		"休息时间结束",
		"休息时间已结束，开始新的工作吧！",
		walk.MsgBoxIconInformation,
	)
}

// 将窗口置顶
func (app *EyeProtectionApp) bringWindowToFront() {
	// 获取窗口句柄
	hwnd := app.mainWindow.Handle()

	// 设置窗口为前台窗口
	win.SetForegroundWindow(hwnd)

	// 如果窗口最小化，则恢复
	if win.IsIconic(hwnd) {
		win.ShowWindow(hwnd, win.SW_RESTORE)
	}

	// 强制激活窗口
	win.SetActiveWindow(hwnd)

	// 设置窗口为顶层窗口
	win.SetWindowPos(
		hwnd,
		win.HWND_TOPMOST,
		0, 0, 0, 0,
		win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_SHOWWINDOW,
	)

	// 短暂延迟后恢复正常窗口层级
	go func() {
		time.Sleep(100 * time.Millisecond)
		win.SetWindowPos(
			hwnd,
			win.HWND_NOTOPMOST,
			0, 0, 0, 0,
			win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_SHOWWINDOW,
		)
	}()
}
