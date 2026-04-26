package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/vetalpaprotsky/gyla/game"
)

const (
	cardWidth  = 50
	cardHeight = 70
	cardGap    = 5

	screenWidth  = 640
	screenHeight = 480
)

var (
	colorWhite     = color.RGBA{255, 255, 255, 255}
	colorBlack     = color.RGBA{0, 0, 0, 255}
	colorRed       = color.RGBA{200, 30, 30, 255}
	colorCardBack  = color.RGBA{40, 60, 150, 255}
	colorPlayable  = color.RGBA{220, 40, 40, 255}
	colorTrumpMark = color.RGBA{255, 215, 0, 255}
	colorBg        = color.RGBA{30, 100, 50, 255}
	colorScoreBg   = color.RGBA{0, 0, 0, 160}
	colorLabelBg   = color.RGBA{0, 0, 0, 180}
	colorGray      = color.RGBA{180, 180, 180, 255}
)

// rankString converts a game.Rank to its display string.
func rankString(r game.Rank) string {
	switch r {
	case game.AceRank:
		return "A"
	case game.KingRank:
		return "K"
	case game.QueenRank:
		return "Q"
	case game.JackRank:
		return "J"
	case game.TenRank:
		return "10"
	case game.NineRank:
		return "9"
	case game.EightRank:
		return "8"
	case game.SevenRank:
		return "7"
	case game.SixRank:
		return "6"
	default:
		return "?"
	}
}

// suitString converts a game.Suit to its display symbol.
func suitString(s game.Suit) string {
	switch s {
	case game.ClubsSuit:
		return "\u2663"
	case game.SpadesSuit:
		return "\u2660"
	case game.HeartsSuit:
		return "\u2665"
	case game.DiamondsSuit:
		return "\u2666"
	default:
		return "?"
	}
}

// suitColor returns red for Hearts/Diamonds, dark gray for Clubs/Spades.
func suitColor(s game.Suit) color.RGBA {
	if s == game.HeartsSuit || s == game.DiamondsSuit {
		return colorRed
	}
	return color.RGBA{40, 40, 50, 255}
}

// cardLabel returns the combined rank+suit label, e.g. "A♠".
func cardLabel(c game.Card) string {
	return rankString(c.Rank) + suitString(c.Suit)
}

// drawSuitSymbol draws a suit symbol as a small colored square at (x, y).
// This is a simple placeholder — a colored rectangle representing the suit.
func drawSuitSymbol(screen *ebiten.Image, s game.Suit, x, y, size int) {
	clr := suitColor(s)
	switch s {
	case game.HeartsSuit:
		// Heart: draw a red triangle-ish shape (top two bumps + bottom point)
		// Simplified: filled red square with a small notch
		sym := ebiten.NewImage(size, size)
		sym.Fill(clr)
		// Cut top-center notch to hint at heart shape
		notch := ebiten.NewImage(size/5, size/4)
		notch.Fill(colorWhite)
		opN := &ebiten.DrawImageOptions{}
		opN.GeoM.Translate(float64(size/2-size/10), 0)
		sym.DrawImage(notch, opN)
		// Cut bottom-left and bottom-right corners
		corner := ebiten.NewImage(size/4, size/4)
		corner.Fill(colorWhite)
		opC := &ebiten.DrawImageOptions{}
		sym.DrawImage(corner, opC)
		opC2 := &ebiten.DrawImageOptions{}
		opC2.GeoM.Translate(float64(size-size/4), 0)
		sym.DrawImage(corner, opC2)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(sym, op)

	case game.DiamondsSuit:
		// Diamond: draw a red rotated square (simplified as a smaller centered square)
		sym := ebiten.NewImage(size, size)
		// Draw diamond as centered smaller filled square with cut corners
		inner := ebiten.NewImage(size*2/3, size*2/3)
		inner.Fill(clr)
		opI := &ebiten.DrawImageOptions{}
		opI.GeoM.Translate(float64(size/6), float64(size/6))
		sym.DrawImage(inner, opI)
		// Vertical bar through center for diamond shape
		vbar := ebiten.NewImage(size/4, size)
		vbar.Fill(clr)
		opV := &ebiten.DrawImageOptions{}
		opV.GeoM.Translate(float64(size/2-size/8), 0)
		sym.DrawImage(vbar, opV)
		// Horizontal bar
		hbar := ebiten.NewImage(size, size/4)
		hbar.Fill(clr)
		opH := &ebiten.DrawImageOptions{}
		opH.GeoM.Translate(0, float64(size/2-size/8))
		sym.DrawImage(hbar, opH)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(sym, op)

	case game.SpadesSuit:
		// Spade: dark filled shape — triangle on top, stem below
		sym := ebiten.NewImage(size, size)
		// Main body
		body := ebiten.NewImage(size*2/3, size*2/3)
		body.Fill(clr)
		opB := &ebiten.DrawImageOptions{}
		opB.GeoM.Translate(float64(size/6), 0)
		sym.DrawImage(body, opB)
		// Stem
		stem := ebiten.NewImage(size/4, size/3)
		stem.Fill(clr)
		opS := &ebiten.DrawImageOptions{}
		opS.GeoM.Translate(float64(size/2-size/8), float64(size*2/3))
		sym.DrawImage(stem, opS)
		// Top point
		top := ebiten.NewImage(size/3, size/3)
		top.Fill(clr)
		opT := &ebiten.DrawImageOptions{}
		opT.GeoM.Translate(float64(size/3), 0)
		sym.DrawImage(top, opT)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(sym, op)

	case game.ClubsSuit:
		// Club: three circles + stem (simplified as three small squares + stem)
		sym := ebiten.NewImage(size, size)
		dot := size / 3
		// Top circle
		c1 := ebiten.NewImage(dot, dot)
		c1.Fill(clr)
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Translate(float64(size/2-dot/2), 0)
		sym.DrawImage(c1, op1)
		// Left circle
		c2 := ebiten.NewImage(dot, dot)
		c2.Fill(clr)
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(0, float64(dot))
		sym.DrawImage(c2, op2)
		// Right circle
		c3 := ebiten.NewImage(dot, dot)
		c3.Fill(clr)
		op3 := &ebiten.DrawImageOptions{}
		op3.GeoM.Translate(float64(size-dot), float64(dot))
		sym.DrawImage(c3, op3)
		// Stem
		stem := ebiten.NewImage(size/4, size/3)
		stem.Fill(clr)
		opSt := &ebiten.DrawImageOptions{}
		opSt.GeoM.Translate(float64(size/2-size/8), float64(size*2/3))
		sym.DrawImage(stem, opSt)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(sym, op)
	}
}

// drawFaceUpCard draws a face-up card at position (x, y) on the screen.
// Cards are white with colored rank text and a suit symbol drawn on them.
func drawFaceUpCard(screen *ebiten.Image, c game.Card, isPlayable bool, x, y int) {
	// If playable, draw green border behind the card.
	if isPlayable {
		border := ebiten.NewImage(cardWidth+4, cardHeight+4)
		border.Fill(colorPlayable)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x-2), float64(y-2))
		screen.DrawImage(border, op)
	}

	// Thin dark border around the card.
	borderImg := ebiten.NewImage(cardWidth+2, cardHeight+2)
	borderImg.Fill(colorBlack)
	opBrd := &ebiten.DrawImageOptions{}
	opBrd.GeoM.Translate(float64(x-1), float64(y-1))
	screen.DrawImage(borderImg, opBrd)

	// White card background.
	cardImg := ebiten.NewImage(cardWidth, cardHeight)
	cardImg.Fill(colorWhite)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(cardImg, op)

	// Trump indicator: gold square in top-right corner.
	if c.IsTrump {
		trumpDot := ebiten.NewImage(8, 8)
		trumpDot.Fill(colorTrumpMark)
		opT := &ebiten.DrawImageOptions{}
		opT.GeoM.Translate(float64(x+cardWidth-10), float64(y+2))
		screen.DrawImage(trumpDot, opT)
	}

	// Draw rank text on a colored background strip so it's readable.
	clr := suitColor(c.Suit)
	label := rankString(c.Rank)
	rankBg := ebiten.NewImage(cardWidth-4, 14)
	rankBg.Fill(clr)
	opRB := &ebiten.DrawImageOptions{}
	opRB.GeoM.Translate(float64(x+2), float64(y+2))
	screen.DrawImage(rankBg, opRB)
	ebitenutil.DebugPrintAt(screen, label, x+4, y+3)

	// Draw suit symbol in the center of the card.
	symbolSize := 18
	symX := x + cardWidth/2 - symbolSize/2
	symY := y + cardHeight/2 - symbolSize/2 + 4
	drawSuitSymbol(screen, c.Suit, symX, symY, symbolSize)

	// Small suit text at bottom-left.
	suitBg := ebiten.NewImage(14, 14)
	suitBg.Fill(clr)
	opSB := &ebiten.DrawImageOptions{}
	opSB.GeoM.Translate(float64(x+2), float64(y+cardHeight-16))
	screen.DrawImage(suitBg, opSB)
	ebitenutil.DebugPrintAt(screen, suitString(c.Suit), x+4, y+cardHeight-15)
}

// drawFaceDownCard draws a card back (blue rectangle) at position (x, y).
func drawFaceDownCard(screen *ebiten.Image, x, y, w, h int) {
	// Dark border.
	borderImg := ebiten.NewImage(w+2, h+2)
	borderImg.Fill(colorBlack)
	opBrd := &ebiten.DrawImageOptions{}
	opBrd.GeoM.Translate(float64(x-1), float64(y-1))
	screen.DrawImage(borderImg, opBrd)

	cardImg := ebiten.NewImage(w, h)
	cardImg.Fill(colorCardBack)

	// Inner pattern: slightly lighter rectangle.
	inner := ebiten.NewImage(w-6, h-6)
	inner.Fill(color.RGBA{60, 80, 180, 255})
	opIn := &ebiten.DrawImageOptions{}
	opIn.GeoM.Translate(3, 3)
	cardImg.DrawImage(inner, opIn)

	// Diamond pattern in center.
	diamond := ebiten.NewImage(8, 8)
	diamond.Fill(color.RGBA{80, 100, 200, 255})
	opD := &ebiten.DrawImageOptions{}
	opD.GeoM.Translate(float64(w/2-4), float64(h/2-4))
	cardImg.DrawImage(diamond, opD)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(cardImg, op)
}

// drawGame renders the full game UI given the current game view.
func drawGame(screen *ebiten.Image, view *game.GameView) {
	// Fill the background with green (card table).
	screen.Fill(colorBg)

	if view == nil {
		ebitenutil.DebugPrintAt(screen, "Waiting for game...", 270, 230)
		return
	}

	round := &view.Round
	hand := round.Hand

	// --- Player's hand (bottom center, face up) ---
	playerCards := hand.Cards
	totalWidth := len(playerCards)*(cardWidth+cardGap) - cardGap
	if totalWidth < 0 {
		totalWidth = 0
	}
	startX := (screenWidth - totalWidth) / 2
	handY := screenHeight - cardHeight - 30

	for i, hc := range playerCards {
		cx := startX + i*(cardWidth+cardGap)
		drawFaceUpCard(screen, hc.Card, hc.IsPlayable, cx, handY)
	}

	// Player name below hand.
	playerNameX := (screenWidth - len(view.You.PlayerName)*6) / 2
	ebitenutil.DebugPrintAt(screen, view.You.PlayerName, playerNameX, screenHeight-22)

	// --- Teammate's hand (top center, face down) ---
	teammateCount := round.TeammateHand
	if teammateCount > 0 {
		smallW := 30
		smallH := 42
		tmTotalW := teammateCount*(smallW+cardGap) - cardGap
		tmStartX := (screenWidth - tmTotalW) / 2
		tmY := 25

		for i := 0; i < teammateCount; i++ {
			cx := tmStartX + i*(smallW+cardGap)
			drawFaceDownCard(screen, cx, tmY, smallW, smallH)
		}
	}
	tmNameX := (screenWidth - len(view.Teammate.PlayerName)*6) / 2
	ebitenutil.DebugPrintAt(screen, view.Teammate.PlayerName, tmNameX, 10)

	// --- Left opponent (left side, face down, stacked vertically) ---
	leftCount := round.LeftOpponentHand
	if leftCount > 0 {
		smallW := 30
		smallH := 42
		lx := 15
		lStartY := (screenHeight - leftCount*(smallH+cardGap) + cardGap) / 2

		for i := 0; i < leftCount; i++ {
			cy := lStartY + i*(smallH+cardGap)
			drawFaceDownCard(screen, lx, cy, smallW, smallH)
		}
	}
	ebitenutil.DebugPrintAt(screen, view.LeftOpponent.PlayerName, 10, screenHeight/2-60)

	// --- Right opponent (right side, face down, stacked vertically) ---
	rightCount := round.RightOpponentHand
	if rightCount > 0 {
		smallW := 30
		smallH := 42
		rx := screenWidth - smallW - 15
		rStartY := (screenHeight - rightCount*(smallH+cardGap) + cardGap) / 2

		for i := 0; i < rightCount; i++ {
			cy := rStartY + i*(smallH+cardGap)
			drawFaceDownCard(screen, rx, cy, smallW, smallH)
		}
	}
	rpName := view.RightOpponent.PlayerName
	ebitenutil.DebugPrintAt(screen, rpName, screenWidth-len(rpName)*6-10, screenHeight/2-60)

	// --- Trick area (center) ---
	drawTrickArea(screen, view)

	// --- Trump indicator / selector (top-left) ---
	if round.Trump != 0 {
		trumpLabel := "Trump: " + suitString(round.Trump)
		labelBg := ebiten.NewImage(len(trumpLabel)*7+8, 18)
		labelBg.Fill(colorLabelBg)
		opLbl := &ebiten.DrawImageOptions{}
		opLbl.GeoM.Translate(8, 0)
		screen.DrawImage(labelBg, opLbl)
		ebitenutil.DebugPrintAt(screen, trumpLabel, 12, 2)
		// Also draw the suit symbol next to the label.
		drawSuitSymbol(screen, round.Trump, len(trumpLabel)*7+12, 1, 14)
	} else if view.NextAction.Name == "assign_trump" && view.NextAction.Player == view.You.Player {
		labelBg := ebiten.NewImage(100, 18)
		labelBg.Fill(colorLabelBg)
		opLbl := &ebiten.DrawImageOptions{}
		opLbl.GeoM.Translate(8, 0)
		screen.DrawImage(labelBg, opLbl)
		ebitenutil.DebugPrintAt(screen, "Choose trump:", 12, 2)
		drawTrumpSelector(screen)
	} else if round.Trump == 0 {
		labelBg := ebiten.NewImage(140, 18)
		labelBg.Fill(colorLabelBg)
		opLbl := &ebiten.DrawImageOptions{}
		opLbl.GeoM.Translate(8, 0)
		screen.DrawImage(labelBg, opLbl)
		ebitenutil.DebugPrintAt(screen, "Waiting for trump...", 12, 2)
	}

	// --- Score display (top-right) ---
	drawScore(screen, view)

	// --- Round / trick info ---
	info := fmt.Sprintf("Round %d  Trick %d", round.Number, currentTrickNumber(round))
	infoBg := ebiten.NewImage(len(info)*7+8, 16)
	infoBg.Fill(colorLabelBg)
	opInfo := &ebiten.DrawImageOptions{}
	opInfo.GeoM.Translate(float64(screenWidth/2-len(info)*7/2-4), float64(screenHeight/2+63))
	screen.DrawImage(infoBg, opInfo)
	ebitenutil.DebugPrintAt(screen, info, screenWidth/2-len(info)*3, screenHeight/2+65)

	// --- Next action hint ---
	if view.NextAction.Name != "" {
		hint := fmt.Sprintf("Next: %s", view.NextAction.Name)
		hintBg := ebiten.NewImage(len(hint)*7+8, 16)
		hintBg.Fill(colorLabelBg)
		opHint := &ebiten.DrawImageOptions{}
		opHint.GeoM.Translate(float64(screenWidth/2-len(hint)*7/2-4), float64(screenHeight/2+81))
		screen.DrawImage(hintBg, opHint)
		ebitenutil.DebugPrintAt(screen, hint, screenWidth/2-len(hint)*3, screenHeight/2+83)
	}
}

const (
	trumpBtnSize = 30
	trumpBtnGap  = 10
	trumpBtnX    = 10
	trumpBtnY    = 22
)

// drawTrumpSelector draws four clickable suit buttons for trump selection.
func drawTrumpSelector(screen *ebiten.Image) {
	suits := [4]game.Suit{game.ClubsSuit, game.SpadesSuit, game.HeartsSuit, game.DiamondsSuit}

	for i, s := range suits {
		x := trumpBtnX + i*(trumpBtnSize+trumpBtnGap)
		y := trumpBtnY

		// Border
		borderImg := ebiten.NewImage(trumpBtnSize+2, trumpBtnSize+2)
		borderImg.Fill(colorBlack)
		opBorder := &ebiten.DrawImageOptions{}
		opBorder.GeoM.Translate(float64(x-1), float64(y-1))
		screen.DrawImage(borderImg, opBorder)

		// White button background.
		btn := ebiten.NewImage(trumpBtnSize, trumpBtnSize)
		btn.Fill(colorWhite)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(btn, op)

		// Draw suit symbol centered in the button.
		symSize := 16
		drawSuitSymbol(screen, s, x+trumpBtnSize/2-symSize/2, y+trumpBtnSize/2-symSize/2, symSize)
	}
}

// trumpButtonHit checks if a mouse click at (mx, my) hits one of the trump
// selector buttons. Returns the suit and true if hit, zero and false otherwise.
func trumpButtonHit(mx, my int) (game.Suit, bool) {
	suits := [4]game.Suit{game.ClubsSuit, game.SpadesSuit, game.HeartsSuit, game.DiamondsSuit}

	for i, s := range suits {
		x := trumpBtnX + i*(trumpBtnSize+trumpBtnGap)
		y := trumpBtnY

		if mx >= x && mx <= x+trumpBtnSize && my >= y && my <= y+trumpBtnSize {
			return s, true
		}
	}

	return 0, false
}

// handCardHit checks if a mouse click at (mx, my) hits one of the player's
// hand cards. Returns the HandCard and true if hit.
func handCardHit(mx, my int, view *game.GameView) (game.HandCard, bool) {
	playerCards := view.Round.Hand.Cards
	totalWidth := len(playerCards)*(cardWidth+cardGap) - cardGap
	if totalWidth < 0 {
		totalWidth = 0
	}
	startX := (screenWidth - totalWidth) / 2
	handY := screenHeight - cardHeight - 30

	for i, hc := range playerCards {
		cx := startX + i*(cardWidth+cardGap)
		if mx >= cx && mx <= cx+cardWidth && my >= handY && my <= handY+cardHeight {
			if hc.IsPlayable {
				return hc, true
			}
			return hc, false
		}
	}

	return game.HandCard{}, false
}

// trickCardPosition returns the (x, y) position for a played card in the trick area,
// based on which player played it relative to "you".
func trickCardPosition(player game.Player, you game.Player) (int, int) {
	centerX := screenWidth/2 - cardWidth/2
	centerY := screenHeight/2 - cardHeight/2

	switch {
	case player == you:
		return centerX, centerY + 40
	case player == leftOpponent(you):
		return centerX - 60, centerY
	case player == teammate(you):
		return centerX, centerY - 40
	case player == rightOpponent(you):
		return centerX + 60, centerY
	default:
		return centerX, centerY
	}
}

func leftOpponent(p game.Player) game.Player {
	switch p {
	case game.Player1:
		return game.Player4
	case game.Player2:
		return game.Player1
	case game.Player3:
		return game.Player2
	case game.Player4:
		return game.Player3
	}
	return game.Player1
}

func teammate(p game.Player) game.Player {
	switch p {
	case game.Player1:
		return game.Player3
	case game.Player2:
		return game.Player4
	case game.Player3:
		return game.Player1
	case game.Player4:
		return game.Player2
	}
	return game.Player1
}

func rightOpponent(p game.Player) game.Player {
	switch p {
	case game.Player1:
		return game.Player2
	case game.Player2:
		return game.Player3
	case game.Player3:
		return game.Player4
	case game.Player4:
		return game.Player1
	}
	return game.Player1
}

// drawTrickArea draws the cards played in the current trick as face-up cards.
func drawTrickArea(screen *ebiten.Image, view *game.GameView) {
	round := &view.Round
	if len(round.Tricks) == 0 {
		return
	}

	currentTrick := round.Tricks[len(round.Tricks)-1]
	you := view.You.Player

	for _, pc := range currentTrick.PlayedCards {
		x, y := trickCardPosition(pc.Player, you)
		drawFaceUpCard(screen, pc.Card, false, x, y)
	}
}

// drawScore draws the score in the top-right corner.
func drawScore(screen *ebiten.Image, view *game.GameView) {
	stats := view.Stats
	t1 := stats.Points[game.Team1]
	t2 := stats.Points[game.Team2]

	yourTeam := view.You.TeamName
	oppTeam := view.LeftOpponent.TeamName

	scoreBg := ebiten.NewImage(130, 34)
	scoreBg.Fill(colorScoreBg)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screenWidth-140), 2)
	screen.DrawImage(scoreBg, op)

	line1 := fmt.Sprintf("%s: %d", yourTeam, teamPoints(view.You.Team, t1, t2))
	line2 := fmt.Sprintf("%s: %d", oppTeam, teamPoints(view.LeftOpponent.Team, t1, t2))
	ebitenutil.DebugPrintAt(screen, line1, screenWidth-135, 4)
	ebitenutil.DebugPrintAt(screen, line2, screenWidth-135, 18)
}

func teamPoints(t game.Team, t1, t2 int) int {
	if t == game.Team1 {
		return t1
	}
	return t2
}

// currentTrickNumber returns the current trick number from the round view.
func currentTrickNumber(r *game.RoundView) int {
	if len(r.Tricks) == 0 {
		return 0
	}
	return r.Tricks[len(r.Tricks)-1].Number
}
