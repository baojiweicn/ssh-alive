package utils

import (
	"github.com/peterh/liner"
)

// Liner 封装 *liner.State, 提供更简便的操作
type Liner struct {
	State   *liner.State
	History *LineHistory

	tmode liner.ModeApplier
	lmode liner.ModeApplier

	paused bool
}

// NewLiner 返回 *Liner, 默认设置允许 Ctrl+C 结束
func NewLiner() *Liner {
	pl := &Liner{}
	pl.tmode, _ = liner.TerminalMode()

	line := liner.NewLiner()
	pl.lmode, _ = liner.TerminalMode()

	line.SetMultiLineMode(true)
	line.SetCtrlCAborts(true)

	pl.State = line

	return pl
}

// Pause 暂停服务
func (pl *Liner) Pause() error {
	if pl.paused {
		panic("Liner already paused")
	}

	pl.paused = true
	pl.DoWriteHistory()

	return pl.tmode.ApplyMode()
}

// Resume 恢复服务
func (pl *Liner) Resume() error {
	if !pl.paused {
		panic("Liner is not paused")
	}

	pl.paused = false

	return pl.lmode.ApplyMode()
}

// Close 关闭服务
func (pl *Liner) Close() (err error) {
	err = pl.State.Close()
	if err != nil {
		return err
	}

	if pl.History != nil && pl.History.historyFile != nil {
		return pl.History.historyFile.Close()
	}

	return nil
}