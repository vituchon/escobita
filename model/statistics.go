package model

// TODO : needs to implement what happens in draw, right now only one winner can be and would be "the first" that the algorithm find
type PlayerStatictics struct {
	CardsTakenCount int  `json:"cardsTakenCount"`
	EscobitasCount  int  `json:"escobitasCount"`
	SeventiesScore  int  `json:"seventiesScore"`
	HasGoldSeven    bool `json:"hasGoldSeven"`
	GoldCardsCount  int  `json:"goldCardsCount"`
}

type StaticticsByPlayer map[Player]PlayerStatictics

func (match Match) CalculateStaticticsByPlayer() StaticticsByPlayer {
	staticticsByPlayer := make(StaticticsByPlayer)
	for _, player := range match.Players {
		staticticsByPlayer[player] = match.doCalculateStatictics(player)
	}
	return staticticsByPlayer
}

func (match Match) CalculateStatictics(player Player) PlayerStatictics {
	return match.doCalculateStatictics(player)
}

func (match Match) doCalculateStatictics(player Player) PlayerStatictics {
	cardsTakenCount := countCardsTaken(player, match)
	escobitasCount := countEscobitas(player, match)
	playerTakenCards := match.Cards.ByPlayer[player].Taken
	seventiesScore := calculateSeventiesScore(playerTakenCards)
	hasGoldSeven := DetermineIsGoldenSevenIsUsed(playerTakenCards)
	goldCardsCount := CountGoldenSuitCards(playerTakenCards)

	ps := PlayerStatictics{
		CardsTakenCount: cardsTakenCount,
		EscobitasCount:  escobitasCount,
		SeventiesScore:  seventiesScore,
		HasGoldSeven:    hasGoldSeven,
		GoldCardsCount:  goldCardsCount,
	}
	return ps
}

func countCardsTaken(player Player, match Match) int {
	return len(match.Cards.ByPlayer[player].Taken)
}

var boolToInt map[bool]int = map[bool]int{
	true:  1,
	false: 0,
}

func countEscobitas(player Player, match Match) int {
	escobitasCount := 0
	//fmt.Printf("match.ActionsByPlayer[player]=%+v", match.ActionsByPlayer[player])
	for _, action := range match.ActionsByPlayer[player] {
		escobitasCount += boolToInt[action.IsEscobita()]
	}
	return escobitasCount
}

func calculateSeventiesScore(cards Deck) int {
	splittedByRank := cards.SplitByRank()
	seventyCardBySuit := make(map[Suit]*Card)
	for suit, cards := range splittedByRank {
		seventyCard := cards.getLeftCloserToRank(7)
		if seventyCard != nil {
			seventyCardBySuit[suit] = seventyCard
		}
	}

	score := 0
	var factors []int = []int{1, 2, 4, 8, 16, 32, 64}
	for _, seventyCard := range seventyCardBySuit {
		score += factors[seventyCard.Rank-1] // ranks starts at 1
	}
	return score
}

func DetermineIsGoldenSevenIsUsed(cards []Card) bool {
	for _, card := range cards {
		if card.IsGoldenSeven() {
			return true
		}
	}
	return false
}

func CountSevenRankCards(cards []Card) int {
	count := 0
	for _, card := range cards {
		count += boolToInt[card.IsSevenRank()]
	}
	return count
}

func CountGoldenSuitCards(cards []Card) int {
	count := 0
	for _, card := range cards {
		count += boolToInt[card.IsGoldenSuit()]
	}
	return count
}

type Tracker struct {
	who   *Player
	count int
}

func (staticticsByPlayer StaticticsByPlayer) calculateMostCardsPlayer() *Player {
	var tracker Tracker = Tracker{nil, 0}
	for player, statictics := range staticticsByPlayer {
		if statictics.CardsTakenCount > tracker.count {
			playerCopy := player // "player" variable is "re used" with a new value in each loop, so a copy is required
			tracker.count = statictics.CardsTakenCount
			tracker.who = &playerCopy
		}
	}
	return tracker.who
}

func (staticticsByPlayer StaticticsByPlayer) calculateSeventiesPlayer() *Player {
	var tracker Tracker = Tracker{nil, 0}
	for player, statictics := range staticticsByPlayer {
		if statictics.SeventiesScore > tracker.count {
			playerCopy := player // "player" variable is "re used" with a new value in each loop, so a copy is required
			tracker.count = statictics.SeventiesScore
			tracker.who = &playerCopy
		}
	}
	return tracker.who
}

func (staticticsByPlayer StaticticsByPlayer) calculateMostGoldCardsPlayer() *Player {
	var tracker Tracker = Tracker{nil, 0}
	for player, statictics := range staticticsByPlayer {
		if statictics.GoldCardsCount > tracker.count {
			playerCopy := player // "player" variable is "re used" with a new value in each loop, so a copy is required
			tracker.count = statictics.GoldCardsCount
			tracker.who = &playerCopy
		}
	}
	return tracker.who
}

type PlayerScoreSummary struct {
	Score      int              `json:"score"`
	Statictics PlayerStatictics `json:"statictics"`
}

type ScoreSummaryByPlayer map[Player]PlayerScoreSummary

func (staticticsByPlayer StaticticsByPlayer) BuildScoreSummaryByPlayer() ScoreSummaryByPlayer {
	scoreSummaryByPlayer := make(ScoreSummaryByPlayer)
	mostCardsPlayer := staticticsByPlayer.calculateMostCardsPlayer()
	seventiesPlayer := staticticsByPlayer.calculateSeventiesPlayer()
	mostGoldCardsPlayer := staticticsByPlayer.calculateMostGoldCardsPlayer()

	//fmt.Printf("mostCardsPlayer = %p\n", mostCardsPlayer)
	//fmt.Printf("seventiesPlayer = %p\n", seventiesPlayer)
	//fmt.Printf("mostGoldCardsPlayer = %p\n", mostGoldCardsPlayer)
	for player, statictics := range staticticsByPlayer {
		score := 0
		if mostCardsPlayer != nil && player == *mostCardsPlayer {
			score += 1
		}
		if seventiesPlayer != nil && player == *seventiesPlayer {
			score += 1
		}
		if mostGoldCardsPlayer != nil && player == *mostGoldCardsPlayer {
			score += 1
		}
		if statictics.HasGoldSeven {
			score += 1
		}
		score += statictics.EscobitasCount
		scoreSummaryByPlayer[player] = PlayerScoreSummary{
			Score:      score,
			Statictics: statictics,
		}
	}
	return scoreSummaryByPlayer
}
