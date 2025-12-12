// FlexBox風レイアウトシステム
package parts

type FlexAlign int

const (
	FlexStart FlexAlign = iota
	FlexCenter
	FlexEnd
	FlexSpaceBetween // MainAxisのみ
	FlexSpaceAround  // MainAxisのみ
	FlexStretch      // CrossAxisのみ
)

// FlexLayoutHはFlexBoxの水平方向レイアウトを調整する
func FlexLayoutH(mainAlign FlexAlign, crossAlign FlexAlign, gap int) func(*Grouping) {
	return func(gr *Grouping) {
		con := gr.Control.GetControlBase()

		totalGrow := 0.0
		totalShrink := 0.0
		totalFlexBasis := 0

		// 可視の子コントロールを抽出、合計basisとgrow、shrinkを計算
		var visibleChildren []Control
		for _, c := range gr.Children {
			cb := c.GetControlBase()
			if cb.Visible {
				visibleChildren = append(visibleChildren, c)

				totalFlexBasis += cb.FlexBasisWidth
				totalGrow += cb.FlexGrow
				totalShrink += cb.FlexShrink
			}
		}

		if len(visibleChildren) == 0 {
			return
		}

		// 合計gap計算
		totalGap := gap * (len(visibleChildren) - 1)

		// 残りスペース計算。マイナスの時は足りてないことを意味する
		remainingSpace := con.Width - totalFlexBasis - totalGap

		// サイズ決定
		currentMainSizeTotal := 0
		for _, c := range visibleChildren {
			cb := c.GetControlBase()

			size := cb.FlexBasisWidth
			if remainingSpace > 0 && totalGrow > 0 {
				size += int(float64(remainingSpace) * (cb.FlexGrow / totalGrow))
			}
			if remainingSpace < 0 && totalShrink > 0 {
				size += int(float64(remainingSpace) * (cb.FlexShrink / totalShrink))
			}

			cb.Width = size
			currentMainSizeTotal += size
		}

		// 位置決定(MainAxis)
		totalSpace := con.Width - currentMainSizeTotal
		currentMainPos := 0
		tgap := gap

		// ギャップと開始位置計算
		switch mainAlign {
		case FlexCenter:
			currentMainPos = (totalSpace - totalGap) / 2
		case FlexEnd:
			currentMainPos = totalSpace - totalGap
		case FlexSpaceBetween:
			if len(visibleChildren) > 1 {
				tgap = totalSpace / (len(visibleChildren) - 1)
			}
		case FlexSpaceAround:
			tgap = totalSpace / len(visibleChildren)
			currentMainPos = tgap / 2
		}

		// 水平方向の位置計算
		for _, c := range visibleChildren {
			cb := c.GetControlBase()
			cb.X = currentMainPos
			currentMainPos += cb.Width + tgap
		}

		// 位置とサイズ決定(CrossAxis)
		for _, c := range visibleChildren {
			cb := c.GetControlBase()

			// CrossAxis位置の決定
			switch crossAlign {
			case FlexStart:
				cb.Y = 0
			case FlexStretch:
				cb.Y = 0
				cb.Height = con.Height
			case FlexCenter:
				cb.Y = (con.Height - cb.Height) / 2
			case FlexEnd:
				cb.Y = con.Height - cb.Height
			}

			// Layouter実装時の再帰呼び出し
			if ali, ok := c.(Layouter); ok {
				ali.Layout()
			}
		}
	}
}

// FlexLayoutVはFlexBoxの垂直方向レイアウトを調整する
func FlexLayoutV(mainAlign FlexAlign, crossAlign FlexAlign, gap int) func(*Grouping) {
	return func(gr *Grouping) {
		con := gr.Control.GetControlBase()

		totalGrow := 0.0
		totalShrink := 0.0
		totalFlexBasis := 0

		// 可視の子コントロールを抽出、合計basisとgrow、shrinkを計算
		var visibleChildren []Control
		for _, c := range gr.Children {
			cb := c.GetControlBase()
			if cb.Visible {
				visibleChildren = append(visibleChildren, c)

				totalFlexBasis += cb.FlexBasisHeight
				totalGrow += cb.FlexGrow
				totalShrink += cb.FlexShrink
			}
		}

		if len(visibleChildren) == 0 {
			return
		}

		// 合計gap計算
		totalGap := gap * (len(visibleChildren) - 1)

		// 残りスペース計算。マイナスの時は足りてないことを意味する
		remainingSpace := con.Height - totalFlexBasis - totalGap

		// サイズ決定
		currentMainSizeTotal := 0
		for _, c := range visibleChildren {
			cb := c.GetControlBase()

			size := cb.FlexBasisHeight
			if remainingSpace > 0 && totalGrow > 0 {
				size += int(float64(remainingSpace) * (cb.FlexGrow / totalGrow))
			}
			if remainingSpace < 0 && totalShrink > 0 {
				size += int(float64(remainingSpace) * (cb.FlexShrink / totalShrink))
			}

			cb.Height = size
			currentMainSizeTotal += size
		}

		// 位置決定(MainAxis)
		totalSpace := con.Height - currentMainSizeTotal
		currentMainPos := 0
		tgap := gap

		// ギャップと開始位置計算
		switch mainAlign {
		case FlexCenter:
			currentMainPos = (totalSpace - totalGap) / 2
		case FlexEnd:
			currentMainPos = totalSpace - totalGap
		case FlexSpaceBetween:
			if len(visibleChildren) > 1 {
				tgap = totalSpace / (len(visibleChildren) - 1)
			}
		case FlexSpaceAround:
			tgap = totalSpace / len(visibleChildren)
			currentMainPos = tgap / 2
		}

		// 垂直方向の位置計算
		for _, c := range visibleChildren {
			cb := c.GetControlBase()
			cb.Y = currentMainPos
			currentMainPos += cb.Height + tgap
		}

		// 位置とサイズ決定(CrossAxis)
		for _, c := range visibleChildren {
			cb := c.GetControlBase()

			// CrossAxis位置の決定
			switch crossAlign {
			case FlexStart:
				cb.X = 0
			case FlexStretch:
				cb.X = 0
				cb.Width = con.Width
			case FlexCenter:
				cb.X = (con.Width - cb.Width) / 2
			case FlexEnd:
				cb.X = con.Width - cb.Width
			}

			// Layouter実装時の再帰呼び出し
			if ali, ok := c.(Layouter); ok {
				ali.Layout()
			}
		}
	}
}
