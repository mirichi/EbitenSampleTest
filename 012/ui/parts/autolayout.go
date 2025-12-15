package parts

// DefaultLayoutはGroupingのデフォルトのレイアウト処理
// AutoResizableな子コントロールのサイズを親に合わせて調整する
func DefaultLayout(g *Grouping) {
	for _, c := range g.Children {
		cb := c.GetControlBase()

		if cb.Visible {
			// AutoResizableがあった場合
			if cb.AutoResizable {
				// Groupingコントロールのサイズに合わせたサイズに更新する
				pcb := g.Control.GetControlBase()
				cb.Width = pcb.Width
				cb.Height = pcb.Height
			}

			// 配下のレイアウト更新
			if al, ok := c.(Layouter); ok {
				al.Layout()
			}
		}
	}
}

// オートレイアウト処理(水平方向)
func AutoLayoutFitH(margin int) func(*Grouping) {
	return func(gr *Grouping) {
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
			arc.GetControlBase().Width = con.Width - margin*2 - total
		}

		// コントロールサイズをもとに左詰めで位置を設定
		x := margin
		for _, c := range gr.Children {
			ccon := c.GetControlBase()
			if ccon.Visible {
				ccon.X = x
				ccon.Y = margin
				ccon.Height = con.Height - margin*2
				x += ccon.Width

				// AutoLayoutInterfaceを実装していたらオートレイアウト実行
				if ali, ok := c.(Layouter); ok {
					ali.Layout()
				}
			}
		}
	}
}

// オートレイアウト処理(垂直方向)
func AutoLayoutFitV(margin int) func(*Grouping) {
	return func(gr *Grouping) {
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
			arc.GetControlBase().Height = con.Height - margin*2 - total
		}

		// コントロールサイズをもとに上詰めで位置を設定
		y := margin
		for _, c := range gr.Children {
			ccon := c.GetControlBase()
			if ccon.Visible {
				ccon.Y = y
				ccon.X = margin
				ccon.Width = con.Width - margin*2
				y += ccon.Height

				// AutoLayoutInterfaceを実装していたらオートレイアウト実行
				if ali, ok := c.(Layouter); ok {
					ali.Layout()
				}
			}
		}
	}
}
