package ui

type AutoLayoutInterface interface {
	Layout()
}

// AutoLayoutVはGroupingとセットで使用してコントロールを縦に並べる機能
type AutoLayoutV struct {
	grouping *Grouping
}

// AutoLayoutV生成
func NewAutoLayoutV(g *Grouping) *AutoLayoutV {
	a := &AutoLayoutV{grouping: g}
	return a
}

// オートレイアウト処理
// リサイズ時に呼ぶことを想定している
func (a *AutoLayoutV) Layout() {
	gr := a.grouping
	con := gr.Control.GetControlBase()

	// 範囲取得。この中にコントロールを配置する
	maxWidth := con.Width
	maxHeight := con.Height

	// コントロールの数
	count := 0

	// コントロールサイズの合計算出
	total := 0
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			total += ccon.Height
			count++
		}
	}

	// コントロールサイズを除いたあまりサイズから間隔を算出
	r := (maxHeight - total) / (count + 1)

	y := 0
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			ccon.X = maxWidth/2 - ccon.Width/2
			y += r
			ccon.Y = y
			y += ccon.Height
		}
	}
}

// AutoLayoutFitはGrouping全体にフィットするように自動サイズ調整する機能
// AutoResizableを埋め込んだコントロールだけリサイズ可能
type AutoLayoutFitH struct {
	grouping *Grouping
}

// AutoLayoutFit生成
func NewAutoLayoutFitH(g *Grouping) *AutoLayoutFitH {
	a := &AutoLayoutFitH{grouping: g}
	return a
}

// オートレイアウト処理(水平方向)
func (a *AutoLayoutFitH) Layout() {
	gr := a.grouping
	con := gr.Control.GetControlBase()

	// コントロールの数
	count := 0

	// 固定サイズのコントロールサイズ合計算出
	total := 0
	var arc Control = nil
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			// AutoResizeが設定されていないコントロールを合計する
			if ccon.AutoResizable {
				arc = c
			} else {
				total += ccon.Width
				count++
			}
		}
	}

	// コントロールサイズを除いたあまりサイズをAutoResizableコントロールに設定
	// このため現状ではAutoResizableは1つに限定
	if arc != nil {
		arc.GetControlBase().Width = con.Width - total

		// AutoLayoutInterfaceを実装していたらオートレイアウト実行
		if ali, ok := arc.(AutoLayoutInterface); ok {
			ali.Layout()
		}
	}

	// コントロールサイズをもとに左詰めで位置を設定
	x := 0
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			ccon.X = x
			ccon.Height = con.Height // 縦方向をどうするかは悩み中
			y += ccon.Width
		}
	}
}

// AutoLayoutFitはGrouping全体にフィットするように自動サイズ調整する機能
// AutoResizableを埋め込んだコントロールだけリサイズ可能
type AutoLayoutFitV struct {
	grouping *Grouping
}

// AutoLayoutFit生成
func NewAutoLayoutFitV(g *Grouping) *AutoLayoutFitV {
	a := &AutoLayoutFitV{grouping: g}
	return a
}

// オートレイアウト処理(垂直方向)
func (a *AutoLayoutFitV) Layout() {
	gr := a.grouping
	con := gr.Control.GetControlBase()

	// コントロールの数
	count := 0

	// 固定サイズのコントロールサイズ合計算出
	total := 0
	var arc Control = nil
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			if ccon.AutoResizable {
				arc = c
			} else {
				total += ccon.Height
				count++
			}
		}
	}

	// コントロールサイズを除いたあまりサイズをAutoResizableコントロールに設定
	// このため現状ではAutoResizableは1つに限定
	if arc != nil {
		arc.GetControlBase().Height = con.Height - total

		// AutoLayoutInterfaceを実装していたらオートレイアウト実行
		if ali, ok := arc.(AutoLayoutInterface); ok {
			ali.Layout()
		}
	}

	// コントロールサイズをもとに左詰めで位置を設定
	y := 0
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			ccon.Y = y
			ccon.Width = con.Width // 横方向をどうするかは悩み中
			y += ccon.Height
		}
	}
}
