package model

type PlayerStatictics struct {
	CardsTakenCount int
	EscobitasCount  int
	SeventiesScore  int
	HasGoldSeven    bool
}

type StaticticsByPlayer map[Player]PlayerStatictics

var boolToInt map[bool]int = map[bool]int{
	true:  1,
	false: 0,
}

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
	seventiesScore := calculateSeventiesScore(player, match)
	hasGoldSeven := hasGoldenSeven(player, match)

	return PlayerStatictics{
		CardsTakenCount: cardsTakenCount,
		EscobitasCount:  escobitasCount,
		SeventiesScore:  seventiesScore,
		HasGoldSeven:    hasGoldSeven,
	}
}

func countEscobitas(player Player, match Match) int {
	escobitasCount := 0
	for _, action := range match.ActionsByPlayer[player] {
		escobitasCount += boolToInt[action.IsEscobita()]
	}
	return escobitasCount
}

func countCardsTaken(player Player, match Match) int {
	return len(match.MatchCards.PerPlayer[player].Taken)
}

var factors []int = []int{2, 4, 8, 16, 32, 64, 128}

func calculateSeventiesScore(player Player, match Match) int {
	score := 0
	for _, card := range match.MatchCards.PerPlayer[player].Taken {
		if card.Rank <= 7 {
			score += factors[card.Rank]
		}
	}
	return score
}

func hasGoldenSeven(player Player, match Match) bool {
	for _, card := range match.MatchCards.PerPlayer[player].Taken {
		if card.Rank == 7 && card.Suit == GOLD {
			return true
		}
	}
	return false
}

type Tracker struct {
	who *Player
	max int
}

func (staticticsByPlayer StaticticsByPlayer) calculateMostCardsPlayer() Player {
	var tracker Tracker = Tracker{nil, 0}
	for player, statictics := range staticticsByPlayer {
		if statictics.CardsTakenCount > tracker.max {
			tracker.max = statictics.CardsTakenCount
			tracker.who = &player
		}
	}
	return *tracker.who
}

func (staticticsByPlayer StaticticsByPlayer) calculateSeventiesPlayer() Player {
	var tracker Tracker = Tracker{nil, 0}
	for player, statictics := range staticticsByPlayer {
		if statictics.SeventiesScore > tracker.max {
			tracker.max = statictics.SeventiesScore
			tracker.who = &player
		}
	}
	return *tracker.who
}

func sumValues(cards []Card) int {
	total := 0
	for _, card := range cards {
		total += determineValue(card)
	}
	return total
}

func determineValue(card Card) int {
	if card.Rank < 8 {
		return card.Rank
	} else {
		return card.Rank - 2
	}
}
