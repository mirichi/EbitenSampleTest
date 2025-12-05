package parts

type AutoLayoutInterface interface {
	Layout()
}

// Layoutはオートレイアウト処理を実行する
// AutoLayoutVでは親コントロールのサイズに合わせて、子コントロールを等間隔に垂直配置する
// リサイズ時に呼び出されることを想定する
func AutoLayoutV(gr *Grouping) {
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

// オートレイアウト処理(水平方向)
func AutoLayoutFitH(gr *Grouping) {
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
	}

	// コントロールサイズをもとに左詰めで位置を設定
	x := 0
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			ccon.X = x
			ccon.Height = con.Height // TODO: 縦方向のサイズ調整ポリシーを検討する
			x += ccon.Width

			// AutoLayoutInterfaceを実装していたらオートレイアウト実行
			if ali, ok := c.(AutoLayoutInterface); ok {
				ali.Layout()
			}
		}
	}
}

// オートレイアウト処理(垂直方向)
func AutoLayoutFitV(gr *Grouping) {
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
	}

	// コントロールサイズをもとに左詰めで位置を設定
	y := 0
	for _, c := range gr.Children {
		ccon := c.GetControlBase()
		if ccon.Visible {
			ccon.Y = y
			ccon.Width = con.Width // TODO: 横方向のサイズ調整ポリシーを検討する
			y += ccon.Height

			// AutoLayoutInterfaceを実装していたらオートレイアウト実行
			if ali, ok := c.(AutoLayoutInterface); ok {
				ali.Layout()
			}
		}
	}
}
